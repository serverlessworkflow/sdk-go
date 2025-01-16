// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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

func TestDoTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"do": [
			{"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}},
			{"task2": {"call": "openapi", "with": {"document": {"name": "doc1"}, "operationId": "op1"}}}
		]
	}`

	var doTask DoTask
	err := json.Unmarshal([]byte(jsonData), &doTask)
	assert.NoError(t, err)

	task1 := doTask.Do.Key("task1").AsCallHTTPTask()
	assert.NotNil(t, task1)
	assert.Equal(t, "http", task1.Call)
	assert.Equal(t, "GET", task1.With.Method)
	assert.Equal(t, "http://example.com", task1.With.Endpoint.String())

	task2 := doTask.Do.Key("task2").AsCallOpenAPITask()
	assert.NotNil(t, task2)
	assert.Equal(t, "openapi", task2.Call)
	assert.Equal(t, "doc1", task2.With.Document.Name)
	assert.Equal(t, "op1", task2.With.OperationID)
}

func TestDoTask_MarshalJSON(t *testing.T) {
	doTask := DoTask{
		TaskBase: TaskBase{},
		Do: &TaskList{
			{Key: "task1", Task: &CallHTTP{
				Call: "http",
				With: HTTPArguments{
					Method:   "GET",
					Endpoint: NewEndpoint("http://example.com"),
				},
			}},
			{Key: "task2", Task: &CallOpenAPI{
				Call: "openapi",
				With: OpenAPIArguments{
					Document:    &ExternalResource{Name: "doc1", Endpoint: NewEndpoint("http://example.com")},
					OperationID: "op1",
				},
			}},
		},
	}

	data, err := json.Marshal(doTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"do": [
			{"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}},
			{"task2": {"call": "openapi", "with": {"document": {"name": "doc1", "endpoint": "http://example.com"}, "operationId": "op1"}}}
		]
	}`, string(data))
}

func TestDoTask_Validation(t *testing.T) {
	doTask := DoTask{
		TaskBase: TaskBase{},
		Do: &TaskList{
			{Key: "task1", Task: &CallHTTP{
				Call: "http",
				With: HTTPArguments{
					Method:   "GET",
					Endpoint: NewEndpoint("http://example.com"),
				},
			}},
			{Key: "task2", Task: &CallOpenAPI{
				Call: "openapi",
				With: OpenAPIArguments{
					Document:    &ExternalResource{Name: "doc1"}, //missing endpoint
					OperationID: "op1",
				},
			}},
		},
	}

	err := validate.Struct(doTask)
	assert.Error(t, err)
}
