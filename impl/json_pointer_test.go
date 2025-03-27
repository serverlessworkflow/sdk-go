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

package impl

import (
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

// TestGenerateJSONPointer_SimpleTask tests a simple workflow task.
func TestGenerateJSONPointer_SimpleTask(t *testing.T) {
	workflow := &model.Workflow{
		Document: model.Document{Name: "simple-workflow"},
		Do: &model.TaskList{
			&model.TaskItem{Key: "task1", Task: &model.SetTask{Set: map[string]interface{}{"value": 10}}},
			&model.TaskItem{Key: "task2", Task: &model.SetTask{Set: map[string]interface{}{"double": "${ .value * 2 }"}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "task2")
	assert.NoError(t, err)
	assert.Equal(t, "/do/1/task2", jsonPointer)
}

// TestGenerateJSONPointer_SimpleTask tests a simple workflow task.
func TestGenerateJSONPointer_Document(t *testing.T) {
	workflow := &model.Workflow{
		Document: model.Document{Name: "simple-workflow"},
		Do: &model.TaskList{
			&model.TaskItem{Key: "task1", Task: &model.SetTask{Set: map[string]interface{}{"value": 10}}},
			&model.TaskItem{Key: "task2", Task: &model.SetTask{Set: map[string]interface{}{"double": "${ .value * 2 }"}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, workflow.Document)
	assert.NoError(t, err)
	assert.Equal(t, "/document", jsonPointer)
}

func TestGenerateJSONPointer_ForkTask(t *testing.T) {
	workflow := &model.Workflow{
		Document: model.Document{Name: "fork-example"},
		Do: &model.TaskList{
			&model.TaskItem{
				Key: "raiseAlarm",
				Task: &model.ForkTask{
					Fork: model.ForkTaskConfiguration{
						Compete: true,
						Branches: &model.TaskList{
							{Key: "callNurse", Task: &model.CallHTTP{Call: "http", With: model.HTTPArguments{Method: "put", Endpoint: model.NewEndpoint("https://hospital.com/api/alert/nurses")}}},
							{Key: "callDoctor", Task: &model.CallHTTP{Call: "http", With: model.HTTPArguments{Method: "put", Endpoint: model.NewEndpoint("https://hospital.com/api/alert/doctor")}}},
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
	workflow := &model.Workflow{
		Document: model.Document{Name: "deep-nested"},
		Do: &model.TaskList{
			&model.TaskItem{
				Key: "step1",
				Task: &model.ForkTask{
					Fork: model.ForkTaskConfiguration{
						Compete: false,
						Branches: &model.TaskList{
							{
								Key: "branchA",
								Task: &model.ForkTask{
									Fork: model.ForkTaskConfiguration{
										Branches: &model.TaskList{
											{
												Key:  "deepTask",
												Task: &model.SetTask{Set: map[string]interface{}{"result": "done"}},
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
	workflow := &model.Workflow{
		Document: model.Document{Name: "nonexistent-test"},
		Do: &model.TaskList{
			&model.TaskItem{Key: "taskA", Task: &model.SetTask{Set: map[string]interface{}{"value": 5}}},
		},
	}

	_, err := GenerateJSONPointer(workflow, "taskX")
	assert.Error(t, err)
}

// TestGenerateJSONPointer_MixedTaskTypes verifies a workflow with different task types.
func TestGenerateJSONPointer_MixedTaskTypes(t *testing.T) {
	workflow := &model.Workflow{
		Document: model.Document{Name: "mixed-tasks"},
		Do: &model.TaskList{
			&model.TaskItem{Key: "compute", Task: &model.SetTask{Set: map[string]interface{}{"result": 42}}},
			&model.TaskItem{Key: "notify", Task: &model.CallHTTP{Call: "http", With: model.HTTPArguments{Method: "post", Endpoint: model.NewEndpoint("https://api.notify.com")}}},
		},
	}

	jsonPointer, err := GenerateJSONPointer(workflow, "notify")
	assert.NoError(t, err)
	assert.Equal(t, "/do/1/notify", jsonPointer)
}
