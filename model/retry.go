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

package model

import (
	"reflect"

	"github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidation(
		RetryStructLevelValidation,
		Retry{},
	)
}

// Retry ...
type Retry struct {
	// Unique retry strategy name
	Name string `json:"name" validate:"required"`
	// Time delay between retry attempts (ISO 8601 duration format)
	Delay string `json:"delay,omitempty"`
	// Maximum time delay between retry attempts (ISO 8601 duration format)
	MaxDelay string `json:"maxDelay,omitempty"`
	// Static value by which the delay increases during each attempt (ISO 8601 time format)
	Increment string `json:"increment,omitempty"`
	// Numeric value, if specified the delay between retries is multiplied by this value.
	Multiplier *floatstr.Float32OrString `json:"multiplier,omitempty" validate:"omitempty,min=1"`
	// Maximum number of retry attempts.
	MaxAttempts intstr.IntOrString `json:"maxAttempts" validate:"required"`
	// If float type, maximum amount of random time added or subtracted from the delay between each retry relative to total delay (between 0 and 1). If string type, absolute maximum amount of random time added or subtracted from the delay between each retry (ISO 8601 duration format)
	Jitter floatstr.Float32OrString `json:"jitter,omitempty" validate:"omitempty,min=0,max=1"`
}

// RetryStructLevelValidation custom validator for Retry Struct
func RetryStructLevelValidation(structLevel validator.StructLevel) {
	retryObj := structLevel.Current().Interface().(Retry)

	if retryObj.Delay != "" {
		err := validateISO8601TimeDuration(retryObj.Delay)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(retryObj.Delay), "Delay", "delay", "reqiso8601duration", "")
		}
	}

	if retryObj.MaxDelay != "" {
		err := validateISO8601TimeDuration(retryObj.MaxDelay)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(retryObj.MaxDelay), "MaxDelay", "maxDelay", "reqiso8601duration", "")
		}
	}

	if retryObj.Increment != "" {
		err := validateISO8601TimeDuration(retryObj.Increment)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(retryObj.Increment), "Increment", "increment", "reqiso8601duration", "")
		}
	}

	if retryObj.Jitter.Type == floatstr.String && retryObj.Jitter.StrVal != "" {
		err := validateISO8601TimeDuration(retryObj.Jitter.StrVal)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(retryObj.Jitter.StrVal), "Jitter", "jitter", "reqiso8601duration", "")
		}
	}
}
