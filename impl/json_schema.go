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
	"errors"
	"fmt"

	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/xeipuuv/gojsonschema"
)

// ValidateJSONSchema validates the provided data against a model.Schema.
func ValidateJSONSchema(data interface{}, schema *model.Schema) error {
	if schema == nil {
		return nil
	}

	schema.ApplyDefaults()

	if schema.Format != model.DefaultSchema {
		return fmt.Errorf("unsupported schema format: '%s'", schema.Format)
	}

	var schemaJSON string
	if schema.Document != nil {
		documentBytes, err := json.Marshal(schema.Document)
		if err != nil {
			return fmt.Errorf("failed to marshal schema document to JSON: %w", err)
		}
		schemaJSON = string(documentBytes)
	} else if schema.Resource != nil {
		// TODO: Handle external resource references (not implemented here)
		return errors.New("external resources are not yet supported")
	} else {
		return errors.New("schema must have either a 'Document' or 'Resource'")
	}

	schemaLoader := gojsonschema.NewStringLoader(schemaJSON)
	dataLoader := gojsonschema.NewGoLoader(data)

	result, err := gojsonschema.Validate(schemaLoader, dataLoader)
	if err != nil {
		// TODO: use model.Error
		return fmt.Errorf("failed to validate JSON schema: %w", err)
	}

	if !result.Valid() {
		var validationErrors string
		for _, err := range result.Errors() {
			validationErrors += fmt.Sprintf("- %s\n", err.String())
		}
		return fmt.Errorf("JSON schema validation failed:\n%s", validationErrors)
	}

	return nil
}
