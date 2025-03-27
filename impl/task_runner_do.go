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
	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"time"
)

// NewTaskRunner creates a TaskRunner instance based on the task type.
func NewTaskRunner(taskName string, task model.Task, taskSupport TaskSupport) (TaskRunner, error) {
	switch t := task.(type) {
	case *model.SetTask:
		return NewSetTaskRunner(taskName, t, taskSupport)
	case *model.RaiseTask:
		return NewRaiseTaskRunner(taskName, t, taskSupport)
	case *model.DoTask:
		return NewDoTaskRunner(t.Do, taskSupport)
	case *model.ForTask:
		return NewForTaskRunner(taskName, t, taskSupport)
	default:
		return nil, fmt.Errorf("unsupported task type '%T' for task '%s'", t, taskName)
	}
}

func NewDoTaskRunner(taskList *model.TaskList, taskSupport TaskSupport) (*DoTaskRunner, error) {
	return &DoTaskRunner{
		TaskList:    taskList,
		TaskSupport: taskSupport,
	}, nil
}

type DoTaskRunner struct {
	TaskList    *model.TaskList
	TaskSupport TaskSupport
}

func (d *DoTaskRunner) Run(input interface{}) (output interface{}, err error) {
	if d.TaskList == nil {
		return input, nil
	}
	return d.runTasks(input, d.TaskList)
}

func (d *DoTaskRunner) GetTaskName() string {
	return ""
}

// runTasks runs all defined tasks sequentially.
func (d *DoTaskRunner) runTasks(input interface{}, tasks *model.TaskList) (output interface{}, err error) {
	output = input
	if tasks == nil {
		return output, nil
	}

	idx := 0
	currentTask := (*tasks)[idx]

	for currentTask != nil {
		if err = d.TaskSupport.SetTaskDef(currentTask); err != nil {
			return nil, err
		}
		if err = d.TaskSupport.SetTaskReferenceFromName(currentTask.Key); err != nil {
			return nil, err
		}

		if shouldRun, err := d.shouldRunTask(input, currentTask); err != nil {
			return output, err
		} else if !shouldRun {
			idx, currentTask = tasks.Next(idx)
			continue
		}

		d.TaskSupport.SetTaskStatus(currentTask.Key, ctx.PendingStatus)
		runner, err := NewTaskRunner(currentTask.Key, currentTask.Task, d.TaskSupport)
		if err != nil {
			return output, err
		}

		d.TaskSupport.SetTaskStatus(currentTask.Key, ctx.RunningStatus)
		if output, err = d.runTask(input, runner, currentTask.Task.GetBase()); err != nil {
			d.TaskSupport.SetTaskStatus(currentTask.Key, ctx.FaultedStatus)
			return output, err
		}

		d.TaskSupport.SetTaskStatus(currentTask.Key, ctx.CompletedStatus)
		input = deepCloneValue(output)
		idx, currentTask = tasks.Next(idx)
	}

	return output, nil
}

func (d *DoTaskRunner) shouldRunTask(input interface{}, task *model.TaskItem) (bool, error) {
	if task.GetBase().If != nil {
		output, err := traverseAndEvaluateBool(task.GetBase().If.String(), input, d.TaskSupport.GetContext())
		if err != nil {
			return false, model.NewErrExpression(err, task.Key)
		}
		return output, nil
	}
	return true, nil
}

// runTask executes an individual task.
func (d *DoTaskRunner) runTask(input interface{}, runner TaskRunner, task *model.TaskBase) (output interface{}, err error) {
	taskName := runner.GetTaskName()

	d.TaskSupport.SetTaskStartedAt(time.Now())
	d.TaskSupport.SetTaskRawInput(input)
	d.TaskSupport.SetTaskName(taskName)

	if task.Input != nil {
		if input, err = d.processTaskInput(task, input, taskName); err != nil {
			return nil, err
		}
	}

	output, err = runner.Run(input)
	if err != nil {
		return nil, err
	}

	d.TaskSupport.SetTaskRawOutput(output)

	if output, err = d.processTaskOutput(task, output, taskName); err != nil {
		return nil, err
	}

	if err = d.processTaskExport(task, output, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

// processTaskInput processes task input validation and transformation.
func (d *DoTaskRunner) processTaskInput(task *model.TaskBase, taskInput interface{}, taskName string) (output interface{}, err error) {
	if task.Input == nil {
		return taskInput, nil
	}

	if err = validateSchema(taskInput, task.Input.Schema, taskName); err != nil {
		return nil, err
	}

	if output, err = traverseAndEvaluate(task.Input.From, taskInput, taskName, d.TaskSupport.GetContext()); err != nil {
		return nil, err
	}

	return output, nil
}

// processTaskOutput processes task output validation and transformation.
func (d *DoTaskRunner) processTaskOutput(task *model.TaskBase, taskOutput interface{}, taskName string) (output interface{}, err error) {
	if task.Output == nil {
		return taskOutput, nil
	}

	if output, err = traverseAndEvaluate(task.Output.As, taskOutput, taskName, d.TaskSupport.GetContext()); err != nil {
		return nil, err
	}

	if err = validateSchema(output, task.Output.Schema, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

func (d *DoTaskRunner) processTaskExport(task *model.TaskBase, taskOutput interface{}, taskName string) (err error) {
	if task.Export == nil {
		return nil
	}

	output, err := traverseAndEvaluate(task.Export.As, taskOutput, taskName, d.TaskSupport.GetContext())
	if err != nil {
		return err
	}

	if err = validateSchema(output, task.Export.Schema, taskName); err != nil {
		return nil
	}

	d.TaskSupport.SetWorkflowInstanceCtx(output)

	return nil
}
