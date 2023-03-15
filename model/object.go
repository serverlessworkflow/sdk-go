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
	"math"
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
//
// +kubebuilder:validation:Type=object
type Object struct {
	Type     Type            `json:"type,inline"`
	IntVal   int32           `json:"intVal,inline"`
	StrVal   string          `json:"strVal,inline"`
	RawValue json.RawMessage `json:"rawValue,inline"`
}

type Type int64

const (
	Integer Type = iota
	String
	Raw
)

func FromInt(val int) Object {
	if val > math.MaxInt32 || val < math.MinInt32 {
		fmt.Println(fmt.Errorf("value: %d overflows int32", val))
	}
	return Object{Type: Integer, IntVal: int32(val)}
}

func FromString(val string) Object {
	return Object{Type: String, StrVal: val}
}

func FromRaw(val interface{}) Object {
	custom, err := json.Marshal(val)
	if err != nil {
		er := fmt.Errorf("failed to parse value to Raw: %w", err)
		fmt.Println(er.Error())
		return Object{}
	}
	return Object{Type: Raw, RawValue: custom}
}

// UnmarshalJSON ...
func (obj *Object) UnmarshalJSON(data []byte) error {
	if data[0] == '"' {
		obj.Type = String
		return json.Unmarshal(data, &obj.StrVal)
	} else if data[0] == '{' {
		obj.Type = Raw
		return json.Unmarshal(data, &obj.RawValue)
	}
	obj.Type = Integer
	return json.Unmarshal(data, &obj.IntVal)
}

// MarshalJSON marshal the given json object into the respective Object subtype.
func (obj Object) MarshalJSON() ([]byte, error) {
	switch obj.Type {
	case String:
		return []byte(fmt.Sprintf(`%q`, obj.StrVal)), nil
	case Integer:
		return []byte(fmt.Sprintf(`%d`, obj.IntVal)), nil
	case Raw:
		val, _ := json.Marshal(obj.RawValue)
		return val, nil
	default:
		return []byte(fmt.Sprintf("%+v", obj)), nil
	}
}
