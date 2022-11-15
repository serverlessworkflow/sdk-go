// Copyright 2021 The Serverless Workflow Specification Authors
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

const (
	// StateTypeDelay ...
	StateTypeDelay = "delay"
	// StateTypeEvent ...
	StateTypeEvent = "event"
	// StateTypeOperation ...
	StateTypeOperation = "operation"
	// StateTypeParallel ...
	StateTypeParallel = "parallel"
	// StateTypeSwitch ...
	StateTypeSwitch = "switch"
	// StateTypeForEach ...
	StateTypeForEach = "foreach"
	// StateTypeInject ...
	StateTypeInject = "inject"
	// StateTypeCallback ...
	StateTypeCallback = "callback"
	// StateTypeSleep ...
	StateTypeSleep = "sleep"
)

func getActionsModelMapping(stateType string, s map[string]interface{}) (State, bool) {
	switch stateType {
	case StateTypeDelay:
		return &DelayState{}, true
	case StateTypeEvent:
		return &EventState{}, true
	case StateTypeOperation:
		return &OperationState{}, true
	case StateTypeParallel:
		return &ParallelState{}, true
	case StateTypeSwitch:
		if _, ok := s["dataConditions"]; ok {
			return &DataBasedSwitchState{}, true
		}
		return &EventBasedSwitchState{}, true
	case StateTypeInject:
		return &InjectState{}, true
	case StateTypeForEach:
		return &ForEachState{}, true
	case StateTypeCallback:
		return &CallbackState{}, true
	case StateTypeSleep:
		return &SleepState{}, true
	}
	return nil, false
}

// StateType ...
type StateType string

// State definition for a Workflow state
type State interface {
	GetID() string
	GetName() string
	GetType() StateType
	GetOnErrors() []OnError
	GetTransition() *Transition
	GetStateDataFilter() *StateDataFilter
	GetCompensatedBy() string
	GetUsedForCompensation() bool
	GetEnd() *End
	GetMetadata() *Metadata
}

// BaseState ...
type BaseState struct {
	// Unique State id
	ID string `json:"id,omitempty"`
	// State name
	Name string `json:"name" validate:"required"`
	// State type
	Type StateType `json:"type" validate:"required"`
	// States error handling and retries definitions
	OnErrors []OnError `json:"onErrors,omitempty"  validate:"omitempty,dive"`
	// Next transition of the workflow after the time delay
	Transition *Transition `json:"transition,omitempty"`
	// State data filter
	StateDataFilter *StateDataFilter `json:"stateDataFilter,omitempty"`
	// Unique Name of a workflow state which is responsible for compensation of this state
	CompensatedBy string `json:"compensatedBy,omitempty"`
	// If true, this state is used to compensate another state. Default is false
	UsedForCompensation bool `json:"usedForCompensation,omitempty"`
	// State end definition
	End      *End      `json:"end,omitempty"`
	Metadata *Metadata `json:"metadata,omitempty"`
}

// GetOnErrors ...
func (s *BaseState) GetOnErrors() []OnError { return s.OnErrors }

// GetCompensatedBy ...
func (s *BaseState) GetCompensatedBy() string { return s.CompensatedBy }

// GetTransition ...
func (s *BaseState) GetTransition() *Transition { return s.Transition }

// GetUsedForCompensation ...
func (s *BaseState) GetUsedForCompensation() bool { return s.UsedForCompensation }

// GetEnd ...
func (s *BaseState) GetEnd() *End { return s.End }

// GetID ...
func (s *BaseState) GetID() string { return s.ID }

// GetName ...
func (s *BaseState) GetName() string { return s.Name }

// GetType ...
func (s *BaseState) GetType() StateType { return s.Type }

// GetStateDataFilter ...
func (s *BaseState) GetStateDataFilter() *StateDataFilter { return s.StateDataFilter }

// GetMetadata ...
func (s *BaseState) GetMetadata() *Metadata { return s.Metadata }

// BaseSwitchState ...
type BaseSwitchState struct {
	BaseState
	// Default transition of the workflow if there is no matching data conditions. Can include a transition or end definition
	DefaultCondition DefaultCondition `json:"defaultCondition,omitempty"`
}

// EventBasedSwitchState Permits transitions to other states based on events
type EventBasedSwitchState struct {
	BaseSwitchState
	// Defines conditions evaluated against events
	EventConditions []EventCondition `json:"eventConditions" validate:"required,min=1,dive"`
	// State specific timeouts
	Timeouts *EventBasedSwitchStateTimeout `json:"timeouts,omitempty"`
}

