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

func TestWaitTask_MarshalJSON(t *testing.T) {
	waitTask := &WaitTask{
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
		Wait: NewDurationExpr("P1DT1H"),
	}

	data, err := json.Marshal(waitTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"wait": "P1DT1H"
	}`, string(data))
}

func TestWaitTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"wait": "P1DT1H"
	}`

	waitTask := &WaitTask{}
	err := json.Unmarshal([]byte(jsonData), waitTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, waitTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, waitTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, waitTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, waitTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, waitTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, waitTask.Metadata)
	assert.Equal(t, NewDurationExpr("P1DT1H"), waitTask.Wait)
}

func TestWaitTask_Validation(t *testing.T) {
	// Valid WaitTask
	waitTask := &WaitTask{
		TaskBase: TaskBase{},
		Wait:     NewDurationExpr("P1DT1H"),
	}
	assert.NoError(t, validate.Struct(waitTask))

	// Invalid WaitTask (empty wait)
	invalidWaitTask := &WaitTask{
		TaskBase: TaskBase{},
	}
	assert.Error(t, validate.Struct(invalidWaitTask))
}
