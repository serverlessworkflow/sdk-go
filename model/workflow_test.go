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
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/util"
	"github.com/stretchr/testify/assert"
)

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
					SpecVersion:    "0.8",
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
			desp: "start empty and states empty",
			data: `{"states": []}`,
			expect: Workflow{
				BaseWorkflow: BaseWorkflow{
					SpecVersion:    "0.8",
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
			err:    `continueAs.workflowId must be string`,
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
			err:    `end must be bool or object`,
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
			desp: `object fail`,
			data: `{"terminate": "true"}`,
			expect: End{
				Terminate: true,
			},
			err: `end.terminate must be bool`,
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
				assert.Equal(t, tc.err, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, v)
		})
	}
}

func TestWorkflowExecTimeoutUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect WorkflowExecTimeout
		err    string
	}

	testCases := []testCase{
		{
			desp: "string success",
			data: `"PT15M"`,
			expect: WorkflowExecTimeout{
				Duration: "PT15M",
			},
			err: ``,
		},
		{
			desp: "string fail",
			data: `PT15M`,
			expect: WorkflowExecTimeout{
				Duration: "PT15M",
			},
			err: `invalid character 'P' looking for beginning of value`,
		},
		{
			desp: `object success`,
			data: `{"duration": "PT15M"}`,
			expect: WorkflowExecTimeout{
				Duration: "PT15M",
			},
			err: ``,
		},
		{
			desp: `object fail`,
			data: `{"duration": PT15M}`,
			expect: WorkflowExecTimeout{
				Duration: "PT15M",
			},
			err: `invalid character 'P' looking for beginning of value`,
		},
		{
			desp: `object key invalid`,
			data: `{"duration_invalid": "PT15M"}`,
			expect: WorkflowExecTimeout{
				Duration: "unlimited",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v WorkflowExecTimeout
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

func TestStartUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect Start
		err    string
	}

	testCases := []testCase{
		{
			desp: "string success",
			data: `"start state"`,
			expect: Start{
				StateName: "start state",
			},
			err: ``,
		},
		{
			desp: "string fail",
			data: `start state`,
			expect: Start{
				StateName: "start state",
			},
			err: `invalid character 's' looking for beginning of value`,
		},
		{
			desp: `object success`,
			data: `{"stateName": "start state"}`,
			expect: Start{
				StateName: "start state",
			},
			err: ``,
		},
		{
			desp: `object fail`,
			data: `{"stateName": start state}`,
			expect: Start{
				StateName: "start state",
			},
			err: `invalid character 's' looking for beginning of value`,
		},
		{
			desp: `object key invalid`,
			data: `{"stateName_invalid": "start state"}`,
			expect: Start{
				StateName: "",
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v Start
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

func TestCronUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect Cron
		err    string
	}

	testCases := []testCase{
		{
			desp: "string success",
			data: `"0 15,30,45 * ? * *"`,
			expect: Cron{
				Expression: "0 15,30,45 * ? * *",
			},
			err: ``,
		},
		{
			desp: "string fail",
			data: `0 15,30,45 * ? * *`,
			expect: Cron{
				Expression: "0 15,30,45 * ? * *",
			},
			err: `invalid character '1' after top-level value`,
		},
		{
			desp: `object success`,
			data: `{"expression": "0 15,30,45 * ? * *"}`,
			expect: Cron{
				Expression: "0 15,30,45 * ? * *",
			},
			err: ``,
		},
		{
			desp: `object fail`,
			data: `{"expression": "0 15,30,45 * ? * *}`,
			expect: Cron{
				Expression: "0 15,30,45 * ? * *",
			},
			err: `unexpected end of JSON input`,
		},
		{
			desp:   `object key invalid`,
			data:   `{"expression_invalid": "0 15,30,45 * ? * *"}`,
			expect: Cron{},
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v Cron
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

func TestTransitionUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect Transition
		err    string
	}

	testCases := []testCase{
		{
			desp: "string success",
			data: `"next state"`,
			expect: Transition{
				NextState: "next state",
			},
			err: ``,
		},
		{
			desp: `object success`,
			data: `{"nextState": "next state"}`,
			expect: Transition{
				NextState: "next state",
			},
			err: ``,
		},
		{
			desp: `object fail`,
			data: `{"nextState": "next state}`,
			expect: Transition{
				NextState: "next state",
			},
			err: `unexpected end of JSON input`,
		},
		{
			desp:   `object key invalid`,
			data:   `{"nextState_invalid": "next state"}`,
			expect: Transition{},
			err:    ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v Transition
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

func TestDataInputSchemaUnmarshalJSON(t *testing.T) {

	var schemaName Object
	err := json.Unmarshal([]byte("{\"key\": \"value\"}"), &schemaName)
	if !assert.NoError(t, err) {
		return
	}

	type testCase struct {
		desp   string
		data   string
		expect DataInputSchema
		err    string
	}

	testCases := []testCase{
		{
			desp: "string success",
			data: "{\"key\": \"value\"}",
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: ``,
		},
		{
			desp: "string fail",
			data: "{\"key\": }",
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: `invalid character '}' looking for beginning of value`,
		},
		{
			desp: `object success (without quotes)`,
			data: `{"key": "value"}`,
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: ``,
		},
		{
			desp: `schema object success`,
			data: `{"schema": "{\"key\": \"value\"}"}`,
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: ``,
		},
		{
			desp: `schema object success (without quotes)`,
			data: `{"schema": {"key": "value"}}`,
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: ``,
		},
		{
			desp: `schema object fail`,
			data: `{"schema": "schema name}`,
			expect: DataInputSchema{
				Schema:                 &schemaName,
				FailOnValidationErrors: true,
			},
			err: `unexpected end of JSON input`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v DataInputSchema
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err, tc.desp)
				assert.Regexp(t, tc.err, err, tc.desp)
				return
			}

			assert.NoError(t, err, tc.desp)
			assert.Equal(t, tc.expect.Schema, v.Schema, tc.desp)
			assert.Equal(t, tc.expect.FailOnValidationErrors, v.FailOnValidationErrors, tc.desp)
		})
	}
}

func TestConstantsUnmarshalJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/test.json":
			_, err := rw.Write([]byte(`{"testkey":"testvalue"}`))
			assert.NoError(t, err)
		default:
			t.Failed()
		}
	}))
	defer server.Close()
	util.HttpClient = *server.Client()

	type testCase struct {
		desp   string
		data   string
		expect Constants
		err    string
	}
	testCases := []testCase{
		{
			desp: "object success",
			data: `{"testkey":"testvalue}`,
			expect: Constants{
				Data: ConstantsData{
					"testkey": []byte(`"testvalue"`),
				},
			},
			err: `unexpected end of JSON input`,
		},
		{
			desp: "object success",
			data: `[]`,
			expect: Constants{
				Data: ConstantsData{
					"testkey": []byte(`"testvalue"`),
				},
			},
			// TODO: improve message: field is empty
			err: `constants must be string or object`,
		},
		{
			desp: "object success",
			data: `{"testkey":"testvalue"}`,
			expect: Constants{
				Data: ConstantsData{
					"testkey": []byte(`"testvalue"`),
				},
			},
			err: ``,
		},
		{
			desp: "file success",
			data: fmt.Sprintf(`"%s/test.json"`, server.URL),
			expect: Constants{
				Data: ConstantsData{
					"testkey": []byte(`"testvalue"`),
				},
			},
			err: ``,
		},
		{
			desp: "file success",
			data: `"uri_invalid"`,
			expect: Constants{
				Data: ConstantsData{
					"testkey": []byte(`"testvalue"`),
				},
			},
			err: `file not found: "uri_invalid"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v Constants
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
