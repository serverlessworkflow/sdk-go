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
	"encoding/json"
	"fmt"

	"github.com/serverlessworkflow/sdk-go/v3/model"

	"sigs.k8s.io/yaml"
)

// New initializes a new WorkflowBuilder instance.
func New() *model.WorkflowBuilder {
	return model.NewWorkflowBuilder()
}

// Yaml generates YAML output from the WorkflowBuilder using custom MarshalYAML implementations.
func Yaml(builder *model.WorkflowBuilder) ([]byte, error) {
	workflow, err := Object(builder)
	if err != nil {
		return nil, fmt.Errorf("failed to build workflow object: %w", err)
	}
	return yaml.Marshal(workflow)
}

// Json generates JSON output from the WorkflowBuilder.
func Json(builder *model.WorkflowBuilder) ([]byte, error) {
	workflow, err := Object(builder)
	if err != nil {
		return nil, fmt.Errorf("failed to build workflow object: %w", err)
	}
	return json.MarshalIndent(workflow, "", "  ")
}

// Object builds and validates the Workflow object from the builder.
func Object(builder *model.WorkflowBuilder) (*model.Workflow, error) {
	workflow := builder.Build()

	// Validate the workflow object
	if err := model.GetValidator().Struct(workflow); err != nil {
		return nil, fmt.Errorf("workflow validation failed: %w", err)
	}

	return workflow, nil
}

// Validate validates any given object using the Workflow model validator.
func Validate(object interface{}) error {
	if err := model.GetValidator().Struct(object); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return nil
}
