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
	"fmt"
	"testing"
)

func TestAuthenticationPolicy(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expected   string
		expectsErr bool
	}{
		{
			name: "Valid Basic Authentication Inline",
			input: `{
				"basic": {
					"username": "john",
					"password": "12345"
				}
			}`,
			expected:   `{"basic":{"username":"john","password":"12345"}}`,
			expectsErr: false,
		},
		{
			name: "Valid Digest Authentication Inline",
			input: `{
				"digest": {
					"username": "digestUser",
					"password": "digestPass"
				}
			}`,
			expected:   `{"digest":{"username":"digestUser","password":"digestPass"}}`,
			expectsErr: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var authPolicy AuthenticationPolicy

			// Unmarshal
			err := json.Unmarshal([]byte(tc.input), &authPolicy)
			if err == nil {
				if authPolicy.Basic != nil {
					err = validate.Struct(authPolicy.Basic)
				}
				if authPolicy.Bearer != nil {
					err = validate.Struct(authPolicy.Bearer)
				}
				if authPolicy.Digest != nil {
					err = validate.Struct(authPolicy.Digest)
				}
				if authPolicy.OAuth2 != nil {
					err = validate.Struct(authPolicy.OAuth2)
				}
			}

			if tc.expectsErr {
				if err == nil {
					t.Errorf("Expected an error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}

				// Marshal
				marshaled, err := json.Marshal(authPolicy)
				if err != nil {
					t.Errorf("Failed to marshal: %v", err)
				}

				if string(marshaled) != tc.expected {
					t.Errorf("Expected %s but got %s", tc.expected, marshaled)
				}

				fmt.Printf("Test '%s' passed. Marshaled output: %s\n", tc.name, marshaled)
			}
		})
	}
}
