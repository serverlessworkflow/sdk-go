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

	"github.com/serverlessworkflow/sdk-go/v3/model"
)

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

func (s *SetTaskRunner) Run(input interface{}, taskSupport TaskSupport) (output interface{}, err error) {
	setObject := utils.DeepClone(s.Task.Set)
	result, err := expr.TraverseAndEvaluateObj(model.NewObjectOrRuntimeExpr(setObject), input, s.TaskName, taskSupport.GetContext())
	if err != nil {
		return nil, err
	}

	output, ok := result.(map[string]interface{})
	if !ok {
		return nil, model.NewErrRuntime(fmt.Errorf("expected output to be a map[string]interface{}, but got a different type. Got: %v", result), s.TaskName)
	}

	return output, nil
}
