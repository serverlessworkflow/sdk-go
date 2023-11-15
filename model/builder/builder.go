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

	"sigs.k8s.io/yaml"

	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func New() *model.WorkflowBuilder {
	return model.NewWorkflowBuilder()
}

func AsObject(builder *model.WorkflowBuilder) *model.Workflow {
	workflow := builder.Build()
	return &workflow
}

func AsJson(builder *model.WorkflowBuilder) ([]byte, error) {
	workflow := builder.Build()
	return json.Marshal(workflow)
}

func AsYaml(builder *model.WorkflowBuilder) ([]byte, error) {
	data, err := AsJson(builder)
	if err != nil {
		return nil, err
	}

	return yaml.JSONToYAML(data)
}
