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

	val "github.com/finbox-in/serverlessworkflow-sdk-go/validator"
)

func TestWorkflowRefUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect WorkflowRef
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal object test",
			data: `{"workflowId": "1", "version": "2", "invoke": "async", "onParentComplete": "continue"}`,
			expect: WorkflowRef{
				WorkflowID:       "1",
				Version:          "2",
				Invoke:           InvokeKindAsync,
				OnParentComplete: "continue",
			},
			err: ``,
		},
		{
			desp: "normal object test & defaults",
			data: `{"workflowId": "1"}`,
			expect: WorkflowRef{
				WorkflowID:       "1",
				Version:          "",
				Invoke:           InvokeKindSync,
				OnParentComplete: "terminate",
			},
			err: ``,
		},
		{
			desp: "normal string test",
			data: `"1"`,
			expect: WorkflowRef{
				WorkflowID:       "1",
				Version:          "",
				Invoke:           InvokeKindSync,
				OnParentComplete: "terminate",
			},
			err: ``,
		},
		{
			desp:   "empty data",
			data:   ` `,
			expect: WorkflowRef{},
			err:    `unexpected end of JSON input`,
		},
		{
			desp:   "invalid string format",
			data:   `"1`,
			expect: WorkflowRef{},
			err:    `unexpected end of JSON input`,
		},
		{
			desp:   "invalid json format",
			data:   `{"workflowId": 1, "version": "2", "invoke": "async", "onParentComplete": "continue"}`,
			expect: WorkflowRef{},
			err:    `json: cannot unmarshal number into Go struct field workflowRefForUnmarshal.workflowId of type string`,
		},
		{
			desp:   "invalid string or object",
			data:   `1`,
			expect: WorkflowRef{},
			err:    `subFlowRef value '1' is not supported, it must be an object or string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v WorkflowRef
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

func TestWorkflowRefValidate(t *testing.T) {
	type testCase struct {
		desp        string
		workflowRef WorkflowRef
		err         string
	}
	testCases := []testCase{
		{
			desp: "all field & defaults",
			workflowRef: WorkflowRef{
				WorkflowID:       "1",
				Version:          "2",
				Invoke:           InvokeKindSync,
				OnParentComplete: "terminate",
			},
			err: ``,
		},
		{
			desp: "all field",
			workflowRef: WorkflowRef{
				WorkflowID:       "1",
				Version:          "2",
				Invoke:           InvokeKindAsync,
				OnParentComplete: "continue",
			},
			err: ``,
		},
		{
			desp: "missing workflowId",
			workflowRef: WorkflowRef{
				WorkflowID:       "",
				Version:          "2",
				Invoke:           InvokeKindSync,
				OnParentComplete: "terminate",
			},
			err: `Key: 'WorkflowRef.WorkflowID' Error:Field validation for 'WorkflowID' failed on the 'required' tag`,
		},
		{
			desp: "invalid invoke",
			workflowRef: WorkflowRef{
				WorkflowID:       "1",
				Version:          "2",
				Invoke:           "sync1",
				OnParentComplete: "terminate",
			},
			err: `Key: 'WorkflowRef.Invoke' Error:Field validation for 'Invoke' failed on the 'oneof' tag`,
		},
		{
			desp: "invalid onParentComplete",
			workflowRef: WorkflowRef{
				WorkflowID:       "1",
				Version:          "2",
				Invoke:           InvokeKindSync,
				OnParentComplete: "terminate1",
			},
			err: `Key: 'WorkflowRef.OnParentComplete' Error:Field validation for 'OnParentComplete' failed on the 'oneof' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.workflowRef)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
