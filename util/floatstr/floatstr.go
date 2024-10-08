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

package floatstr

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

// Float32OrString is a type that can hold a float32 or a string.
// implementation borrowed from apimachinary intstr package: https://github.com/kubernetes/apimachinery/blob/master/pkg/util/intstr/intstr.go
type Float32OrString struct {
	Type     Type    `json:"type,omitempty"`
	FloatVal float32 `json:"floatVal,omitempty"`
	StrVal   string  `json:"strVal,omitempty"`
}

// Type represents the stored type of Float32OrString.
type Type int64

const (
	// Float ...
	Float Type = iota // The Float32OrString holds a float.
	// String ...
	String // The Float32OrString holds a string.
)

// FromFloat creates an Float32OrString object with a float32 value. It is
// your responsibility not to call this method with a value greater
// than float32.
func FromFloat(val float32) Float32OrString {
	return Float32OrString{Type: Float, FloatVal: val}
}

// FromString creates a Float32OrString object with a string value.
func FromString(val string) Float32OrString {
	return Float32OrString{Type: String, StrVal: val}
}

// Parse the given string and try to convert it to a float32 before
// setting it as a string value.
func Parse(val string) Float32OrString {
	f, err := strconv.ParseFloat(val, 32)
	if err != nil {
		return FromString(val)
	}
	return FromFloat(float32(f))
}

// UnmarshalJSON implements the json.Unmarshaller interface.
func (floatstr *Float32OrString) UnmarshalJSON(value []byte) error {
	if value[0] == '"' {
		floatstr.Type = String
		return json.Unmarshal(value, &floatstr.StrVal)
	}
	floatstr.Type = Float
	return json.Unmarshal(value, &floatstr.FloatVal)
}

// MarshalJSON implements the json.Marshaller interface.
func (floatstr *Float32OrString) MarshalJSON() ([]byte, error) {
	switch floatstr.Type {
	case Float:
		return json.Marshal(floatstr.FloatVal)
	case String:
		return json.Marshal(floatstr.StrVal)
	default:
		return []byte{}, fmt.Errorf("impossible Float32OrString.Type")
	}
}

// String returns the string value, or the float value.
func (floatstr *Float32OrString) String() string {
	if floatstr == nil {
		return "<nil>"
	}
	if floatstr.Type == String {
		return floatstr.StrVal
	}
	return strconv.FormatFloat(float64(floatstr.FloatValue()), 'E', -1, 32)
}

// FloatValue returns the FloatVal if type float32, or if
// it is a String, will attempt a conversion to float32,
// returning 0 if a parsing error occurs.
func (floatstr *Float32OrString) FloatValue() float32 {
	if floatstr.Type == String {
		f, _ := strconv.ParseFloat(floatstr.StrVal, 32)
		return float32(f)
	}
	return floatstr.FloatVal
}

func init() {
	val.GetValidator().RegisterCustomTypeFunc(func(fl reflect.Value) interface{} {
		if fl.Kind() != reflect.Struct {
			return errors.New("invalid type: expected Float32OrString")
		}

		// Get the Float32OrString value
		_, ok := fl.Interface().(Float32OrString)
		if !ok {
			return fmt.Errorf("invalid type: expected Float32OrString")
		}

		return nil
	}, Float32OrString{})
}

func ValidateFloat32OrString(sl validator.StructLevel) {
	// Get the current struct being validated.
	current := sl.Current()

	for i := 0; i < current.NumField(); i++ {
		field := current.Type().Field(i)
		value := current.Field(i)

		// Check if the field is a pointer and handle nil pointers.
		if value.Kind() == reflect.Ptr {
			if value.IsNil() {
				continue // Skip nil pointers.
			}
			value = value.Elem() // Dereference the pointer.
		}

		// Check if the field is of type Float32OrString.
		if value.Type() == reflect.TypeOf(Float32OrString{}) {
			// Extract validation tags from the field.
			tags := field.Tag.Get("validate")

			// Split tags and look for min/max.
			tagList := strings.Split(tags, ",")
			for _, tag := range tagList {
				if strings.HasPrefix(tag, "min=") {
					minVal, err := strconv.ParseFloat(strings.TrimPrefix(tag, "min="), 32)
					if err != nil {
						sl.ReportError(value.Interface(), field.Name, field.Name, "min", "")
						continue
					}

					if value.FieldByName("FloatVal").Float() < minVal {
						sl.ReportError(value.Interface(), field.Name, field.Name, "min", "")
					}
				}

				if strings.HasPrefix(tag, "max=") {
					maxVal, err := strconv.ParseFloat(strings.TrimPrefix(tag, "max="), 32)
					if err != nil {
						sl.ReportError(value.Interface(), field.Name, field.Name, "max", "")
						continue
					}

					if value.FieldByName("FloatVal").Float() > maxVal {
						sl.ReportError(value.Interface(), field.Name, field.Name, "max", "")
					}
				}
			}
		}
	}
}
