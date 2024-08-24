// Copyright 2024 The Serverless Workflow Specification Authors
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

	"github.com/serverlessworkflow/sdk-go/v4/validate"
	"sigs.k8s.io/yaml"
)

func Validate(builder *WorkflowBuilder) error {
	data, err := Json(builder)
	if err != nil {
		return err
	}

	err = validate.FromJSONSource(data)
	if err != nil {
		return err
	}

	return nil
}

func Json(builder *WorkflowBuilder) ([]byte, error) {
	data, err := json.MarshalIndent(builder.Node(), "", "  ")
	if err != nil {
		return nil, err
	}

	return data, nil
}

func Yaml(builder *WorkflowBuilder) ([]byte, error) {
	data, err := Json(builder)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(data)
}
