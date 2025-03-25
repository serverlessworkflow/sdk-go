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
	"regexp"
)

var _ Object = &ObjectOrString{}
var _ Object = &ObjectOrRuntimeExpr{}
var _ Object = &RuntimeExpression{}
var _ Object = &URITemplateOrRuntimeExpr{}
var _ Object = &StringOrRuntimeExpr{}
var _ Object = &JsonPointerOrRuntimeExpression{}

type Object interface {
	String() string
	GetValue() interface{}
}

// ObjectOrString is a type that can hold either a string or an object.
type ObjectOrString struct {
	Value interface{} `validate:"object_or_string"`
}

func (o *ObjectOrString) String() string {
	return fmt.Sprintf("%v", o.Value)
}

func (o *ObjectOrString) GetValue() interface{} {
	return o.Value
}

// UnmarshalJSON unmarshals data into either a string or an object.
func (o *ObjectOrString) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		o.Value = asString
		return nil
	}

	var asObject map[string]interface{}
	if err := json.Unmarshal(data, &asObject); err == nil {
		o.Value = asObject
		return nil
	}

	return errors.New("ObjectOrString must be a string or an object")
}

// MarshalJSON marshals ObjectOrString into JSON.
func (o *ObjectOrString) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.Value)
}

// ObjectOrRuntimeExpr is a type that can hold either a RuntimeExpression or an object.
type ObjectOrRuntimeExpr struct {
	Value interface{} `json:"-" validate:"object_or_runtime_expr"` // Custom validation tag.
}

func NewObjectOrRuntimeExpr(value interface{}) *ObjectOrRuntimeExpr {
	return &ObjectOrRuntimeExpr{
		Value: value,
	}
}

func (o *ObjectOrRuntimeExpr) String() string {
	return fmt.Sprintf("%v", o.Value)
}

func (o *ObjectOrRuntimeExpr) GetValue() interface{} {
	return o.Value
}

func (o *ObjectOrRuntimeExpr) AsStringOrMap() interface{} {
	switch o.Value.(type) {
	case map[string]interface{}:
		return o.Value.(map[string]interface{})
	case string:
		return o.Value.(string)
	case RuntimeExpression:
		return o.Value.(RuntimeExpression).Value
	}
	return nil
}

// UnmarshalJSON unmarshals data into either a RuntimeExpression or an object.
func (o *ObjectOrRuntimeExpr) UnmarshalJSON(data []byte) error {
	// Attempt to decode as a RuntimeExpression
	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(data, &runtimeExpr); err == nil && runtimeExpr.IsValid() {
		o.Value = runtimeExpr
		return nil
	}

	// Attempt to decode as a generic object
	var asObject map[string]interface{}
	if err := json.Unmarshal(data, &asObject); err == nil {
		o.Value = asObject
		return nil
	}

	// If neither succeeds, return an error
	return fmt.Errorf("ObjectOrRuntimeExpr must be a runtime expression or an object")
}

// MarshalJSON marshals ObjectOrRuntimeExpr into JSON.
func (o *ObjectOrRuntimeExpr) MarshalJSON() ([]byte, error) {
	switch v := o.Value.(type) {
	case RuntimeExpression:
		return json.Marshal(v.String())
	case map[string]interface{}:
		return json.Marshal(v)
	default:
		return nil, fmt.Errorf("ObjectOrRuntimeExpr contains unsupported type")
	}
}

// Validate validates the ObjectOrRuntimeExpr using the custom validation logic.
func (o *ObjectOrRuntimeExpr) Validate() error {
	switch v := o.Value.(type) {
	case RuntimeExpression:
		if !v.IsValid() {
			return fmt.Errorf("invalid runtime expression: %s", v.Value)
		}
	case map[string]interface{}:
		if len(v) == 0 {
			return fmt.Errorf("object cannot be empty")
		}
	default:
		return fmt.Errorf("unsupported value type for ObjectOrRuntimeExpr")
	}
	return nil
}

func NewStringOrRuntimeExpr(value string) *StringOrRuntimeExpr {
	return &StringOrRuntimeExpr{
		Value: value,
	}
}

// StringOrRuntimeExpr is a type that can hold either a RuntimeExpression or a string.
type StringOrRuntimeExpr struct {
	Value interface{} `json:"-" validate:"string_or_runtime_expr"` // Custom validation tag.
}

func (s *StringOrRuntimeExpr) AsObjectOrRuntimeExpr() *ObjectOrRuntimeExpr {
	return &ObjectOrRuntimeExpr{Value: s.Value}
}

// UnmarshalJSON unmarshals data into either a RuntimeExpression or a string.
func (s *StringOrRuntimeExpr) UnmarshalJSON(data []byte) error {
	// Attempt to decode as a RuntimeExpression
	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(data, &runtimeExpr); err == nil && runtimeExpr.IsValid() {
		s.Value = runtimeExpr
		return nil
	}

	// Attempt to decode as a string
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		s.Value = asString
		return nil
	}

	// If neither succeeds, return an error
	return fmt.Errorf("StringOrRuntimeExpr must be a runtime expression or a string")
}

