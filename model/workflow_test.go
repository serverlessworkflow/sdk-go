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
	"testing"

	validator "github.com/go-playground/validator/v10"

	"github.com/stretchr/testify/assert"
)

func TestDocument_JSONMarshal(t *testing.T) {
	doc := Document{
		DSL:       "1.0.0",
		Namespace: "example-namespace",
		Name:      "example-name",
		Version:   "1.0.0",
		Title:     "Example Workflow",
		Summary:   "This is a sample workflow document.",
		Tags: map[string]string{
			"env":  "prod",
			"team": "workflow",
		},
		Metadata: map[string]interface{}{
			"author":  "John Doe",
			"created": "2025-01-01",
		},
	}

	data, err := json.Marshal(doc)
	assert.NoError(t, err)

	expectedJSON := `{
		"dsl": "1.0.0",
		"namespace": "example-namespace",
		"name": "example-name",
		"version": "1.0.0",
		"title": "Example Workflow",
		"summary": "This is a sample workflow document.",
		"tags": {
			"env": "prod",
			"team": "workflow"
		},
		"metadata": {
			"author": "John Doe",
			"created": "2025-01-01"
		}
	}`

	// Use JSON comparison to avoid formatting mismatches
	var expected, actual map[string]interface{}
	assert.NoError(t, json.Unmarshal([]byte(expectedJSON), &expected))
	assert.NoError(t, json.Unmarshal(data, &actual))
	assert.Equal(t, expected, actual)
}

func TestDocument_JSONUnmarshal(t *testing.T) {
	inputJSON := `{
		"dsl": "1.0.0",
		"namespace": "example-namespace",
		"name": "example-name",
		"version": "1.0.0",
		"title": "Example Workflow",
		"summary": "This is a sample workflow document.",
		"tags": {
			"env": "prod",
			"team": "workflow"
		},
		"metadata": {
			"author": "John Doe",
			"created": "2025-01-01"
		}
	}`

	var doc Document
	err := json.Unmarshal([]byte(inputJSON), &doc)
	assert.NoError(t, err)

	expected := Document{
		DSL:       "1.0.0",
		Namespace: "example-namespace",
		Name:      "example-name",
		Version:   "1.0.0",
		Title:     "Example Workflow",
		Summary:   "This is a sample workflow document.",
		Tags: map[string]string{
			"env":  "prod",
			"team": "workflow",
		},
		Metadata: map[string]interface{}{
			"author":  "John Doe",
			"created": "2025-01-01",
		},
	}

	assert.Equal(t, expected, doc)
}

func TestDocument_JSONUnmarshal_InvalidJSON(t *testing.T) {
	invalidJSON := `{
		"dsl": "1.0.0",
		"namespace": "example-namespace",
		"name": "example-name",
		"version": "1.0.0",
		"tags": {
			"env": "prod",
			"team": "workflow"
		"metadata": {
			"author": "John Doe",
			"created": "2025-01-01"
		}
	}` // Missing closing brace for "tags"

	var doc Document
	err := json.Unmarshal([]byte(invalidJSON), &doc)
	assert.Error(t, err)
}

func TestDocument_Validation_MissingRequiredField(t *testing.T) {
	inputJSON := `{
		"namespace": "example-namespace",
		"name": "example-name",
		"version": "1.0.0"
	}` // Missing "dsl"

	var doc Document
	err := json.Unmarshal([]byte(inputJSON), &doc)
	assert.NoError(t, err) // JSON is valid for unmarshalling

	// Validate the struct
	err = validate.Struct(doc)
	assert.Error(t, err)

	// Assert that the error is specifically about the missing "dsl" field
	assert.Contains(t, err.Error(), "Key: 'Document.DSL' Error:Field validation for 'DSL' failed on the 'required' tag")
}

