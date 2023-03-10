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

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
	"github.com/stretchr/testify/assert"
)

func TestSwitchStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp string
		obj  State
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal & eventConditions",
			obj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "switch",
				},
				SwitchState: &SwitchState{
					DefaultCondition: DefaultCondition{
						Transition: &Transition{
							NextState: "1",
						},
					},
					EventConditions: []EventCondition{
						{
							EventRef: "1",
							Transition: &Transition{
								NextState: "2",
							},
						},
					},
				},
			},
			err: ``,
		},
		{
			desp: "normal & dataConditions",
			obj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "switch",
				},
				SwitchState: &SwitchState{
					DefaultCondition: DefaultCondition{
						Transition: &Transition{
							NextState: "1",
						},
					},
					DataConditions: []DataCondition{
						{
							Condition: "1",
							Transition: &Transition{
								NextState: "2",
							},
						},
					},
				},
			},
			err: ``,
		},
		{
			desp: "missing eventConditions & dataConditions",
			obj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "switch",
				},
				SwitchState: &SwitchState{
					DefaultCondition: DefaultCondition{
						Transition: &Transition{
							NextState: "1",
						},
					},
				},
			},
			err: `Key: 'State.SwitchState.DataConditions' Error:Field validation for 'DataConditions' failed on the 'required' tag`,
		},
		{
			desp: "exclusive eventConditions & dataConditions",
			obj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "switch",
				},
				SwitchState: &SwitchState{
					DefaultCondition: DefaultCondition{
						Transition: &Transition{
							NextState: "1",
						},
					},
					EventConditions: []EventCondition{
						{
							EventRef: "1",
							Transition: &Transition{
								NextState: "2",
							},
						},
					},
					DataConditions: []DataCondition{
						{
							Condition: "1",
							Transition: &Transition{
								NextState: "2",
							},
						},
					},
				},
			},
			err: `Key: 'State.SwitchState.DataConditions' Error:Field validation for 'DataConditions' failed on the 'exclusive' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.obj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestDefaultConditionStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp string
		obj  DefaultCondition
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal & end",
			obj: DefaultCondition{
				End: &End{
					Terminate: true,
				},
			},
			err: ``,
		},
		{
			desp: "normal & transition",
			obj: DefaultCondition{
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: ``,
		},
		{
			desp: "missing end & transition",
			obj:  DefaultCondition{},
			err:  `DefaultCondition.Transition' Error:Field validation for 'Transition' failed on the 'required' tag`,
		},
		{
			desp: "exclusive end & transition",
			obj: DefaultCondition{
				End: &End{
					Terminate: true,
				},
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: `Key: 'DefaultCondition.Transition' Error:Field validation for 'Transition' failed on the 'exclusive' tag`,
		},
	}
	for _, tc := range testCases[2:] {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.obj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestEventConditionStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp string
		obj  EventCondition
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal & end",
			obj: EventCondition{
				EventRef: "1",
				End: &End{
					Terminate: true,
				},
			},
			err: ``,
		},
		{
			desp: "normal & transition",
			obj: EventCondition{
				EventRef: "1",
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: ``,
		},
		{
			desp: "missing end & transition",
			obj: EventCondition{
				EventRef: "1",
			},
			err: `Key: 'EventCondition.Transition' Error:Field validation for 'Transition' failed on the 'required' tag`,
		},
		{
			desp: "exclusive end & transition",
			obj: EventCondition{
				EventRef: "1",
				End: &End{
					Terminate: true,
				},
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: `Key: 'EventCondition.Transition' Error:Field validation for 'Transition' failed on the 'exclusive' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.obj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestDataConditionStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp string
		obj  DataCondition
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal & end",
			obj: DataCondition{
				Condition: "1",
				End: &End{
					Terminate: true,
				},
			},
			err: ``,
		},
		{
			desp: "normal & transition",
			obj: DataCondition{
				Condition: "1",
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: ``,
		},
		{
			desp: "missing end & transition",
			obj: DataCondition{
				Condition: "1",
			},
			err: `Key: 'DataCondition.Transition' Error:Field validation for 'Transition' failed on the 'required' tag`,
		},
		{
			desp: "exclusive end & transition",
			obj: DataCondition{
				Condition: "1",
				End: &End{
					Terminate: true,
				},
				Transition: &Transition{
					NextState: "1",
				},
			},
			err: `Key: 'DataCondition.Transition' Error:Field validation for 'Transition' failed on the 'exclusive' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.obj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
