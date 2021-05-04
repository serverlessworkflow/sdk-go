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

import "k8s.io/apimachinery/pkg/util/intstr"

const (
	StateTypeDelay     = "delay"
	StateTypeEvent     = "event"
	StateTypeOperation = "operation"
	StateTypeParallel  = "parallel"
	StateTypeSwitch    = "switch"
	StateTypeForEach   = "foreach"
	StateTypeSubflow   = "subflow"
	StateTypeInject    = "inject"
	StateTypeCallback  = "callback"

	CompletionTypeAnd  = "and"
	CompletionTypeXor  = "xor"
	CompletionTypeNOfM = "n_of_m"
)

type StateType string

// Option types on how to complete branch execution.
type CompletionType string

// State definition for a Workflow state
type State interface {
	GetID() string
	GetName() string
	GetType() StateType
	GetOnErrors() []Error
	GetTransition() Transition
	GetStateDataFilter() StateDataFilter
	GetCompensatedBy() string
	GetUsedForCompensation() bool
	GetEnd() End
	GetMetadata() Metadata
}

type BaseState struct {
	// Unique State id
	ID string `json:"id,omitempty"`
	// State name
	Name string `json:"name"`
	// State type
	Type StateType `json:"type"`
	// States error handling and retries definitions
	OnErrors []Error `json:"onErrors,omitempty"`
	// Next transition of the workflow after the time delay
	Transition Transition `json:"transition,omitempty"`
	// State data filter
	StateDataFilter StateDataFilter `json:"stateDataFilter,omitempty"`
	// Unique Name of a workflow state which is responsible for compensation of this state
	CompensatedBy string `json:"compensatedBy,omitempty"`
	// If true, this state is used to compensate another state. Default is false
	UsedForCompensation bool `json:"usedForCompensation,omitempty"`
	// State end definition
	End      End      `json:"end,omitempty"`
	Metadata Metadata `json:"metadata,omitempty"`
}

func (s *BaseState) GetOnErrors() []Error { return s.OnErrors }

func (s *BaseState) GetCompensatedBy() string { return s.CompensatedBy }

func (s *BaseState) GetTransition() Transition { return s.Transition }

func (s *BaseState) GetUsedForCompensation() bool { return s.UsedForCompensation }

func (s *BaseState) GetEnd() End { return s.End }

func (s *BaseState) GetID() string { return s.ID }

func (s *BaseState) GetName() string { return s.Name }

func (s *BaseState) GetType() StateType { return s.Type }

func (s *BaseState) GetStateDataFilter() StateDataFilter { return s.StateDataFilter }

func (s *BaseState) GetMetadata() Metadata { return s.Metadata }

// Causes the workflow execution to delay for a specified duration
type DelayState struct {
	BaseState
	// Amount of time (ISO 8601 format) to delay
	TimeDelay string `json:"timeDelay"`
}

// This state is used to wait for events from event sources, then consumes them and invoke one or more actions to run in sequence or parallel
type EventState struct {
	BaseState
	// If true consuming one of the defined events causes its associated actions to be performed. If false all of the defined events must be consumed in order for actions to be performed
	Exclusive bool `json:"exclusive,omitempty"`
	// Define the events to be consumed and optional actions to be performed
	OnEvents []OnEvents `json:"onEvents"`
	// Time period to wait for incoming events (ISO 8601 format)
	Timeout string `json:"timeout,omitempty"`
}

// Defines actions be performed. Does not wait for incoming events
type OperationState struct {
	BaseState
	// Specifies whether actions are performed in sequence or in parallel
	ActionMode ActionMode `json:"actionMode,omitempty"`
	// Actions to be performed
	Actions []Action `json:"actions"`
}

// Consists of a number of states that are executed in parallel
type ParallelState struct {
	BaseState
	// Branch Definitions
	Branches []Branch `json:"branches"`
	// Option types on how to complete branch execution.
	CompletionType CompletionType `json:"completionType,omitempty"`
	// Used when completionType is set to 'n_of_m' to specify the 'N' value
	N intstr.IntOrString `json:"n,omitempty"`
}

