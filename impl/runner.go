package impl

import (
	"context"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ WorkflowRunner = &workflowRunnerImpl{}

type WorkflowRunner interface {
	GetWorkflowDef() *model.Workflow
	Run(input interface{}) (output interface{}, err error)
	GetContext() *WorkflowContext
}

func NewDefaultRunner(workflow *model.Workflow) WorkflowRunner {
	wfContext := &WorkflowContext{}
	wfContext.SetStatus(PendingStatus)
	// TODO: based on the workflow definition, the context might change.
	ctx := WithWorkflowContext(context.Background(), wfContext)
	return &workflowRunnerImpl{
		Workflow:  workflow,
		Context:   ctx,
		RunnerCtx: wfContext,
	}
}

type workflowRunnerImpl struct {
	Workflow  *model.Workflow
	Context   context.Context
	RunnerCtx *WorkflowContext
}

func (wr *workflowRunnerImpl) GetContext() *WorkflowContext {
	return wr.RunnerCtx
}

func (wr *workflowRunnerImpl) GetTaskContext() TaskContext {
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

	wr.RunnerCtx.SetInput(input)
	// Run tasks sequentially
	wr.RunnerCtx.SetStatus(RunningStatus)
	doRunner, err := NewDoTaskRunner(wr.Workflow.Do, wr)
	if err != nil {
		return nil, err
	}
	output, err = doRunner.Run(wr.RunnerCtx.GetInput())
	if err != nil {
		return nil, err
	}

	// Process output
	if output, err = wr.processWorkflowOutput(output); err != nil {
		return nil, err
	}

	wr.RunnerCtx.SetOutput(output)
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
			wr.RunnerCtx.SetInstanceCtx(input)
		}
	}

	wr.RunnerCtx.SetInput(input)
	wr.RunnerCtx.SetOutput(input)
	return input, nil
}

// processWorkflowOutput applies output transformations.
func (wr *workflowRunnerImpl) processWorkflowOutput(output interface{}) (interface{}, error) {
	if wr.Workflow.Output != nil {
		var err error
		if output, err = traverseAndEvaluate(wr.Workflow.Output.As, output, "/"); err != nil {
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
