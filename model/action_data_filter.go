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
	"bytes"
	"encoding/json"
	"fmt"
)

// ActionDataFilter used to filter action data results.
type ActionDataFilter struct {
	// Workflow expression that selects state data that the state action can use
	FromStateData string `json:"fromStateData,omitempty"`

	// UseResults represent where action data results is added/merged to state data. If it's false, results & toStateData should be ignored.
	// Defaults to true.
	UseResults bool `json:"useResults,omitempty"`

	// Workflow expression that filters the actions' data results
	Results string `json:"results,omitempty"`
	// Workflow expression that selects a state data element to which the action results should be added/merged into. If not specified, denote, the top-level state data element
	ToStateData string `json:"toStateData,omitempty"`
}

type actionDataFilterForUnmarshal ActionDataFilter

// UnmarshalJSON implements json.Unmarshaler
func (f *ActionDataFilter) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	v := actionDataFilterForUnmarshal{
		UseResults: true,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		// TODO: replace the error message with correct type's name
		return err
	}

	*f = ActionDataFilter(v)
	return nil
}
