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

func TestCallbackStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp             string
		callbackStateObj State
		err              string
	}
	testCases := []testCase{
		{
			desp: "normal",
			callbackStateObj: State{
				BaseState: BaseState{
					Name: "callbackTest",
					Type: StateTypeCallback,
				},
				CallbackState: &CallbackState{
					Action: Action{
						ID:   "1",
						Name: "action1",
					},
					EventRef: "refExample",
				},
			},
			err: ``,
		},
		{
			desp: "missing required EventRef",
			callbackStateObj: State{
				BaseState: BaseState{
					Name: "callbackTest",
					Type: StateTypeCallback,
				},
				CallbackState: &CallbackState{
					Action: Action{
						ID:   "1",
						Name: "action1",
					},
				},
			},
			err: `Key: 'State.CallbackState.EventRef' Error:Field validation for 'EventRef' failed on the 'required' tag`,
		},
		// TODO need to register custom types
		//{
		//	desp: "missing required Action",
		//	callbackStateObj: State{
		//		BaseState: BaseState{
		//			Name: "callbackTest",
		//			Type: StateTypeCallback,
		//		},
		//		CallbackState: &CallbackState{
		//			EventRef: "refExample",
		//		},
		//	},
		//	err: `Key: 'State.CallbackState.Action' Error:Field validation for 'Action' failed on the 'required' tag`,
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(&tc.callbackStateObj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
