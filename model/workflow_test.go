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

var workflowStructDefault = Workflow{
	BaseWorkflow: BaseWorkflow{
		ID:          "id",
		SpecVersion: "0.8",
		Auth: AuthArray{
			{
				Name: "auth name",
			},
		},
		Start: &Start{
			StateName: "name state",
		},
	},
	States: []State{{
		BaseState: BaseState{
			Name: "name state",
			Type: StateTypeOperation,
		},
		OperationState: &OperationState{
			ActionMode: "sequential",
			Actions: []Action{
				{},
			},
		},
	}},
}

func TestValidationAsStructLevelValidation(t *testing.T) {
	type testCase[T any] struct {
		name     string
		instance T
		err      string
	}
	testCases := []testCase[any]{
		{
			name:     "workflow success",
			instance: workflowStructDefault,
		},
		{
			name: "workflow auth.name repeat",
			instance: func() Workflow {
				w := workflowStructDefault
				w.Auth = append(w.Auth, w.Auth[0])
				return w
			}(),
			err: `Key: 'Workflow.[]Auth.Name' Error:Field validation for '[]Auth.Name' failed on the 'reqnameunique' tag`,
		},
		{
			name: "workflow id exclude key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = ""
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow key exclude id",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = "key"
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow id and key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = "key"
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow without id and key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = ""
				return w
			}(),
			err: `Key: 'Workflow.BaseWorkflow.ID' Error:Field validation for 'ID' failed on the 'required_without' tag
Key: 'Workflow.BaseWorkflow.Key' Error:Field validation for 'Key' failed on the 'required_without' tag`,
		},
		{
			name: "workflow start",
			instance: func() Workflow {
				w := workflowStructDefault
				w.Start = &Start{
					StateName: "start state not found",
				}
				return w
			}(),
			err: `Key: 'Workflow.Start' Error:Field validation for 'Start' failed on the 'startnotexists' tag`,
		},
		{
			name: "valid ContinueAs",
			instance: ContinueAs{
				WorkflowID: "another-test",
				Version:    "2",
				Data:       FromString("${ del(.customerCount) }"),
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
			instance: ContinueAs{
				WorkflowID: "test",
				Version:    "1",
				Data:       FromString("${ del(.customerCount) }"),
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration: "invalid",
				},
			},
			err: `Key: 'ContinueAs.workflowExecTimeout' Error:Field validation for 'workflowExecTimeout' failed on the 'iso8601duration' tag`,
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

func TestWorkflowStartUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect Workflow
		err    string
	}
	testCases := []testCase{
		{
			desp: "start string",
			data: `{"start": "start state name"}`,
			expect: Workflow{
				BaseWorkflow: BaseWorkflow{
					ExpressionLang: "jq",
					Start: &Start{
						StateName: "start state name",
					},
				},
				States: []State{},
			},
			err: ``,
		},
		{
			desp: "start empty and use the first state",
			data: `{"states": [{"name": "start state name", "type": "operation"}]}`,
			expect: Workflow{
				BaseWorkflow: BaseWorkflow{
					ExpressionLang: "jq",
					Start: &Start{
						StateName: "start state name",
					},
				},
				States: []State{
					{
						BaseState: BaseState{
							Name: "start state name",
							Type: StateTypeOperation,
						},
						OperationState: &OperationState{
							ActionMode: "sequential",
						},
					},
				},
			},
			err: ``,
		},
		{
			desp: "start empty, and states empty",
			data: `{"states": []}`,
			expect: Workflow{
				BaseWorkflow: BaseWorkflow{
					ExpressionLang: "jq",
				},
				States: []State{},
			},
			err: ``,
		},
	}

	for _, tc := range testCases[1:] {
		t.Run(tc.desp, func(t *testing.T) {
			var v Workflow
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
				Data:       FromString("3"),
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
				Data:       Object{},
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

func TestEndUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect End
		err    string
	}
	testCases := []testCase{
		{
			desp: "bool success",
			data: `true`,
			expect: End{
				Terminate: true,
			},
			err: ``,
		},
		{
			desp:   "string fail",
			data:   `"true"`,
			expect: End{},
			err:    `json: cannot unmarshal string into Go value of type bool`,
		},
		{
			desp: `object success`,
			data: `{"terminate": true}`,
			expect: End{
				Terminate: true,
			},
			err: ``,
		},
		{
			desp:   `object key invalid`,
			data:   `{"terminate_parameter_invalid": true}`,
			expect: End{},
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v End
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
