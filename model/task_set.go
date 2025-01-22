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

// SetTask represents a task used to set data.
type SetTask struct {
	TaskBase `json:",inline"`       // Inline TaskBase fields
	Set      map[string]interface{} `json:"set" validate:"required,min=1,dive"`
}

// MarshalJSON for SetTask to ensure proper serialization.
func (st *SetTask) MarshalJSON() ([]byte, error) {
	type Alias SetTask
	return json.Marshal((*Alias)(st))
}

// UnmarshalJSON for SetTask to ensure proper deserialization.
func (st *SetTask) UnmarshalJSON(data []byte) error {
	type Alias SetTask
	alias := (*Alias)(st)
	return json.Unmarshal(data, alias)
}
