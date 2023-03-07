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
	"reflect"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"

	validator "github.com/go-playground/validator/v10"
)

// EventKind defines this event as either `consumed` or `produced`
type EventKind string

const (
	// EventKindConsumed means the event continuation of workflow instance execution
	EventKindConsumed EventKind = "consumed"

	// EventKindProduced means the event was created during workflow instance execution
	EventKindProduced EventKind = "produced"
)

func init() {
	val.GetValidator().RegisterStructValidation(EventStructLevelValidation, Event{})
}

// EventStructLevelValidation custom validator for event kind consumed
func EventStructLevelValidation(structLevel validator.StructLevel) {
	event := structLevel.Current().Interface().(Event)
	if event.Kind == EventKindConsumed && len(event.Type) == 0 {
		structLevel.ReportError(reflect.ValueOf(event.Type), "Type", "type", "reqtypeconsumed", "")
	}
}

// Event used to define events and their correlations
type Event struct {
	Common `json:",inline"`
	// Unique event name
	Name string `json:"name" validate:"required"`
	// CloudEvent source
	Source string `json:"source,omitempty"`
	// CloudEvent type
	Type string `json:"type" validate:"required"`
	// Defines the CloudEvent as either 'consumed' or 'produced' by the workflow.
	// Defaults to `consumed`
	Kind EventKind `json:"kind,omitempty"`
	// If `true`, only the Event payload is accessible to consuming Workflow states. If `false`, both event payload and context attributes should be accessible"
	// Defaults to true
	DataOnly bool `json:"dataOnly,omitempty"`
	// CloudEvent correlation definitions
	Correlation []Correlation `json:"correlation,omitempty" validate:"omitempty,dive"`
}

type eventForUnmarshal Event

// UnmarshalJSON unmarshal Event object from json bytes
func (e *Event) UnmarshalJSON(data []byte) error {
	v := eventForUnmarshal{
		DataOnly: true,
		Kind:     EventKindConsumed,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*e = Event(v)
	return nil
}

// Correlation define event correlation rules for an event. Only used for `consumed` events
type Correlation struct {
	// CloudEvent Extension Context Attribute name
	ContextAttributeName string `json:"contextAttributeName" validate:"required"`
	// CloudEvent Extension Context Attribute value
	ContextAttributeValue string `json:"contextAttributeValue,omitempty"`
}

// EventRef defining invocation of a function via event
type EventRef struct {
	// Reference to the unique name of a 'produced' event definition
	TriggerEventRef string `json:"triggerEventRef" validate:"required"`
	// Reference to the unique name of a 'consumed' event definition
	ResultEventRef string `json:"resultEventRef" validate:"required"`
	// ResultEventTimeout defines maximum amount of time (ISO 8601 format) to wait for the result event. If not defined it be set to the actionExecutionTimeout
	ResultEventTimeout string `json:"resultEventTimeout,omitempty" validate:"omitempty,iso8601duration"`
	// If string type, an expression which selects parts of the states data output to become the data (payload) of the event referenced by 'triggerEventRef'.
	// If object type, a custom object to become the data (payload) of the event referenced by 'triggerEventRef'.
	Data *Object `json:"data,omitempty"`
	// Add additional extension context attributes to the produced event
	ContextAttributes map[string]Object `json:"contextAttributes,omitempty"`
	// Invoke specifies if the subflow should be invoked sync or async.
	// Defaults to sync.
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneof=async sync"`
}

type eventRefForUnmarshal EventRef

// UnmarshalJSON implements json.Unmarshaler
func (e *EventRef) UnmarshalJSON(data []byte) error {
	v := eventRefForUnmarshal{
		Invoke: InvokeKindSync,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return err
	}

	*e = EventRef(v)
	return nil
}
