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
	"github.com/stretchr/testify/assert"
	"testing"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func TestSleepStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp  string
		state State
		err   string
	}
	testCases := []testCase{
		{
			desp: "normal duration",
			state: State{
				BaseState: BaseState{
					Name: "1",
					Type: "sleep",
					End: &End{
						Terminate: true,
					},
				},
				SleepState: &SleepState{
					Duration: "PT10S",
				},
			},
			err: ``,
		},
		{
			desp: "invalid duration",
			state: State{
				BaseState: BaseState{
					Name: "1",
					Type: "sleep",
				},
				SleepState: &SleepState{
					Duration: "T10S",
				},
			},
			err: `Key: 'State.SleepState.Duration' Error:Field validation for 'Duration' failed on the 'iso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.state)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestSleepStateToString(t *testing.T) {

	single := "2023-05-22T10:00:35+02:00"
	total := "2023-05-22T10:01:35+02:00"
	stateExecTimeout := StateExecTimeout{
		Single: single,
		Total:  total,
	}

	timeout := SleepStateTimeout{
		StateExecTimeout: &stateExecTimeout,
	}

	sleepState := SleepState{
		Duration: "duration",
		Timeouts: &timeout,
	}
	value := sleepState.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ Duration:duration, Timeouts:&{StateExecTimeout:{ Single:2023-05-22T10:00:35+02:00, Total:2023-05-22T10:01:35+02:00}} }", value)
}