// Defines a sub-workflow to be executed
type SubflowState struct {
	BaseState
	// Workflow execution must wait for sub-workflow to finish before continuing
	WaitForCompletion bool `json:"waitForCompletion,omitempty"`
	// Sub-workflow unique id
	WorkflowID string `json:"workflowId"`
	// SubFlow state repeat exec definition
	Repeat Repeat `json:"repeat,omitempty"`
}

type InjectState struct {
	BaseState
	// JSON object which can be set as states data input and can be manipulated via filters
	Data map[string]interface{} `json:"data"`
}

type ForEachState struct {
	BaseState
	// Workflow expression selecting an array element of the states data
	InputCollection string `json:"inputCollection"`
	// Workflow expression specifying an array element of the states data to add the results of each iteration
	OutputCollection string `json:"outputCollection,omitempty"`
	// Name of the iteration parameter that can be referenced in actions/workflow. For each parallel iteration, this param should contain an unique element of the inputCollection array
	IterationParam string `json:"iterationParam"`
	// Specifies how upper bound on how many iterations may run in parallel
	Max intstr.IntOrString `json:"max,omitempty"`
	// Actions to be executed for each of the elements of inputCollection
	Actions []Action `json:"actions,omitempty"`
	// Unique Id of a workflow to be executed for each of the elements of inputCollection
	WorkflowID string `json:"workflowId,omitempty"`
}

type CallbackState struct {
	BaseState
	// Defines the action to be executed
	Action Action `json:"action"`
	// References an unique callback event name in the defined workflow events
	EventRef string `json:"eventRef"`
	// Time period to wait for incoming events (ISO 8601 format)
	Timeout string `json:"timeout"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

type BaseSwitchState struct {
	BaseState
	// Default transition of the workflow if there is no matching data conditions. Can include a transition or end definition
	Default DefaultDef `json:"default,omitempty"`
}

// Permits transitions to other states based on events
type EventBasedSwitchState struct {
	BaseSwitchState
	// Defines conditions evaluated against events
	EventConditions []EventCondition `json:"eventConditions"`
}

type EventCondition interface {
	GetName() string
	GetEventRef() string
	GetEventDataFilter() EventDataFilter
	GetMetadata() Metadata
}

type BaseEventCondition struct {
	// Event condition name
	Name string `json:"name,omitempty"`
	// References an unique event name in the defined workflow events
	EventRef string `json:"eventRef"`
	// Event data filter definition
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
	Metadata        Metadata        `json:"metadata,omitempty"`
}

func (e *BaseEventCondition) GetEventRef() string { return e.EventRef }

func (e *BaseEventCondition) GetEventDataFilter() EventDataFilter { return e.EventDataFilter }

func (e *BaseEventCondition) GetMetadata() Metadata { return e.Metadata }

func (e *BaseEventCondition) GetName() string { return e.Name }

// Switch state data event condition
type TransitionEventCondition struct {
	BaseEventCondition
	// Next transition of the workflow if there is valid matches
	Transition Transition `json:"transition"`
}

// Switch state data event condition
type EndEventCondition struct {
	BaseEventCondition
	// Explicit transition to end
	End End `json:"end"`
}

// Permits transitions to other states based on data conditions
type DataBasedSwitchState struct {
	BaseSwitchState
	DataConditions []DataCondition `json:"dataConditions"`
}

type DataCondition interface {
	GetName() string
	GetCondition() string
	GetMetadata() Metadata
}

type BaseDataCondition struct {
	// Data condition name
	Name string `json:"name,omitempty"`
	// Workflow expression evaluated against state data. Must evaluate to true or false
	Condition string   `json:"condition"`
	Metadata  Metadata `json:"metadata,omitempty"`
}

func (b *BaseDataCondition) GetName() string { return b.Name }

func (b *BaseDataCondition) GetCondition() string { return b.Condition }

func (b *BaseDataCondition) GetMetadata() Metadata { return b.Metadata }

type TransitionDataCondition struct {
	BaseDataCondition
	// Workflow transition if condition is evaluated to true
	Transition Transition `json:"transition"`
}

type EndDataCondition struct {
	BaseDataCondition
	// Workflow end definition
	End End `json:"end"`
}
