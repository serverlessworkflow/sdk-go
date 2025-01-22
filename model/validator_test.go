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
	"testing"
)

func TestRegexValidators(t *testing.T) {
	testCases := []struct {
		name     string
		validate func(string) bool
		input    string
		expected bool
	}{
		// ISO 8601 Duration Tests
		{"ISO 8601 Duration Valid 1", isISO8601DurationValid, "P2Y", true},
		{"ISO 8601 Duration Valid 2", isISO8601DurationValid, "P1DT12H30M", true},
		{"ISO 8601 Duration Valid 3", isISO8601DurationValid, "P1Y2M3D", true},
		{"ISO 8601 Duration Valid 4", isISO8601DurationValid, "P1Y2M3D4H", false},
		{"ISO 8601 Duration Valid 5", isISO8601DurationValid, "P1Y", true},
		{"ISO 8601 Duration Valid 6", isISO8601DurationValid, "PT1H", true},
		{"ISO 8601 Duration Valid 7", isISO8601DurationValid, "P1Y2M3D4H5M6S", false},
		{"ISO 8601 Duration Invalid 1", isISO8601DurationValid, "P", false},
		{"ISO 8601 Duration Invalid 2", isISO8601DurationValid, "P1Y2M3D4H5M6S7", false},
		{"ISO 8601 Duration Invalid 3", isISO8601DurationValid, "1Y", false},

		// Semantic Versioning Tests
		{"Semantic Version Valid 1", isSemanticVersionValid, "1.0.0", true},
		{"Semantic Version Valid 2", isSemanticVersionValid, "1.2.3", true},
		{"Semantic Version Valid 3", isSemanticVersionValid, "1.2.3-beta", true},
		{"Semantic Version Valid 4", isSemanticVersionValid, "1.2.3-beta.1", true},
		{"Semantic Version Valid 5", isSemanticVersionValid, "1.2.3-beta.1+build.123", true},
		{"Semantic Version Invalid 1", isSemanticVersionValid, "v1.2.3", false},
		{"Semantic Version Invalid 2", isSemanticVersionValid, "1.2", false},
		{"Semantic Version Invalid 3", isSemanticVersionValid, "1.2.3-beta.x", true},

		// RFC 1123 Hostname Tests
		{"RFC 1123 Hostname Valid 1", isHostnameValid, "example.com", true},
		{"RFC 1123 Hostname Valid 2", isHostnameValid, "my-hostname", true},
		{"RFC 1123 Hostname Valid 3", isHostnameValid, "subdomain.example.com", true},
		{"RFC 1123 Hostname Invalid 1", isHostnameValid, "127.0.0.1", false},
		{"RFC 1123 Hostname Invalid 2", isHostnameValid, "example.com.", false},
		{"RFC 1123 Hostname Invalid 3", isHostnameValid, "example..com", false},
		{"RFC 1123 Hostname Invalid 4", isHostnameValid, "example.com-", false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.validate(tc.input)
			if result != tc.expected {
				t.Errorf("Validation failed for '%s': input='%s', expected=%v, got=%v", tc.name, tc.input, tc.expected, result)
			}
		})
	}
}
