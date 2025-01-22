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
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

func TestTaskList_UnmarshalJSON(t *testing.T) {
	jsonData := `[
		{"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}},
		{"task2": {"do": [{"task3": {"call": "openapi", "with": {"document": {"name": "doc1"}, "operationId": "op1"}}}]}}
	]`

	var taskList TaskList
	err := json.Unmarshal([]byte(jsonData), &taskList)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(taskList))

	task1 := taskList.Key("task1").AsCallHTTPTask()
	assert.NotNil(t, task1)
	assert.Equal(t, "http", task1.Call)
	assert.Equal(t, "GET", task1.With.Method)
	assert.Equal(t, "http://example.com", task1.With.Endpoint.URITemplate.String())

	task2 := taskList.Key("task2").AsDoTask()
	assert.NotNil(t, task2)
	assert.Equal(t, 1, len(*task2.Do))

	task3 := task2.Do.Key("task3").AsCallOpenAPITask()
	assert.NotNil(t, task3)
	assert.Equal(t, "openapi", task3.Call)
	assert.Equal(t, "doc1", task3.With.Document.Name)
	assert.Equal(t, "op1", task3.With.OperationID)
}

func TestTaskList_MarshalJSON(t *testing.T) {
	taskList := TaskList{
		{Key: "task1", Task: &CallHTTP{
			Call: "http",
			With: HTTPArguments{
				Method:   "GET",
				Endpoint: &Endpoint{URITemplate: &LiteralUri{Value: "http://example.com"}},
			},
		}},
		{Key: "task2", Task: &DoTask{
			Do: &TaskList{
				{Key: "task3", Task: &CallOpenAPI{
					Call: "openapi",
					With: OpenAPIArguments{
						Document:    &ExternalResource{Name: "doc1", Endpoint: NewEndpoint("http://example.com")},
						OperationID: "op1",
					},
				}},
			},
		}},
	}

	data, err := json.Marshal(taskList)
	assert.NoError(t, err)
	assert.JSONEq(t, `[
		{"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}},
		{"task2": {"do": [{"task3": {"call": "openapi", "with": {"document": {"name": "doc1", "endpoint": "http://example.com"}, "operationId": "op1"}}}]}}
	]`, string(data))
}

func TestTaskList_Validation(t *testing.T) {
	taskList := TaskList{
		{Key: "task1", Task: &CallHTTP{
			Call: "http",
			With: HTTPArguments{
				Method:   "GET",
				Endpoint: NewEndpoint("http://example.com"),
			},
		}},
		{Key: "task2", Task: &DoTask{
			Do: &TaskList{
				{Key: "task3", Task: &CallOpenAPI{
					Call: "openapi",
					With: OpenAPIArguments{
						Document:    &ExternalResource{Name: "doc1", Endpoint: NewEndpoint("http://example.com")},
						OperationID: "op1",
					},
				}},
			},
		}},
	}

	// Validate each TaskItem explicitly
	for _, taskItem := range taskList {
		err := validate.Struct(taskItem)
		if err != nil {
			var validationErrors validator.ValidationErrors
			if errors.As(err, &validationErrors) {
				for _, validationErr := range validationErrors {
					t.Errorf("Validation failed on field '%s' with tag '%s'", validationErr.Field(), validationErr.Tag())
				}
			} else {
				t.Errorf("Unexpected error: %v", err)
			}
		}
	}

}
