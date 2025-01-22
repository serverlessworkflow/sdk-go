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

func TestEmitTask_MarshalJSON(t *testing.T) {
	emitTask := &EmitTask{
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
		Emit: EmitTaskConfiguration{
			Event: EmitEventDefinition{
				With: &EventProperties{
					ID:              "event-id",
					Source:          &URITemplateOrRuntimeExpr{Value: "http://example.com/source"},
					Type:            "example.event.type",
					Time:            &StringOrRuntimeExpr{Value: "2023-01-01T00:00:00Z"},
					Subject:         "example.subject",
					DataContentType: "application/json",
					DataSchema:      &URITemplateOrRuntimeExpr{Value: "http://example.com/schema"},
					Additional: map[string]interface{}{
						"extra": "value",
					},
				},
			},
		},
	}

	data, err := json.Marshal(emitTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"emit": {
			"event": {
				"with": {
					"id": "event-id",
					"source": "http://example.com/source",
					"type": "example.event.type",
					"time": "2023-01-01T00:00:00Z",
					"subject": "example.subject",
					"datacontenttype": "application/json",
					"dataschema": "http://example.com/schema",
					"extra": "value"
				}
			}
		}
	}`, string(data))
}

func TestEmitTask_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"emit": {
			"event": {
				"with": {
					"id": "event-id",
					"source": "http://example.com/source",
					"type": "example.event.type",
					"time": "2023-01-01T00:00:00Z",
					"subject": "example.subject",
					"datacontenttype": "application/json",
					"dataschema": "http://example.com/schema",
					"extra": "value"
				}
			}
		}
	}`

	var emitTask EmitTask
	err := json.Unmarshal([]byte(jsonData), &emitTask)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, emitTask.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, emitTask.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, emitTask.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, emitTask.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, emitTask.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, emitTask.Metadata)
	assert.Equal(t, "event-id", emitTask.Emit.Event.With.ID)
	assert.Equal(t, "http://example.com/source", emitTask.Emit.Event.With.Source.String())
	assert.Equal(t, "example.event.type", emitTask.Emit.Event.With.Type)
	assert.Equal(t, "2023-01-01T00:00:00Z", emitTask.Emit.Event.With.Time.String())
	assert.Equal(t, "example.subject", emitTask.Emit.Event.With.Subject)
	assert.Equal(t, "application/json", emitTask.Emit.Event.With.DataContentType)
	assert.Equal(t, "http://example.com/schema", emitTask.Emit.Event.With.DataSchema.String())
	assert.Equal(t, map[string]interface{}{"extra": "value"}, emitTask.Emit.Event.With.Additional)
}

func TestListenTask_MarshalJSON_WithUntilCondition(t *testing.T) {
	listenTask := ListenTask{
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
		Listen: ListenTaskConfiguration{
			To: &EventConsumptionStrategy{
				Any: []*EventFilter{
					{
						With: &EventProperties{
							Type:   "example.event.type",
							Source: &URITemplateOrRuntimeExpr{Value: "http://example.com/source"},
						},
					},
				},
				Until: &EventConsumptionUntil{
					Condition: NewRuntimeExpression("workflow.data.condition == true"),
				},
			},
		},
	}

	data, err := json.Marshal(listenTask)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"listen": {
			"to": {
				"any": [
					{
						"with": {
							"type": "example.event.type",
							"source": "http://example.com/source"
						}
					}
				],
				"until": "workflow.data.condition == true"
			}
		}
	}`, string(data))
}

func TestEventConsumptionUntil_MarshalJSON(t *testing.T) {
	tests := []struct {
		name      string
		until     *EventConsumptionUntil
		expected  string
		shouldErr bool
	}{
		{
			name: "Until Disabled",
			until: &EventConsumptionUntil{
				IsDisabled: true,
			},
			expected:  `false`,
			shouldErr: false,
		},
		{
			name: "Until Condition Set",
			until: &EventConsumptionUntil{
				Condition: &RuntimeExpression{Value: "workflow.data.condition == true"},
			},
			expected:  `"workflow.data.condition == true"`,
			shouldErr: false,
		},
		{
			name: "Until Nested Strategy",
			until: &EventConsumptionUntil{
				Strategy: &EventConsumptionStrategy{
					One: &EventFilter{
						With: &EventProperties{Type: "example.event.type"},
					},
				},
			},
			expected:  `{"one":{"with":{"type":"example.event.type"}}}`,
			shouldErr: false,
		},
		{
			name:      "Until Nil",
			until:     &EventConsumptionUntil{},
			expected:  `null`,
			shouldErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.until)
			if test.shouldErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, test.expected, string(data))
			}
		})
	}
}
