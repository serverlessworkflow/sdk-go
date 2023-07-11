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

func TestActionStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase[Action]{
		{
			Desp:  "action empty",
			Model: Action{},
			Err: `Key: 'Action.FunctionRef' Error:Field validation for 'FunctionRef' failed on the 'exclusive' tag
Key: 'Action.EventRef' Error:Field validation for 'EventRef' failed on the 'exclusive' tag
Key: 'Action.SubFlowRef' Error:Field validation for 'SubFlowRef' failed on the 'exclusive' tag`,
		},
		{
			Desp: "action functionRef and eventRef",
			Model: Action{
				FunctionRef: &FunctionRef{
					RefName: "function 1",
					Invoke:  InvokeKindSync,
				},
				EventRef: &EventRef{
					TriggerEventRef: "event 1",
					ResultEventRef:  "event 1",
					Invoke:          InvokeKindAsync,
				},
			},
			Err: `Key: 'Action.FunctionRef' Error:Field validation for 'FunctionRef' failed on the 'exclusive' tag
Key: 'Action.EventRef' Error:Field validation for 'EventRef' failed on the 'exclusive' tag
Key: 'Action.SubFlowRef' Error:Field validation for 'SubFlowRef' failed on the 'exclusive' tag`,
		},
		{
			Desp: "action eventRef and subFlowRef",
			Model: Action{
				EventRef: &EventRef{
					TriggerEventRef: "event 1",
					ResultEventRef:  "event 1",
					Invoke:          InvokeKindAsync,
				},
				SubFlowRef: &WorkflowRef{
					WorkflowID:       "teste",
					Invoke:           InvokeKindAsync,
					OnParentComplete: OnParentCompleteTypeTerminate,
				},
			},
			Err: `Key: 'Action.FunctionRef' Error:Field validation for 'FunctionRef' failed on the 'exclusive' tag
Key: 'Action.EventRef' Error:Field validation for 'EventRef' failed on the 'exclusive' tag
Key: 'Action.SubFlowRef' Error:Field validation for 'SubFlowRef' failed on the 'exclusive' tag`,
		},
		{
			Desp: "action functionRef",
			Model: Action{
				FunctionRef: &FunctionRef{
					RefName: "function 1",
					Invoke:  InvokeKindSync,
				},
			},
			Err: ``,
		},
		{
			Desp: "action eventRef",
			Model: Action{
				EventRef: &EventRef{
					TriggerEventRef: "event 1",
					ResultEventRef:  "event 1",
					Invoke:          InvokeKindAsync,
				},
			},
			Err: ``,
		},
		{
			Desp: "action subFlowRef",
			Model: Action{
				SubFlowRef: &WorkflowRef{
					WorkflowID:       "teste",
					Invoke:           InvokeKindAsync,
					OnParentComplete: "terminate",
				},
			},
			Err: ``,
		},
	}

	StructLevelValidationCtx(t, NewValidatorContext(&workflowStructDefault), testCases)
}
