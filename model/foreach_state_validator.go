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

func init() {
	val.GetValidator().RegisterStructValidationCtx(forEachStateStructLevelValidation, ForEachState{})
}

// ForEachStateStructLevelValidation custom validator for ForEachState
func forEachStateStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	stateObj := structLevel.Current().Interface().(ForEachState)

	if stateObj.Mode != ForEachModeTypeParallel {
		return
	}

	if stateObj.BatchSize == nil {
		return
	}

	if !val.ValidateGt0IntStr(stateObj.BatchSize) {
		structLevel.ReportError(reflect.ValueOf(stateObj.BatchSize), "BatchSize", "batchSize", "gt0", "")
	}
}
