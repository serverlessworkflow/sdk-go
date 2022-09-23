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
	"context"
	"reflect"

	validator "github.com/go-playground/validator/v10"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

// SleepState suspends workflow execution for a given time duration.
type SleepState struct {
	BaseState

	// Duration (ISO 8601 duration format) to sleep
	Duration string `json:"duration" validate:"required"`
	// Timeouts State specific timeouts
	Timeouts *SleepStateTimeout `json:"timeouts,omitempty"`
}

// SleepStateTimeout defines timeout settings for sleep state
type SleepStateTimeout struct {
	StateExecTimeout StateExecTimeout `json:"stateExecTimeout,omitempty"`
}

// SleepStateStructLevelValidation custom validator for SleepState
func SleepStateStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	sleepState := structLevel.Current().Interface().(SleepState)

	err := validateISO8601TimeDuration(sleepState.Duration)
	if err != nil {
		structLevel.ReportError(reflect.ValueOf(sleepState.Duration), "Duration", "duration", "reqiso8601duration", "")
	}
}

func init() {
	val.GetValidator().RegisterStructValidationCtx(
		SleepStateStructLevelValidation,
		SleepState{},
	)
}
