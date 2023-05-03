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

// EventKind defines this event as either `consumed` or `produced`
type EventKind string

const (
	// EventKindConsumed means the event continuation of workflow instance execution
	EventKindConsumed EventKind = "consumed"

	// EventKindProduced means the event was created during workflow instance execution
	EventKindProduced EventKind = "produced"
)

// Event used to define events and their correlations
type Event struct {
	Common `json:",inline"`
	// Unique event name.
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// CloudEvent source.
	// +optional
	Source string `json:"source,omitempty"`
	// CloudEvent type.
	// +kubebuilder:validation:Required
	Type string `json:"type" validate:"required"`
	// Defines the CloudEvent as either 'consumed' or 'produced' by the workflow. Defaults to `consumed`.
	// +kubebuilder:validation:Enum=consumed;produced
	// +kubebuilder:default=consumed
	Kind EventKind `json:"kind,omitempty"`
	// If `true`, only the Event payload is accessible to consuming Workflow states. If `false`, both event payload
	// and context attributes should be accessible. Defaults to true.
	// +optional
	DataOnly bool `json:"dataOnly,omitempty"`
	// Define event correlation rules for this event. Only used for consumed events.
	// +optional
	Correlation []Correlation `json:"correlation,omitempty" validate:"omitempty,dive"`
}

type eventUnmarshal Event

// UnmarshalJSON unmarshal Event object from json bytes
func (e *Event) UnmarshalJSON(data []byte) error {
	e.ApplyDefault()
	return unmarshalObject("event", data, (*eventUnmarshal)(e))
}

// ApplyDefault set the default values for Event
func (e *Event) ApplyDefault() {
	e.DataOnly = true
	e.Kind = EventKindConsumed
}

// Correlation define event correlation rules for an event. Only used for `consumed` events
type Correlation struct {
	// CloudEvent Extension Context Attribute name
	// +kubebuilder:validation:Required
	ContextAttributeName string `json:"contextAttributeName" validate:"required"`
	// CloudEvent Extension Context Attribute value
	// +optional
	ContextAttributeValue string `json:"contextAttributeValue,omitempty"`
}

// EventRef defining invocation of a function via event
type EventRef struct {
	// Reference to the unique name of a 'produced' event definition,
	// +kubebuilder:validation:Required
	TriggerEventRef string `json:"triggerEventRef" validate:"required"`
	// Reference to the unique name of a 'consumed' event definition
	// +kubebuilder:validation:Required
	ResultEventRef string `json:"resultEventRef" validate:"required"`
	// Maximum amount of time (ISO 8601 format) to wait for the result event. If not defined it be set to the
	// actionExecutionTimeout
	// +optional
	ResultEventTimeout string `json:"resultEventTimeout,omitempty" validate:"omitempty,iso8601duration"`
	// If string type, an expression which selects parts of the states data output to become the data (payload)
	// of the event referenced by triggerEventRef. If object type, a custom object to become the data (payload)
	// of the event referenced by triggerEventRef.
	// +optional
	Data *Object `json:"data,omitempty"`
	// Add additional extension context attributes to the produced event.
	// +optional
	ContextAttributes map[string]Object `json:"contextAttributes,omitempty"`
	// Specifies if the function should be invoked sync or async. Default is sync.
	// +kubebuilder:validation:Enum=async;sync
	// +kubebuilder:default=sync
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneofkind"`
}

type eventRefUnmarshal EventRef

// UnmarshalJSON implements json.Unmarshaler
func (e *EventRef) UnmarshalJSON(data []byte) error {
	e.ApplyDefault()
	return unmarshalObject("eventRef", data, (*eventRefUnmarshal)(e))
}

// ApplyDefault set the default values for Event Ref
func (e *EventRef) ApplyDefault() {
	e.Invoke = InvokeKindSync
}
