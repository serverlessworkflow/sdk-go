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

	"sigs.k8s.io/yaml"
)

// WorkflowBuilder helps construct and serialize a Workflow object.
type WorkflowBuilder struct {
	workflow *Workflow
}

// NewWorkflowBuilder initializes a new WorkflowBuilder.
func NewWorkflowBuilder() *WorkflowBuilder {
	return &WorkflowBuilder{
		workflow: &Workflow{
			Document: Document{},
			Do:       &TaskList{},
		},
	}
}

// SetDocument sets the Document fields in the Workflow.
func (wb *WorkflowBuilder) SetDocument(dsl, namespace, name, version string) *WorkflowBuilder {
	wb.workflow.Document.DSL = dsl
	wb.workflow.Document.Namespace = namespace
	wb.workflow.Document.Name = name
	wb.workflow.Document.Version = version
	return wb
}

// AddTask adds a TaskItem to the Workflow's Do list.
func (wb *WorkflowBuilder) AddTask(key string, task Task) *WorkflowBuilder {
	*wb.workflow.Do = append(*wb.workflow.Do, &TaskItem{
		Key:  key,
		Task: task,
	})
	return wb
}

// SetInput sets the Input for the Workflow.
func (wb *WorkflowBuilder) SetInput(input *Input) *WorkflowBuilder {
	wb.workflow.Input = input
	return wb
}

// SetOutput sets the Output for the Workflow.
func (wb *WorkflowBuilder) SetOutput(output *Output) *WorkflowBuilder {
	wb.workflow.Output = output
	return wb
}

// SetTimeout sets the Timeout for the Workflow.
func (wb *WorkflowBuilder) SetTimeout(timeout *TimeoutOrReference) *WorkflowBuilder {
	wb.workflow.Timeout = timeout
	return wb
}

// SetUse sets the Use section for the Workflow.
func (wb *WorkflowBuilder) SetUse(use *Use) *WorkflowBuilder {
	wb.workflow.Use = use
	return wb
}

// SetSchedule sets the Schedule for the Workflow.
func (wb *WorkflowBuilder) SetSchedule(schedule *Schedule) *WorkflowBuilder {
	wb.workflow.Schedule = schedule
	return wb
}

// Build returns the constructed Workflow object.
func (wb *WorkflowBuilder) Build() *Workflow {
	return wb.workflow
}

// ToYAML serializes the Workflow to YAML format.
func (wb *WorkflowBuilder) ToYAML() ([]byte, error) {
	return yaml.Marshal(wb.workflow)
}

// ToJSON serializes the Workflow to JSON format.
func (wb *WorkflowBuilder) ToJSON() ([]byte, error) {
	return json.MarshalIndent(wb.workflow, "", "  ")
}
