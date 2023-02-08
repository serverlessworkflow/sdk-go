// Copyright 2021 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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

func TestUnmarshalJSONMultipleAuthProperties(t *testing.T) {
	t.Run("BearerAuthProperties", func(t *testing.T) {
		a1JSON := `{
		"name": "a1",
		"scheme": "bearer",
		"properties": {
			"token": "token1"
		}
	}`
		a2JSON := `{
		"name": "a2",
		"scheme": "bearer",
		"properties": {
			"token": "token2"
		}
	}`

		var a1 Auth
		err := json.Unmarshal([]byte(a1JSON), &a1)
		assert.NoError(t, err)

		var a2 Auth
		err = json.Unmarshal([]byte(a2JSON), &a2)
		assert.NoError(t, err)

		a1Properties := a1.Properties.Bearer
		a2Properties := a2.Properties.Bearer

		assert.Equal(t, "token1", a1Properties.Token)
		assert.Equal(t, "token2", a2Properties.Token)
		assert.NotEqual(t, a1Properties, a2Properties)
	})

	t.Run("OAuth2AuthProperties", func(t *testing.T) {
		a1JSON := `{
	"name": "a1",
	"scheme": "oauth2",
	"properties": {
		"clientSecret": "secret1"
	}
}`

		a2JSON := `{
	"name": "a2",
	"scheme": "oauth2",
	"properties": {
		"clientSecret": "secret2"
	}
}`

		var a1 Auth
		err := json.Unmarshal([]byte(a1JSON), &a1)
		assert.NoError(t, err)

		var a2 Auth
		err = json.Unmarshal([]byte(a2JSON), &a2)
		assert.NoError(t, err)

		a1Properties := a1.Properties.OAuth2
		a2Properties := a2.Properties.OAuth2

		assert.Equal(t, "secret1", a1Properties.ClientSecret)
		assert.Equal(t, "secret2", a2Properties.ClientSecret)
		assert.NotEqual(t, a1Properties, a2Properties)
	})
}