func TestSchemaValidation(t *testing.T) {

	tests := []struct {
		name      string
		jsonInput string
		valid     bool
	}{
		// Valid Cases
		{
			name: "Valid Inline Schema",
			jsonInput: `{
				"document": "{\"key\":\"value\"}"
			}`,
			valid: true,
		},
		{
			name: "Valid External Schema",
			jsonInput: `{
				"resource": {
					"name": "external-schema",
					"endpoint": {
						"uri": "http://example.com/schema"
					}
				}
			}`,
			valid: true,
		},
		{
			name: "Valid External Schema Without Name",
			jsonInput: `{
				"resource": {
					"endpoint": {
						"uri": "http://example.com/schema"
					}
				}
			}`,
			valid: true,
		},
		{
			name: "Valid Inline Schema with Format",
			jsonInput: `{
				"format": "yaml",
				"document": "{\"key\":\"value\"}"
			}`,
			valid: true,
		},
		{
			name: "Valid External Schema with Format",
			jsonInput: `{
				"format": "xml",
				"resource": {
					"name": "external-schema",
					"endpoint": {
						"uri": "http://example.com/schema"
					}
				}
			}`,
			valid: true,
		},
		// Invalid Cases
		{
			name: "Invalid Both Document and Resource",
			jsonInput: `{
				"document": "{\"key\":\"value\"}",
				"resource": {
					"endpoint": {
						"uri": "http://example.com/schema"
					}
				}
			}`,
			valid: false,
		},
		{
			name: "Invalid Missing Both Document and Resource",
			jsonInput: `{
				"format": "json"
			}`,
			valid: false,
		},
		{
			name: "Invalid Resource Without Endpoint",
			jsonInput: `{
				"resource": {
					"name": "external-schema"
				}
			}`,
			valid: false,
		},
		{
			name: "Invalid Resource with Invalid URL",
			jsonInput: `{
				"resource": {
					"name": "external-schema",
					"endpoint": {
						"uri": "not-a-valid-url"
					}
				}
			}`,
			valid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var schema Schema
			err := json.Unmarshal([]byte(tt.jsonInput), &schema)
			if tt.valid {
				// Assert no unmarshalling error
				assert.NoError(t, err)

				// Validate the struct
				err = validate.Struct(schema)
				assert.NoError(t, err, "Expected valid schema but got validation error: %v", err)
			} else {
				// Assert unmarshalling or validation error
				if err == nil {
					err = validate.Struct(schema)
				}
				assert.Error(t, err, "Expected validation error but got none")
			}
		})
	}
}

type InputTestCase struct {
	Name      string
	Input     Input
	ShouldErr bool
}

func TestInputValidation(t *testing.T) {
	cases := []InputTestCase{
		{
			Name: "Valid input with Schema and From (object)",
			Input: Input{
				Schema: &Schema{
					Format: "json",
					Document: func() *string {
						doc := "example schema"
						return &doc
					}(),
				},
				From: &ObjectOrRuntimeExpr{
					Value: map[string]interface{}{
						"key": "value",
					},
				},
			},
			ShouldErr: false,
		},
		{
			Name: "Invalid input with Schema and From (expr)",
			Input: Input{
				Schema: &Schema{
					Format: "json",
				},
				From: &ObjectOrRuntimeExpr{
					Value: "example input",
				},
			},
			ShouldErr: true,
		},
		{
			Name: "Valid input with Schema and From (expr)",
			Input: Input{
				Schema: &Schema{
					Format: "json",
				},
				From: &ObjectOrRuntimeExpr{
					Value: "${ expression }",
				},
			},
			ShouldErr: true,
		},
		{
			Name: "Invalid input with Empty From (expr)",
			Input: Input{
				From: &ObjectOrRuntimeExpr{
					Value: "",
				},
			},
			ShouldErr: true,
		},
		{
			Name: "Invalid input with Empty From (object)",
			Input: Input{
				From: &ObjectOrRuntimeExpr{
					Value: map[string]interface{}{},
				},
			},
			ShouldErr: true,
		},
		{
			Name: "Invalid input with Unsupported From Type",
			Input: Input{
				From: &ObjectOrRuntimeExpr{
					Value: 123,
				},
			},
			ShouldErr: true,
		},
		{
			Name: "Valid input with Schema Only",
			Input: Input{
				Schema: &Schema{
					Format: "json",
				},
			},
			ShouldErr: false,
		},
		{
			Name:      "input with Neither Schema Nor From",
			Input:     Input{},
			ShouldErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validate.Struct(tc.Input)
			if tc.ShouldErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}
		})
	}
}

