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

var stateTransitionDefault = State{
	BaseState: BaseState{
		Name: "name state",
		Type: StateTypeOperation,
		Transition: &Transition{
			NextState: "next name state",
		},
	},
	OperationState: &OperationState{
		ActionMode: "sequential",
		Actions: []Action{{
			FunctionRef: &FunctionRef{
				RefName: "function 1",
				Invoke:  InvokeKindAsync,
			},
		}},
	},
}

var stateEndDefault = State{
	BaseState: BaseState{
		Name: "name state",
		Type: StateTypeOperation,
		End: &End{
			Terminate: true,
		},
	},
	OperationState: &OperationState{
		ActionMode: "sequential",
		Actions: []Action{{
			FunctionRef: &FunctionRef{
				RefName: "test",
				Invoke:  InvokeKindAsync,
			},
		}},
	},
}

var switchStateTransitionDefault = State{
	BaseState: BaseState{
		Name: "name state",
		Type: StateTypeSwitch,
	},
	SwitchState: &SwitchState{
		DataConditions: []DataCondition{
			{
				Condition: "${ .applicant | .age >= 18 }",
				Transition: &Transition{
					NextState: "nex state",
				},
			},
		},
		DefaultCondition: DefaultCondition{
			Transition: &Transition{
				NextState: "nex state",
			},
		},
	},
}

func TestStateStructLevelValidation(t *testing.T) {
	// type testCase struct {
	// 	name     string
	// 	instance State
	// 	err      string
	// }

	testCases := []ValidationCase[State]{
		{
			Desp:  "state transition success",
			Model: stateTransitionDefault,
			Err:   ``,
		},
		{
			Desp:  "state end success",
			Model: stateEndDefault,
			Err:   ``,
		},
		{
			Desp:  "switch state success",
			Model: switchStateTransitionDefault,
			Err:   ``,
		},
		{
			Desp: "state end and transition",
			Model: func() State {
				s := stateTransitionDefault
				s.End = stateEndDefault.End
				return s
			}(),
			Err: `Key: 'State.BaseState.Transition' Error:Field validation for 'Transition' failed on the 'exclusive' tag`,
		},
		{
			Desp: "basestate without end and transition",
			Model: func() State {
				s := stateTransitionDefault
				s.Transition = nil
				return s
			}(),
			Err: `Key: 'State.BaseState.Transition' Error:Field validation for 'Transition' failed on the 'required' tag`,
		},
	}

	workflow := &Workflow{
		Functions: Functions{{
			Name: "function 1",
		}},
	}
	StructLevelValidationCtx(t, NewValidatorContext(workflow), testCases)
}
