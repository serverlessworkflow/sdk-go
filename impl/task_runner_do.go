// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package impl

import (
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/impl/expr"
	"github.com/serverlessworkflow/sdk-go/v3/impl/utils"
	"time"

	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

// NewTaskRunner creates a TaskRunner instance based on the task type.
func NewTaskRunner(taskName string, task model.Task, workflowDef *model.Workflow) (TaskRunner, error) {
	switch t := task.(type) {
	case *model.SetTask:
		return NewSetTaskRunner(taskName, t)
	case *model.RaiseTask:
		return NewRaiseTaskRunner(taskName, t, workflowDef)
	case *model.DoTask:
		return NewDoTaskRunner(t.Do)
	case *model.ForTask:
		return NewForTaskRunner(taskName, t)
	case *model.CallHTTP:
		return NewCallHttpRunner(taskName, t)
	case *model.ForkTask:
		return NewForkTaskRunner(taskName, t, workflowDef)
	default:
		return nil, fmt.Errorf("unsupported task type '%T' for task '%s'", t, taskName)
	}
}

func NewDoTaskRunner(taskList *model.TaskList) (*DoTaskRunner, error) {
	return &DoTaskRunner{
		TaskList: taskList,
	}, nil
}

type DoTaskRunner struct {
	TaskList *model.TaskList
}

func (d *DoTaskRunner) Run(input interface{}, taskSupport TaskSupport) (output interface{}, err error) {
	if d.TaskList == nil {
		return input, nil
	}
	return d.runTasks(input, taskSupport)
}

func (d *DoTaskRunner) GetTaskName() string {
	return ""
}

// runTasks runs all defined tasks sequentially.
func (d *DoTaskRunner) runTasks(input interface{}, taskSupport TaskSupport) (output interface{}, err error) {
	output = input
	if d.TaskList == nil {
		return output, nil
	}

	idx := 0
	currentTask := (*d.TaskList)[idx]

	for currentTask != nil {
		if err = taskSupport.SetTaskDef(currentTask); err != nil {
			return nil, err
		}
		if err = taskSupport.SetTaskReferenceFromName(currentTask.Key); err != nil {
			return nil, err
		}

		if shouldRun, err := d.shouldRunTask(input, taskSupport, currentTask); err != nil {
			return output, err
		} else if !shouldRun {
			idx, currentTask = d.TaskList.Next(idx)
			continue
		}

		taskSupport.SetTaskStatus(currentTask.Key, ctx.PendingStatus)

		// Check if this task is a SwitchTask and handle it
		if switchTask, ok := currentTask.Task.(*model.SwitchTask); ok {
			flowDirective, err := d.evaluateSwitchTask(input, taskSupport, currentTask.Key, switchTask)
			if err != nil {
				taskSupport.SetTaskStatus(currentTask.Key, ctx.FaultedStatus)
				return output, err
			}
			taskSupport.SetTaskStatus(currentTask.Key, ctx.CompletedStatus)

			// Process FlowDirective: update idx/currentTask accordingly
			idx, currentTask = d.TaskList.KeyAndIndex(flowDirective.Value)
			if currentTask == nil {
				return nil, fmt.Errorf("flow directive target '%s' not found", flowDirective.Value)
			}
			continue
		}

		runner, err := NewTaskRunner(currentTask.Key, currentTask.Task, taskSupport.GetWorkflowDef())
		if err != nil {
			return output, err
		}

		taskSupport.SetTaskStatus(currentTask.Key, ctx.RunningStatus)
		if output, err = d.runTask(input, taskSupport, runner, currentTask.Task.GetBase()); err != nil {
			taskSupport.SetTaskStatus(currentTask.Key, ctx.FaultedStatus)
			return output, err
		}

		taskSupport.SetTaskStatus(currentTask.Key, ctx.CompletedStatus)
		input = utils.DeepCloneValue(output)
		idx, currentTask = d.TaskList.Next(idx)
	}

	return output, nil
}

func (d *DoTaskRunner) shouldRunTask(input interface{}, taskSupport TaskSupport, task *model.TaskItem) (bool, error) {
	if task.GetBase().If != nil {
		output, err := expr.TraverseAndEvaluateBool(task.GetBase().If.String(), input, taskSupport.GetContext())
		if err != nil {
			return false, model.NewErrExpression(err, task.Key)
		}
		return output, nil
	}
	return true, nil
}

func (d *DoTaskRunner) evaluateSwitchTask(input interface{}, taskSupport TaskSupport, taskKey string, switchTask *model.SwitchTask) (*model.FlowDirective, error) {
	var defaultThen *model.FlowDirective
	for _, switchItem := range switchTask.Switch {
		for _, switchCase := range switchItem {
			if switchCase.When == nil {
				defaultThen = switchCase.Then
				continue
			}
			result, err := expr.TraverseAndEvaluateBool(model.NormalizeExpr(switchCase.When.String()), input, taskSupport.GetContext())
			if err != nil {
				return nil, model.NewErrExpression(err, taskKey)
			}
			if result {
				if switchCase.Then == nil {
					return nil, model.NewErrExpression(fmt.Errorf("missing 'then' directive in matched switch case"), taskKey)
				}
				return switchCase.Then, nil
			}
		}
	}
	if defaultThen != nil {
		return defaultThen, nil
	}
	return nil, model.NewErrExpression(fmt.Errorf("no matching switch case"), taskKey)
}

// runTask executes an individual task.
func (d *DoTaskRunner) runTask(input interface{}, taskSupport TaskSupport, runner TaskRunner, task *model.TaskBase) (output interface{}, err error) {
	taskName := runner.GetTaskName()

	taskSupport.SetTaskStartedAt(time.Now())
	taskSupport.SetTaskRawInput(input)
	taskSupport.SetTaskName(taskName)

	if task.Input != nil {
		if input, err = d.processTaskInput(task, input, taskSupport, taskName); err != nil {
			return nil, err
		}
	}

	output, err = runner.Run(input, taskSupport)
	if err != nil {
		return nil, err
	}

	taskSupport.SetTaskRawOutput(output)

	if output, err = d.processTaskOutput(task, output, taskSupport, taskName); err != nil {
		return nil, err
	}

	if err = d.processTaskExport(task, output, taskSupport, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

// processTaskInput processes task input validation and transformation.
func (d *DoTaskRunner) processTaskInput(task *model.TaskBase, taskInput interface{}, taskSupport TaskSupport, taskName string) (output interface{}, err error) {
	if task.Input == nil {
		return taskInput, nil
	}

	if err = utils.ValidateSchema(taskInput, task.Input.Schema, taskName); err != nil {
		return nil, err
	}

	if output, err = expr.TraverseAndEvaluateObj(task.Input.From, taskInput, taskName, taskSupport.GetContext()); err != nil {
		return nil, err
	}

	return output, nil
}

// processTaskOutput processes task output validation and transformation.
func (d *DoTaskRunner) processTaskOutput(task *model.TaskBase, taskOutput interface{}, taskSupport TaskSupport, taskName string) (output interface{}, err error) {
	if task.Output == nil {
		return taskOutput, nil
	}

	if output, err = expr.TraverseAndEvaluateObj(task.Output.As, taskOutput, taskName, taskSupport.GetContext()); err != nil {
		return nil, err
	}

	if err = utils.ValidateSchema(output, task.Output.Schema, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

func (d *DoTaskRunner) processTaskExport(task *model.TaskBase, taskOutput interface{}, taskSupport TaskSupport, taskName string) (err error) {
	if task.Export == nil {
		return nil
	}

	output, err := expr.TraverseAndEvaluateObj(task.Export.As, taskOutput, taskName, taskSupport.GetContext())
	if err != nil {
		return err
	}

	if err = utils.ValidateSchema(output, task.Export.Schema, taskName); err != nil {
		return nil
	}

	taskSupport.SetWorkflowInstanceCtx(output)

	return nil
}
