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

// ForTask represents a task configuration to iterate over a collection.
type ForTask struct {
	TaskBase `json:",inline"`     // Inline TaskBase fields
	For      ForTaskConfiguration `json:"for" validate:"required"`
	While    string               `json:"while,omitempty"`
	Do       *TaskList            `json:"do" validate:"required,dive"`
}

func (f *ForTask) GetBase() *TaskBase {
	return &f.TaskBase
}

// ForTaskConfiguration defines the loop configuration for iterating over a collection.
type ForTaskConfiguration struct {
	Each string `json:"each,omitempty"`         // Variable name for the current item
	In   string `json:"in" validate:"required"` // Runtime expression for the collection
	At   string `json:"at,omitempty"`           // Variable name for the current index
}
