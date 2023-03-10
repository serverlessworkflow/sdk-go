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

// SwitchState is workflow's gateways: direct transitions onf a workflow based on certain conditions.
type SwitchState struct {
	// TODO: don't use BaseState for this, there are a few fields that SwitchState don't need.

	// Default transition of the workflow if there is no matching data conditions. Can include a transition or end definition
	// Required
	DefaultCondition DefaultCondition `json:"defaultCondition"`
	// Defines conditions evaluated against events
	EventConditions []EventCondition `json:"eventConditions" validate:"omitempty,min=1,dive"`
	// Defines conditions evaluated against data
	DataConditions []DataCondition `json:"dataConditions" validate:"omitempty,min=1,dive"`
	// SwitchState specific timeouts
	Timeouts *SwitchStateTimeout `json:"timeouts,omitempty"`
}

func (s *SwitchState) MarshalJSON() ([]byte, error) {
	type Alias SwitchState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *SwitchStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(s),
		Timeouts: s.Timeouts,
	})
	return custom, err
}

// DefaultCondition Can be either a transition or end definition
type DefaultCondition struct {
	Transition *Transition `json:"transition,omitempty"`
	End        *End        `json:"end,omitempty"`
}

// SwitchStateTimeout defines the specific timeout settings for switch state
type SwitchStateTimeout struct {
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`

	// EventTimeout specify the expire value to transitions to defaultCondition
	// when event-based conditions do not arrive.
	// NOTE: this is only available for EventConditions
	EventTimeout string `json:"eventTimeout,omitempty" validate:"omitempty,iso8601duration"`
}

// EventCondition specify events which the switch state must wait for.
type EventCondition struct {
	// Event condition name
	Name string `json:"name,omitempty"`
	// References a unique event name in the defined workflow events
	EventRef string `json:"eventRef" validate:"required"`
	// Event data filter definition
	EventDataFilter *EventDataFilter `json:"eventDataFilter,omitempty"`
	Metadata        Metadata         `json:"metadata,omitempty"`
	// Explicit transition to end
	End *End `json:"end" validate:"omitempty"`
	// Workflow transition if condition is evaluated to true
	Transition *Transition `json:"transition" validate:"omitempty"`
}

// DataCondition specify a data-based condition statement which causes a transition to another workflow state
// if evaluated to true.
type DataCondition struct {
	// Data condition name
	Name string `json:"name,omitempty"`
	// Workflow expression evaluated against state data. Must evaluate to true or false
	Condition string   `json:"condition" validate:"required"`
	Metadata  Metadata `json:"metadata,omitempty"`

	// Explicit transition to end
	End *End `json:"end" validate:"omitempty"`
	// Workflow transition if condition is evaluated to true
	Transition *Transition `json:"transition" validate:"omitempty"`
}
