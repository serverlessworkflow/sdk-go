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
)

// OperationState defines a set of actions to be performed in sequence or in parallel.
type OperationState struct {
	// Specifies whether actions are performed in sequence or in parallel, defaults to sequential
	ActionMode ActionMode `json:"actionMode,omitempty" validate:"required,oneof=sequential parallel"`
	// Actions to be performed
	// +listType=atomic
	// +optional
	Actions []Action `json:"actions" validate:"required,min=1,dive"`
	// State specific timeouts
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

type operationStateForUnmarshal OperationState

// UnmarshalJSON unmarshal OperationState object from json bytes
func (o *OperationState) UnmarshalJSON(data []byte) error {

	v := operationStateForUnmarshal{
		ActionMode: ActionModeSequential,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*o = OperationState(v)
	return nil
}

// OperationStateTimeout defines the specific timeout settings for operation state
type OperationStateTimeout struct {
	StateExecTimeout  *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string            `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
