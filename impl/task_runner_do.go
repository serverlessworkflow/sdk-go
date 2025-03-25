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

	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskRunner = &DoTaskRunner{}

type TaskSupport interface {
	GetTaskContext() TaskContext
	GetWorkflowDef() *model.Workflow
}

// TODO: refactor to receive a resolver handler instead of the workflow runner

// NewTaskRunner creates a TaskRunner instance based on the task type.
func NewTaskRunner(taskName string, task model.Task, taskSupport TaskSupport) (TaskRunner, error) {
	switch t := task.(type) {
	case *model.SetTask:
		return NewSetTaskRunner(taskName, t)
	case *model.RaiseTask:
		return NewRaiseTaskRunner(taskName, t, taskSupport.GetWorkflowDef())
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
	return d.executeTasks(input, d.TaskList)
}

func (d *DoTaskRunner) GetTaskName() string {
	return ""
}

// executeTasks runs all defined tasks sequentially.
func (d *DoTaskRunner) executeTasks(input interface{}, tasks *model.TaskList) (output interface{}, err error) {
	output = input
	if tasks == nil {
		return output, nil
	}

	idx := 0
	currentTask := (*tasks)[idx]
	ctx := d.TaskSupport.GetTaskContext()

	for currentTask != nil {
		if shouldRun, err := d.shouldRunTask(input, currentTask); err != nil {
			return output, err
		} else if !shouldRun {
			idx, currentTask = tasks.Next(idx)
			continue
		}

		ctx.SetTaskStatus(currentTask.Key, PendingStatus)
		runner, err := NewTaskRunner(currentTask.Key, currentTask.Task, d.TaskSupport)
		if err != nil {
			return output, err
		}

		ctx.SetTaskStatus(currentTask.Key, RunningStatus)
		if output, err = d.runTask(input, runner, currentTask.Task.GetBase()); err != nil {
			ctx.SetTaskStatus(currentTask.Key, FaultedStatus)
			return output, err
		}

		ctx.SetTaskStatus(currentTask.Key, CompletedStatus)
		input = deepCloneValue(output)
		idx, currentTask = tasks.Next(idx)
	}

	return output, nil
}

func (d *DoTaskRunner) shouldRunTask(input interface{}, task *model.TaskItem) (bool, error) {
	if task.GetBase().If != nil {
		output, err := expr.TraverseAndEvaluate(task.GetBase().If.String(), input)
		if err != nil {
			return false, model.NewErrExpression(err, task.Key)
		}
		if result, ok := output.(bool); ok && !result {
			return false, nil
		}
	}
	return true, nil
}

// runTask executes an individual task.
func (d *DoTaskRunner) runTask(input interface{}, runner TaskRunner, task *model.TaskBase) (output interface{}, err error) {
	taskName := runner.GetTaskName()

	if task.Input != nil {
		if input, err = d.processTaskInput(task, input, taskName); err != nil {
			return nil, err
		}
	}

	output, err = runner.Run(input)
	if err != nil {
		return nil, err
	}

	if output, err = d.processTaskOutput(task, output, taskName); err != nil {
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

	if output, err = traverseAndEvaluate(task.Input.From, taskInput, taskName); err != nil {
		return nil, err
	}

	return output, nil
}

// processTaskOutput processes task output validation and transformation.
func (d *DoTaskRunner) processTaskOutput(task *model.TaskBase, taskOutput interface{}, taskName string) (output interface{}, err error) {
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
