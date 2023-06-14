// Copyright 2022 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"testing"
)

func TestBaseStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "repeat name",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States = []State{model.States[0], model.States[0]}
				return *model
			},
			Err: `workflow.states has duplicate "name"`,
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.CompensatedBy = "invalid state compensate by"
				return *model
			},
			Err: `workflow.states[0].compensatedBy don't exist "invalid state compensate by"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.States = make(States, 0, 2)

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	operationState2 := buildOperationState(baseWorkflow, "next state")
	buildEndByState(operationState2, true, false)
	action2 := buildActionByOperationState(operationState2, "action 2")
	buildFunctionRef(baseWorkflow, action2, "function 2")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				return *baseWorkflow.DeepCopy()
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.End = nil
				return *model
			},
			Err: `workflow.states[0].transition is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildTransitionByState(&model.States[0], &model.States[1], false)

				return *model
			},
			Err: `workflow.states[0].transition exclusive`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