func TestFlowDirectiveValidation(t *testing.T) {
	cases := []struct {
		Name      string
		Input     FlowDirective
		IsEnum    bool // Expected result for IsEnum method.
		ShouldErr bool // Expected result for validation.
	}{
		{
			Name:      "Valid Enum: continue",
			Input:     FlowDirective{Value: "continue"},
			IsEnum:    true,
			ShouldErr: false,
		},
		{
			Name:      "Valid Enum: exit",
			Input:     FlowDirective{Value: "exit"},
			IsEnum:    true,
			ShouldErr: false,
		},
		{
			Name:      "Valid Enum: end",
			Input:     FlowDirective{Value: "end"},
			IsEnum:    true,
			ShouldErr: false,
		},
		{
			Name:      "Valid Free-form String",
			Input:     FlowDirective{Value: "custom-directive"},
			IsEnum:    false,
			ShouldErr: false,
		},
		{
			Name:      "Invalid Empty String",
			Input:     FlowDirective{Value: ""},
			IsEnum:    false,
			ShouldErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			// Validate the struct
			err := validate.Var(tc.Input.Value, "required")
			if tc.ShouldErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}

			// Check IsEnum result
			assert.Equal(t, tc.IsEnum, tc.Input.IsEnum(), "unexpected IsEnum result")
		})
	}
}

func TestUse_MarshalJSON(t *testing.T) {
	use := Use{
		Authentications: map[string]*AuthenticationPolicy{
			"auth1": NewBasicAuth("alice", "secret"),
		},
		Errors: map[string]*Error{
			"error1": {Type: NewUriTemplate("http://example.com/errors"), Status: 404},
		},
		Extensions: ExtensionList{
			{Key: "ext1", Extension: &Extension{Extend: "call"}},
			{Key: "ext2", Extension: &Extension{Extend: "emit"}},
			{Key: "ext3", Extension: &Extension{Extend: "for"}},
		},
		Functions: NamedTaskMap{
			"func1": &CallHTTP{Call: "http", With: HTTPArguments{Endpoint: NewEndpoint("http://example.com/"), Method: "GET"}},
		},
		Retries: map[string]*RetryPolicy{
			"retry1": {
				Delay: NewDurationExpr("PT5S"),
				Limit: RetryLimit{Attempt: &RetryLimitAttempt{Count: 3}},
			},
		},
		Secrets:  []string{"secret1", "secret2"},
		Timeouts: map[string]*Timeout{"timeout1": {After: NewDurationExpr("PT1M")}},
		Catalogs: map[string]*Catalog{
			"catalog1": {Endpoint: NewEndpoint("http://example.com")},
		},
	}

	data, err := json.Marshal(use)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"authentications": {"auth1": { "basic": {"username": "alice", "password": "secret"}}},
		"errors": {"error1": {"type": "http://example.com/errors", "status": 404}},
		"extensions": [
			{"ext1": {"extend": "call"}},
			{"ext2": {"extend": "emit"}},
			{"ext3": {"extend": "for"}}
		],
		"functions": {"func1": {"call": "http", "with": {"endpoint": "http://example.com/", "method": "GET"}}},
		"retries": {"retry1": {"delay": "PT5S", "limit": {"attempt": {"count": 3}}}},
		"secrets": ["secret1", "secret2"],
		"timeouts": {"timeout1": {"after": "PT1M"}},
		"catalogs": {"catalog1": {"endpoint": "http://example.com"}}
	}`, string(data))
}

func TestUse_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"authentications": {"auth1": { "basic": {"username": "alice", "password": "secret"}}},
		"errors": {"error1": {"type": "http://example.com/errors", "status": 404}},
		"extensions": [{"ext1": {"extend": "call"}}],
		"functions": {"func1": {"call": "http", "with": {"endpoint": "http://example.com"}}},
		"retries": {"retry1": {"delay": "PT5S", "limit": {"attempt": {"count": 3}}}},
		"secrets": ["secret1", "secret2"],
		"timeouts": {"timeout1": {"after": "PT1M"}},
		"catalogs": {"catalog1": {"endpoint": "http://example.com"}}
	}`

	var use Use
	err := json.Unmarshal([]byte(jsonData), &use)
	assert.NoError(t, err)

	assert.NotNil(t, use.Authentications["auth1"])
	assert.Equal(t, "alice", use.Authentications["auth1"].Basic.Username)
	assert.Equal(t, "secret", use.Authentications["auth1"].Basic.Password)

	assert.NotNil(t, use.Errors["error1"])
	assert.Equal(t, "http://example.com/errors", use.Errors["error1"].Type.String())
	assert.Equal(t, 404, use.Errors["error1"].Status)

	assert.NotNil(t, use.Extensions.Key("ext1"))
	assert.Equal(t, "call", use.Extensions.Key("ext1").Extend)

	assert.NotNil(t, use.Functions["func1"])
	assert.IsType(t, &CallHTTP{With: HTTPArguments{Endpoint: NewEndpoint("http://example.com")}}, use.Functions["func1"])

	assert.NotNil(t, use.Retries["retry1"])
	assert.Equal(t, "PT5S", use.Retries["retry1"].Delay.AsExpression())
	assert.Equal(t, 3, use.Retries["retry1"].Limit.Attempt.Count)

	assert.Equal(t, []string{"secret1", "secret2"}, use.Secrets)

	assert.NotNil(t, use.Timeouts["timeout1"])
	assert.Equal(t, "PT1M", use.Timeouts["timeout1"].After.AsExpression())

	assert.NotNil(t, use.Catalogs["catalog1"])
	assert.Equal(t, "http://example.com", use.Catalogs["catalog1"].Endpoint.URITemplate.String())
}

