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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConditionUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect DefaultCondition
		err    string
	}

	testCases := []testCase{
		{
			desp: "json nextState success",
			data: `{"transition": {"nextState": "next state"}}`,
			expect: DefaultCondition{
				Transition: &Transition{
					NextState: "next state",
				},
			},
			err: ``,
		},
		{
			desp: "invalid json nextState",
			data: `{"transition": {"nextState": "next state}}`,
			err:  `unexpected end of JSON input`,
		},
		{
			desp: "invalid json nextState type",
			data: `{"transition": {"nextState": true}}`,
			err:  `transition.nextState must be string`,
		},
		{
			desp: "transition json success",
			data: `{"transition": "next state"}`,
			expect: DefaultCondition{
				Transition: &Transition{
					NextState: "next state",
				},
			},
			err: ``,
		},
		{
			desp: "invalid json transition",
			data: `{"transition": "next state}`,
			err:  `unexpected end of JSON input`,
		},
		{
			desp: "invalid json transition type",
			data: `{"transition": true}`,
			err:  `transition must be string or object`,
		},
		{
			desp: "string success",
			data: `"next state"`,
			expect: DefaultCondition{
				Transition: &Transition{
					NextState: "next state",
				},
			},
			err: ``,
		},
		{
			desp: "invalid string syntax",
			data: `"next state`,
			err:  `unexpected end of JSON input`,
		},
		{
			desp: "invalid type",
			data: `123`,
			err:  `defaultCondition must be string or object`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v DefaultCondition
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, v)
		})
	}
}

func TestSwitchStateToString(t *testing.T) {
	single := "2023-03-22T10:00:35+02:00"
	total := "2023-03-22T10:01:35+02:00"
	stateExecTimeout := StateExecTimeout{
		Single: single,
		Total:  total,
	}

	timeout := SwitchStateTimeout{
		StateExecTimeout: &stateExecTimeout,
	}

	transition := Transition{
		NextState:     "next",
		ProduceEvents: []ProduceEvent{},
		Compensate:    false,
	}

	end := End{
		Terminate:     true,
		Compensate:    false,
		ProduceEvents: []ProduceEvent{},
		ContinueAs:    &ContinueAs{},
	}

	defaultCOndition := DefaultCondition{
		Transition: &transition,
		End:        &end,
	}

	switchState := SwitchState{
		Timeouts:         &timeout,
		DefaultCondition: defaultCOndition,
		EventConditions:  []EventCondition{},
		DataConditions:   []DataCondition{},
	}
	value := switchState.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ DataConditions:[], EventConditions:[], DefaultCondition:{ Transition:[next, [], false], End:&{Terminate:true ProduceEvents:[] Compensate:false ContinueAs:{ WorkflowID:, Version:, WorkflowExecTimeout:[, false, ], Data:{ Type:0, IntVal:0, StrVal:, RawValue:[] } }} }, Timeouts:{ StateExecTimeout:{ Single:2023-03-22T10:00:35+02:00, Total:2023-03-22T10:01:35+02:00}, EventTimeout: } }", value)
}
