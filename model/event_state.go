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
	"fmt"
)

// EventState used to wait for events from event sources, then consumes them and invoke one or more actions to run in sequence or parallel
type EventState struct {
	// TODO: EventState doesn't have usedForCompensation field.

	// If true consuming one of the defined events causes its associated actions to be performed.
	// If false all the defined events must be consumed in order for actions to be performed
	// Defaults to true
	Exclusive bool `json:"exclusive,omitempty"`
	// Define the events to be consumed and optional actions to be performed
	// +optional
	OnEvents []OnEvents `json:"onEvents" validate:"required,min=1,dive"`
	// State specific timeouts
	Timeouts *EventStateTimeout `json:"timeouts,omitempty"`
}

func (e *EventState) MarshalJSON() ([]byte, error) {
	type Alias EventState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *EventStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(e),
		Timeouts: e.Timeouts,
	})
	return custom, err
}

type eventStateForUnmarshal EventState

// UnmarshalJSON unmarshal EventState object from json bytes
func (e *EventState) UnmarshalJSON(data []byte) error {
	v := eventStateForUnmarshal{
		Exclusive: true,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("eventState value '%s' is not supported, it must be an object or string", string(data))
	}

	*e = EventState(v)
	return nil
}

// OnEvents define which actions are be performed for the one or more events.
type OnEvents struct {
	// References one or more unique event names in the defined workflow events
	EventRefs []string `json:"eventRefs" validate:"required,min=1"`
	// Specifies how actions are to be performed (in sequence or parallel)
	// Defaults to sequential
	ActionMode ActionMode `json:"actionMode,omitempty" validate:"required,oneof=sequential parallel"`
	// Actions to be performed if expression matches
	// +listType=atomic
	// +optional
	Actions []Action `json:"actions,omitempty" validate:"omitempty,dive"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

type onEventsForUnmarshal OnEvents

// UnmarshalJSON unmarshal OnEvents object from json bytes
func (o *OnEvents) UnmarshalJSON(data []byte) error {
	v := onEventsForUnmarshal{
		ActionMode: ActionModeSequential,
	}

	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("onEvents value '%s' is not supported, it must be an object or string", string(data))
	}

	*o = OnEvents(v)

	return nil
}

// EventStateTimeout defines timeout settings for event state
type EventStateTimeout struct {
	StateExecTimeout  *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string            `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
	EventTimeout      string            `json:"eventTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
