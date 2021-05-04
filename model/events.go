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

const (
	EventKindConsumed EventKind = "consumed"
	EventKindProduced EventKind = "produced"
)

type EventKind string

type Event struct {
	Common
	// Unique event name
	Name string `json:"name"`
	// CloudEvent source
	Source string `json:"source,omitempty"`
	// CloudEvent type
	Type string `json:"type"`
	// Defines the CloudEvent as either 'consumed' or 'produced' by the workflow. Default is 'consumed'
	Kind EventKind `json:"kind,omitempty"`
	// CloudEvent correlation definitions
	Correlation []Correlation `json:"correlation,omitempty"`
}

type Correlation struct {
	// CloudEvent Extension Context Attribute name
	ContextAttributeName string `json:"contextAttributeName"`
	// CloudEvent Extension Context Attribute value
	ContextAttributeValue string `json:"contextAttributeValue,omitempty"`
}

type EventRef struct {
	// Reference to the unique name of a 'produced' event definition
	TriggerEventRef string `json:"triggerEventRef"`
	// Reference to the unique name of a 'consumed' event definition
	ResultEventRef string `json:"resultEventRef"`
	// TODO: create StringOrMap structure
	// If string type, an expression which selects parts of the states data output to become the data (payload) of the event referenced by 'triggerEventRef'. If object type, a custom object to become the data (payload) of the event referenced by 'triggerEventRef'.
	Data interface{} `json:"data,omitempty"`
	// Add additional extension context attributes to the produced event
	ContextAttributes map[string]interface{} `json:"contextAttributes,omitempty"`
}
