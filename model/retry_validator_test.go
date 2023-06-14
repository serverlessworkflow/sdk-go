// Copyright 2021 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"
)

func TestRetryStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildRetryRef(baseWorkflow, action1, "retry 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Retries[0].Delay = "PT5S"
				model.Retries[0].MaxDelay = "PT5S"
				model.Retries[0].Increment = "PT5S"
				model.Retries[0].Jitter = floatstr.FromString("PT5S")
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Retries[0].Name = ""
				model.States[0].OperationState.Actions[0].RetryRef = ""
				return *model
			},
			Err: `workflow.retries[0].name is required`,
		},
		{
			Desp: "repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Retries = append(model.Retries, model.Retries[0])
				return *model
			},
			Err: `workflow.retries has duplicate "name"`,
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].RetryRef = "invalid retry"
				return *model
			},
			Err: `workflow.states[0].actions[0].retryRef don't exist "invalid retry"`,
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Retries[0].Delay = "P5S"
				model.Retries[0].MaxDelay = "P5S"
				model.Retries[0].Increment = "P5S"
				model.Retries[0].Jitter = floatstr.FromString("P5S")

				return *model
			},
			Err: `workflow.retries[0].delay invalid iso8601 duration "P5S"
workflow.retries[0].maxDelay invalid iso8601 duration "P5S"
workflow.retries[0].increment invalid iso8601 duration "P5S"
workflow.retries[0].jitter invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
