// Copyright 2022 The Serverless Workflow Specification Authors
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
)

func buildOperationState(workflow *Workflow, name string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeOperation,
		},
		OperationState: &OperationState{
			ActionMode: ActionModeSequential,
		},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func buildOperationStateTimeout(state *State) *OperationStateTimeout {
	state.OperationState.Timeouts = &OperationStateTimeout{
		ActionExecTimeout: "PT5S",
	}
	return state.OperationState.Timeouts
}

func TestOperationStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions = []Action{}
				return *model
			},
			Err: ``,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.ActionMode = ActionModeParallel + "invalid"
				return *model
			},
			Err: `workflow.states[0].actionMode need by one of [sequential parallel]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestOperationStateTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	operationStateTimeout := buildOperationStateTimeout(operationState)
	buildStateExecTimeoutByOperationStateTimeout(operationStateTimeout)

	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "omitempty",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Timeouts.ActionExecTimeout = ""
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Timeouts.ActionExecTimeout = "P5S"
				return *model
			},
			Err: `workflow.states[0].timeouts.actionExecTimeout invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
