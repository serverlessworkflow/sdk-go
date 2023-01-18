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

// StateType ...
type StateType string

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

func getActionsModelMapping(stateType string) (State, bool) {
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
		return &SwitchState{}, true
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
	// DeepCopyState fixes undefined (type State has no field or method DeepCopyState)
	DeepCopyState() State
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
