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

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func TestSleepValidate(t *testing.T) {
	type testCase struct {
		desp  string
		sleep Sleep
		err   string
	}
	testCases := []testCase{
		{
			desp: "all field empty",
			sleep: Sleep{
				Before: "",
				After:  "",
			},
			err: ``,
		},
		{
			desp: "only before field",
			sleep: Sleep{
				Before: "PT5M",
				After:  "",
			},
			err: ``,
		},
		{
			desp: "only after field",
			sleep: Sleep{
				Before: "",
				After:  "PT5M",
			},
			err: ``,
		},
		{
			desp: "all field",
			sleep: Sleep{
				Before: "PT5M",
				After:  "PT5M",
			},
			err: ``,
		},
		{
			desp: "invalid before value",
			sleep: Sleep{
				Before: "T5M",
				After:  "PT5M",
			},
			err: `Key: 'Sleep.Before' Error:Field validation for 'Before' failed on the 'iso8601duration' tag`,
		},
		{
			desp: "invalid after value",
			sleep: Sleep{
				Before: "PT5M",
				After:  "T5M",
			},
			err: `Key: 'Sleep.After' Error:Field validation for 'After' failed on the 'iso8601duration' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.sleep)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestFunctionRefUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect FunctionRef
		err    string
	}

	testCases := []testCase{
		{
			desp:   "invalid object refName",
			data:   `{"refName": 1}`,
			expect: FunctionRef{},
			err:    "functionRef.refName must be an string",
		},
		{
			desp: "object with refName",
			data: `{"refName": "function name"}`,
			expect: FunctionRef{
				RefName: "function name",
				Invoke:  InvokeKindSync,
			},
			err: ``,
		},
		{
			desp: "object with refName and Invoke",
			data: `{"refName": "function name", "invoke": "async"}`,
			expect: FunctionRef{
				RefName: "function name",
				Invoke:  InvokeKindAsync,
			},
			err: ``,
		},
		{
			desp: "refName string",
			data: `"function name"`,
			expect: FunctionRef{
				RefName: "function name",
				Invoke:  InvokeKindSync,
			},
			err: ``,
		},
	}

	for _, tc := range testCases[:1] {
		t.Run(tc.desp, func(t *testing.T) {
			var v FunctionRef
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, v)
		})
	}
}
