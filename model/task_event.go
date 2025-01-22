// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

// EmitTask represents the configuration for emitting events.
type EmitTask struct {
	TaskBase `json:",inline"`      // Inline TaskBase fields
	Emit     EmitTaskConfiguration `json:"emit" validate:"required"`
}

func (e *EmitTask) MarshalJSON() ([]byte, error) {
	type Alias EmitTask // Prevent recursion
	return json.Marshal((*Alias)(e))
}

// ListenTask represents a task configuration to listen to events.
type ListenTask struct {
	TaskBase `json:",inline"`        // Inline TaskBase fields
	Listen   ListenTaskConfiguration `json:"listen" validate:"required"`
}

type ListenTaskConfiguration struct {
	To *EventConsumptionStrategy `json:"to" validate:"required"`
}

// MarshalJSON for ListenTask to ensure proper serialization.
func (lt *ListenTask) MarshalJSON() ([]byte, error) {
	type Alias ListenTask
	return json.Marshal((*Alias)(lt))
}

// UnmarshalJSON for ListenTask to ensure proper deserialization.
func (lt *ListenTask) UnmarshalJSON(data []byte) error {
	type Alias ListenTask
	alias := (*Alias)(lt)
	return json.Unmarshal(data, alias)
}

type EmitTaskConfiguration struct {
	Event EmitEventDefinition `json:"event" validate:"required"`
}

type EmitEventDefinition struct {
	With *EventProperties `json:"with" validate:"required"`
}

type EventProperties struct {
	ID              string                    `json:"id,omitempty"`
	Source          *URITemplateOrRuntimeExpr `json:"source,omitempty" validate:"omitempty"` // URI template or runtime expression
	Type            string                    `json:"type,omitempty"`
	Time            *StringOrRuntimeExpr      `json:"time,omitempty" validate:"omitempty,string_or_runtime_expr"` // ISO 8601 date-time string or runtime expression
	Subject         string                    `json:"subject,omitempty"`
	DataContentType string                    `json:"datacontenttype,omitempty"`
	DataSchema      *URITemplateOrRuntimeExpr `json:"dataschema,omitempty" validate:"omitempty"` // URI template or runtime expression
	Additional      map[string]interface{}    `json:"-"`
}

// UnmarshalJSON implements custom unmarshaling for EventProperties.
func (e *EventProperties) UnmarshalJSON(data []byte) error {
	type Alias EventProperties // Prevent recursion
	alias := &struct {
		Additional map[string]interface{} `json:"-"` // Inline the additional properties
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	// Decode the entire JSON into a map to capture additional properties
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal EventProperties: %w", err)
	}

	// Unmarshal known fields into the alias
	if err := json.Unmarshal(data, alias); err != nil {
		return fmt.Errorf("failed to unmarshal EventProperties fields: %w", err)
	}

	// Validate fields requiring custom unmarshaling
	if e.Source != nil && e.Source.Value == nil {
		return fmt.Errorf("invalid Source: must be a valid URI template or runtime expression")
	}

	if e.DataSchema != nil && e.DataSchema.Value == nil {
		return fmt.Errorf("invalid DataSchema: must be a valid URI template or runtime expression")
	}

	// Extract additional properties by removing known keys
	for key := range raw {
		switch key {
		case "id", "source", "type", "time", "subject", "datacontenttype", "dataschema":
			delete(raw, key)
		}
	}

	e.Additional = raw
	return nil
}

// MarshalJSON implements custom marshaling for EventProperties.
func (e *EventProperties) MarshalJSON() ([]byte, error) {
	// Create a map for known fields
	known := make(map[string]interface{})

	if e.ID != "" {
		known["id"] = e.ID
	}
	if e.Source != nil {
		known["source"] = e.Source
	}
	if e.Type != "" {
		known["type"] = e.Type
	}
	if e.Time != nil {
		known["time"] = e.Time
	}
	if e.Subject != "" {
		known["subject"] = e.Subject
	}
	if e.DataContentType != "" {
		known["datacontenttype"] = e.DataContentType
	}
	if e.DataSchema != nil {
		known["dataschema"] = e.DataSchema
	}

	// Merge additional properties
	for key, value := range e.Additional {
		known[key] = value
	}

	return json.Marshal(known)
}

