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

	validator "github.com/go-playground/validator/v10"
	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidation(retryStructLevelValidation, Retry{})
}

// RetryStructLevelValidation custom validator for Retry Struct
func retryStructLevelValidation(structLevel validator.StructLevel) {
	retryObj := structLevel.Current().Interface().(Retry)

	if retryObj.Jitter.Type == floatstr.String && retryObj.Jitter.StrVal != "" {
		err := val.ValidateISO8601TimeDuration(retryObj.Jitter.StrVal)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(retryObj.Jitter.StrVal), "Jitter", "jitter", "iso8601duration", "")
		}
	}
}
