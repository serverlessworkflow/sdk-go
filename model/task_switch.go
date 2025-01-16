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

import "encoding/json"

// SwitchTask represents a task configuration for conditional branching.
type SwitchTask struct {
	TaskBase `json:",inline"` // Inline TaskBase fields
	Switch   []SwitchItem     `json:"switch" validate:"required,min=1,dive,switch_item"`
}

type SwitchItem map[string]SwitchCase

// SwitchCase defines a condition and the corresponding outcome for a switch task.
type SwitchCase struct {
	When *RuntimeExpression `json:"when,omitempty"`
	Then *FlowDirective     `json:"then" validate:"required"`
}

// MarshalJSON for SwitchTask to ensure proper serialization.
func (st *SwitchTask) MarshalJSON() ([]byte, error) {
	type Alias SwitchTask
	return json.Marshal((*Alias)(st))
}

// UnmarshalJSON for SwitchTask to ensure proper deserialization.
func (st *SwitchTask) UnmarshalJSON(data []byte) error {
	type Alias SwitchTask
	alias := (*Alias)(st)
	return json.Unmarshal(data, alias)
}
