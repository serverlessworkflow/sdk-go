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

	"github.com/stretchr/testify/assert"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func TestDelayStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp          string
		delayStateObj State
		err           string
	}
	testCases := []testCase{
		{
			desp: "normal",
			delayStateObj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "delay",
					End: &End{
						Terminate: true,
					},
				},
				DelayState: &DelayState{
					TimeDelay: "PT5S",
				},
			},
			err: ``,
		},
		{
			desp: "missing required timeDelay",
			delayStateObj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "delay",
				},
				DelayState: &DelayState{
					TimeDelay: "",
				},
			},
			err: `Key: 'State.DelayState.TimeDelay' Error:Field validation for 'TimeDelay' failed on the 'required' tag`,
		},
		{
			desp: "invalid timeDelay duration",
			delayStateObj: State{
				BaseState: BaseState{
					Name: "1",
					Type: "delay",
				},
				DelayState: &DelayState{
					TimeDelay: "P5S",
				},
			},
			err: `Key: 'State.DelayState.TimeDelay' Error:Field validation for 'TimeDelay' failed on the 'iso8601duration' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.delayStateObj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
