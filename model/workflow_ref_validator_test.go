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

import "testing"

func TestWorkflowRefStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(&baseWorkflow.States[0], true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	baseWorkflow.States[0].OperationState.Actions[0].FunctionRef = nil
	baseWorkflow.States[0].OperationState.Actions[0].SubFlowRef = &WorkflowRef{
		WorkflowID:       "workflowID",
		Invoke:           InvokeKindSync,
		OnParentComplete: OnParentCompleteTypeTerminate,
	}

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
				model.States[0].OperationState.Actions[0].SubFlowRef.WorkflowID = ""
				model.States[0].OperationState.Actions[0].SubFlowRef.Invoke = ""
				model.States[0].OperationState.Actions[0].SubFlowRef.OnParentComplete = ""
				return *model
			},
			Err: `workflow.states[0].actions[0].subFlowRef.workflowID is required
workflow.states[0].actions[0].subFlowRef.invoke is required
workflow.states[0].actions[0].subFlowRef.onParentComplete is required`,
		},
		{
			Desp: "invalid type",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].SubFlowRef.Invoke = "invalid invoce"
				model.States[0].OperationState.Actions[0].SubFlowRef.OnParentComplete = "invalid parent complete"
				return *model
			},
			Err: `workflow.states[0].actions[0].subFlowRef.invoke need by one of [sync async]
workflow.states[0].actions[0].subFlowRef.onParentComplete need by one of [terminate continue]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
