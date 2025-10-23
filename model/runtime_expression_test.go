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

func TestRuntimeExpressionUnmarshalJSON(t *testing.T) {
	tests := []struct {
		Name      string
		JSONInput string
		Expected  string
		ExpectErr bool
	}{
		{
			Name:      "Valid RuntimeExpression",
			JSONInput: `{ "expression": "${runtime.value}" }`,
			Expected:  "${runtime.value}",
			ExpectErr: false,
		},
		{
			Name:      "Invalid RuntimeExpression",
			JSONInput: `{ "expression": "1234invalid_runtime" }`,
			Expected:  "",
			ExpectErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.Name, func(t *testing.T) {
			var acme *RuntimeExpressionAcme
			err := json.Unmarshal([]byte(tc.JSONInput), &acme)

			if tc.ExpectErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.Expected, acme.Expression.Value)
			}

			// Test marshalling
			if !tc.ExpectErr {
				output, err := json.Marshal(acme)
				assert.NoError(t, err)
				assert.JSONEq(t, tc.JSONInput, string(output))
			}
		})
	}
}

// EndpointAcme represents a struct using URITemplate.
type RuntimeExpressionAcme struct {
	Expression RuntimeExpression `json:"expression"`
}

func TestIsStrictExpr(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{
			name:       "StrictExpr with braces",
			expression: "${.some.path}",
			want:       true,
		},
		{
			name:       "Missing closing brace",
			expression: "${.some.path",
			want:       false,
		},
		{
			name:       "Missing opening brace",
			expression: ".some.path}",
			want:       false,
		},
		{
			name:       "Empty string",
			expression: "",
			want:       false,
		},
		{
			name:       "No braces at all",
			expression: ".some.path",
			want:       false,
		},
		{
			name:       "With spaces but still correct",
			expression: "${  .some.path   }",
			want:       true,
		},
		{
			name:       "Only braces",
			expression: "${}",
			want:       true, // Technically matches prefix+suffix
		},
		{
			name:       "With single quote",
			expression: "echo 'hello, I love ${ .project }",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsStrictExpr(tc.expression)
			if got != tc.want {
				t.Errorf("IsStrictExpr(%q) = %v, want %v", tc.expression, got, tc.want)
			}
		})
	}
}

func TestSanitize(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       string
	}{
		{
			name:       "Remove braces and replace single quotes",
			expression: "${ 'some.path' }",
			want:       `"some.path"`,
		},
		{
			name:       "Already sanitized string, no braces",
			expression: ".some.path",
			want:       ".some.path",
		},
		{
			name:       "Multiple single quotes",
			expression: "${ 'foo' + 'bar' }",
			want:       `"foo" + "bar"`,
		},
		{
			name:       "Only braces with spaces",
			expression: "${    }",
			want:       "",
		},
		{
			name:       "No braces, just single quotes to be replaced",
			expression: "'some.path'",
			want:       `"some.path"`,
		},
		{
			name:       "Nothing to sanitize",
			expression: "",
			want:       "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeExpr(tc.expression)
			if got != tc.want {
				t.Errorf("Sanitize(%q) = %q, want %q", tc.expression, got, tc.want)
			}
		})
	}
}

func TestIsValid(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       bool
	}{
		{
			name:       "Valid expression - simple path",
			expression: "${ .foo }",
			want:       true,
		},
		{
			name:       "Valid expression - array slice",
			expression: "${ .arr[0] }",
			want:       true,
		},
		{
			name:       "Invalid syntax",
			expression: "${ .foo( }",
			want:       false,
		},
		{
			name:       "No braces but valid JQ (directly provided)",
			expression: ".bar",
			want:       true,
		},
		{
			name:       "Empty expression",
			expression: "",
			want:       true, // empty is parseable but yields an empty query
		},
		{
			name:       "Invalid bracket usage",
			expression: "${ .arr[ }",
			want:       false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := IsValidExpr(tc.expression)
			if got != tc.want {
				t.Errorf("IsValid(%q) = %v, want %v", tc.expression, got, tc.want)
			}
		})
	}
}