func TestUse_Validation(t *testing.T) {
	use := &Use{
		Authentications: map[string]*AuthenticationPolicy{
			"auth1": NewBasicAuth("alice", "secret"),
		},
		Errors: map[string]*Error{
			"error1": {Type: &URITemplateOrRuntimeExpr{&LiteralUri{"http://example.com/errors"}}, Status: 404},
		},
		Extensions: ExtensionList{},
		Functions: map[string]Task{
			"func1": &CallHTTP{Call: "http", With: HTTPArguments{Endpoint: NewEndpoint("http://example.com"), Method: "GET"}},
		},
		Retries: map[string]*RetryPolicy{
			"retry1": {
				Delay: NewDurationExpr("PT5S"),
				Limit: RetryLimit{Attempt: &RetryLimitAttempt{Count: 3}},
			},
		},
		Secrets:  []string{"secret1", "secret2"},
		Timeouts: map[string]*Timeout{"timeout1": {After: NewDurationExpr("PT1M")}},
		Catalogs: map[string]*Catalog{
			"catalog1": {Endpoint: NewEndpoint("http://example.com")},
		},
	}

	err := validate.Struct(use)
	assert.NoError(t, err)

	// Test with missing required fields
	use.Catalogs["catalog1"].Endpoint = nil
	err = validate.Struct(use)
	assert.Error(t, err)

	var validationErrors validator.ValidationErrors
	if errors.As(err, &validationErrors) {
		for _, validationErr := range validationErrors {
			t.Logf("Validation failed on field '%s' with tag '%s'", validationErr.Namespace(), validationErr.Tag())
		}

		assert.Contains(t, validationErrors.Error(), "Catalogs[catalog1].Endpoint")
		assert.Contains(t, validationErrors.Error(), "required")
	}
}
