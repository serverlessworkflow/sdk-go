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
	// StateTypeSubflow ...
	StateTypeSubflow = "subflow"
	// StateTypeInject ...
	StateTypeInject = "inject"
	// StateTypeCallback ...
	StateTypeCallback = "callback"

	// CompletionTypeAnd ..
	CompletionTypeAnd = "and"
	// CompletionTypeXor ...
	CompletionTypeXor = "xor"
	// CompletionTypeNOfM ...
	CompletionTypeNOfM = "n_of_m"
)

// StateType ...
type StateType string

// CompletionType Option types on how to complete branch execution.
type CompletionType string

// State definition for a Workflow state
type State interface {
	GetID() string
	GetName() string
	GetType() StateType
	GetOnErrors() []Error
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
	OnErrors []Error `json:"onErrors,omitempty"  validate:"omitempty,dive"`
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
func (s *BaseState) GetOnErrors() []Error { return s.OnErrors }

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

// DelayState Causes the workflow execution to delay for a specified duration
type DelayState struct {
	BaseState
	// Amount of time (ISO 8601 format) to delay
	TimeDelay string `json:"timeDelay" validate:"required"`
}

// EventState This state is used to wait for events from event sources, then consumes them and invoke one or more actions to run in sequence or parallel
type EventState struct {
	BaseState
	// If true consuming one of the defined events causes its associated actions to be performed. If false all of the defined events must be consumed in order for actions to be performed
	Exclusive bool `json:"exclusive,omitempty"`
	// Define the events to be consumed and optional actions to be performed
	OnEvents []OnEvents `json:"onEvents" validate:"required,min=1,dive"`
	// Time period to wait for incoming events (ISO 8601 format)
	Timeout string `json:"timeout,omitempty"`
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

	if eventStateMap["exclusive"] == nil {
		e.Exclusive = true
	} else {
		e.Exclusive = eventStateMap["exclusive"].(bool)
	}
	e.Timeout = requiresNotNilOrEmpty(eventStateMap["timeout"])

	eventStateRaw := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &eventStateRaw); err != nil {
		return err
	}
	if err := json.Unmarshal(eventStateRaw["onEvents"], &e.OnEvents); err != nil {
		return err
	}

	return nil
}

// OperationState Defines actions be performed. Does not wait for incoming events
type OperationState struct {
	BaseState
	// Specifies whether actions are performed in sequence or in parallel
	ActionMode ActionMode `json:"actionMode,omitempty"`
	// Actions to be performed
	Actions []Action `json:"actions" validate:"required,min=1,dive"`
}

// ParallelState Consists of a number of states that are executed in parallel
type ParallelState struct {
	BaseState
	// Branch Definitions
	Branches []Branch `json:"branches" validate:"required,min=1,dive"`
	// Option types on how to complete branch execution.
	CompletionType CompletionType `json:"completionType,omitempty"`
	// Used when completionType is set to 'n_of_m' to specify the 'N' value
	N intstr.IntOrString `json:"n,omitempty"`
}

// SubflowState Defines a sub-workflow to be executed
type SubflowState struct {
	BaseState
	// Workflow execution must wait for sub-workflow to finish before continuing
	WaitForCompletion bool `json:"waitForCompletion,omitempty"`
	// Sub-workflow unique id
	WorkflowID string `json:"workflowId" validate:"required"`
	// SubFlow state repeat exec definition
	Repeat Repeat `json:"repeat,omitempty"`
}

// InjectState ...
type InjectState struct {
	BaseState
	// JSON object which can be set as states data input and can be manipulated via filters
	Data map[string]interface{} `json:"data" validate:"required,min=1"`
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
	Max intstr.IntOrString `json:"max,omitempty"`
	// Actions to be executed for each of the elements of inputCollection
	Actions []Action `json:"actions,omitempty"`
	// Unique Id of a workflow to be executed for each of the elements of inputCollection
	WorkflowID string `json:"workflowId,omitempty"`
}

// CallbackState ...
type CallbackState struct {
	BaseState
	// Defines the action to be executed
	Action Action `json:"action" validate:"required"`
	// References an unique callback event name in the defined workflow events
	EventRef string `json:"eventRef" validate:"required"`
	// Time period to wait for incoming events (ISO 8601 format)
	Timeout string `json:"timeout" validate:"required"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

// BaseSwitchState ...
type BaseSwitchState struct {
	BaseState
	// Default transition of the workflow if there is no matching data conditions. Can include a transition or end definition
	Default DefaultDef `json:"default,omitempty"`
}

// EventBasedSwitchState Permits transitions to other states based on events
type EventBasedSwitchState struct {
	BaseSwitchState
	// Defines conditions evaluated against events
	EventConditions []EventCondition `json:"eventConditions" validate:"required,min=1,dive"`
}

// UnmarshalJSON implementation for json Unmarshal function for the Eventbasedswitch type
func (j *EventBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	eventBasedSwitch := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &eventBasedSwitch)
	if err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	err = json.Unmarshal(eventBasedSwitch["eventConditions"], &rawConditions)
	if err != nil {
		return err
	}

	j.EventConditions = make([]EventCondition, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		err = json.Unmarshal(rawCondition, &mapConditions)
		if err != nil {
			return err
		}
		var condition EventCondition
		if _, ok := mapConditions["end"]; ok {
			condition = &EndEventCondition{}
		} else {
			condition = &TransitionEventCondition{}
		}
		err := json.Unmarshal(rawCondition, condition)
		if err != nil {
			return err
		}
		j.EventConditions[i] = condition
	}
	return nil
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
	// References an unique event name in the defined workflow events
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
	DataConditions []DataCondition `json:"dataConditions" validate:"required,min=1,dive"`
}

// UnmarshalJSON implementation for json Unmarshal function for the Databasedswitch type
func (j *DataBasedSwitchState) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &j.BaseSwitchState); err != nil {
		return err
	}
	dataBasedSwitch := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &dataBasedSwitch)
	if err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	err = json.Unmarshal(dataBasedSwitch["dataConditions"], &rawConditions)
	if err != nil {
		return err
	}

	j.DataConditions = make([]DataCondition, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		err = json.Unmarshal(rawCondition, &mapConditions)
		if err != nil {
			return err
		}
		var condition DataCondition
		if _, ok := mapConditions["end"]; ok {
			condition = &EndDataCondition{}
		} else {
			condition = &TransitionDataCondition{}
		}
		err := json.Unmarshal(rawCondition, condition)
		if err != nil {
			return err
		}
		j.DataConditions[i] = condition
	}
	return nil
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
