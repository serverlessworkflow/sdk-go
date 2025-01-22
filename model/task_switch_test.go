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

func TestSwitchTask_MarshalJSON(t *testing.T) {
	switchTask := &SwitchTask{
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
		Switch: []SwitchItem{
			{
				"case1": SwitchCase{
					When: &RuntimeExpression{Value: "${condition1}"},
					Then: &FlowDirective{Value: "next"},
				},
			},
			{
				"case2": SwitchCase{
					When: &RuntimeExpression{Value: "${condition2}"},
					Then: &FlowDirective{Value: "end"},
				},
			},
		},
	}

	data, err := json.Marshal(switchTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"switch": [
			{
				"case1": {
					"when": "${condition1}",
					"then": "next"
				}
			},
			{
				"case2": {
					"when": "${condition2}",
					"then": "end"
				}
			}
		]
	}`, string(data))
}

func TestSwitchTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"switch": [
			{
				"case1": {
					"when": "${condition1}",
					"then": "next"
				}
			},
			{
				"case2": {
					"when": "${condition2}",
					"then": "end"
				}
			}
		]
	}`

	var switchTask SwitchTask
	err := json.Unmarshal([]byte(jsonData), &switchTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, switchTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, switchTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, switchTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, switchTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, switchTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, switchTask.Metadata)
	assert.Equal(t, 2, len(switchTask.Switch))
	assert.Equal(t, &RuntimeExpression{Value: "${condition1}"}, switchTask.Switch[0]["case1"].When)
	assert.Equal(t, &FlowDirective{Value: "next"}, switchTask.Switch[0]["case1"].Then)
	assert.Equal(t, &RuntimeExpression{Value: "${condition2}"}, switchTask.Switch[1]["case2"].When)
	assert.Equal(t, &FlowDirective{Value: "end"}, switchTask.Switch[1]["case2"].Then)
}

func TestSwitchTask_Validation(t *testing.T) {
	// Valid SwitchTask
	switchTask := SwitchTask{
		TaskBase: TaskBase{},
		Switch: []SwitchItem{
			{
				"case1": SwitchCase{
					When: &RuntimeExpression{Value: "${condition1}"},
					Then: &FlowDirective{Value: "next"},
				},
			},
		},
	}
	assert.NoError(t, validate.Struct(switchTask))

	// Invalid SwitchTask (empty switch)
	invalidSwitchTask := SwitchTask{
		TaskBase: TaskBase{},
		Switch:   []SwitchItem{},
	}
	assert.Error(t, validate.Struct(invalidSwitchTask))

	// Invalid SwitchTask (SwitchItem with multiple keys)
	invalidSwitchItemTask := SwitchTask{
		TaskBase: TaskBase{},
		Switch: []SwitchItem{
			{
				"case1": SwitchCase{When: &RuntimeExpression{Value: "${condition1}"}, Then: &FlowDirective{Value: "next"}},
				"case2": SwitchCase{When: &RuntimeExpression{Value: "${condition2}"}, Then: &FlowDirective{Value: "end"}},
			},
		},
	}
	assert.Error(t, validate.Struct(invalidSwitchItemTask))
}
