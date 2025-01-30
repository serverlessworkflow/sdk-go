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

// ForkTask represents a task configuration to execute multiple tasks concurrently.
type ForkTask struct {
	TaskBase `json:",inline"`      // Inline TaskBase fields
	Fork     ForkTaskConfiguration `json:"fork" validate:"required"`
}

func (f *ForkTask) GetBase() *TaskBase {
	return &f.TaskBase
}

// ForkTaskConfiguration defines the configuration for the branches to perform concurrently.
type ForkTaskConfiguration struct {
	Branches *TaskList `json:"branches" validate:"required,dive"`
	Compete  bool      `json:"compete,omitempty"`
}
