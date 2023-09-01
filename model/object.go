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
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"strconv"
)

type Type int8

const (
	Null Type = iota
	String
	Int
	Float
	Map
	Slice
	Bool
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
//
// +kubebuilder:validation:Type=object
type Object struct {
	Type        Type   `json:"type,inline"`
	StringValue string `json:"strVal,inline"`
	IntValue    int32  `json:"intVal,inline"`
	FloatValue  float64
	MapValue    map[string]Object
	SliceValue  []Object
	BoolValue   bool `json:"boolValue,inline"`
}

// UnmarshalJSON implements json.Unmarshaler
func (obj *Object) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)

	if data[0] == '"' {
		obj.Type = String
		return json.Unmarshal(data, &obj.StringValue)
	} else if data[0] == 't' || data[0] == 'f' {
		obj.Type = Bool
		return json.Unmarshal(data, &obj.BoolValue)
	} else if data[0] == 'n' {
		obj.Type = Null
		return nil
	} else if data[0] == '{' {
		obj.Type = Map
		return json.Unmarshal(data, &obj.MapValue)
	} else if data[0] == '[' {
		obj.Type = Slice
		return json.Unmarshal(data, &obj.SliceValue)
	}

	number := string(data)
	intValue, err := strconv.ParseInt(number, 10, 32)
	if err == nil {
		obj.Type = Int
		obj.IntValue = int32(intValue)
		return nil
	}

	floatValue, err := strconv.ParseFloat(number, 64)
	if err == nil {
		obj.Type = Float
		obj.FloatValue = floatValue
		return nil
	}

	return fmt.Errorf("json invalid number %q", number)
}

// MarshalJSON marshal the given json object into the respective Object subtype.
func (obj Object) MarshalJSON() ([]byte, error) {
	switch obj.Type {
	case String:
		return []byte(fmt.Sprintf(`%q`, obj.StringValue)), nil
	case Int:
		return []byte(fmt.Sprintf(`%d`, obj.IntValue)), nil
	case Float:
		return []byte(fmt.Sprintf(`%f`, obj.FloatValue)), nil
	case Map:
		return json.Marshal(obj.MapValue)
	case Slice:
		return json.Marshal(obj.SliceValue)
	case Bool:
		return []byte(fmt.Sprintf(`%t`, obj.BoolValue)), nil
	case Null:
		return []byte("null"), nil
	default:
		panic("object invalid type")
	}
}

func FromString(val string) Object {
	return Object{Type: String, StringValue: val}
}

func FromInt(val int) Object {
	if val > math.MaxInt32 || val < math.MinInt32 {
		fmt.Println(fmt.Errorf("value: %d overflows int32", val))
	}
	return Object{Type: Int, IntValue: int32(val)}
}

func FromFloat(val float64) Object {
	if val > math.MaxFloat64 || val < -math.MaxFloat64 {
		fmt.Println(fmt.Errorf("value: %f overflows float64", val))
	}
	return Object{Type: Float, FloatValue: float64(val)}
}

func FromMap(mapValue map[string]any) Object {
	mapValueObject := make(map[string]Object, len(mapValue))
	for key, value := range mapValue {
		mapValueObject[key] = FromInterface(value)
	}
	return Object{Type: Map, MapValue: mapValueObject}
}

func FromSlice(sliceValue []any) Object {
	sliceValueObject := make([]Object, len(sliceValue))
	for key, value := range sliceValue {
		sliceValueObject[key] = FromInterface(value)
	}
	return Object{Type: Slice, SliceValue: sliceValueObject}
}

func FromBool(val bool) Object {
	return Object{Type: Bool, BoolValue: val}
}

func FromNull() Object {
	return Object{Type: Null}
}

func FromInterface(value any) Object {
	switch v := value.(type) {
	case string:
		return FromString(v)
	case int:
		return FromInt(v)
	case int32:
		return FromInt(int(v))
	case float64:
		return FromFloat(v)
	case map[string]any:
		return FromMap(v)
	case []any:
		return FromSlice(v)
	case bool:
		return FromBool(v)
	case nil:
		return FromNull()
	}
	panic("invalid type")
}

func ToInterface(object Object) any {
	switch object.Type {
	case String:
		return object.StringValue
	case Int:
		return object.IntValue
	case Float:
		return object.FloatValue
	case Map:
		mapInterface := make(map[string]any, len(object.MapValue))
		for key, value := range object.MapValue {
			mapInterface[key] = ToInterface(value)
		}
		return mapInterface
	case Slice:
		sliceInterface := make([]any, len(object.SliceValue))
		for key, value := range object.SliceValue {
			sliceInterface[key] = ToInterface(value)
		}
		return sliceInterface
	case Bool:
		return object.BoolValue
	case Null:
		return nil
	}
	panic("invalid type")
}
