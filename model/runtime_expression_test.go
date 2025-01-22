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
