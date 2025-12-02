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

package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunTask_MarshalJSON(t *testing.T) {
	runTask := RunTask{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Run: RunTaskConfiguration{
			Await: boolPtr(true),
			Container: &Container{
				Image:   "example-image",
				Name:    "example-name",
				Command: "example-command",
				Ports: map[string]interface{}{
					"8080": "80",
				},
				Environment: map[string]string{
					"ENV_VAR": "value",
				},
				Input: "example-input",
				Arguments: []string{
					"arg1",
					"arg2",
				},
				Lifetime: &ContainerLifetime{
					Cleanup: "eventually",
					After:   NewDurationExpr("20s"),
				},
			},
		},
	}

	data, err := json.Marshal(runTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"container": {
				"image": "example-image",
				"name": "example-name",
				"command": "example-command",
				"ports": {"8080": "80"},
				"environment": {"ENV_VAR": "value"},
				"stdin": "example-input",
				"arguments": ["arg1","arg2"],
				"lifetime": {
					"cleanup": "eventually",
					"after": "20s"
				}
			}
		}
	}`, string(data))
}

func TestRunTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"container": {
				"image": "example-image",
				"name": "example-name",
				"command": "example-command",
				"ports": {"8080": "80"},
				"environment": {"ENV_VAR": "value"},
				"stdin": "example-input",
				"arguments": ["arg1","arg2"],
				"lifetime": {
					"cleanup": "eventually",
					"after": {
						"seconds": 20
					}
				}
			}
		}
	}`

	var runTask RunTask
	err := json.Unmarshal([]byte(jsonData), &runTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, runTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, runTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, runTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, runTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, runTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, runTask.Metadata)
	assert.Equal(t, true, *runTask.Run.Await)
	assert.Equal(t, "example-image", runTask.Run.Container.Image)
	assert.Equal(t, "example-command", runTask.Run.Container.Command)
	assert.Equal(t, map[string]interface{}{"8080": "80"}, runTask.Run.Container.Ports)
	assert.Equal(t, map[string]string{"ENV_VAR": "value"}, runTask.Run.Container.Environment)
	assert.Equal(t, "example-name", runTask.Run.Container.Name)
	assert.Equal(t, "example-input", runTask.Run.Container.Input)
	assert.Equal(t, []string{"arg1", "arg2"}, runTask.Run.Container.Arguments)
	assert.Equal(t, "eventually", runTask.Run.Container.Lifetime.Cleanup)
	assert.Equal(t, &DurationInline{Seconds: 20}, runTask.Run.Container.Lifetime.After.AsInline())
}

func TestRunTaskScriptArgsMap_MarshalJSON(t *testing.T) {
	runTask := RunTask{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Run: RunTaskConfiguration{
			Await: boolPtr(true),
			Script: &Script{
				Language: "python",
				Arguments: &RunArguments{
					Value: map[string]interface{}{
						"arg1": "value1",
					},
				},
				Environment: map[string]string{
					"ENV_VAR": "value",
				},
				InlineCode: stringPtr("print('Hello, World!')"),
				Input:      "example-input",
			},
		},
	}

	data, err := json.Marshal(runTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"script": {
				"language": "python",
				"arguments": {"arg1": "value1"},
				"environment": {"ENV_VAR": "value"},
				"code": "print('Hello, World!')",
				"stdin": "example-input"
			}
		}
	}`, string(data))
}

func TestRunTaskScriptArgsMap_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"script": {
				"language": "python",
				"arguments": {"arg1": "value1"},
				"environment": {"ENV_VAR": "value"},
				"code": "print('Hello, World!')",
				"stdin": "example-input"
			}
		}
	}`

	var runTask RunTask
	err := json.Unmarshal([]byte(jsonData), &runTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, runTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, runTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, runTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, runTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, runTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, runTask.Metadata)
	assert.Equal(t, true, *runTask.Run.Await)
	assert.Equal(t, "python", runTask.Run.Script.Language)
	assert.Equal(t, map[string]interface{}{"arg1": "value1"}, runTask.Run.Script.Arguments.AsMap())
	assert.Equal(t, map[string]string{"ENV_VAR": "value"}, runTask.Run.Script.Environment)
	assert.Equal(t, "print('Hello, World!')", *runTask.Run.Script.InlineCode)
	assert.Equal(t, "example-input", runTask.Run.Script.Input)
}

func TestRunTaskScriptArgArray_MarshalJSON(t *testing.T) {
	runTask := RunTask{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Run: RunTaskConfiguration{
			Await: boolPtr(true),
			Script: &Script{
				Language: "python",
				Arguments: &RunArguments{
					Value: []string{
						"arg1=value1",
					},
				},
				Environment: map[string]string{
					"ENV_VAR": "value",
				},
				InlineCode: stringPtr("print('Hello, World!')"),
				Input:      "example-input",
			},
		},
	}

	data, err := json.Marshal(runTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"script": {
				"language": "python",
				"arguments": ["arg1=value1"],
				"environment": {"ENV_VAR": "value"},
				"code": "print('Hello, World!')",
				"stdin": "example-input"
			}
		}
	}`, string(data))
}

func TestRunTaskScriptArgsArray_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"run": {
			"await": true,
			"script": {
				"language": "python",
				"arguments": ["arg1=value1"],
				"environment": {"ENV_VAR": "value"},
				"code": "print('Hello, World!')",
				"stdin": "example-input"
			}
		}
	}`

	var runTask RunTask
	err := json.Unmarshal([]byte(jsonData), &runTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, runTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, runTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, runTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, runTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, runTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, runTask.Metadata)
	assert.Equal(t, true, *runTask.Run.Await)
	assert.Equal(t, "python", runTask.Run.Script.Language)
	assert.Equal(t, []string{"arg1=value1"}, runTask.Run.Script.Arguments.AsSlice())
	assert.Equal(t, map[string]string{"ENV_VAR": "value"}, runTask.Run.Script.Environment)
	assert.Equal(t, "print('Hello, World!')", *runTask.Run.Script.InlineCode)
	assert.Equal(t, "example-input", runTask.Run.Script.Input)
}

func boolPtr(b bool) *bool {
	return &b
}

func stringPtr(s string) *string {
	return &s
}
