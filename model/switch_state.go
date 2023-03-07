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
	"context"
	"encoding/json"
	"reflect"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"

	validator "github.com/go-playground/validator/v10"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(SwitchStateStructLevelValidation, SwitchState{})
	val.GetValidator().RegisterStructValidationCtx(DefaultConditionStructLevelValidation, DefaultCondition{})
	val.GetValidator().RegisterStructValidationCtx(EventConditionStructLevelValidation, EventCondition{})
	val.GetValidator().RegisterStructValidationCtx(DataConditionStructLevelValidation, DataCondition{})
}

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

// SwitchStateStructLevelValidation custom validator for SwitchState
func SwitchStateStructLevelValidation(ctx context.Context, structLevel validator.StructLevel) {
	switchState := structLevel.Current().Interface().(SwitchState)
	switch {
	case len(switchState.DataConditions) == 0 && len(switchState.EventConditions) == 0:
		structLevel.ReportError(reflect.ValueOf(switchState), "DataConditions", "dataConditions", "required", "must have one of dataConditions, eventConditions")
	case len(switchState.DataConditions) > 0 && len(switchState.EventConditions) > 0:
		structLevel.ReportError(reflect.ValueOf(switchState), "DataConditions", "dataConditions", "exclusive", "must have one of dataConditions, eventConditions")
	}
}

// DefaultCondition Can be either a transition or end definition
type DefaultCondition struct {
	Transition *Transition `json:"transition,omitempty"`
	End        *End        `json:"end,omitempty"`
}

// DefaultConditionStructLevelValidation custom validator for DefaultCondition
func DefaultConditionStructLevelValidation(ctx context.Context, structLevel validator.StructLevel) {
	defaultCondition := structLevel.Current().Interface().(DefaultCondition)
	switch {
	case defaultCondition.End == nil && defaultCondition.Transition == nil:
		structLevel.ReportError(reflect.ValueOf(defaultCondition), "Transition", "transition", "required", "must have one of transition, end")
	case defaultCondition.Transition != nil && defaultCondition.End != nil:
		structLevel.ReportError(reflect.ValueOf(defaultCondition), "Transition", "transition", "exclusive", "must have one of transition, end")
	}
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

// EventConditionStructLevelValidation custom validator for EventCondition
func EventConditionStructLevelValidation(ctx context.Context, structLevel validator.StructLevel) {
	eventCondition := structLevel.Current().Interface().(EventCondition)
	switch {
	case eventCondition.End == nil && eventCondition.Transition == nil:
		structLevel.ReportError(reflect.ValueOf(eventCondition), "Transition", "transition", "required", "must have one of transition, end")
	case eventCondition.Transition != nil && eventCondition.End != nil:
		structLevel.ReportError(reflect.ValueOf(eventCondition), "Transition", "transition", "exclusive", "must have one of transition, end")
	}
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

// DataConditionStructLevelValidation custom validator for DataCondition
func DataConditionStructLevelValidation(ctx context.Context, structLevel validator.StructLevel) {
	dataCondition := structLevel.Current().Interface().(DataCondition)
	switch {
	case dataCondition.End == nil && dataCondition.Transition == nil:
		structLevel.ReportError(reflect.ValueOf(dataCondition), "Transition", "transition", "required", "must have one of transition, end")
	case dataCondition.Transition != nil && dataCondition.End != nil:
		structLevel.ReportError(reflect.ValueOf(dataCondition), "Transition", "transition", "exclusive", "must have one of transition, end")
	}
}
