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

func Test_unmarshal(t *testing.T) {
	testCases := []struct {
		name   string
		json   string
		object Object
		any    any
		err    string
	}{
		{
			name:   "string",
			json:   "\"value\"",
			object: FromString("value"),
			any:    any("value"),
		},
		{
			name:   "int",
			json:   "123",
			object: FromInt(123),
			any:    any(int32(123)),
		},
		{
			name:   "float",
			json:   "123.123",
			object: FromFloat(123.123),
			any:    any(123.123),
		},
		{
			name:   "map",
			json:   "{\"key\": \"value\", \"key2\": 123}",
			object: FromMap(map[string]any{"key": "value", "key2": 123}),
			any:    any(map[string]any{"key": "value", "key2": int32(123)}),
		},
		{
			name:   "slice",
			json:   "[\"key\", 123]",
			object: FromSlice([]any{"key", 123}),
			any:    any([]any{"key", int32(123)}),
		},
		{
			name:   "bool true",
			json:   "true",
			object: FromBool(true),
			any:    any(true),
		},
		{
			name:   "bool false",
			json:   "false",
			object: FromBool(false),
			any:    any(false),
		},
		{
			name:   "null",
			json:   "null",
			object: FromNull(),
			any:    nil,
		},
		{
			name: "string invalid",
			json: "\"invalid",
			err:  "unexpected end of JSON input",
		},
		{
			name: "number invalid",
			json: "123a",
			err:  "invalid character 'a' after top-level value",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			o := Object{}
			err := json.Unmarshal([]byte(tc.json), &o)
			if tc.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.object, o)
				assert.Equal(t, ToInterface(tc.object), tc.any)
			} else {
				assert.Equal(t, tc.err, err.Error())
			}
		})
	}
}

func Test_marshal(t *testing.T) {
	testCases := []struct {
		name   string
		json   string
		object Object
		err    string
	}{
		{
			name:   "string",
			json:   "\"value\"",
			object: FromString("value"),
		},
		{
			name:   "int",
			json:   "123",
			object: FromInt(123),
		},
		{
			name:   "float",
			json:   "123.123000",
			object: FromFloat(123.123),
		},
		{
			name:   "map",
			json:   "{\"key\":\"value\",\"key2\":123}",
			object: FromMap(map[string]any{"key": "value", "key2": 123}),
		},
		{
			name:   "slice",
			json:   "[\"key\",123]",
			object: FromSlice([]any{"key", 123}),
		},
		{
			name:   "bool true",
			json:   "true",
			object: FromBool(true),
		},
		{
			name:   "bool false",
			json:   "false",
			object: FromBool(false),
		},
		{
			name:   "null",
			json:   "null",
			object: FromNull(),
		},
		{
			name: "interface",
			json: "[\"value\",123,123.123000,[1],{\"key\":1.100000},true,false,null]",
			object: FromInterface([]any{
				"value",
				123,
				123.123,
				[]any{1},
				map[string]any{"key": 1.1},
				true,
				false,
				nil,
			}),
		},
		{
			name:   "fromraw rawmessage slice",
			json:   "[\"x\",2,{\"k\":true},null]",
			object: FromRaw(json.RawMessage([]byte(`["x",2,{"k":true},null]`))),
		},
		{
			name:   "fromraw primitive string",
			json:   "\"hello\"",
			object: FromRaw("hello"),
		},
		{
			name:   "fromraw primitive int",
			json:   "42",
			object: FromRaw(42),
		},
		{
			name:   "fromraw slice and map",
			json:   "[\"x\",2,{\"k\":true},null]",
			object: FromRaw([]any{"x", 2, map[string]any{"k": true}, nil}),
		},
		{
			name:   "fromraw map",
			json:   "{\"a\":1,\"b\":[true,false],\"c\":null}",
			object: FromRaw(map[string]any{"a": 1, "b": []any{true, false}, "c": nil}),
		},
		{
			name:   "fromraw null",
			json:   "null",
			object: FromRaw(nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			json, err := json.Marshal(tc.object)
			if tc.err == "" {
				assert.NoError(t, err)
				assert.Equal(t, tc.json, string(json))
			} else {
				assert.Equal(t, tc.err, err.Error())
			}
		})
	}
}
