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

func TestSetTask_MarshalJSON(t *testing.T) {
	setTask := SetTask{
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
		Set: map[string]interface{}{
			"key1": "value1",
			"key2": 42,
		},
	}

	data, err := json.Marshal(setTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"set": {
			"key1": "value1",
			"key2": 42
		}
	}`, string(data))
}

func TestSetTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"set": {
			"key1": "value1",
			"key2": 42
		}
	}`

	var setTask SetTask
	err := json.Unmarshal([]byte(jsonData), &setTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, setTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, setTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, setTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, setTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, setTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, setTask.Metadata)
	expectedSet := map[string]interface{}{
		"key1": "value1",
		"key2": float64(42), // Match JSON unmarshaling behavior
	}
	assert.Equal(t, expectedSet, setTask.Set)
}

func TestSetTask_Validation(t *testing.T) {
	// Valid SetTask
	setTask := SetTask{
		TaskBase: TaskBase{},
		Set: map[string]interface{}{
			"key": "value",
		},
	}
	assert.NoError(t, validate.Struct(setTask))

	// Invalid SetTask (empty set)
	invalidSetTask := SetTask{
		TaskBase: TaskBase{},
		Set:      map[string]interface{}{},
	}
	assert.Error(t, validate.Struct(invalidSetTask))
}
