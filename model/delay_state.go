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
	"reflect"

	"github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidation(
		DelayStateStructLevelValidation,
		DelayState{},
	)
}

// DelayState Causes the workflow execution to delay for a specified duration
type DelayState struct {
	BaseState
	// Amount of time (ISO 8601 format) to delay
	TimeDelay string `json:"timeDelay" validate:"required"`
}

// DelayStateStructLevelValidation custom validator for DelayState Struct
func DelayStateStructLevelValidation(structLevel validator.StructLevel) {
	delayStateObj := structLevel.Current().Interface().(DelayState)

	err := validateISO8601TimeDuration(delayStateObj.TimeDelay)
	if err != nil {
		structLevel.ReportError(reflect.ValueOf(delayStateObj.TimeDelay), "TimeDelay", "timeDelay", "reqiso8601duration", "")
	}
}
