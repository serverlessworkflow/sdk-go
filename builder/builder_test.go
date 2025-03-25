// Copyright 2023 The Serverless Workflow Specification Authors
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

package builder

import (
	"errors"
	"testing"

	validator "github.com/go-playground/validator/v10"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/serverlessworkflow/sdk-go/v3/test"

	"github.com/stretchr/testify/assert"
)

func TestBuilder_Yaml(t *testing.T) {
	builder := New().
		SetDocument("1.0.0", "examples", "example-workflow", "1.0.0").
		AddTask("task1", &model.CallHTTP{
			TaskBase: model.TaskBase{
				If: &model.RuntimeExpression{Value: "${condition}"},
			},
			Call: "http",
			With: model.HTTPArguments{
				Method:   "GET",
				Endpoint: model.NewEndpoint("http://example.com"),
			},
		})

	// Generate YAML from the builder
	yamlData, err := Yaml(builder)
	assert.NoError(t, err)

	// Define the expected YAML structure
	expectedYAML := `document:
  dsl: 1.0.0
  namespace: examples
  name: example-workflow
  version: 1.0.0
do:
- task1:
    call: http
    if: ${condition}
    with:
      method: GET
      endpoint: http://example.com
`

	// Use assertYAMLEq to compare YAML structures
	test.AssertYAMLEq(t, expectedYAML, string(yamlData))
}

func TestBuilder_Json(t *testing.T) {
	builder := New().
		SetDocument("1.0.0", "examples", "example-workflow", "1.0.0").
		AddTask("task1", &model.CallHTTP{
			TaskBase: model.TaskBase{
				If: &model.RuntimeExpression{Value: "${condition}"},
			},
			Call: "http",
			With: model.HTTPArguments{
				Method:   "GET",
				Endpoint: model.NewEndpoint("http://example.com"),
			},
		})

	jsonData, err := Json(builder)
	assert.NoError(t, err)

	expectedJSON := `{
  "document": {
    "dsl": "1.0.0",
    "namespace": "examples",
    "name": "example-workflow",
    "version": "1.0.0"
  },
  "do": [
    {
      "task1": {
        "call": "http",
        "if": "${condition}",
        "with": {
          "method": "GET",
          "endpoint": "http://example.com"
        }
      }
    }
  ]
}`
	assert.JSONEq(t, expectedJSON, string(jsonData))
}

func TestBuilder_Object(t *testing.T) {
	builder := New().
		SetDocument("1.0.0", "examples", "example-workflow", "1.0.0").
		AddTask("task1", &model.CallHTTP{
			TaskBase: model.TaskBase{
				If: &model.RuntimeExpression{Value: "${condition}"},
			},
			Call: "http",
			With: model.HTTPArguments{
				Method:   "GET",
				Endpoint: model.NewEndpoint("http://example.com"),
			},
		})

	workflow, err := Object(builder)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)

	assert.Equal(t, "1.0.0", workflow.Document.DSL)
	assert.Equal(t, "examples", workflow.Document.Namespace)
	assert.Equal(t, "example-workflow", workflow.Document.Name)
	assert.Equal(t, "1.0.0", workflow.Document.Version)
	assert.Len(t, *workflow.Do, 1)
	assert.Equal(t, "http", (*workflow.Do)[0].Task.(*model.CallHTTP).Call)
}

func TestBuilder_Validate(t *testing.T) {
	workflow := &model.Workflow{
		Document: model.Document{
			DSL:       "1.0.0",
			Namespace: "examples",
			Name:      "example-workflow",
			Version:   "1.0.0",
		},
		Do: &model.TaskList{
			&model.TaskItem{
				Key: "task1",
				Task: &model.CallHTTP{
					Call: "http",
					With: model.HTTPArguments{
						Method:   "GET",
						Endpoint: model.NewEndpoint("http://example.com"),
					},
				},
			},
		},
	}

	err := Validate(workflow)
	assert.NoError(t, err)

	// Test validation failure
	workflow.Do = &model.TaskList{
		&model.TaskItem{
			Key: "task2",
			Task: &model.CallHTTP{
				Call: "http",
				With: model.HTTPArguments{
					Method: "GET", // Missing Endpoint
				},
			},
		},
	}
	err = Validate(workflow)
	assert.Error(t, err)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		t.Logf("Validation errors: %v", validationErrors)
		assert.Contains(t, validationErrors.Error(), "Do[0].Task.With.Endpoint")
		assert.Contains(t, validationErrors.Error(), "required")
	}
}
