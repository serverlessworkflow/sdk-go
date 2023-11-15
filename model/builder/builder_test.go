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
		d := `{"id":"id test","key":"key test","version":"","start":{"stateName":"","schedule":{"cron":{"expression":""}}},"dataInputSchema":{"schema":"","failOnValidationErrors":false},"specVersion":"","constants":{},"timeouts":{"workflowExecTimeout":{"duration":""},"stateExecTimeout":{"total":""}},"states":[],"functions":[{"name":"function name","operation":""},{"name":"function name2","operation":""}]}`
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
		d := `constants: {}
dataInputSchema:
  failOnValidationErrors: false
  schema: ""
functions:
- name: function name
  operation: ""
- name: function name2
  operation: ""
id: id test
key: key test
specVersion: ""
start:
  schedule:
    cron:
      expression: ""
  stateName: ""
states: []
timeouts:
  stateExecTimeout:
    total: ""
  workflowExecTimeout:
    duration: ""
version: ""
`
		assert.Equal(t, d, string(data))
	}
}
