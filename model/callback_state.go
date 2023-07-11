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

// CallbackState executes a function and waits for callback event that indicates completion of the task.
type CallbackState struct {
	// Defines the action to be executed.
	// +kubebuilder:validation:Required
	Action Action `json:"action"`
	// References a unique callback event name in the defined workflow events.
	// +kubebuilder:validation:Required
	EventRef string `json:"eventRef" validate:"required"`
	// Time period to wait for incoming events (ISO 8601 format)
	// +optional
	Timeouts *CallbackStateTimeout `json:"timeouts,omitempty"`
	// Event data filter definition.
	// +optional
	EventDataFilter *EventDataFilter `json:"eventDataFilter,omitempty"`
}

func (c *CallbackState) MarshalJSON() ([]byte, error) {
	type Alias CallbackState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *CallbackStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(c),
		Timeouts: c.Timeouts,
	})
	return custom, err
}

// CallbackStateTimeout defines timeout settings for callback state
type CallbackStateTimeout struct {
	// Default workflow state execution timeout (ISO 8601 duration format)
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// Default single actions definition execution timeout (ISO 8601 duration format)
	// +optional
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
	// Default timeout for consuming defined events (ISO 8601 duration format)
	// +optional
	EventTimeout string `json:"eventTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
