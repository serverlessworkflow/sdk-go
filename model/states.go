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
	"fmt"
	"strings"
)

// StateType ...
type StateType string

const (
	// StateTypeDelay ...
	StateTypeDelay StateType = "delay"
	// StateTypeEvent ...
	StateTypeEvent StateType = "event"
	// StateTypeOperation ...
	StateTypeOperation StateType = "operation"
	// StateTypeParallel ...
	StateTypeParallel StateType = "parallel"
	// StateTypeSwitch ...
	StateTypeSwitch StateType = "switch"
	// StateTypeForEach ...
	StateTypeForEach StateType = "foreach"
	// StateTypeInject ...
	StateTypeInject StateType = "inject"
	// StateTypeCallback ...
	StateTypeCallback StateType = "callback"
	// StateTypeSleep ...
	StateTypeSleep StateType = "sleep"
)

// BaseState ...
type BaseState struct {
	// Unique State id.
	// +optional
	ID string `json:"id,omitempty"`
	// State name.
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// stateType can be any of delay, callback, event, foreach, inject, operation, parallel, sleep, switch
	// +kubebuilder:validation:Enum:=delay;callback;event;foreach;inject;operation;parallel;sleep;switch
	// +kubebuilder:validation:Required
	Type StateType `json:"type" validate:"required"`
	// States error handling and retries definitions.
	// +optional
	OnErrors []OnError `json:"onErrors,omitempty"  validate:"omitempty,dive"`
	// Next transition of the workflow after the time delay.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Transition *Transition `json:"transition,omitempty"`
	// State data filter.
	// +optional
	StateDataFilter *StateDataFilter `json:"stateDataFilter,omitempty"`
	// Unique Name of a workflow state which is responsible for compensation of this state.
	// +optional
	CompensatedBy string `json:"compensatedBy,omitempty"`
	// If true, this state is used to compensate another state. Default is false.
	// +optional
	UsedForCompensation bool `json:"usedForCompensation,omitempty"`
	// State end definition.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	End *End `json:"end,omitempty"`
	// Metadata information.
	// +optional
	Metadata *Metadata `json:"metadata,omitempty"`
}

func (b *BaseState) MarshalJSON() ([]byte, error) {
	type Alias BaseState
	if b == nil {
		return []byte("null"), nil
	}
	cus, err := json.Marshal(struct {
		*Alias
	}{
		Alias: (*Alias)(b),
	})
	return cus, err
}

type State struct {
	BaseState `json:",inline"`
	// delayState Causes the workflow execution to delay for a specified duration.
	// +optional
	*DelayState `json:"delayState,omitempty"`
	// event states await one or more events and perform actions when they are received. If defined as the
	// workflow starting state, the event state definition controls when the workflow instances should be created.
	// +optional
	*EventState `json:"eventState,omitempty"`
	// operationState defines a set of actions to be performed in sequence or in parallel.
	// +optional
	*OperationState `json:"operationState,omitempty"`
	// parallelState Consists of a number of states that are executed in parallel.
	// +optional
	*ParallelState `json:"parallelState,omitempty"`
	// switchState is workflow's gateways: direct transitions onf a workflow based on certain conditions.
	// +optional
	*SwitchState `json:"switchState,omitempty"`
	// forEachState used to execute actions for each element of a data set.
	// +optional
	*ForEachState `json:"forEachState,omitempty"`
	// injectState used to inject static data into state data input.
	// +optional
	*InjectState `json:"injectState,omitempty"`
	// callbackState executes a function and waits for callback event that indicates completion of the task.
	// +optional
	*CallbackState `json:"callbackState,omitempty"`
	// sleepState suspends workflow execution for a given time duration.
	// +optional
	*SleepState `json:"sleepState,omitempty"`
}

func (s *State) MarshalJSON() ([]byte, error) {
	if s == nil {
		return []byte("null"), nil
	}
	r := []byte("")
	var errs error

	if s.DelayState != nil {
		r, errs = s.DelayState.MarshalJSON()
	}

	if s.EventState != nil {
		r, errs = s.EventState.MarshalJSON()
	}

	if s.OperationState != nil {
		r, errs = s.OperationState.MarshalJSON()
	}

	if s.ParallelState != nil {
		r, errs = s.ParallelState.MarshalJSON()
	}

	if s.SwitchState != nil {
		r, errs = s.SwitchState.MarshalJSON()
	}

	if s.ForEachState != nil {
		r, errs = s.ForEachState.MarshalJSON()
	}

	if s.InjectState != nil {
		r, errs = s.InjectState.MarshalJSON()
	}

	if s.CallbackState != nil {
		r, errs = s.CallbackState.MarshalJSON()
	}

	if s.SleepState != nil {
		r, errs = s.SleepState.MarshalJSON()
	}

	b, err := s.BaseState.MarshalJSON()
	if err != nil {
		return nil, err
	}

	//remove }{ as BaseState and the State Type needs to be merged together
	partialResult := append(b, r...)
	result := strings.Replace(string(partialResult), "}{", ",", 1)
	return []byte(result), errs
}

func (s *State) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &s.BaseState); err != nil {
		return err
	}

	mapState := map[string]interface{}{}
	if err := json.Unmarshal(data, &mapState); err != nil {
		return err
	}

	switch mapState["type"] {
	case string(StateTypeDelay):
		state := &DelayState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.DelayState = state

	case string(StateTypeEvent):
		state := &EventState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.EventState = state

	case string(StateTypeOperation):
		state := &OperationState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.OperationState = state

	case string(StateTypeParallel):
		state := &ParallelState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.ParallelState = state

	case string(StateTypeSwitch):
		state := &SwitchState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.SwitchState = state

	case string(StateTypeForEach):
		state := &ForEachState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.ForEachState = state

	case string(StateTypeInject):
		state := &InjectState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.InjectState = state

	case string(StateTypeCallback):
		state := &CallbackState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.CallbackState = state

	case string(StateTypeSleep):
		state := &SleepState{}
		if err := json.Unmarshal(data, state); err != nil {
			return err
		}
		s.SleepState = state
	case nil:
		return fmt.Errorf("state parameter 'type' not defined")
	default:
		return fmt.Errorf("state type %v not supported", mapState["type"])
	}
	return nil
}
