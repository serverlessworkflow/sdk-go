package impl

import (
	"context"
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ WorkflowRunner = &workflowRunnerImpl{}

type WorkflowRunner interface {
	GetWorkflowDef() *model.Workflow
	Run(input interface{}) (output interface{}, err error)
	GetContext() *WorkflowRunnerContext
}

func NewDefaultRunner(workflow *model.Workflow) WorkflowRunner {
	wfContext := &WorkflowRunnerContext{}
	wfContext.SetStatus(PendingStatus)
	// TODO: based on the workflow definition, the context might change.
	ctx := WithRunnerContext(context.Background(), wfContext)
	return &workflowRunnerImpl{
		Workflow:  workflow,
		Context:   ctx,
		RunnerCtx: wfContext,
	}
}

type workflowRunnerImpl struct {
	Workflow  *model.Workflow
	Context   context.Context
	RunnerCtx *WorkflowRunnerContext
}

func (wr *workflowRunnerImpl) GetContext() *WorkflowRunnerContext {
	return wr.RunnerCtx
}

func (wr *workflowRunnerImpl) GetWorkflowDef() *model.Workflow {
	return wr.Workflow
}

// Run executes the workflow synchronously.
func (wr *workflowRunnerImpl) Run(input interface{}) (output interface{}, err error) {
	defer func() {
		if err != nil {
			wr.RunnerCtx.SetStatus(FaultedStatus)
			err = wr.wrapWorkflowError(err, "/")
		}
	}()

	// Process input
	if input, err = wr.processWorkflowInput(input); err != nil {
		return nil, err
	}

	// Run tasks sequentially
	wr.RunnerCtx.SetStatus(RunningStatus)
	if err = wr.executeTasks(wr.Workflow.Do); err != nil {
		return nil, err
	}

	output = wr.RunnerCtx.GetOutput()

	// Process output
	if output, err = wr.processWorkflowOutput(output); err != nil {
		return nil, err
	}

	wr.RunnerCtx.SetStatus(CompletedStatus)
	return output, nil
}

// wrapWorkflowError ensures workflow errors have a proper instance reference.
func (wr *workflowRunnerImpl) wrapWorkflowError(err error, taskName string) error {
	if knownErr := model.AsError(err); knownErr != nil {
		return knownErr.WithInstanceRef(wr.Workflow, taskName)
	}
	return model.NewErrRuntime(err, taskName)
}

// processWorkflowInput validates and transforms input if needed.
func (wr *workflowRunnerImpl) processWorkflowInput(input interface{}) (interface{}, error) {
	if wr.Workflow.Input != nil {
		var err error
		if err = validateSchema(input, wr.Workflow.Input.Schema, "/"); err != nil {
			return nil, err
		}

		if wr.Workflow.Input.From != nil {
			if input, err = traverseAndEvaluate(wr.Workflow.Input.From, input, "/"); err != nil {
				return nil, err
			}
			wr.RunnerCtx.SetWorkflowCtx(input)
		}
	}

	wr.RunnerCtx.SetInput(input)
	wr.RunnerCtx.SetOutput(input)
	return input, nil
}

// executeTasks runs all defined tasks sequentially.
func (wr *workflowRunnerImpl) executeTasks(tasks *model.TaskList) error {
	if tasks == nil {
		return nil
	}

	idx := 0
	currentTask := (*tasks)[idx]

	for currentTask != nil {
		wr.RunnerCtx.SetInput(wr.RunnerCtx.GetOutput())
		if shouldRun, err := wr.shouldRunTask(currentTask); err != nil {
			return err
		} else if !shouldRun {
			wr.RunnerCtx.SetOutput(wr.RunnerCtx.GetInput())
			idx, currentTask = tasks.Next(idx)
			continue
		}

		wr.RunnerCtx.SetTaskStatus(currentTask.Key, PendingStatus)
		runner, err := NewTaskRunner(currentTask.Key, currentTask.Task, wr)
		if err != nil {
			return err
		}

		wr.RunnerCtx.SetTaskStatus(currentTask.Key, RunningStatus)
		var output interface{}
		if output, err = wr.runTask(runner, currentTask.Task.GetBase()); err != nil {
			wr.RunnerCtx.SetTaskStatus(currentTask.Key, FaultedStatus)
			return err
		}
		// TODO: make sure that `output` is a map[string]interface{}, so compatible to JSON traversal.

		wr.RunnerCtx.SetTaskStatus(currentTask.Key, CompletedStatus)
		wr.RunnerCtx.SetOutput(output)

		idx, currentTask = tasks.Next(idx)
	}

	return nil
}

func (wr *workflowRunnerImpl) shouldRunTask(task *model.TaskItem) (bool, error) {
	if task.GetBase().If != nil {
		output, err := expr.TraverseAndEvaluate(task.GetBase().If.String(), wr.RunnerCtx.GetInput())
		if err != nil {
			return false, model.NewErrExpression(err, task.Key)
		}
		if result, ok := output.(bool); ok && !result {
			return false, nil
		}
	}
	return true, nil
}

// processWorkflowOutput applies output transformations.
func (wr *workflowRunnerImpl) processWorkflowOutput(output interface{}) (interface{}, error) {
	if wr.Workflow.Output != nil {
		var err error
		if output, err = traverseAndEvaluate(wr.Workflow.Output.As, wr.RunnerCtx.GetOutput(), "/"); err != nil {
			return nil, err
		}

		if err = validateSchema(output, wr.Workflow.Output.Schema, "/"); err != nil {
			return nil, err
		}
	}

	wr.RunnerCtx.SetOutput(output)
	return output, nil
}

// ----------------- Task funcs ------------------- //

// TODO: refactor to receive a resolver handler instead of the workflow runner

// NewTaskRunner creates a TaskRunner instance based on the task type.
func NewTaskRunner(taskName string, task model.Task, wr *workflowRunnerImpl) (TaskRunner, error) {
	switch t := task.(type) {
	case *model.SetTask:
		return NewSetTaskRunner(taskName, t)
	case *model.RaiseTask:
		if err := wr.resolveErrorDefinition(t); err != nil {
			return nil, err
		}
		return NewRaiseTaskRunner(taskName, t)
	default:
		return nil, fmt.Errorf("unsupported task type '%T' for task '%s'", t, taskName)
	}
}

// TODO: can e refactored to a definition resolver callable from the context
func (wr *workflowRunnerImpl) resolveErrorDefinition(t *model.RaiseTask) error {
	if t.Raise.Error.Ref != nil {
		notFoundErr := model.NewErrValidation(fmt.Errorf("%v error definition not found in 'uses'", t.Raise.Error.Ref), "")
		if wr.Workflow.Use != nil && wr.Workflow.Use.Errors != nil {
			definition, ok := wr.Workflow.Use.Errors[*t.Raise.Error.Ref]
			if !ok {
				return notFoundErr
			}
			t.Raise.Error.Definition = definition
			return nil
		}
		return notFoundErr
	}
	return nil
}

// runTask executes an individual task.
func (wr *workflowRunnerImpl) runTask(runner TaskRunner, task *model.TaskBase) (output interface{}, err error) {
	taskInput := wr.RunnerCtx.GetInput()
	taskName := runner.GetTaskName()

	defer func() {
		if err != nil {
			err = wr.wrapWorkflowError(err, taskName)
		}
	}()

	if task.Input != nil {
		if taskInput, err = wr.validateAndEvaluateTaskInput(task, taskInput, taskName); err != nil {
			return nil, err
		}
	}

	output, err = runner.Run(taskInput)
	if err != nil {
		return nil, err
	}

	if output, err = wr.validateAndEvaluateTaskOutput(task, output, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

// validateAndEvaluateTaskInput processes task input validation and transformation.
func (wr *workflowRunnerImpl) validateAndEvaluateTaskInput(task *model.TaskBase, taskInput interface{}, taskName string) (output interface{}, err error) {
	if task.Input == nil {
		return taskInput, nil
	}

	if err = validateSchema(taskInput, task.Input.Schema, taskName); err != nil {
		return nil, err
	}

	if output, err = traverseAndEvaluate(task.Input.From, taskInput, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

// validateAndEvaluateTaskOutput processes task output validation and transformation.
func (wr *workflowRunnerImpl) validateAndEvaluateTaskOutput(task *model.TaskBase, taskOutput interface{}, taskName string) (output interface{}, err error) {
	if task.Output == nil {
		return taskOutput, nil
	}

	if output, err = traverseAndEvaluate(task.Output.As, taskOutput, taskName); err != nil {
		return nil, err
	}

	if err = validateSchema(output, task.Output.Schema, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

func validateSchema(data interface{}, schema *model.Schema, taskName string) error {
	if schema != nil {
		if err := ValidateJSONSchema(data, schema); err != nil {
			return model.NewErrValidation(err, taskName)
		}
	}
	return nil
}

func traverseAndEvaluate(runtimeExpr *model.ObjectOrRuntimeExpr, input interface{}, taskName string) (output interface{}, err error) {
	if runtimeExpr == nil {
		return input, nil
	}
	output, err = expr.TraverseAndEvaluate(runtimeExpr.AsStringOrMap(), input)
	if err != nil {
		return nil, model.NewErrExpression(err, taskName)
	}
	return output, nil
}
