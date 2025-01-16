// Copyright 2020 The Serverless Workflow Specification Authors
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

package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFromYAMLSource(t *testing.T) {
	source := []byte(`
document:
  dsl: 1.0.0
  namespace: examples
  name: example-workflow
  version: 1.0.0
do:
  - task1:
      call: http
      with:
        method: GET
        endpoint: http://example.com
`)
	workflow, err := FromYAMLSource(source)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, "example-workflow", workflow.Document.Name)
}

func TestFromJSONSource(t *testing.T) {
	source := []byte(`{
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
				"with": {
					"method": "GET",
					"endpoint": "http://example.com"
				}
			}
		}
	]
}`)
	workflow, err := FromJSONSource(source)
	assert.NoError(t, err)
	assert.NotNil(t, workflow)
	assert.Equal(t, "example-workflow", workflow.Document.Name)
}

func TestFromFile(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Valid YAML File",
			filePath:    "testdata/valid_workflow.yaml",
			expectError: false,
		},
		{
			name:        "Invalid YAML File",
			filePath:    "testdata/invalid_workflow.yaml",
			expectError: true,
		},
		{
			name:        "Unsupported File Extension",
			filePath:    "testdata/unsupported_workflow.txt",
			expectError: true,
		},
		{
			name:        "Non-existent File",
			filePath:    "testdata/nonexistent_workflow.yaml",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workflow, err := FromFile(tt.filePath)
			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, workflow)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, workflow)
				assert.Equal(t, "example-workflow", workflow.Document.Name)
			}
		})
	}
}

func TestCheckFilePath(t *testing.T) {
	tests := []struct {
		name        string
		filePath    string
		expectError bool
	}{
		{
			name:        "Valid YAML File Path",
			filePath:    "testdata/valid_workflow.yaml",
			expectError: false,
		},
		{
			name:        "Unsupported File Extension",
			filePath:    "testdata/unsupported_workflow.txt",
			expectError: true,
		},
		{
			name:        "Directory Path",
			filePath:    "testdata",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := checkFilePath(tt.filePath)
			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
