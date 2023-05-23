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

import "fmt"

// EventDataFilter used to filter consumed event payloads.
type EventDataFilter struct {
	// If set to false, event payload is not added/merged to state data. In this case 'data' and 'toStateData'
	// should be ignored. Default is true.
	// +optional
	UseData bool `json:"useData,omitempty"`
	// Workflow expression that filters of the event data (payload).
	// +optional
	Data string `json:"data,omitempty"`
	// Workflow expression that selects a state data element to which the action results should be added/merged into.
	// If not specified denotes the top-level state data element
	// +optional
	ToStateData string `json:"toStateData,omitempty"`
}

func (f EventDataFilter) String() string {
	return fmt.Sprintf("{ UseData:%t, Data:%s, ToStateData:%s }", f.UseData, f.Data, f.ToStateData)
}

type eventDataFilterUnmarshal EventDataFilter

// UnmarshalJSON implements json.Unmarshaler
func (f *EventDataFilter) UnmarshalJSON(data []byte) error {
	f.ApplyDefault()
	return unmarshalObject("eventDataFilter", data, (*eventDataFilterUnmarshal)(f))
}

// ApplyDefault set the default values for Event Data Filter
func (f *EventDataFilter) ApplyDefault() {
	f.UseData = true
}
