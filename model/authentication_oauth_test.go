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

func TestOAuth2AuthenticationPolicyValidation(t *testing.T) {
	testCases := []struct {
		name       string
		policy     OAuth2AuthenticationPolicy
		shouldPass bool
	}{
		{
			name: "Valid: Use set",
			policy: OAuth2AuthenticationPolicy{
				Use: "mysecret",
			},
			shouldPass: true,
		},
		{
			name: "Valid: Properties set",
			policy: OAuth2AuthenticationPolicy{
				Properties: &OAuth2AuthenticationProperties{
					Grant:     ClientCredentialsGrant,
					Scopes:    []string{"scope1", "scope2"},
					Authority: &LiteralUri{Value: "https://auth.example.com"},
				},
			},
			shouldPass: true,
		},
		{
			name: "Invalid: Both Use and Properties set",
			policy: OAuth2AuthenticationPolicy{
				Use: "mysecret",
				Properties: &OAuth2AuthenticationProperties{
					Grant:     ClientCredentialsGrant,
					Scopes:    []string{"scope1", "scope2"},
					Authority: &LiteralUri{Value: "https://auth.example.com"},
				},
			},
			shouldPass: false,
		},
		{
			name:       "Invalid: Neither Use nor Properties set",
			policy:     OAuth2AuthenticationPolicy{},
			shouldPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validate.Struct(tc.policy)
			if tc.shouldPass {
				if err != nil {
					t.Errorf("Expected validation to pass, but got error: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("Expected validation to fail, but it passed")
				}
			}
		})
	}
}

func TestAuthenticationOAuth2Policy(t *testing.T) {
	testCases := []struct {
		name       string
		input      string
		expected   string
		expectsErr bool
	}{
		{
			name: "Valid OAuth2 Authentication Inline",
			input: `{
				"oauth2": {
					"authority": "https://auth.example.com",
					"grant": "client_credentials",
					"scopes": ["scope1", "scope2"]
				}
			}`,
			expected:   `{"oauth2":{"authority":"https://auth.example.com","grant":"client_credentials","scopes":["scope1","scope2"]}}`,
			expectsErr: false,
		},
		{
			name: "Valid OAuth2 Authentication Use",
			input: `{
				"oauth2": {
					"use": "mysecret"
				}
			}`,
			expected:   `{"oauth2":{"use":"mysecret"}}`,
			expectsErr: false,
		},
		{
			name: "Invalid OAuth2: Both properties and use set",
			input: `{
				"oauth2": {
					"authority": "https://auth.example.com",
					"grant": "client_credentials",
					"use": "mysecret"
				}
			}`,
			expectsErr: true,
		},
		{
			name: "Invalid OAuth2: Missing required fields",
			input: `{
				"oauth2": {}
			}`,
			expectsErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var authPolicy AuthenticationPolicy

			// Unmarshal
			err := json.Unmarshal([]byte(tc.input), &authPolicy)
			if err == nil {
				err = validate.Struct(authPolicy)
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
