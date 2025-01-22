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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoint_UnmarshalJSON(t *testing.T) {
	t.Run("Valid RuntimeExpression", func(t *testing.T) {
		input := `"${example}"`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.NoError(t, err, "Unmarshal should not return an error")
		assert.NotNil(t, endpoint.RuntimeExpression, "RuntimeExpression should be set")
		assert.Equal(t, "${example}", endpoint.RuntimeExpression.Value, "RuntimeExpression value should match")
	})

	t.Run("Invalid RuntimeExpression", func(t *testing.T) {
		input := `"123invalid-expression"`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.Error(t, err, "Unmarshal should return an error for invalid runtime expression")
		assert.Nil(t, endpoint.RuntimeExpression, "RuntimeExpression should not be set")
	})

	t.Run("Invalid LiteralUriTemplate", func(t *testing.T) {
		uriTemplate := &LiteralUriTemplate{Value: "example.com/{id}"}
		assert.False(t, LiteralUriPattern.MatchString(uriTemplate.Value), "LiteralUriTemplate should not match URI pattern")
	})

	t.Run("Valid URITemplate", func(t *testing.T) {
		input := `"http://example.com/{id}"`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.NoError(t, err, "Unmarshal should not return an error")
		assert.NotNil(t, endpoint.URITemplate, "URITemplate should be set")
	})

	t.Run("Valid EndpointConfiguration", func(t *testing.T) {
		input := `{
			"uri": "http://example.com/{id}",
			"authentication": {
				"basic": { "username": "admin", "password": "admin" }
			}
		}`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.NoError(t, err, "Unmarshal should not return an error")
		assert.NotNil(t, endpoint.EndpointConfig, "EndpointConfig should be set")
		assert.Equal(t, "admin", endpoint.EndpointConfig.Authentication.AuthenticationPolicy.Basic.Username, "Authentication Username should match")
		assert.Equal(t, "admin", endpoint.EndpointConfig.Authentication.AuthenticationPolicy.Basic.Password, "Authentication Password should match")
	})

	t.Run("Invalid JSON Structure", func(t *testing.T) {
		input := `{"invalid": "data"}`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.Error(t, err, "Unmarshal should return an error for invalid JSON structure")
	})

	t.Run("Empty Input", func(t *testing.T) {
		input := `{}`
		var endpoint Endpoint
		err := json.Unmarshal([]byte(input), &endpoint)

		assert.NoError(t, err, "Unmarshal should not return an error for empty input")
		assert.Nil(t, endpoint.RuntimeExpression, "RuntimeExpression should not be set")
		assert.Nil(t, endpoint.URITemplate, "URITemplate should not be set")
		assert.Nil(t, endpoint.EndpointConfig, "EndpointConfig should not be set")
	})
}

func TestEndpoint_MarshalJSON(t *testing.T) {
	t.Run("Marshal RuntimeExpression", func(t *testing.T) {
		endpoint := &Endpoint{
			RuntimeExpression: &RuntimeExpression{Value: "${example}"},
		}

		data, err := json.Marshal(endpoint)
		assert.NoError(t, err, "Marshal should not return an error")
		assert.JSONEq(t, `"${example}"`, string(data), "Output JSON should match")
	})

	t.Run("Marshal URITemplate", func(t *testing.T) {
		endpoint := &Endpoint{
			URITemplate: &LiteralUriTemplate{Value: "http://example.com/{id}"},
		}

		data, err := json.Marshal(endpoint)
		assert.NoError(t, err, "Marshal should not return an error")
		assert.JSONEq(t, `"http://example.com/{id}"`, string(data), "Output JSON should match")
	})

	t.Run("Marshal EndpointConfiguration", func(t *testing.T) {
		endpoint := &Endpoint{
			EndpointConfig: &EndpointConfiguration{
				URI: &LiteralUriTemplate{Value: "http://example.com/{id}"},
				Authentication: &ReferenceableAuthenticationPolicy{AuthenticationPolicy: &AuthenticationPolicy{Basic: &BasicAuthenticationPolicy{
					Username: "john",
					Password: "secret",
				}}},
			},
		}

		data, err := json.Marshal(endpoint)
		assert.NoError(t, err, "Marshal should not return an error")
		expected := `{
			"uri": "http://example.com/{id}",
			"authentication": {
				"basic": { "username": "john", "password": "secret" }
			}
		}`
		assert.JSONEq(t, expected, string(data), "Output JSON should match")
	})

	t.Run("Marshal Empty Endpoint", func(t *testing.T) {
		endpoint := Endpoint{}

		data, err := json.Marshal(endpoint)
		assert.NoError(t, err, "Marshal should not return an error")
		assert.JSONEq(t, `{}`, string(data), "Output JSON should be empty")
	})
}
