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

func TestExtension_UnmarshalJSON(t *testing.T) {
	jsonData := `{
        "extend": "call",
        "when": "${condition}",
        "before": [
            {"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}}
        ],
        "after": [
            {"task2": {"call": "openapi", "with": {"document": {"name": "doc1"}, "operationId": "op1"}}}
        ]
    }`

	var extension Extension
	err := json.Unmarshal([]byte(jsonData), &extension)
	assert.NoError(t, err)
	assert.Equal(t, "call", extension.Extend)
	assert.Equal(t, NewExpr("${condition}"), extension.When)

	task1 := extension.Before.Key("task1").AsCallHTTPTask()
	assert.NotNil(t, task1)
	assert.Equal(t, "http", task1.Call)
	assert.Equal(t, "GET", task1.With.Method)
	assert.Equal(t, "http://example.com", task1.With.Endpoint.String())

	// Check if task2 exists before accessing its fields
	task2 := extension.After.Key("task2")
	assert.NotNil(t, task2, "task2 should not be nil")
	openAPITask := task2.AsCallOpenAPITask()
	assert.NotNil(t, openAPITask)
	assert.Equal(t, "openapi", openAPITask.Call)
	assert.Equal(t, "doc1", openAPITask.With.Document.Name)
	assert.Equal(t, "op1", openAPITask.With.OperationID)
}

func TestExtension_MarshalJSON(t *testing.T) {
	extension := Extension{
		Extend: "call",
		When:   NewExpr("${condition}"),
		Before: &TaskList{
			{Key: "task1", Task: &CallHTTP{
				Call: "http",
				With: HTTPArguments{
					Method:   "GET",
					Endpoint: NewEndpoint("http://example.com"),
				},
			}},
		},
		After: &TaskList{
			{Key: "task2", Task: &CallOpenAPI{
				Call: "openapi",
				With: OpenAPIArguments{
					Document:    &ExternalResource{Name: "doc1", Endpoint: NewEndpoint("http://example.com")},
					OperationID: "op1",
				},
			}},
		},
	}

	data, err := json.Marshal(extension)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"extend": "call",
		"when": "${condition}",
		"before": [
			{"task1": {"call": "http", "with": {"method": "GET", "endpoint": "http://example.com"}}}
		],
		"after": [
			{"task2": {"call": "openapi", "with": {"document": {"name": "doc1", "endpoint": "http://example.com"}, "operationId": "op1"}}}
		]
	}`, string(data))
}

func TestExtension_Validation(t *testing.T) {
	extension := Extension{
		Extend: "call",
		When:   NewExpr("${condition}"),
		Before: &TaskList{
			{Key: "task1", Task: &CallHTTP{
				Call: "http",
				With: HTTPArguments{
					Method:   "GET",
					Endpoint: NewEndpoint("http://example.com"),
				},
			}},
		},
		After: &TaskList{
			{Key: "task2", Task: &CallOpenAPI{
				Call: "openapi",
				With: OpenAPIArguments{
					Document: &ExternalResource{
						Name: "doc1", // Missing Endpoint
					},
					OperationID: "op1",
				},
			}},
		},
	}

	err := validate.Struct(extension)
	assert.Error(t, err)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, validationErr := range validationErrors {
			t.Logf("Validation failed on field '%s' with tag '%s': %s",
				validationErr.StructNamespace(), validationErr.Tag(), validationErr.Param())
		}

		// Assert on specific validation errors
		assert.Contains(t, validationErrors.Error(), "After[0].Task.With.Document.Endpoint")
		assert.Contains(t, validationErrors.Error(), "required")
	} else {
		t.Errorf("Unexpected error type: %v", err)
	}
}
