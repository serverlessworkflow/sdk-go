// Copyright 2022 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
)

// InjectState used to inject static data into state data input.
type InjectState struct {
	// JSON object which can be set as state's data input and can be manipulated via filter
	// +kubebuilder:validation:MinProperties=1
	Data map[string]Object `json:"data" validate:"required,min=1"`
	// State specific timeouts
	// +optional
	Timeouts *InjectStateTimeout `json:"timeouts,omitempty"`
}

func (i *InjectState) MarshalJSON() ([]byte, error) {
	type Alias InjectState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *InjectStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(i),
		Timeouts: i.Timeouts,
	})
	return custom, err
}

// InjectStateTimeout defines timeout settings for inject state
type InjectStateTimeout struct {
	// Default workflow state execution timeout (ISO 8601 duration format)
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
}
