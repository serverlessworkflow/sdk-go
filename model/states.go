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

	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	//StateTypeDelay ...
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

	// CompletionTypeAllOf ...
	CompletionTypeAllOf CompletionType = "allOf"
	// CompletionTypeAtLeast ...
	CompletionTypeAtLeast CompletionType = "atLeast"

	// ForEachModeTypeSequential ...
	ForEachModeTypeSequential ForEachModeType = "sequential"
	// ForEachModeTypeParallel ...
	ForEachModeTypeParallel ForEachModeType = "parallel"
)

// StateType ...
type StateType string

// CompletionType Option types on how to complete branch execution.
type CompletionType string

// ForEachModeType Specifies how iterations are to be performed (sequentially or in parallel)
type ForEachModeType string

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
	ID string `json:"id,omitempty" validate:"omitempty,min=1"`
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
	CompensatedBy string `json:"compensatedBy,omitempty" validate:"omitempty,min=1"`
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

// EventState This state is used to wait for events from event sources, then consumes them and invoke one or more actions to run in sequence or parallel
type EventState struct {
	BaseState
	// If true consuming one of the defined events causes its associated actions to be performed. If false all of the defined events must be consumed in order for actions to be performed
	Exclusive bool `json:"exclusive,omitempty"`
	// Define the events to be consumed and optional actions to be performed
	OnEvents []OnEvents `json:"onEvents" validate:"required,min=1,dive"`
	// State specific timeouts
	Timeout *EventStateTimeout `json:"timeouts,omitempty"`
}

// UnmarshalJSON ...
func (e *EventState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &e.BaseState); err != nil {
		return err
	}

	eventStateMap := make(map[string]interface{})
	if err := json.Unmarshal(data, &eventStateMap); err != nil {
		return err
	}

	e.Exclusive = true

	if eventStateMap["exclusive"] != nil {
		exclusiveVal, ok := eventStateMap["exclusive"].(bool)
		if ok {
			e.Exclusive = exclusiveVal
		}
	}

	eventStateRaw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &eventStateRaw); err != nil {
		return err
	}
	if err := json.Unmarshal(eventStateRaw["onEvents"], &e.OnEvents); err != nil {
		return err
	}
	if err := unmarshalKey("timeouts", eventStateRaw, &e.Timeout); err != nil {
		return err
	}

	return nil
}

// EventStateTimeout ...
type EventStateTimeout struct {
	StateExecTimeout  StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string           `json:"actionExecTimeout,omitempty"`
	EventTimeout      string           `json:"eventTimeout,omitempty"`
}

// OperationState Defines actions be performed. Does not wait for incoming events
type OperationState struct {
	BaseState
	// Specifies whether actions are performed in sequence or in parallel
	ActionMode ActionMode `json:"actionMode,omitempty"`
	// Actions to be performed
	Actions []Action `json:"actions" validate:"required,min=1,dive"`
	// State specific timeouts
	Timeouts *OperationStateTimeout `json:"timeouts,omitempty"`
}

// OperationStateTimeout ...
type OperationStateTimeout struct {
	StateExecTimeout  StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string           `json:"actionExecTimeout,omitempty" validate:"omitempty,min=1"`
}

// ParallelState Consists of a number of states that are executed in parallel
type ParallelState struct {
	BaseState
	// Branch Definitions
	Branches []Branch `json:"branches" validate:"required,min=1,dive"`
	// Option types on how to complete branch execution.
	CompletionType CompletionType `json:"completionType,omitempty"`
	// Used when completionType is set to 'atLeast' to specify the minimum number of branches that must complete before the state will transition."
	NumCompleted intstr.IntOrString `json:"numCompleted,omitempty"`
	// State specific timeouts
	Timeouts *ParallelStateTimeout `json:"timeouts,omitempty"`
}

// ParallelStateTimeout ...
type ParallelStateTimeout struct {
	StateExecTimeout  StateExecTimeout `json:"stateExecTimeout,omitempty"`
	BranchExecTimeout string           `json:"branchExecTimeout,omitempty" validate:"omitempty,min=1"`
}

// InjectState ...
type InjectState struct {
	BaseState
	// JSON object which can be set as states data input and can be manipulated via filters
	Data map[string]interface{} `json:"data" validate:"required,min=1"`
	// State specific timeouts
	Timeouts *InjectStateTimeout `json:"timeouts,omitempty"`
}

// InjectStateTimeout ...
type InjectStateTimeout struct {
	StateExecTimeout StateExecTimeout `json:"stateExecTimeout,omitempty"`
}

// ForEachState ...
type ForEachState struct {
	BaseState
	// Workflow expression selecting an array element of the states data
	InputCollection string `json:"inputCollection" validate:"required"`
	// Workflow expression specifying an array element of the states data to add the results of each iteration
	OutputCollection string `json:"outputCollection,omitempty"`
	// Name of the iteration parameter that can be referenced in actions/workflow. For each parallel iteration, this param should contain an unique element of the inputCollection array
	IterationParam string `json:"iterationParam" validate:"required"`
	// Specifies how upper bound on how many iterations may run in parallel
	BatchSize intstr.IntOrString `json:"batchSize,omitempty"`
	// Actions to be executed for each of the elements of inputCollection
	Actions []Action `json:"actions,omitempty"`
	// State specific timeout
	Timeouts *ForEachStateTimeout `json:"timeouts,omitempty"`
	// Mode Specifies how iterations are to be performed (sequentially or in parallel)
	Mode ForEachModeType `json:"mode,omitempty"`
}

// ForEachStateTimeout ...
type ForEachStateTimeout struct {
	StateExecTimeout  StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string           `json:"actionExecTimeout,omitempty"`
}

// CallbackState ...
type CallbackState struct {
	BaseState
	// Defines the action to be executed
	Action Action `json:"action" validate:"required"`
	// References a unique callback event name in the defined workflow events
	EventRef string `json:"eventRef" validate:"required"`
	// Time period to wait for incoming events (ISO 8601 format)
	Timeouts CallbackStateTimeout `json:"timeouts" validate:"required"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

// CallbackStateTimeout ...
type CallbackStateTimeout struct {
	StateExecTimeout  StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string           `json:"actionExecTimeout,omitempty"`
	EventTimeout      string           `json:"eventTimeout,omitempty"`
}

// SleepState ...
type SleepState struct {
	BaseState
	// Duration (ISO 8601 duration format) to sleep
	Duration string `json:"duration" validate:"required"`
	// Timeouts State specific timeouts
	Timeouts *SleepStateTimeout `json:"timeouts,omitempty"`
}

// SleepStateTimeout ...
type SleepStateTimeout struct {
	StateExecTimeout StateExecTimeout `json:"stateExecTimeout,omitempty"`
}

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

// UnmarshalJSON implementation for json Unmarshal function for the Eventbasedswitch type
func (j *EventBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	eventBasedSwitch := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &eventBasedSwitch); err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	if err := unmarshalKey("timeouts", eventBasedSwitch, &j.Timeouts); err != nil {
		return err
	}
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
	StateExecTimeout StateExecTimeout `json:"stateExecTimeout,omitempty"`
	EventTimeout     string           `json:"eventTimeout,omitempty"`
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

// UnmarshalJSON implementation for json Unmarshal function for the Databasedswitch type
func (j *DataBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	dataBasedSwitch := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &dataBasedSwitch); err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	if err := unmarshalKey("timeouts", dataBasedSwitch, &j.Timeouts); err != nil {
		return err
	}
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
	StateExecTimeout StateExecTimeout `json:"stateExecTimeout,omitempty"`
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
