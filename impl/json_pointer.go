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

package impl

import (
	"encoding/json"
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"reflect"
	"strings"
)

func findJsonPointer(data interface{}, target string, path string) (string, bool) {
	switch node := data.(type) {
	case map[string]interface{}:
		for key, value := range node {
			newPath := fmt.Sprintf("%s/%s", path, key)
			if key == target {
				return newPath, true
			}
			if result, found := findJsonPointer(value, target, newPath); found {
				return result, true
			}
		}
	case []interface{}:
		for i, item := range node {
			newPath := fmt.Sprintf("%s/%d", path, i)
			if result, found := findJsonPointer(item, target, newPath); found {
				return result, true
			}
		}
	}
	return "", false
}

// GenerateJSONPointer Function to generate JSON Pointer from a Workflow reference
func GenerateJSONPointer(workflow *model.Workflow, targetNode interface{}) (string, error) {
	// Convert struct to JSON
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %w", err)
	}

	// Convert JSON to a generic map for traversal
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %w", err)
	}

	transformedNode := ""
	switch node := targetNode.(type) {
	case string:
		transformedNode = node
	default:
		transformedNode = strings.ToLower(reflect.TypeOf(targetNode).Name())
	}

	// Search for the target node
	jsonPointer, found := findJsonPointer(jsonMap, transformedNode, "")
	if !found {
		return "", fmt.Errorf("node '%s' not found", targetNode)
	}

	return jsonPointer, nil
}