// EventFilter defines a mechanism to filter events based on predefined criteria.
type EventFilter struct {
	With      *EventProperties       `json:"with" validate:"required"`
	Correlate map[string]Correlation `json:"correlate,omitempty" validate:"omitempty,dive"` // Keyed correlation filters
}

// Correlation defines the mapping of event attributes for correlation.
type Correlation struct {
	From   string `json:"from" validate:"required"` // Runtime expression to extract the correlation value
	Expect string `json:"expect,omitempty"`         // Expected value or runtime expression
}

// EventConsumptionStrategy defines the consumption strategy for events.
type EventConsumptionStrategy struct {
	All   []*EventFilter         `json:"all,omitempty" validate:"omitempty,dive"`
	Any   []*EventFilter         `json:"any,omitempty" validate:"omitempty,dive"`
	One   *EventFilter           `json:"one,omitempty" validate:"omitempty"`
	Until *EventConsumptionUntil `json:"until,omitempty" validate:"omitempty"`
}

// EventConsumptionUntil handles the complex conditions of the "until" field.
type EventConsumptionUntil struct {
	Condition  *RuntimeExpression        `json:"-" validate:"omitempty"`
	Strategy   *EventConsumptionStrategy `json:"-" validate:"omitempty"`
	IsDisabled bool                      `json:"-"` // True when "until: false"
}

// UnmarshalJSON for EventConsumptionUntil to handle the "oneOf" behavior.
func (ecu *EventConsumptionUntil) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal EventConsumptionUntil: %w", err)
	}

	switch v := raw.(type) {
	case bool:
		if !v {
			ecu.IsDisabled = true
		} else {
			return fmt.Errorf("invalid value for 'until': true is not supported")
		}
	case string:
		ecu.Condition = &RuntimeExpression{Value: v}
	case map[string]interface{}:
		strategyData, err := json.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal 'until' strategy: %w", err)
		}
		var strategy EventConsumptionStrategy
		if err := json.Unmarshal(strategyData, &strategy); err != nil {
			return fmt.Errorf("failed to unmarshal 'until' strategy: %w", err)
		}
		ecu.Strategy = &strategy
	default:
		return fmt.Errorf("invalid type for 'until'")
	}

	return nil
}

// MarshalJSON for EventConsumptionUntil to handle proper serialization.
func (ecu *EventConsumptionUntil) MarshalJSON() ([]byte, error) {
	if ecu.IsDisabled {
		return json.Marshal(false)
	}
	if ecu.Condition != nil {
		// Serialize the condition directly
		return json.Marshal(ecu.Condition.Value)
	}
	if ecu.Strategy != nil {
		// Serialize the nested strategy
		return json.Marshal(ecu.Strategy)
	}
	// Return null if nothing is set
	return json.Marshal(nil)
}

// UnmarshalJSON for EventConsumptionStrategy to enforce "oneOf" behavior and handle edge cases.
func (ecs *EventConsumptionStrategy) UnmarshalJSON(data []byte) error {
	temp := struct {
		All   []*EventFilter         `json:"all"`
		Any   []*EventFilter         `json:"any"`
		One   *EventFilter           `json:"one"`
		Until *EventConsumptionUntil `json:"until"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Count non-nil fields (ignoring empty lists for `all` and `any`)
	count := 0
	if len(temp.All) > 0 {
		count++
		ecs.All = temp.All
	}
	if len(temp.Any) > 0 || temp.Until != nil {
		count++
		ecs.Any = temp.Any
		ecs.Until = temp.Until
	}
	if temp.One != nil {
		count++
		ecs.One = temp.One
	}

	// Ensure only one primary field (all, any, one) is set
	if count > 1 {
		return errors.New("invalid EventConsumptionStrategy: only one primary strategy type (all, any, or one) must be specified")
	}

	return nil
}

// MarshalJSON for EventConsumptionStrategy to ensure proper serialization.
func (ecs *EventConsumptionStrategy) MarshalJSON() ([]byte, error) {
	temp := struct {
		All   []*EventFilter         `json:"all,omitempty"`
		Any   []*EventFilter         `json:"any,omitempty"`
		One   *EventFilter           `json:"one,omitempty"`
		Until *EventConsumptionUntil `json:"until,omitempty"`
	}{
		All:   ecs.All,
		Any:   ecs.Any,
		One:   ecs.One,
		Until: ecs.Until,
	}

	return json.Marshal(temp)
}
