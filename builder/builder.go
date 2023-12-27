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
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func New() *model.WorkflowBuilder {
	return model.NewWorkflowBuilder()
}

func Yaml(builder *model.WorkflowBuilder) ([]byte, error) {
	data, err := Json(builder)
	if err != nil {
		return nil, err
	}
	return yaml.JSONToYAML(data)
}

func Json(builder *model.WorkflowBuilder) ([]byte, error) {
	workflow, err := Object(builder)
	if err != nil {
		return nil, err
	}
	return json.Marshal(workflow)
}

func Object(builder *model.WorkflowBuilder) (*model.Workflow, error) {
	workflow := builder.Build()
	ctx := model.NewValidatorContext(&workflow)
	if err := val.GetValidator().StructCtx(ctx, workflow); err != nil {
		return nil, err
	}
	return &workflow, nil
}

func Validate(object interface{}) error {
	ctx := model.NewValidatorContext(object)
	if err := val.GetValidator().StructCtx(ctx, object); err != nil {
		return val.WorkflowError(err)
	}
	return nil
}
