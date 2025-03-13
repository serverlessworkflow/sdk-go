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

	"github.com/serverlessworkflow/sdk-go/v3/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskRunner = &SetTaskRunner{}
var _ TaskRunner = &RaiseTaskRunner{}
var _ TaskRunner = &ForTaskRunner{}

type TaskRunner interface {
	Run(input interface{}) (interface{}, error)
	GetTaskName() string
}

func NewSetTaskRunner(taskName string, task *model.SetTask) (*SetTaskRunner, error) {
	if task == nil || task.Set == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no set configuration provided for SetTask %s", taskName), taskName)
	}
	return &SetTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

type SetTaskRunner struct {
	Task     *model.SetTask
	TaskName string
}

func (s *SetTaskRunner) GetTaskName() string {
	return s.TaskName
}

func (s *SetTaskRunner) Run(input interface{}) (output interface{}, err error) {
	setObject := deepClone(s.Task.Set)
	result, err := expr.TraverseAndEvaluate(setObject, input)
	if err != nil {
		return nil, model.NewErrExpression(err, s.TaskName)
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		return nil, model.NewErrRuntime(fmt.Errorf("expected output to be a map[string]interface{}, but got a different type. Got: %v", result), s.TaskName)
	}

	return output, nil
}

func NewRaiseTaskRunner(taskName string, task *model.RaiseTask, workflowDef *model.Workflow) (*RaiseTaskRunner, error) {
	if err := resolveErrorDefinition(task, workflowDef); err != nil {
		return nil, err
	}
	if task.Raise.Error.Definition == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no raise configuration provided for RaiseTask %s", taskName), taskName)
	}
	return &RaiseTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

// TODO: can e refactored to a definition resolver callable from the context
func resolveErrorDefinition(t *model.RaiseTask, workflowDef *model.Workflow) error {
	if workflowDef != nil && t.Raise.Error.Ref != nil {
		notFoundErr := model.NewErrValidation(fmt.Errorf("%v error definition not found in 'uses'", t.Raise.Error.Ref), "")
		if workflowDef.Use != nil && workflowDef.Use.Errors != nil {
			definition, ok := workflowDef.Use.Errors[*t.Raise.Error.Ref]
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

type RaiseTaskRunner struct {
	Task     *model.RaiseTask
	TaskName string
}

var raiseErrFuncMapping = map[string]func(error, string) *model.Error{
	model.ErrorTypeAuthentication: model.NewErrAuthentication,
	model.ErrorTypeValidation:     model.NewErrValidation,
	model.ErrorTypeCommunication:  model.NewErrCommunication,
	model.ErrorTypeAuthorization:  model.NewErrAuthorization,
	model.ErrorTypeConfiguration:  model.NewErrConfiguration,
	model.ErrorTypeExpression:     model.NewErrExpression,
	model.ErrorTypeRuntime:        model.NewErrRuntime,
	model.ErrorTypeTimeout:        model.NewErrTimeout,
}

func (r *RaiseTaskRunner) Run(input interface{}) (output interface{}, err error) {
	output = input
	// TODO: make this an external func so we can call it after getting the reference? Or we can get the reference from the workflow definition
	var detailResult interface{}
	detailResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Detail.AsObjectOrRuntimeExpr(), input, r.TaskName)
	if err != nil {
		return nil, err
	}

	var titleResult interface{}
	titleResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Title.AsObjectOrRuntimeExpr(), input, r.TaskName)
	if err != nil {
		return nil, err
	}

	instance := &model.JsonPointerOrRuntimeExpression{Value: r.TaskName}

	var raiseErr *model.Error
	if raiseErrF, ok := raiseErrFuncMapping[r.Task.Raise.Error.Definition.Type.String()]; ok {
		raiseErr = raiseErrF(fmt.Errorf("%v", detailResult), instance.String())
	} else {
		raiseErr = r.Task.Raise.Error.Definition
		raiseErr.Detail = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", detailResult))
		raiseErr.Instance = instance
	}

	raiseErr.Title = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", titleResult))
	err = raiseErr

	return output, err
}

func (r *RaiseTaskRunner) GetTaskName() string {
	return r.TaskName
}

func NewForTaskRunner(taskName string, task *model.ForTask, taskSupport TaskSupport) (*ForTaskRunner, error) {
	if task == nil || task.Do == nil {
		return nil, model.NewErrValidation(fmt.Errorf("invalid For task %s", taskName), taskName)
	}

	doRunner, err := NewDoTaskRunner(task.Do, taskSupport)
	if err != nil {
		return nil, err
	}

	return &ForTaskRunner{
		Task:     task,
		TaskName: taskName,
		DoRunner: doRunner,
	}, nil
}

const (
	forTaskDefaultEach = "$item"
	forTaskDefaultAt   = "$index"
)

type ForTaskRunner struct {
	Task     *model.ForTask
	TaskName string
	DoRunner *DoTaskRunner
}

func (f *ForTaskRunner) Run(input interface{}) (interface{}, error) {
	f.sanitizeFor()
	in, err := expr.TraverseAndEvaluate(f.Task.For.In, input)
	if err != nil {
		return nil, err
	}

	var forOutput interface{}
	rv := reflect.ValueOf(in)
	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < rv.Len(); i++ {
			item := rv.Index(i).Interface()

			if forOutput, err = f.processForItem(i, item, forOutput); err != nil {
				return nil, err
			}
		}
	case reflect.Invalid:
		return input, nil
	default:
		if forOutput, err = f.processForItem(0, in, forOutput); err != nil {
			return nil, err
		}
	}

	return forOutput, nil
}

func (f *ForTaskRunner) processForItem(idx int, item interface{}, forOutput interface{}) (interface{}, error) {
	forInput := map[string]interface{}{
		f.Task.For.At:   idx,
		f.Task.For.Each: item,
	}
	if forOutput != nil {
		if outputMap, ok := forOutput.(map[string]interface{}); ok {
			for key, value := range outputMap {
				forInput[key] = value
			}
		} else {
			return nil, fmt.Errorf("task %s item %s at index %d returned a non-json object, impossible to merge context", f.TaskName, f.Task.For.Each, idx)
		}
	}
	var err error
	forOutput, err = f.DoRunner.Run(forInput)
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