// UnmarshalJSON implementation for json Unmarshal function for the EventBasedSwitch type
func (j *EventBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	eventBasedSwitch := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &eventBasedSwitch); err != nil {
		return err
	}

	eventBaseTimeoutsRawMessage, ok := eventBasedSwitch["timeouts"]
	if ok {
		if err := json.Unmarshal(eventBaseTimeoutsRawMessage, &j.Timeouts); err != nil {
			return err
		}
	}

	var rawConditions []json.RawMessage
	if err := json.Unmarshal(eventBasedSwitch["eventConditions"], &rawConditions); err != nil {
		return err
	}

	j.EventConditions = make([]EventCondition, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		if err := json.Unmarshal(rawCondition, &mapConditions); err != nil {
			return err
		}
		var condition EventCondition
		if _, ok := mapConditions["end"]; ok {
			condition = &EndEventCondition{}
		} else {
			condition = &TransitionEventCondition{}
		}
		if err := json.Unmarshal(rawCondition, condition); err != nil {
			return err
		}
		j.EventConditions[i] = condition
	}

	return nil
}

// EventBasedSwitchStateTimeout ...
type EventBasedSwitchStateTimeout struct {
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	EventTimeout     string            `json:"eventTimeout,omitempty"`
}

// EventCondition ...
type EventCondition interface {
	GetName() string
	GetEventRef() string
	GetEventDataFilter() EventDataFilter
	GetMetadata() Metadata
}

// BaseEventCondition ...
type BaseEventCondition struct {
	// Event condition name
	Name string `json:"name,omitempty"`
	// References a unique event name in the defined workflow events
	EventRef string `json:"eventRef" validate:"required"`
	// Event data filter definition
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
	Metadata        Metadata        `json:"metadata,omitempty"`
}

// GetEventRef ...
func (e *BaseEventCondition) GetEventRef() string { return e.EventRef }

// GetEventDataFilter ...
func (e *BaseEventCondition) GetEventDataFilter() EventDataFilter { return e.EventDataFilter }

// GetMetadata ...
func (e *BaseEventCondition) GetMetadata() Metadata { return e.Metadata }

// GetName ...
func (e *BaseEventCondition) GetName() string { return e.Name }

// TransitionEventCondition Switch state data event condition
type TransitionEventCondition struct {
	BaseEventCondition
	// Next transition of the workflow if there is valid matches
	Transition Transition `json:"transition" validate:"required"`
}

// EndEventCondition Switch state data event condition
type EndEventCondition struct {
	BaseEventCondition
	// Explicit transition to end
	End End `json:"end" validate:"required"`
}

// DataBasedSwitchState Permits transitions to other states based on data conditions
type DataBasedSwitchState struct {
	BaseSwitchState
	DataConditions []DataCondition              `json:"dataConditions" validate:"required,min=1,dive"`
	Timeouts       *DataBasedSwitchStateTimeout `json:"timeouts,omitempty"`
}

// UnmarshalJSON implementation for json Unmarshal function for the DataBasedSwitch type
func (j *DataBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	dataBasedSwitch := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &dataBasedSwitch); err != nil {
		return err
	}
	if err := json.Unmarshal(data, &j.Timeouts); err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	if err := json.Unmarshal(dataBasedSwitch["dataConditions"], &rawConditions); err != nil {
		return err
	}
	j.DataConditions = make([]DataCondition, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		if err := json.Unmarshal(rawCondition, &mapConditions); err != nil {
			return err
		}
		var condition DataCondition
		if _, ok := mapConditions["end"]; ok {
			condition = &EndDataCondition{}
		} else {
			condition = &TransitionDataCondition{}
		}
		if err := json.Unmarshal(rawCondition, condition); err != nil {
			return err
		}
		j.DataConditions[i] = condition
	}
	return nil
}

// DataBasedSwitchStateTimeout ...
type DataBasedSwitchStateTimeout struct {
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
}

// DataCondition ...
type DataCondition interface {
	GetName() string
	GetCondition() string
	GetMetadata() Metadata
}

// BaseDataCondition ...
type BaseDataCondition struct {
	// Data condition name
	Name string `json:"name,omitempty"`
	// Workflow expression evaluated against state data. Must evaluate to true or false
	Condition string   `json:"condition" validate:"required"`
	Metadata  Metadata `json:"metadata,omitempty"`
}

// GetName ...
func (b *BaseDataCondition) GetName() string { return b.Name }

// GetCondition ...
func (b *BaseDataCondition) GetCondition() string { return b.Condition }

// GetMetadata ...
func (b *BaseDataCondition) GetMetadata() Metadata { return b.Metadata }

// TransitionDataCondition ...
type TransitionDataCondition struct {
	BaseDataCondition
	// Workflow transition if condition is evaluated to true
	Transition Transition `json:"transition" validate:"required"`
}

// EndDataCondition ...
type EndDataCondition struct {
	BaseDataCondition
	// Workflow end definition
	End End `json:"end" validate:"required"`
}
