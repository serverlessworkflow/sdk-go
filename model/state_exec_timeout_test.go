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

func TestStateExecTimeoutUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp string
		data string

		expect *StateExecTimeout
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal string",
			data: `"PT10S"`,

			expect: &StateExecTimeout{
				Single: "",
				Total:  "PT10S",
			},
			err: ``,
		},
		{
			desp: "normal object with total",
			data: `{
				"total": "PT10S"
			}`,

			expect: &StateExecTimeout{
				Single: "",
				Total:  "PT10S",
			},
			err: ``,
		},
		{
			desp: "normal object with total & single",
			data: `{
				"single": "PT1S",
				"total": "PT10S"
			}`,

			expect: &StateExecTimeout{
				Single: "PT1S",
				Total:  "PT10S",
			},
			err: ``,
		},
		{
			desp: "invalid string or object",
			data: `PT10S`,

			expect: &StateExecTimeout{},
			err:    `stateExecTimeout value 'PT10S' is not supported, it must be an object or string`,
		},
		{
			desp: "invalid total type",
			data: `{
				"single": "PT1S",
				"total": 10
			}`,

			expect: &StateExecTimeout{},
			err:    `json: cannot unmarshal number into Go struct field stateExecTimeoutForUnmarshal.total of type string`,
		},
		{
			desp: "invalid single type",
			data: `{
				"single": 1,
				"total": "PT10S"
			}`,

			expect: &StateExecTimeout{
				Single: "",
				Total:  "PT10S",
			},
			err: `json: cannot unmarshal number into Go struct field stateExecTimeoutForUnmarshal.single of type string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			actual := &StateExecTimeout{}
			err := actual.UnmarshalJSON([]byte(tc.data))

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, actual)
		})
	}
}

func TestStateExecTimeoutStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp    string
		timeout StateExecTimeout
		err     string
	}
	testCases := []testCase{
		{
			desp: "normal total",
			timeout: StateExecTimeout{
				Total: "PT10S",
			},
			err: ``,
		},
		{
			desp: "normal total & single",
			timeout: StateExecTimeout{
				Single: "PT10S",
				Total:  "PT10S",
			},
			err: ``,
		},
		{
			desp: "missing total",
			timeout: StateExecTimeout{
				Single: "PT10S",
				Total:  "",
			},
			err: `Key: 'StateExecTimeout.Total' Error:Field validation for 'Total' failed on the 'required' tag`,
		},
		{
			desp: "invalid total duration",
			timeout: StateExecTimeout{
				Single: "PT10S",
				Total:  "T10S",
			},
			err: `Key: 'StateExecTimeout.Total' Error:Field validation for 'Total' failed on the 'iso8601duration' tag`,
		},
		{
			desp: "invalid single duration",
			timeout: StateExecTimeout{
				Single: "T10S",
				Total:  "PT10S",
			},
			err: `Key: 'StateExecTimeout.Single' Error:Field validation for 'Single' failed on the 'iso8601duration' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.timeout)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
