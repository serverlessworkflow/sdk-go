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

func TestObjectOrRuntimeExpr_UnmarshalJSON(t *testing.T) {
	cases := []struct {
		Name      string
		JSON      string
		Expected  interface{}
		ShouldErr bool
	}{
		{
			Name:      "Unmarshal valid string",
			JSON:      `"${ expression }"`,
			Expected:  RuntimeExpression{Value: "${ expression }"},
			ShouldErr: false,
		},
		{
			Name: "Unmarshal valid object",
			JSON: `{
				"key": "value"
			}`,
			Expected: map[string]interface{}{
				"key": "value",
			},
			ShouldErr: false,
		},
		{
			Name:      "Unmarshal invalid type",
			JSON:      `123`,
			ShouldErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			var obj ObjectOrRuntimeExpr
			err := json.Unmarshal([]byte(tc.JSON), &obj)
			if tc.ShouldErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
				assert.Equal(t, tc.Expected, obj.Value, "unexpected unmarshalled value")
			}
		})
	}
}

func TestURITemplateOrRuntimeExprValidation(t *testing.T) {
	cases := []struct {
		Name      string
		Input     *URITemplateOrRuntimeExpr
		ShouldErr bool
	}{
		{
			Name: "Valid URI template",
			Input: &URITemplateOrRuntimeExpr{
				Value: &LiteralUriTemplate{Value: "http://example.com/{id}"},
			},
			ShouldErr: false,
		},
		{
			Name: "Valid URI",
			Input: &URITemplateOrRuntimeExpr{
				Value: &LiteralUri{Value: "http://example.com"},
			},
			ShouldErr: false,
		},
		{
			Name: "Valid runtime expression",
			Input: &URITemplateOrRuntimeExpr{
				Value: RuntimeExpression{Value: "${expression}"},
			},
			ShouldErr: false,
		},
		{
			Name: "Invalid runtime expression",
			Input: &URITemplateOrRuntimeExpr{
				Value: RuntimeExpression{Value: "123invalid-expression"},
			},
			ShouldErr: true,
		},
		{
			Name: "Invalid URI format",
			Input: &URITemplateOrRuntimeExpr{
				Value: &LiteralUri{Value: "invalid-uri"},
			},
			ShouldErr: true,
		},
		{
			Name: "Unsupported type",
			Input: &URITemplateOrRuntimeExpr{
				Value: 123,
			},
			ShouldErr: true,
		},
		{
			Name: "Valid URI as string",
			Input: &URITemplateOrRuntimeExpr{
				Value: "http://example.com",
			},
			ShouldErr: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validate.Var(tc.Input, "uri_template_or_runtime_expr")
			if tc.ShouldErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}
		})
	}
}

func TestJsonPointerOrRuntimeExpressionValidation(t *testing.T) {
	cases := []struct {
		Name      string
		Input     JsonPointerOrRuntimeExpression
		ShouldErr bool
	}{
		{
			Name: "Valid JSON Pointer",
			Input: JsonPointerOrRuntimeExpression{
				Value: "/valid/json/pointer",
			},
			ShouldErr: false,
		},
		{
			Name: "Valid runtime expression",
			Input: JsonPointerOrRuntimeExpression{
				Value: RuntimeExpression{Value: "${expression}"},
			},
			ShouldErr: false,
		},
		{
			Name: "Invalid JSON Pointer",
			Input: JsonPointerOrRuntimeExpression{
				Value: "invalid-json-pointer",
			},
			ShouldErr: true,
		},
		{
			Name: "Invalid runtime expression",
			Input: JsonPointerOrRuntimeExpression{
				Value: RuntimeExpression{Value: "123invalid-expression"},
			},
			ShouldErr: true,
		},
		{
			Name: "Unsupported type",
			Input: JsonPointerOrRuntimeExpression{
				Value: 123,
			},
			ShouldErr: true,
		},
	}

	for _, tc := range cases {
		t.Run(tc.Name, func(t *testing.T) {
			err := validate.Var(tc.Input, "json_pointer_or_runtime_expr")
			if tc.ShouldErr {
				assert.Error(t, err, "expected an error, but got none")
			} else {
				assert.NoError(t, err, "expected no error, but got one")
			}
		})
	}
}
