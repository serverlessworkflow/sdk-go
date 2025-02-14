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

// Workflow represents the root structure of a workflow.
type Workflow struct {
	Document Document            `json:"document" yaml:"document" validate:"required"`
	Input    *Input              `json:"input,omitempty" yaml:"input,omitempty"`
	Use      *Use                `json:"use,omitempty" yaml:"use"`
	Do       *TaskList           `json:"do" yaml:"do" validate:"required,dive"`
	Timeout  *TimeoutOrReference `json:"timeout,omitempty" yaml:"timeout,omitempty"`
	Output   *Output             `json:"output,omitempty" yaml:"output,omitempty"`
	Schedule *Schedule           `json:"schedule,omitempty" yaml:"schedule,omitempty"`
}

func (w *Workflow) MarshalYAML() (interface{}, error) {
	// Create a map to hold fields
	data := map[string]interface{}{
		"document": w.Document,
	}

	// Conditionally add fields
	if w.Input != nil {
		data["input"] = w.Input
	}
	if w.Use != nil {
		data["use"] = w.Use
	}
	data["do"] = w.Do
	if w.Timeout != nil {
		data["timeout"] = w.Timeout
	}
	if w.Output != nil {
		data["output"] = w.Output
	}
	if w.Schedule != nil {
		data["schedule"] = w.Schedule
	}

	return data, nil
}

// Document holds metadata for the workflow.
type Document struct {
	DSL       string                 `json:"dsl" yaml:"dsl" validate:"required,semver_pattern"`
	Namespace string                 `json:"namespace" yaml:"namespace" validate:"required,hostname_rfc1123"`
	Name      string                 `json:"name" yaml:"name" validate:"required,hostname_rfc1123"`
	Version   string                 `json:"version" yaml:"version" validate:"required,semver_pattern"`
	Title     string                 `json:"title,omitempty" yaml:"title,omitempty"`
	Summary   string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Tags      map[string]string      `json:"tags,omitempty" yaml:"tags,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty" yaml:"metadata,omitempty"`
}

// Input Configures the workflow's input.
type Input struct {
	Schema *Schema              `json:"schema,omitempty" validate:"omitempty"`
	From   *ObjectOrRuntimeExpr `json:"from,omitempty" validate:"omitempty"`
}

// Output Configures the output of a workflow or task.
type Output struct {
	Schema *Schema              `json:"schema,omitempty" validate:"omitempty"`
	As     *ObjectOrRuntimeExpr `json:"as,omitempty" validate:"omitempty"`
}

// Export Set the content of the context.
type Export struct {
	Schema *Schema              `json:"schema,omitempty" validate:"omitempty"`
	As     *ObjectOrRuntimeExpr `json:"as,omitempty" validate:"omitempty"`
}

// Schedule the workflow.
type Schedule struct {
	Every *Duration                 `json:"every,omitempty" validate:"omitempty"`
	Cron  string                    `json:"cron,omitempty" validate:"omitempty"`
	After *Duration                 `json:"after,omitempty" validate:"omitempty"`
	On    *EventConsumptionStrategy `json:"on,omitempty" validate:"omitempty"`
}

const DefaultSchema = "json"

// Schema represents the definition of a schema.
type Schema struct {
	Format   string            `json:"format,omitempty"`
	Document interface{}       `json:"document,omitempty" validate:"omitempty"`
	Resource *ExternalResource `json:"resource,omitempty" validate:"omitempty"`
}

func (s *Schema) ApplyDefaults() {
	if len(s.Format) == 0 {
		s.Format = DefaultSchema
	}
}

// UnmarshalJSON for Schema enforces "oneOf" behavior.
func (s *Schema) UnmarshalJSON(data []byte) error {
	s.ApplyDefaults()

	// Parse into a temporary map for flexibility
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check for "document"
	if doc, ok := raw["document"]; ok {
		// Determine if "document" is a string or an object
		switch doc.(type) {
		case string:
			s.Document = doc
		case map[string]interface{}:
			s.Document = doc
		default:
			return errors.New("invalid Schema: 'document' must be a string or an object")
		}
	}

	// Check for "resource"
	if res, ok := raw["resource"]; ok {
		var resource ExternalResource
		resBytes, err := json.Marshal(res)
		if err != nil {
			return fmt.Errorf("invalid Schema: failed to parse 'resource': %w", err)
		}
		if err := json.Unmarshal(resBytes, &resource); err != nil {
			return fmt.Errorf("invalid Schema: failed to parse 'resource': %w", err)
		}
		s.Resource = &resource
	}

	// Validate "oneOf" logic
	if (s.Document != nil && s.Resource != nil) || (s.Document == nil && s.Resource == nil) {
		return errors.New("invalid Schema: must specify either 'document' or 'resource', but not both")
	}

	return nil
}

// MarshalJSON for Schema marshals the correct field.
func (s *Schema) MarshalJSON() ([]byte, error) {
	s.ApplyDefaults()

	if s.Document != nil {
		return json.Marshal(map[string]interface{}{
			"format":   s.Format,
			"document": s.Document,
		})
	}
	if s.Resource != nil {
		return json.Marshal(map[string]interface{}{
			"format":   s.Format,
			"resource": s.Resource,
		})
	}

	return nil, errors.New("invalid Schema: no valid field to marshal")
}

type ExternalResource struct {
	Name     string    `json:"name,omitempty"`
	Endpoint *Endpoint `json:"endpoint" validate:"required"`
}

type Use struct {
	Authentications map[string]*AuthenticationPolicy `json:"authentications,omitempty" validate:"omitempty,dive"`
	Errors          map[string]*Error                `json:"errors,omitempty" validate:"omitempty,dive"`
	Extensions      ExtensionList                    `json:"extensions,omitempty" validate:"omitempty,dive"`
	Functions       NamedTaskMap                     `json:"functions,omitempty" validate:"omitempty,dive"`
	Retries         map[string]*RetryPolicy          `json:"retries,omitempty" validate:"omitempty,dive"`
	Secrets         []string                         `json:"secrets,omitempty"`
	Timeouts        map[string]*Timeout              `json:"timeouts,omitempty" validate:"omitempty,dive"`
	Catalogs        map[string]*Catalog              `json:"catalogs,omitempty" validate:"omitempty,dive"`
}

type Catalog struct {
	Endpoint *Endpoint `json:"endpoint" validate:"required"`
}

// FlowDirective represents a directive that can be an enumerated or free-form string.
type FlowDirective struct {
	Value string `json:"-" validate:"required"` // Ensure the value is non-empty.
}

type FlowDirectiveType string

const (
	FlowDirectiveContinue FlowDirectiveType = "continue"
	FlowDirectiveExit     FlowDirectiveType = "exit"
	FlowDirectiveEnd      FlowDirectiveType = "end"
)

// Enumerated values for FlowDirective.
var validFlowDirectives = map[string]struct{}{
	"continue": {},
	"exit":     {},
	"end":      {},
}

// IsEnum checks if the FlowDirective matches one of the enumerated values.
func (f *FlowDirective) IsEnum() bool {
	_, exists := validFlowDirectives[f.Value]
	return exists
}

// IsTermination checks if the FlowDirective matches FlowDirectiveExit or FlowDirectiveEnd.
func (f *FlowDirective) IsTermination() bool {
	return f.Value == string(FlowDirectiveExit) || f.Value == string(FlowDirectiveEnd)
}

func (f *FlowDirective) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}
	f.Value = value
	return nil
}

func (f *FlowDirective) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}
