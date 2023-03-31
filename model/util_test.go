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
	"fmt"
	"net/http"
	"net/http/httptest"
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

func Test_getBytesFromFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		switch req.URL.Path {
		case "/test.json":
			_, err := rw.Write([]byte("{}"))
			assert.NoError(t, err)
		default:
			t.Failed()
		}
	}))
	defer server.Close()
	httpClient = *server.Client()

	data, err := getBytesFromFile(server.URL + "/test.json")
	assert.NoError(t, err)
	assert.Equal(t, "{}", string(data))

	data, err = getBytesFromFile("../parser/testdata/eventdefs.yml")
	assert.NoError(t, err)
	assert.Equal(t, "[{\"correlation\":[{\"contextAttributeName\":\"accountId\"}],\"name\":\"PaymentReceivedEvent\",\"source\":\"paymentEventSource\",\"type\":\"payment.receive\"},{\"kind\":\"produced\",\"name\":\"ConfirmationCompletedEvent\",\"type\":\"payment.confirmation\"}]", string(data))
}

func Test_unmarshalObjectOrFile(t *testing.T) {
	t.Run("httptest", func(t *testing.T) {
		type structString struct {
			FieldValue string `json:"fieldValue"`
		}
		type listStructString []structString

		server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			switch req.URL.Path {
			case "/test.json":
				_, err := rw.Write([]byte(`[{"fieldValue": "value"}]`))
				assert.NoError(t, err)
			default:
				t.Failed()
			}
		}))
		defer server.Close()
		httpClient = *server.Client()

		structValue := &structString{}
		data := []byte(`"fieldValue": "value"`)
		err := unmarshalObjectOrFile("structString", data, structValue)
		assert.Error(t, err)
		assert.Equal(t, &structString{}, structValue)

		listStructValue := &listStructString{}
		data = []byte(`[{"fieldValue": "value"}]`)
		err = unmarshalObjectOrFile("listStructString", data, listStructValue)
		assert.NoError(t, err)
		assert.Equal(t, listStructString{{FieldValue: "value"}}, *listStructValue)

		listStructValue = &listStructString{}
		data = []byte(fmt.Sprintf(`"%s/test.json"`, server.URL))
		err = unmarshalObjectOrFile("listStructString", data, listStructValue)
		assert.NoError(t, err)
		assert.Equal(t, listStructString{{FieldValue: "value"}}, *listStructValue)
	})

	t.Run("file://", func(t *testing.T) {
		retries := &Retries{}
		data := []byte(`"file://../parser/testdata/applicationrequestretries.json"`)
		err := unmarshalObjectOrFile("retries", data, retries)
		assert.NoError(t, err)
	})

	t.Run("external url", func(t *testing.T) {
		retries := &Retries{}
		data := []byte(`"https://raw.githubusercontent.com/serverlessworkflow/sdk-net/main/src/ServerlessWorkflow.Sdk.UnitTests/Resources/retries/default.yaml"`)
		err := unmarshalObjectOrFile("retries", data, retries)
		assert.NoError(t, err)
	})

}

