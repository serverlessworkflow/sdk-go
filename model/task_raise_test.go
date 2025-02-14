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

func TestRaiseTask_MarshalJSON(t *testing.T) {
	raiseTask := &RaiseTask{
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
		Raise: RaiseTaskConfiguration{
			Error: RaiseTaskError{
				Definition: &Error{
					Type:   &URITemplateOrRuntimeExpr{Value: "http://example.com/error"},
					Status: 500,
					Title:  NewStringOrRuntimeExpr("Internal Server Error"),
					Detail: NewStringOrRuntimeExpr("An unexpected error occurred."),
				},
			},
		},
	}

	data, err := json.Marshal(raiseTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"raise": {
			"error": {
				"type": "http://example.com/error",
				"status": 500,
				"title": "Internal Server Error",
				"detail": "An unexpected error occurred."
			}
		}
	}`, string(data))
}

func TestRaiseTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"raise": {
			"error": {
				"type": "http://example.com/error",
				"status": 500,
				"title": "Internal Server Error",
				"detail": "An unexpected error occurred."
			}
		}
	}`

	var raiseTask *RaiseTask
	err := json.Unmarshal([]byte(jsonData), &raiseTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, raiseTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, raiseTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, raiseTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, raiseTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, raiseTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, raiseTask.Metadata)
	assert.Equal(t, "http://example.com/error", raiseTask.Raise.Error.Definition.Type.String())
	assert.Equal(t, 500, raiseTask.Raise.Error.Definition.Status)
	assert.Equal(t, "Internal Server Error", raiseTask.Raise.Error.Definition.Title.String())
	assert.Equal(t, "An unexpected error occurred.", raiseTask.Raise.Error.Definition.Detail.String())
}
