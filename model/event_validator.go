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
	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(eventStructLevelValidation), Event{})
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(eventRefStructLevelValidation), EventRef{})
}

// eventStructLevelValidation custom validator for event kind consumed
func eventStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
}

func eventRefStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
	model := structLevel.Current().Interface().(EventRef)
	if model.TriggerEventRef != "" && !ctx.ExistEvent(model.TriggerEventRef) {
		structLevel.ReportError(model.TriggerEventRef, "triggerEventRef", "TriggerEventRef", val.TagExists, "")
	}
	if model.ResultEventRef != "" && !ctx.ExistEvent(model.ResultEventRef) {
		structLevel.ReportError(model.ResultEventRef, "triggerEventRef", "TriggerEventRef", val.TagExists, "")
	}
}
