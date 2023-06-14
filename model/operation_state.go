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

	"github.com/serverlessworkflow/sdk-go/v2/util"
)

// OperationState defines a set of actions to be performed in sequence or in parallel.
type OperationState struct {
	// Specifies whether actions are performed in sequence or in parallel, defaults to sequential.
	// +kubebuilder:validation:Enum=sequential;parallel
	// +kubebuilder:default=sequential
	ActionMode ActionMode `json:"actionMode,omitempty" validate:"required,oneofkind"`
	// Actions to be performed
	// +kubebuilder:validation:MinItems=1
	Actions []Action `json:"actions" validate:"required,min=1,dive"`
	// State specific timeouts
	// +optional
	Timeouts *OperationStateTimeout `json:"timeouts,omitempty"`
}

func (a *OperationState) MarshalJSON() ([]byte, error) {
	type Alias OperationState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *OperationStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(a),
		Timeouts: a.Timeouts,
	})
	return custom, err
}

type operationStateUnmarshal OperationState

// UnmarshalJSON unmarshal OperationState object from json bytes
func (o *OperationState) UnmarshalJSON(data []byte) error {
	o.ApplyDefault()
	return util.UnmarshalObject("operationState", data, (*operationStateUnmarshal)(o))
}

// ApplyDefault set the default values for Operation State
func (o *OperationState) ApplyDefault() {
	o.ActionMode = ActionModeSequential
}

// OperationStateTimeout defines the specific timeout settings for operation state
type OperationStateTimeout struct {
	// Defines workflow state execution timeout.
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// Default single actions definition execution timeout (ISO 8601 duration format)
	// +optional
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
