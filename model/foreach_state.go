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

	"k8s.io/apimachinery/pkg/util/intstr"
)

// ForEachModeType Specifies how iterations are to be performed (sequentially or in parallel)
type ForEachModeType string

func (f ForEachModeType) KindValues() []string {
	return []string{
		string(ForEachModeTypeSequential),
		string(ForEachModeTypeParallel),
	}
}

func (f ForEachModeType) String() string {
	return string(f)
}

const (
	// ForEachModeTypeSequential specifies iterations should be done sequentially.
	ForEachModeTypeSequential ForEachModeType = "sequential"
	// ForEachModeTypeParallel specifies iterations should be done parallel.
	ForEachModeTypeParallel ForEachModeType = "parallel"
)

// ForEachState used to execute actions for each element of a data set.
type ForEachState struct {
	// Workflow expression selecting an array element of the states' data.
	// +kubebuilder:validation:Required
	InputCollection string `json:"inputCollection" validate:"required"`
	// Workflow expression specifying an array element of the states data to add the results of each iteration.
	// +optional
	OutputCollection string `json:"outputCollection,omitempty"`
	// Name of the iteration parameter that can be referenced in actions/workflow. For each parallel iteration,
	// this param should contain a unique element of the inputCollection array.
	// +optional
	IterationParam string `json:"iterationParam,omitempty"`
	// Specifies how many iterations may run in parallel at the same time. Used if mode property is set to
	// parallel (default). If not specified, its value should be the size of the inputCollection.
	// +optional
	BatchSize *intstr.IntOrString `json:"batchSize,omitempty"`
	// Actions to be executed for each of the elements of inputCollection.
	// +kubebuilder:validation:MinItems=1
	Actions []Action `json:"actions,omitempty" validate:"required,min=1,dive"`
	// State specific timeout.
	// +optional
	Timeouts *ForEachStateTimeout `json:"timeouts,omitempty"`
	// Specifies how iterations are to be performed (sequential or in parallel), defaults to parallel.
	// +kubebuilder:validation:Enum=sequential;parallel
	// +kubebuilder:default=parallel
	Mode ForEachModeType `json:"mode,omitempty" validate:"required,oneofkind"`
}

func (f *ForEachState) MarshalJSON() ([]byte, error) {
	type Alias ForEachState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *ForEachStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(f),
		Timeouts: f.Timeouts,
	})
	return custom, err
}

type forEachStateUnmarshal ForEachState

// UnmarshalJSON implements json.Unmarshaler
func (f *ForEachState) UnmarshalJSON(data []byte) error {
	f.ApplyDefault()
	return unmarshalObject("forEachState", data, (*forEachStateUnmarshal)(f))
}

// ApplyDefault set the default values
func (f *ForEachState) ApplyDefault() {
	f.Mode = ForEachModeTypeParallel
}

// ForEachStateTimeout defines timeout settings for foreach state
type ForEachStateTimeout struct {
	// Default workflow state execution timeout (ISO 8601 duration format)
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// Default single actions definition execution timeout (ISO 8601 duration format)
	// +optional
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
