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

func buildActionByOperationState(state *State, name string) *Action {
	action := Action{
		Name: name,
	}

	state.OperationState.Actions = append(state.OperationState.Actions, action)
	return &state.OperationState.Actions[len(state.OperationState.Actions)-1]
}

func buildActionByForEachState(state *State, name string) *Action {
	action := Action{
		Name: name,
	}

	state.ForEachState.Actions = append(state.ForEachState.Actions, action)
	return &state.ForEachState.Actions[len(state.ForEachState.Actions)-1]
}

func buildActionByBranch(branch *Branch, name string) *Action {
	action := Action{
		Name: name,
	}

	branch.Actions = append(branch.Actions, action)
	return &branch.Actions[len(branch.Actions)-1]
}

func buildFunctionRef(workflow *Workflow, action *Action, name string) (*FunctionRef, *Function) {
	function := Function{
		Name:      name,
		Operation: "http://function/function_name",
		Type:      FunctionTypeREST,
	}

	functionRef := FunctionRef{
		RefName: name,
		Invoke:  InvokeKindSync,
	}
	action.FunctionRef = &functionRef

	workflow.Functions = append(workflow.Functions, function)
	return &functionRef, &function
}

func buildRetryRef(workflow *Workflow, action *Action, name string) {
	retry := Retry{
		Name: name,
	}

	workflow.Retries = append(workflow.Retries, retry)
	action.RetryRef = name
}

func buildSleep(action *Action) *Sleep {
	action.Sleep = &Sleep{
		Before: "PT5S",
		After:  "PT5S",
	}
	return action.Sleep
}

func TestActionStructLevelValidation(t *testing.T) {
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
			Desp: "require_without",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].FunctionRef = nil
				return *model
			},
			Err: `Key: 'Workflow.States[0].OperationState.Actions[0].FunctionRef' Error:Field validation for 'FunctionRef' failed on the 'required_without' tag`,
		},
		{
			Desp: "exclude",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildEventRef(model, &model.States[0].OperationState.Actions[0], "event 1", "event2")
				return *model
			},
			Err: `workflow.states[0].actions[0].functionRef exclusive
workflow.states[0].actions[0].eventRef exclusive
workflow.states[0].actions[0].subFlowRef exclusive`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].FunctionRef.Invoke = InvokeKindSync + "invalid"
				return *model
			},
			Err: `workflow.states[0].actions[0].functionRef.invoke need by one of [sync async]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestFunctionRefStructLevelValidation(t *testing.T) {
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
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].FunctionRef.RefName = "invalid function"
				return *model
			},
			Err: `workflow.states[0].actions[0].functionRef.refName don't exist "invalid function"`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestSleepStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildSleep(action1)
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
				model.States[0].OperationState.Actions[0].Sleep.Before = ""
				model.States[0].OperationState.Actions[0].Sleep.After = ""
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].Sleep.Before = "P5S"
				model.States[0].OperationState.Actions[0].Sleep.After = "P5S"
				return *model
			},
			Err: `workflow.states[0].actions[0].sleep.before invalid iso8601 duration "P5S"
workflow.states[0].actions[0].sleep.after invalid iso8601 duration "P5S"`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}
