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

	"github.com/tidwall/gjson"
)

// LiteralUriPattern matches standard URIs without placeholders.
var LiteralUriPattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+\-.]*://[^{}\s]+$`)

// LiteralUriTemplatePattern matches URIs with placeholders.
var LiteralUriTemplatePattern = regexp.MustCompile(`^[A-Za-z][A-Za-z0-9+\-.]*://.*\{.*}.*$`)

// URITemplate represents a URI that can be a literal URI or a URI template.
type URITemplate interface {
	IsURITemplate() bool
	String() string
	GetValue() interface{}
}

// UnmarshalURITemplate is a shared function for unmarshalling URITemplate fields.
func UnmarshalURITemplate(data []byte) (URITemplate, error) {
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal URITemplate: %w", err)
	}

	if LiteralUriTemplatePattern.MatchString(raw) {
		return &LiteralUriTemplate{Value: raw}, nil
	}

	if LiteralUriPattern.MatchString(raw) {
		return &LiteralUri{Value: raw}, nil
	}

	return nil, fmt.Errorf("invalid URI or URI template format: %s", raw)
}

type LiteralUriTemplate struct {
	Value string `json:"-" validate:"required,uri_template_pattern"` // Validate pattern for URI template.
}

func (t *LiteralUriTemplate) IsURITemplate() bool {
	return true
}

func (t *LiteralUriTemplate) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.Value)
}

func (t *LiteralUriTemplate) String() string {
	return t.Value
}

func (t *LiteralUriTemplate) GetValue() interface{} {
	return t.Value
}

type LiteralUri struct {
	Value string `json:"-" validate:"required,uri_pattern"` // Validate pattern for URI.
}

func (u *LiteralUri) IsURITemplate() bool {
	return true
}

func (u *LiteralUri) MarshalJSON() ([]byte, error) {
	return json.Marshal(u.Value)
}

func (u *LiteralUri) String() string {
	return u.Value
}

func (u *LiteralUri) GetValue() interface{} {
	return u.Value
}

type EndpointConfiguration struct {
	RuntimeExpression *RuntimeExpression                 `json:"-"`
	URI               URITemplate                        `json:"uri" validate:"required"`
	Authentication    *ReferenceableAuthenticationPolicy `json:"authentication,omitempty"`
}

// UnmarshalJSON implements custom unmarshalling for EndpointConfiguration.
func (e *EndpointConfiguration) UnmarshalJSON(data []byte) error {
	// Use a temporary structure to unmarshal the JSON
	type Alias EndpointConfiguration
	temp := &struct {
		URI json.RawMessage `json:"uri"`
		*Alias
	}{
		Alias: (*Alias)(e),
	}

	if err := json.Unmarshal(data, &temp); err != nil {
		return fmt.Errorf("failed to unmarshal EndpointConfiguration: %w", err)
	}

	// Unmarshal the URI field into the appropriate URITemplate implementation
	uri, err := UnmarshalURITemplate(temp.URI)
	if err == nil {
		e.URI = uri
		return nil
	}

	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(temp.URI, &runtimeExpr); err == nil && runtimeExpr.IsValid() {
		e.RuntimeExpression = &runtimeExpr
		return nil
	}

	return errors.New("failed to unmarshal EndpointConfiguration: data does not match any known schema")
}

// MarshalJSON implements custom marshalling for Endpoint.
func (e *EndpointConfiguration) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	if e.Authentication != nil {
		m["authentication"] = e.Authentication
	}

	if e.RuntimeExpression != nil {
		m["uri"] = e.RuntimeExpression
	} else if e.URI != nil {
		m["uri"] = e.URI
	}

	// Return an empty JSON object when no fields are set
	return json.Marshal(m)
}

type Endpoint struct {
	RuntimeExpression *RuntimeExpression     `json:"-"`
	URITemplate       URITemplate            `json:"-"`
	EndpointConfig    *EndpointConfiguration `json:"-"`
}

func NewEndpoint(uri string) *Endpoint {
	return &Endpoint{URITemplate: &LiteralUri{Value: uri}}
}

func (e *Endpoint) String() string {
	if e.RuntimeExpression != nil {
		return e.RuntimeExpression.String()
	}
	if e.URITemplate != nil {
		return e.URITemplate.String()
	}
	if e.EndpointConfig != nil {
		return e.EndpointConfig.URI.String()
	}
	return ""
}

// UnmarshalJSON implements custom unmarshalling for Endpoint.
func (e *Endpoint) UnmarshalJSON(data []byte) error {
	if gjson.ValidBytes(data) && gjson.ParseBytes(data).IsObject() && len(gjson.ParseBytes(data).Map()) == 0 {
		// Leave the Endpoint fields unset (nil)
		return nil
	}

	// Then try to unmarshal as URITemplate
	if uriTemplate, err := UnmarshalURITemplate(data); err == nil {
		e.URITemplate = uriTemplate
		return nil
	}

	// First try to unmarshal as RuntimeExpression
	var runtimeExpr RuntimeExpression
	if err := json.Unmarshal(data, &runtimeExpr); err == nil && runtimeExpr.IsValid() {
		e.RuntimeExpression = &runtimeExpr
		return nil
	}

	// Finally, try to unmarshal as EndpointConfiguration
	var endpointConfig EndpointConfiguration
	if err := json.Unmarshal(data, &endpointConfig); err == nil {
		e.EndpointConfig = &endpointConfig
		return nil
	}

	return errors.New("failed to unmarshal Endpoint: data does not match any known schema")
}

// MarshalJSON implements custom marshalling for Endpoint.
func (e *Endpoint) MarshalJSON() ([]byte, error) {
	if e.RuntimeExpression != nil {
		return json.Marshal(e.RuntimeExpression)
	}
	if e.URITemplate != nil {
		return json.Marshal(e.URITemplate)
	}
	if e.EndpointConfig != nil {
		return json.Marshal(e.EndpointConfig)
	}
	// Return an empty JSON object when no fields are set
	return []byte("{}"), nil
}
