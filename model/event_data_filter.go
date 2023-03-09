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

// EventDataFilter used to filter consumed event payloads.
type EventDataFilter struct {
	// UseData represent where event payload is added/merged to state data. If it's false, data & toStateData
	// should be ignored. Defaults to true.
	UseData bool `json:"useData,omitempty"`
	// Workflow expression that filters of the event data (payload)
	// +optional
	Data string `json:"data,omitempty"`
	// Workflow expression that selects a state data element to which the event payload should be added/merged into.
	// If not specified, denotes, the top-level state data element.
	ToStateData string `json:"toStateData,omitempty"`
}

type eventDataFilterForUnmarshal EventDataFilter

func (f *EventDataFilter) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	v := eventDataFilterForUnmarshal{
		UseData: true,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		// TODO: replace the error message with correct type's name
		return err
	}

	*f = EventDataFilter(v)
	return nil
}
