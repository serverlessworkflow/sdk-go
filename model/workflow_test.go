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

func TestContinueAsStructLevelValidation(t *testing.T) {
	type testCase struct {
		name       string
		continueAs ContinueAs
		err        string
	}

	testCases := []testCase{
		{
			name: "valid ContinueAs",
			continueAs: ContinueAs{
				WorkflowID: "another-test",
				Version:    "2",
				Data:       SwObject{Object(String("${ del(.customerCount) }"))},
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration:  "PT1H",
					Interrupt: false,
					RunBefore: "test",
				},
			},
			err: ``,
		},
		{
			name: "invalid WorkflowExecTimeout",
			continueAs: ContinueAs{
				WorkflowID: "test",
				Version:    "1",
				Data:       SwObject{Object(String("${ del(.customerCount) }"))},
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration: "invalid",
				},
			},
			err: `Key: 'ContinueAs.workflowExecTimeout' Error:Field validation for 'workflowExecTimeout' failed on the 'iso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.continueAs)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestContinueAsUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect ContinueAs
		err    string
	}
	testCases := []testCase{
		{
			desp: "string",
			data: `"1"`,
			expect: ContinueAs{
				WorkflowID: "1",
			},
			err: ``,
		},
		{
			desp: "object all field set",
			data: `{"workflowId": "1", "version": "2", "data": "3", "workflowExecTimeout": {"duration": "PT1H", "interrupt": true, "runBefore": "4"}}`,
			expect: ContinueAs{
				WorkflowID: "1",
				Version:    "2",
				Data:       SwObject{Object(String("3"))},
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration:  "PT1H",
					Interrupt: true,
					RunBefore: "4",
				},
			},
			err: ``,
		},
		{
			desp: "object optional field unset",
			data: `{"workflowId": "1"}`,
			expect: ContinueAs{
				WorkflowID: "1",
				Version:    "",
				Data:       SwObject{},
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration:  "",
					Interrupt: false,
					RunBefore: "",
				},
			},
			err: ``,
		},
		{
			desp:   "invalid string format",
			data:   `"{`,
			expect: ContinueAs{},
			err:    `unexpected end of JSON input`,
		},
		{
			desp:   "invalid object format",
			data:   `{"workflowId": 1}`,
			expect: ContinueAs{},
			err:    `json: cannot unmarshal number into Go struct field continueAsForUnmarshal.workflowId of type string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v ContinueAs
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
