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

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
	"github.com/stretchr/testify/assert"
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
				RefName: "test",
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
	type testCase struct {
		name     string
		instance State
		err      string
	}

	testCases := []testCase{
		{
			name:     "state transition success",
			instance: stateTransitionDefault,
			err:      ``,
		},
		{
			name:     "state end success",
			instance: stateEndDefault,
			err:      ``,
		},
		{
			name:     "switch state success",
			instance: switchStateTransitionDefault,
			err:      ``,
		},
		{
			name: "state end and transition",
			instance: func() State {
				s := stateTransitionDefault
				s.End = stateEndDefault.End
				return s
			}(),
			err: `Key: 'State.BaseState.Transition' Error:Field validation for 'Transition' failed on the 'exclusive' tag`,
		},
		{
			name: "basestate without end and transition",
			instance: func() State {
				s := stateTransitionDefault
				s.Transition = nil
				return s
			}(),
			err: `Key: 'State.BaseState.Transition' Error:Field validation for 'Transition' failed on the 'required' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.instance)

			if tc.err != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.err, err.Error())
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}
