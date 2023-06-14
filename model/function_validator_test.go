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

import "testing"

func TestFunctionStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.Functions = Functions{{
		Name:      "function 1",
		Operation: "http://function/action",
	}}

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "workflow function.name repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Functions = append(model.Functions, model.Functions[0])
				return *model
			},
			Err: `workflow.functions has duplicate "name"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
