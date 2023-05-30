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
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/intstr"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(parallelStateStructLevelValidation, ParallelState{})
}

// ParallelStateStructLevelValidation custom validator for ParallelState
func parallelStateStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	parallelStateObj := structLevel.Current().Interface().(ParallelState)

	if parallelStateObj.CompletionType == CompletionTypeAllOf {
		return
	}

	switch parallelStateObj.NumCompleted.Type {
	case intstr.Int:
		if parallelStateObj.NumCompleted.IntVal <= 0 {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", "")
		}
	case intstr.String:
		v, err := strconv.Atoi(parallelStateObj.NumCompleted.StrVal)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", err.Error())
			return
		}

		if v <= 0 {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", "")
		}
	}
}
