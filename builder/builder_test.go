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
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func prepareBuilder() *model.WorkflowBuilder {
	builder := New()
	builder.BaseWorkflow().Key("key test")
	builder.BaseWorkflow().ID("id test")

	function := builder.AddFunctions()
	function.Name("function name")
	function.Operation("http://test")
	function2 := builder.AddFunctions()
	function2.Name("function name2")
	function2.Operation("http://test")

	state1 := builder.AddStates()
	state1.BaseState().Name("state")
	state1.BaseState().Type(model.StateTypeInject)
	end := state1.BaseState().End()
	end.Terminate(true)

	inject := state1.InjectState()
	inject.Data(map[string]model.Object{
		"test": model.FromMap(map[string]any{}),
	})

	return builder
}

func TestBuild(t *testing.T) {
	workflow, err := Object(prepareBuilder())
	if assert.NoError(t, err) {
		assert.Equal(t, "key test", workflow.Key)
		assert.Equal(t, "id test", workflow.ID)
		assert.Equal(t, "0.8", workflow.SpecVersion)
		assert.Equal(t, "jq", workflow.ExpressionLang.String())
		assert.Equal(t, 2, len(workflow.Functions))

		assert.Equal(t, "function name", workflow.Functions[0].Name)
		assert.Equal(t, "function name2", workflow.Functions[1].Name)
	}
}

func TestAsJson(t *testing.T) {
	data, err := Json(prepareBuilder())
	if assert.NoError(t, err) {
		d := `{"id":"id test","key":"key test","version":"","specVersion":"0.8","expressionLang":"jq","states":[{"name":"state","type":"inject","end":{"terminate":true},"data":{"test":{}}}],"functions":[{"name":"function name","operation":"http://test","type":"rest"},{"name":"function name2","operation":"http://test","type":"rest"}]}`
		assert.Equal(t, d, string(data))
	}
}

func TestAsYaml(t *testing.T) {
	data, err := Yaml(prepareBuilder())
	if assert.NoError(t, err) {
		d := `expressionLang: jq
functions:
- name: function name
  operation: http://test
  type: rest
- name: function name2
  operation: http://test
  type: rest
id: id test
key: key test
specVersion: "0.8"
states:
- data:
    test: {}
  end:
    terminate: true
  name: state
  type: inject
version: ""
`

		assert.Equal(t, d, string(data))
	}
}
