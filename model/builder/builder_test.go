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

	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	builder := New()
	builder.BaseWorkflow().Key("key test")
	builder.BaseWorkflow().ID("id test")

	function := builder.AddFunction()
	function.Name("function name")
	function2 := builder.AddFunction()
	function2.Name("function name2")

	workflow := AsObject(builder)

	assert.Equal(t, "key test", workflow.Key)
	assert.Equal(t, "id test", workflow.ID)
	assert.Equal(t, "0.8", workflow.SpecVersion)
	assert.Equal(t, "jq", workflow.ExpressionLang.String())
	assert.Equal(t, 2, len(workflow.Functions))

	assert.Equal(t, "function name", workflow.Functions[0].Name)
	assert.Equal(t, "function name2", workflow.Functions[1].Name)
}

func TestAsJson(t *testing.T) {
	builder := New()
	builder.BaseWorkflow().Key("key test")
	builder.BaseWorkflow().ID("id test")
	function := builder.AddFunction()
	function.Name("function name")
	function2 := builder.AddFunction()
	function2.Name("function name2")

	data, err := AsJson(builder)
	if assert.NoError(t, err) {
		d := `{"id":"id test","key":"key test","version":"","specVersion":"0.8","expressionLang":"jq","states":[],"functions":[{"name":"function name","operation":"","type":"rest"},{"name":"function name2","operation":"","type":"rest"}]}`
		assert.Equal(t, d, string(data))
	}
}

func TestAsYaml(t *testing.T) {
	builder := New()
	builder.BaseWorkflow().Key("key test")
	builder.BaseWorkflow().ID("id test")
	function := builder.AddFunction()
	function.Name("function name")
	function2 := builder.AddFunction()
	function2.Name("function name2")

	data, err := AsYaml(builder)
	if assert.NoError(t, err) {
		d := `expressionLang: jq
functions:
- name: function name
  operation: ""
  type: rest
- name: function name2
  operation: ""
  type: rest
id: id test
key: key test
specVersion: "0.8"
states: []
version: ""
`
		assert.Equal(t, d, string(data))
	}
}
