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

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
)

func AssertYAMLEq(t *testing.T, expected, actual string) {
	var expectedMap, actualMap map[string]interface{}

	// Unmarshal the expected YAML
	err := yaml.Unmarshal([]byte(expected), &expectedMap)
	assert.NoError(t, err, "failed to unmarshal expected YAML")

	// Unmarshal the actual YAML
	err = yaml.Unmarshal([]byte(actual), &actualMap)
	assert.NoError(t, err, "failed to unmarshal actual YAML")

	// Assert equality of the two maps
	assert.Equal(t, expectedMap, actualMap, "YAML structures do not match")
}
