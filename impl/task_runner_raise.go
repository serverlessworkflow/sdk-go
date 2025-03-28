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
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

func NewRaiseTaskRunner(taskName string, task *model.RaiseTask, taskSupport TaskSupport) (*RaiseTaskRunner, error) {
	if err := resolveErrorDefinition(task, taskSupport.GetWorkflowDef()); err != nil {
		return nil, err
	}

	if task.Raise.Error.Definition == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no raise configuration provided for RaiseTask %s", taskName), taskName)
	}
	return &RaiseTaskRunner{
		Task:        task,
		TaskName:    taskName,
		TaskSupport: taskSupport,
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
	Task        *model.RaiseTask
	TaskName    string
	TaskSupport TaskSupport
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
	detailResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Detail.AsObjectOrRuntimeExpr(), input, r.TaskName, r.TaskSupport.GetContext())
	if err != nil {
		return nil, err
	}

	var titleResult interface{}
	titleResult, err = traverseAndEvaluate(r.Task.Raise.Error.Definition.Title.AsObjectOrRuntimeExpr(), input, r.TaskName, r.TaskSupport.GetContext())
	if err != nil {
		return nil, err
	}

	instance := r.TaskSupport.GetTaskReference()

	var raiseErr *model.Error
	if raiseErrF, ok := raiseErrFuncMapping[r.Task.Raise.Error.Definition.Type.String()]; ok {
		raiseErr = raiseErrF(fmt.Errorf("%v", detailResult), instance)
	} else {
		raiseErr = r.Task.Raise.Error.Definition
		raiseErr.Detail = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", detailResult))
		raiseErr.Instance = &model.JsonPointerOrRuntimeExpression{Value: instance}
	}

	raiseErr.Title = model.NewStringOrRuntimeExpr(fmt.Sprintf("%v", titleResult))
	err = raiseErr

	return output, err
}

func (r *RaiseTaskRunner) GetTaskName() string {
	return r.TaskName
}