// MarshalJSON marshals StringOrRuntimeExpr into JSON.
func (s *StringOrRuntimeExpr) MarshalJSON() ([]byte, error) {
	switch v := s.Value.(type) {
	case RuntimeExpression:
		return json.Marshal(v.String())
	case string:
		return json.Marshal(v)
	default:
		return nil, fmt.Errorf("StringOrRuntimeExpr contains unsupported type")
	}
}

func (s *StringOrRuntimeExpr) String() string {
	switch v := s.Value.(type) {
	case RuntimeExpression:
		return v.String()
	case string:
		return v
	default:
		return ""
	}
}

func (s *StringOrRuntimeExpr) GetValue() interface{} {
	return s.Value
}

// URITemplateOrRuntimeExpr represents a type that can be a URITemplate or a RuntimeExpression.
type URITemplateOrRuntimeExpr struct {
	Value interface{} `json:"-" validate:"uri_template_or_runtime_expr"` // Custom validation.
}

func NewUriTemplate(uriTemplate string) *URITemplateOrRuntimeExpr {
	return &URITemplateOrRuntimeExpr{
		Value: uriTemplate,
	}
}

// UnmarshalJSON unmarshals data into either a URITemplate or a RuntimeExpression.
func (u *URITemplateOrRuntimeExpr) UnmarshalJSON(data []byte) error {
	// Attempt to decode as URITemplate
	uriTemplate, err := UnmarshalURITemplate(data)
	if err == nil {
		u.Value = uriTemplate
		return nil
	}

	// Attempt to decode as RuntimeExpression
	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(data, &runtimeExpr); err == nil && runtimeExpr.IsValid() {
		u.Value = runtimeExpr
		return nil
	}

	// Return an error if neither succeeds
	return fmt.Errorf("URITemplateOrRuntimeExpr must be a valid URITemplate or RuntimeExpression")
}

// MarshalJSON marshals URITemplateOrRuntimeExpr into JSON.
func (u *URITemplateOrRuntimeExpr) MarshalJSON() ([]byte, error) {
	switch v := u.Value.(type) {
	case URITemplate:
		return json.Marshal(v.String())
	case RuntimeExpression:
		return json.Marshal(v.String())
	case string:
		// Attempt to marshal as RuntimeExpression
		runtimeExpr := RuntimeExpression{Value: v}
		if runtimeExpr.IsValid() {
			return json.Marshal(runtimeExpr.String())
		}
		// Otherwise, treat as a Literal URI
		uriTemplate, err := UnmarshalURITemplate([]byte(fmt.Sprintf(`"%s"`, v)))
		if err == nil {
			return json.Marshal(uriTemplate.String())
		}
		return nil, fmt.Errorf("invalid string for URITemplateOrRuntimeExpr: %s", v)
	default:
		return nil, fmt.Errorf("unsupported type for URITemplateOrRuntimeExpr: %T", v)
	}
}

func (u *URITemplateOrRuntimeExpr) String() string {
	switch v := u.Value.(type) {
	case URITemplate:
		return v.String()
	case RuntimeExpression:
		return v.String()
	case string:
		return v
	}
	return ""
}

func (u *URITemplateOrRuntimeExpr) GetValue() interface{} {
	return u.Value
}

// JsonPointerOrRuntimeExpression represents a type that can be a JSON Pointer or a RuntimeExpression.
type JsonPointerOrRuntimeExpression struct {
	Value interface{} `json:"-" validate:"json_pointer_or_runtime_expr"` // Custom validation tag.
}

// JSONPointerPattern validates JSON Pointers as per RFC 6901.
var JSONPointerPattern = regexp.MustCompile(`^(/([^/~]|~[01])*)*$`)

// UnmarshalJSON unmarshals data into either a JSON Pointer or a RuntimeExpression.
func (j *JsonPointerOrRuntimeExpression) UnmarshalJSON(data []byte) error {
	// Attempt to decode as a JSON Pointer
	var jsonPointer string
	if err := json.Unmarshal(data, &jsonPointer); err == nil {
		if JSONPointerPattern.MatchString(jsonPointer) {
			j.Value = jsonPointer
			return nil
		}
	}

	// Attempt to decode as RuntimeExpression
	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(data, &runtimeExpr); err == nil {
		if runtimeExpr.IsValid() {
			j.Value = runtimeExpr
			return nil
		}
	}

	// If neither succeeds, return an error
	return fmt.Errorf("JsonPointerOrRuntimeExpression must be a valid JSON Pointer or RuntimeExpression")
}

// MarshalJSON marshals JsonPointerOrRuntimeExpression into JSON.
func (j *JsonPointerOrRuntimeExpression) MarshalJSON() ([]byte, error) {
	switch v := j.Value.(type) {
	case string: // JSON Pointer
		return json.Marshal(v)
	case RuntimeExpression:
		return json.Marshal(v.String())
	default:
		return nil, fmt.Errorf("JsonPointerOrRuntimeExpression contains unsupported type")
	}
}

func (j *JsonPointerOrRuntimeExpression) String() string {
	switch v := j.Value.(type) {
	case RuntimeExpression:
		return v.String()
	case string:
		return v
	default:
		return ""
	}
}

func (j *JsonPointerOrRuntimeExpression) GetValue() interface{} {
	return j.Value
}

func (j *JsonPointerOrRuntimeExpression) IsValid() bool {
	return JSONPointerPattern.MatchString(j.String())
}
