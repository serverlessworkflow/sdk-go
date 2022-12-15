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
	"strconv"
	"strings"
)

// Object is used to allow integration with DeepCopy tool by replacing 'interface' generic type.
// The DeepCopy tool allow us to easily import the Workflow types into a Kubernetes operator,
// which requires the DeepCopy method.
//
// It can marshal and unmarshal any type.
// This object type can be three types:
//   - String	- holds string values
//   - Integer	- holds int32 values, JSON marshal any number to float64 by default, during the marshaling process it is
//     parsed to int32
//   - raw		- holds any not typed value, replaces the interface{} behavior.
type Object struct {
	IObject
}

// IObject interface that can converted into one of the three subtypes
type IObject interface {
	DeepCopyIObject() IObject
}

// raw generic subtype
type raw struct {
	IObject interface{}
}

func (o raw) DeepCopyIObject() IObject {
	return o
}

// Integer int32 type
type Integer int

func (m Integer) DeepCopyIObject() IObject {
	return m
}

// String string type
type String string

func (m String) DeepCopyIObject() IObject {
	return m
}

// MarshalJSON marshal the given json object into the respective Object subtype.
func (obj Object) MarshalJSON() ([]byte, error) {
	switch val := obj.IObject.(type) {
	case String:
		return []byte(fmt.Sprintf(`%q`, val)), nil
	case Integer:
		return []byte(fmt.Sprintf(`%d`, val)), nil
	case raw:
		custom, err := json.Marshal(&struct {
			raw
		}{
			val,
		})
		if err != nil {
			return nil, err
		}

		// remove the field name and the last '}' for marshalling purposes
		st := strings.Replace(string(custom), "{\"IObject\":", "", 1)
		st = strings.TrimSuffix(st, "}")
		return []byte(st), nil
	default:
		return []byte(fmt.Sprintf("%+v", obj.IObject)), nil
	}
}

// UnmarshalJSON ...
func (obj *Object) UnmarshalJSON(data []byte) error {
	var test interface{}
	if err := json.Unmarshal(data, &test); err != nil {
		return err
	}
	switch val := test.(type) {
	case string:
		var strVal String
		if err := json.Unmarshal(data, &strVal); err != nil {
			return err
		}
		obj.IObject = strVal
		return nil

	case map[string]interface{}:
		var cstVal raw
		if err := json.Unmarshal(data, &cstVal.IObject); err != nil {
			return err
		}
		obj.IObject = cstVal
		return nil

	default:
		// json parses all not typed numbers as float64, let's enforce to int32
		if valInt, parseErr := strconv.Atoi(fmt.Sprint(val)); parseErr != nil {
			return fmt.Errorf("falied to parse %d to int32: %s", valInt, parseErr.Error())
		} else {
			var intVal Integer
			if err := json.Unmarshal(data, &intVal); err != nil {
				return err
			}
			obj.IObject = intVal
			return nil
		}
	}
}