func Test_primitiveOrMapType(t *testing.T) {
	type dataMap map[string]json.RawMessage

	t.Run("unmarshal", func(t *testing.T) {
		var valBool bool
		valMap := &dataMap{}
		data := []byte(`"value":true`)
		err := unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.Error(t, err)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`{value":true}`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.Error(t, err)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`value":true}`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.Error(t, err)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`"true"`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.Error(t, err)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`true`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.NoError(t, err)
		assert.Equal(t, &dataMap{}, valMap)
		assert.True(t, valBool)

		valString := ""
		valMap = &dataMap{}
		data = []byte(`"true"`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valString, valMap)
		assert.NoError(t, err)
		assert.Equal(t, &dataMap{}, valMap)
		assert.Equal(t, `true`, valString)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`{"value":true}`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.NoError(t, err)
		assert.NotNil(t, valMap)
		assert.Equal(t, valMap, &dataMap{"value": []byte("true")})
		assert.False(t, valBool)

		valBool = false
		valMap = &dataMap{}
		data = []byte(`{"value": "true"}`)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valMap)
		assert.NoError(t, err)
		assert.NotNil(t, valMap)
		assert.Equal(t, valMap, &dataMap{"value": []byte(`"true"`)})
		assert.False(t, valBool)
	})

	t.Run("test personalized syntaxError error message", func(t *testing.T) {
		type structString struct {
			FieldValue string `json:"fieldValue"`
		}

		var valString string
		valStruct := &structString{}
		data := []byte(`{"fieldValue": "value"`)
		err := unmarshalPrimitiveOrObject("structBool", data, &valString, valStruct)
		assert.Error(t, err)
		assert.Equal(t, "structBool has a syntax error \"unexpected end of JSON input\"", err.Error())

		data = []byte(`{\n  "fieldValue": value\n}`)
		err = unmarshalPrimitiveOrObject("structBool", data, &valString, valStruct)
		assert.Error(t, err)
		assert.Equal(t, "structBool has a syntax error \"invalid character '\\\\\\\\' looking for beginning of object key string\"", err.Error())
		// assert.Equal(t, `structBool value '{"fieldValue": value}' is not supported, it has a syntax error "invalid character 'v' looking for beginning of value"`, err.Error())
	})

	t.Run("test personalized unmarshalTypeError error message", func(t *testing.T) {
		type structBool struct {
			FieldValue bool `json:"fieldValue"`
		}

		var valBool bool
		valStruct := &structBool{}
		data := []byte(`{
  "fieldValue": "true"
}`)
		err := unmarshalPrimitiveOrObject("structBool", data, &valBool, valStruct)
		assert.Error(t, err)
		assert.Equal(t, "structBool.fieldValue must be bool", err.Error())

		valBool = false
		valStruct = &structBool{}
		data = []byte(`"true"`)
		err = unmarshalPrimitiveOrObject("structBool", data, &valBool, valStruct)
		assert.Error(t, err)
		assert.Equal(t, "structBool must be bool or object", err.Error())
	})

	t.Run("check json with spaces", func(t *testing.T) {
		var valBool bool
		valStruct := &dataMap{}
		data := []byte(` {"value": "true"} `)
		err := unmarshalPrimitiveOrObject("dataMap", data, &valBool, valStruct)
		assert.NoError(t, err)

		valBool = false
		valStruct = &dataMap{}
		data = []byte(` true `)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valStruct)
		assert.NoError(t, err)

		valString := ""
		valStruct = &dataMap{}
		data = []byte(` "true" `)
		err = unmarshalPrimitiveOrObject("dataMap", data, &valString, valStruct)
		assert.NoError(t, err)
	})

	t.Run("check tabs", func(t *testing.T) {
		valString := ""
		valStruct := &dataMap{}
		data := []byte(string('\t') + `"true"` + string('\t'))
		err := unmarshalPrimitiveOrObject("dataMap", data, &valString, valStruct)
		assert.NoError(t, err)

		valBool := false
		valStruct = &dataMap{}
		data = []byte(string('\t') + `true` + string('\t'))
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valStruct)
		assert.NoError(t, err)
	})

	t.Run("check breakline", func(t *testing.T) {
		valString := ""
		valStruct := &dataMap{}
		data := []byte(string('\n') + `"true"` + string('\n'))
		err := unmarshalPrimitiveOrObject("dataMap", data, &valString, valStruct)
		assert.NoError(t, err)

		valBool := false
		valStruct = &dataMap{}
		data = []byte(string('\n') + `true` + string('\n'))
		err = unmarshalPrimitiveOrObject("dataMap", data, &valBool, valStruct)
		assert.NoError(t, err)
	})

	t.Run("test recursivity and default value", func(t *testing.T) {
		valStruct := &structBool{}
		data := []byte(`{"fieldValue": false}`)
		err := json.Unmarshal(data, valStruct)
		assert.NoError(t, err)
		assert.False(t, valStruct.FieldValue)
	})
}

type structBool struct {
	FieldValue bool `json:"fieldValue"`
}

func (s *structBool) UnmarshalJSON(data []byte) error {
	s.FieldValue = true
	return unmarshalObject("unmarshalJSON", data, s)
}
