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
	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(validationWrap(actionStructLevelValidationCtx), Action{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(functionRefStructLevelValidation), FunctionRef{})
}

func actionStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	action := structLevel.Current().Interface().(Action)
	values := []bool{
		action.FunctionRef != nil,
		action.EventRef != nil,
		action.SubFlowRef != nil,
	}

	if validationNotExclusiveParamters(values) {
		structLevel.ReportError(action.FunctionRef, "functionRef", "FunctionRef", TagExclusive, "")
		structLevel.ReportError(action.EventRef, "EventRef", "eventRef", TagExclusive, "")
		structLevel.ReportError(action.SubFlowRef, "SubFlowRef", "subFlowRef", TagExclusive, "")
	}

	if action.RetryRef != "" && !ctx.MapRetries.contain(action.RetryRef) {
		structLevel.ReportError(action.RetryRef, "retryRef", "RetryRef", TagExists, "")
	}
}

func functionRefStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
	functionRef := structLevel.Current().Interface().(FunctionRef)
	if !ctx.MapFunctions.contain(functionRef.RefName) {
		structLevel.ReportError(functionRef.RefName, "refName", "RefName", TagExists, functionRef.RefName)
	}
}
