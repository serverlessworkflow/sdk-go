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
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	validator "github.com/go-playground/validator/v10"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

// StateExecTimeout defines workflow state execution timeout
type StateExecTimeout struct {
	// Single state execution timeout, not including retries (ISO 8601 duration format)
	Single string `json:"single,omitempty"`
	// Total state execution timeout, including retries (ISO 8601 duration format)
	Total string `json:"total" validate:"required"`
}

// just define another type to unmarshal object, so the UnmarshalJSON will not called recursively
type stateExecTimeoutForUnmarshal StateExecTimeout

// UnmarshalJSON unmarshal StateExecTimeout object from json bytes
func (s *StateExecTimeout) UnmarshalJSON(data []byte) error {
	// We must trim the leading space, because we use first byte to detect data's type
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		// TODO: Normalize error messages
		return fmt.Errorf("no bytes to unmarshal")
	}

	var err error
	switch data[0] {
	case '"':
		s.Total, err = unmarshalString(data)
		return err
	case '{':
		var v stateExecTimeoutForUnmarshal
		err = json.Unmarshal(data, &v)
		if err != nil {
			// TODO: replace the error message with correct type's name
			return err
		}

		*s = StateExecTimeout(v)
		return nil
	}

	return fmt.Errorf("stateExecTimeout value '%s' not support, it must be an object or string", string(data))
}

// StateExecTimeoutStructLevelValidation custom validator for StateExecTimeout
func StateExecTimeoutStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	timeoutObj := structLevel.Current().Interface().(StateExecTimeout)

	// TODO: use Custom Validation Functions tags for iso8601duration
	err := validateISO8601TimeDuration(timeoutObj.Total)
	if err != nil {
		structLevel.ReportError(reflect.ValueOf(timeoutObj.Total), "Total", "total", "reqiso8601duration", "")
	}

	if timeoutObj.Single != "" {
		err = validateISO8601TimeDuration(timeoutObj.Single)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(timeoutObj.Single), "Single", "single", "reqiso8601duration", "")
		}
	}
}

func init() {
	val.GetValidator().RegisterStructValidationCtx(
		StateExecTimeoutStructLevelValidation,
		StateExecTimeout{},
	)
}
