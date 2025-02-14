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
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGenerateJSONPointer_SimpleTask tests a simple workflow task.
func TestGenerateJSONPointer_SimpleTask(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "simple-workflow"},
		Do: &TaskList{
			&TaskItem{Key: "task1", Task: &SetTask{Set: map[string]interface{}{"value": 10}}},
			&TaskItem{Key: "task2", Task: &SetTask{Set: map[string]interface{}{"double": "${ .value * 2 }"}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "task2")
	assert.NoError(t, err)
	assert.Equal(t, "/do/1/task2", jsonPointer)
}

// TestGenerateJSONPointer_SimpleTask tests a simple workflow task.
func TestGenerateJSONPointer_Document(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "simple-workflow"},
		Do: &TaskList{
			&TaskItem{Key: "task1", Task: &SetTask{Set: map[string]interface{}{"value": 10}}},
			&TaskItem{Key: "task2", Task: &SetTask{Set: map[string]interface{}{"double": "${ .value * 2 }"}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, workflow.Document)
	assert.NoError(t, err)
	assert.Equal(t, "/document", jsonPointer)
}

func TestGenerateJSONPointer_ForkTask(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "fork-example"},
		Do: &TaskList{
			&TaskItem{
				Key: "raiseAlarm",
				Task: &ForkTask{
					Fork: ForkTaskConfiguration{
						Compete: true,
						Branches: &TaskList{
							{Key: "callNurse", Task: &CallHTTP{Call: "http", With: HTTPArguments{Method: "put", Endpoint: NewEndpoint("https://hospital.com/api/alert/nurses")}}},
							{Key: "callDoctor", Task: &CallHTTP{Call: "http", With: HTTPArguments{Method: "put", Endpoint: NewEndpoint("https://hospital.com/api/alert/doctor")}}},
						},
					},
				},
			},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "callDoctor")
	assert.NoError(t, err)
	assert.Equal(t, "/do/0/raiseAlarm/fork/branches/1/callDoctor", jsonPointer)
}

// TestGenerateJSONPointer_DeepNestedTask tests multiple nested task levels.
func TestGenerateJSONPointer_DeepNestedTask(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "deep-nested"},
		Do: &TaskList{
			&TaskItem{
				Key: "step1",
				Task: &ForkTask{
					Fork: ForkTaskConfiguration{
						Compete: false,
						Branches: &TaskList{
							{
								Key: "branchA",
								Task: &ForkTask{
									Fork: ForkTaskConfiguration{
										Branches: &TaskList{
											{
												Key:  "deepTask",
												Task: &SetTask{Set: map[string]interface{}{"result": "done"}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "deepTask")
	assert.NoError(t, err)
	assert.Equal(t, "/do/0/step1/fork/branches/0/branchA/fork/branches/0/deepTask", jsonPointer)
}

// TestGenerateJSONPointer_NonExistentTask checks for a task that doesn't exist.
func TestGenerateJSONPointer_NonExistentTask(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "nonexistent-test"},
		Do: &TaskList{
			&TaskItem{Key: "taskA", Task: &SetTask{Set: map[string]interface{}{"value": 5}}},
		},
	}

	_, err := GenerateJSONPointer(workflow, "taskX")
	assert.Error(t, err)
}

// TestGenerateJSONPointer_MixedTaskTypes verifies a workflow with different task types.
func TestGenerateJSONPointer_MixedTaskTypes(t *testing.T) {
	workflow := &Workflow{
		Document: Document{Name: "mixed-tasks"},
		Do: &TaskList{
			&TaskItem{Key: "compute", Task: &SetTask{Set: map[string]interface{}{"result": 42}}},
			&TaskItem{Key: "notify", Task: &CallHTTP{Call: "http", With: HTTPArguments{Method: "post", Endpoint: NewEndpoint("https://api.notify.com")}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "notify")
	assert.NoError(t, err)
	assert.Equal(t, "/do/1/notify", jsonPointer)
}
