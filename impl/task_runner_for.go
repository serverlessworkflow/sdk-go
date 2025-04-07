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
	"reflect"
	"strings"

	"github.com/serverlessworkflow/sdk-go/v3/impl/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

const (
	forTaskDefaultEach = "$item"
	forTaskDefaultAt   = "$index"
)

func NewForTaskRunner(taskName string, task *model.ForTask, workflowDef *model.Workflow) (*ForTaskRunner, error) {
	if task == nil || task.Do == nil {
		return nil, model.NewErrValidation(fmt.Errorf("invalid For task %s", taskName), taskName)
	}

	doRunner, err := NewDoTaskRunner(task.Do, workflowDef)
	if err != nil {
		return nil, err
	}

	return &ForTaskRunner{
		Task:     task,
		TaskName: taskName,
		DoRunner: doRunner,
	}, nil
}

type ForTaskRunner struct {
	Task     *model.ForTask
	TaskName string
	DoRunner *DoTaskRunner
}

func (f *ForTaskRunner) Run(input interface{}, taskSupport TaskSupport) (interface{}, error) {
	defer func() {
		// clear local variables
		taskSupport.RemoveLocalExprVars(f.Task.For.Each, f.Task.For.At)
	}()
	f.sanitizeFor()
	in, err := expr.TraverseAndEvaluate(f.Task.For.In, input, taskSupport.GetContext())
	if err != nil {
		return nil, err
	}

	forOutput := input
	rv := reflect.ValueOf(in)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()

			if forOutput, err = f.processForItem(i, item, taskSupport, forOutput); err != nil {
				return nil, err
			}
			if f.Task.While != "" {
				whileIsTrue, err := traverseAndEvaluateBool(f.Task.While, forOutput, taskSupport.GetContext())
				if err != nil {
					return nil, err
				}
				if !whileIsTrue {
					break
				}
			}
		}
	case reflect.Invalid:
		return input, nil
	default:
		if forOutput, err = f.processForItem(0, in, taskSupport, forOutput); err != nil {
			return nil, err
		}
	}

	return forOutput, nil
}

func (f *ForTaskRunner) processForItem(idx int, item interface{}, taskSupport TaskSupport, forOutput interface{}) (interface{}, error) {
	forVars := map[string]interface{}{
		f.Task.For.At:   idx,
		f.Task.For.Each: item,
	}
	// Instead of Set, we Add since other tasks in this very same context might be adding variables to the context
	taskSupport.AddLocalExprVars(forVars)
	// output from previous iterations are merged together
	var err error
	forOutput, err = f.DoRunner.Run(forOutput, taskSupport)
	if err != nil {
		return nil, err
	}

	return forOutput, nil
}

func (f *ForTaskRunner) sanitizeFor() {
	f.Task.For.Each = strings.TrimSpace(f.Task.For.Each)
	f.Task.For.At = strings.TrimSpace(f.Task.For.At)

	if f.Task.For.Each == "" {
		f.Task.For.Each = forTaskDefaultEach
	}
	if f.Task.For.At == "" {
		f.Task.For.At = forTaskDefaultAt
	}

	if !strings.HasPrefix(f.Task.For.Each, "$") {
		f.Task.For.Each = "$" + f.Task.For.Each
	}
	if !strings.HasPrefix(f.Task.For.At, "$") {
		f.Task.For.At = "$" + f.Task.For.At
	}
}

func (f *ForTaskRunner) GetTaskName() string {
	return f.TaskName
}
