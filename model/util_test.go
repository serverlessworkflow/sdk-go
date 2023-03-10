// Copyright 2022 The Serverless Workflow Specification Authors
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

func TestIncludePaths(t *testing.T) {
	assert.NotNil(t, IncludePaths())
	assert.True(t, len(IncludePaths()) > 0)

	// update include paths
	paths := []string{"/root", "/path"}
	SetIncludePaths(paths)
	assert.Equal(t, IncludePaths(), paths)

	assert.PanicsWithError(t, "1 must be an absolute file path", assert.PanicTestFunc(func() {
		SetIncludePaths([]string{"1"})
	}))
}

func Test_primitiveOrMapType(t *testing.T) {
	type dataMap map[string]json.RawMessage

	t.Run("unmarshal", func(t *testing.T) {
		data := []byte(`"value":true`)
		_, _, err := primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.Error(t, err)

		data = []byte(`{value":true}`)
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.Error(t, err)

		data = []byte(`value":true}`)
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.Error(t, err)

		data = []byte(`"true"`)
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.Error(t, err)

		data = []byte(`true`)
		valMap, valBool, err := primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)
		assert.Nil(t, valMap)
		assert.True(t, valBool)

		data = []byte(`"true"`)
		valMap, valString, err := primitiveOrStruct[string, dataMap]("dataMap", data)
		assert.NoError(t, err)
		assert.Nil(t, valMap)
		assert.Equal(t, `true`, valString)

		data = []byte(`{"value":true}`)
		valMap, valBool, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)
		assert.NotNil(t, valMap)
		assert.Equal(t, valMap, &dataMap{"value": []byte("true")})
		assert.False(t, valBool)

		data = []byte(`{"value": "true"}`)
		valMap, valBool, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)
		assert.NotNil(t, valMap)
		assert.Equal(t, valMap, &dataMap{"value": []byte(`"true"`)})
		assert.False(t, valBool)
	})

	t.Run("test personalized syntaxError error message", func(t *testing.T) {
		type structBool struct {
			FieldValue string `json:"fieldValue"`
		}

		data := []byte(` {"fieldValue": "value" `)
		_, _, err := primitiveOrStruct[string, structBool]("structBool", data)
		assert.Error(t, err)
		assert.Equal(t, `structBool value '{"fieldValue": "value"' is not supported, it has a syntax error "unexpected end of JSON input"`, err.Error())

		data = []byte(` {"fieldValue": value} `)
		_, _, err = primitiveOrStruct[string, structBool]("structBool", data)
		assert.Error(t, err)
		assert.Equal(t, `structBool value '{"fieldValue": value}' is not supported, it has a syntax error "invalid character 'v' looking for beginning of value"`, err.Error())
	})

	t.Run("test personalized unmarshalTypeError error message", func(t *testing.T) {
		type structBool struct {
			FieldValue bool `json:"fieldValue"`
		}

		data := []byte(` {"fieldValue": "true"} `)
		_, _, err := primitiveOrStruct[bool, structBool]("structBool", data)
		assert.Error(t, err)
		assert.Equal(t, `structBool value '{"fieldValue": "true"}' is not supported, the value field fieldValue must be bool`, err.Error())

		data = []byte(` "true" `)
		_, _, err = primitiveOrStruct[bool, structBool]("structBool", data)
		assert.Error(t, err)
		assert.Equal(t, `structBool value '"true"' is not supported, it must be an object or bool`, err.Error())
	})

	t.Run("check json with spaces", func(t *testing.T) {
		data := []byte(` {"value": "true"} `)
		_, _, err := primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)

		data = []byte(` true `)
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)

		data = []byte(` "true" `)
		_, _, err = primitiveOrStruct[string, dataMap]("dataMap", data)
		assert.NoError(t, err)
	})

	t.Run("check tabs", func(t *testing.T) {
		data := []byte(string('\t') + `"true"` + string('\t'))
		_, _, err := primitiveOrStruct[string, dataMap]("dataMap", data)
		assert.NoError(t, err)

		data = []byte(string('\t') + `true` + string('\t'))
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)
	})

	t.Run("check breakline", func(t *testing.T) {
		data := []byte(string('\n') + `"true"` + string('\n'))
		_, _, err := primitiveOrStruct[string, dataMap]("dataMap", data)
		assert.NoError(t, err)

		data = []byte(string('\n') + `true` + string('\n'))
		_, _, err = primitiveOrStruct[bool, dataMap]("dataMap", data)
		assert.NoError(t, err)
	})
}
