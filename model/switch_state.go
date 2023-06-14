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
	"strings"

	"github.com/serverlessworkflow/sdk-go/v2/util"
)

type EventConditions []EventCondition

// SwitchState is workflow's gateways: direct transitions onf a workflow based on certain conditions.
type SwitchState struct {
	// TODO: don't use BaseState for this, there are a few fields that SwitchState don't need.

	// Default transition of the workflow if there is no matching data conditions. Can include a transition or
	// end definition.
	DefaultCondition DefaultCondition `json:"defaultCondition"`
	// Defines conditions evaluated against events.
	// +optional
	EventConditions EventConditions `json:"eventConditions" validate:"dive"`
	// Defines conditions evaluated against data
	// +optional
	DataConditions []DataCondition `json:"dataConditions" validate:"dive"`
	// SwitchState specific timeouts
	// +optional
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

	// Avoid marshal empty objects as null.
	st := strings.Replace(string(custom), "\"eventConditions\":null,", "", 1)
	st = strings.Replace(st, "\"dataConditions\":null,", "", 1)
	st = strings.Replace(st, "\"end\":null,", "", -1)
	return []byte(st), err
}

// DefaultCondition Can be either a transition or end definition
type DefaultCondition struct {
	// Serverless workflow states can have one or more incoming and outgoing transitions (from/to other states).
	// Each state can define a transition definition that is used to determine which state to transition to next.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Transition *Transition `json:"transition,omitempty"`
	// 	If this state an end state
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	End *End `json:"end,omitempty"`
}

type defaultConditionUnmarshal DefaultCondition

// UnmarshalJSON implements json.Unmarshaler
func (e *DefaultCondition) UnmarshalJSON(data []byte) error {
	var nextState string
	err := util.UnmarshalPrimitiveOrObject("defaultCondition", data, &nextState, (*defaultConditionUnmarshal)(e))
	if err != nil {
		return err
	}

	if nextState != "" {
		e.Transition = &Transition{NextState: nextState}
	}

	return err
}

// SwitchStateTimeout defines the specific timeout settings for switch state
type SwitchStateTimeout struct {
	// Default workflow state execution timeout (ISO 8601 duration format)
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// Specify the expire value to transitions to defaultCondition. When event-based conditions do not arrive.
	// NOTE: this is only available for EventConditions
	// +optional
	EventTimeout string `json:"eventTimeout,omitempty" validate:"omitempty,iso8601duration"`
}

// EventCondition specify events which the switch state must wait for.
type EventCondition struct {
	// Event condition name.
	// +optional
	Name string `json:"name,omitempty"`
	// References a unique event name in the defined workflow events.
	// +kubebuilder:validation:Required
	EventRef string `json:"eventRef" validate:"required"`
	// Event data filter definition.
	// +optional
	EventDataFilter *EventDataFilter `json:"eventDataFilter,omitempty"`
	// Metadata information.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Metadata Metadata `json:"metadata,omitempty"`
	// TODO End or Transition needs to be exclusive tag, one or another should be set.
	// Explicit transition to end
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	End *End `json:"end" validate:"omitempty"`
	// Workflow transition if condition is evaluated to true
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Transition *Transition `json:"transition" validate:"omitempty"`
}

// DataCondition specify a data-based condition statement which causes a transition to another workflow state
// if evaluated to true.
type DataCondition struct {
	// Data condition name.
	// +optional
	Name string `json:"name,omitempty"`
	// Workflow expression evaluated against state data. Must evaluate to true or false.
	// +kubebuilder:validation:Required
	Condition string `json:"condition" validate:"required"`
	// Metadata information.
	// +optional
	Metadata Metadata `json:"metadata,omitempty"`
	// TODO End or Transition needs to be exclusive tag, one or another should be set.
	// Explicit transition to end
	End *End `json:"end" validate:"omitempty"`
	// Workflow transition if condition is evaluated to true
	Transition *Transition `json:"transition,omitempty" validate:"omitempty"`
}
