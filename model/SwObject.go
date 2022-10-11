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

type SwObject struct {
	Object
}

type Object interface {
	DeepCopyObject() Object
}

type custom struct {
	Object interface{}
}

func (o custom) DeepCopyObject() Object {
	return o
}

type String string
type Integer int

func (m Integer) DeepCopyObject() Object {
	return m
}

func (m String) DeepCopyObject() Object {
	return m
}

func (sw SwObject) MarshalJSON() ([]byte, error) {
	switch val := sw.Object.(type) {
	case String:
		return []byte(fmt.Sprintf(`%q`, val)), nil
	case Integer:
		return []byte(fmt.Sprintf(`%d`, val)), nil
	case custom:
		custom, err := json.Marshal(&struct {
			custom
		}{
			val,
		})
		if err != nil {
			return nil, err
		}

		// remove the field name and the last '}' for marshalling purposes
		st := strings.Replace(string(custom), "{\"Object\":", "", 1)
		st = strings.TrimSuffix(st, "}")
		return []byte(st), nil
	default:
		return []byte(fmt.Sprintf("%+v", sw.Object)), nil
	}
}

func (sw *SwObject) UnmarshalJSON(data []byte) error {

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
		sw.Object = strVal
		return nil

	case map[string]interface{}:
		var cstVal custom
		if err := json.Unmarshal(data, &cstVal.Object); err != nil {
			return err
		}
		sw.Object = cstVal
		return nil

	default:
		// json parses all not typed numbers as float64, let's enforce to int32
		if valInt, perr := strconv.Atoi(fmt.Sprint(val)); perr != nil {
			return fmt.Errorf("falied to parse %d to int32: %s", valInt, perr.Error())
		} else {
			var intVal Integer
			if err := json.Unmarshal(data, &intVal); err != nil {
				return err
			}
			sw.Object = intVal
			return nil
		}
	}
}
