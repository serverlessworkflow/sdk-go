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
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(eventStateStructLevelValidationCtx), EventState{})
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(onEventsStructLevelValidationCtx), OnEvents{})
}

func eventStateStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	// EventRefs
}

func onEventsStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	onEvent := structLevel.Current().Interface().(OnEvents)
	for _, eventRef := range onEvent.EventRefs {
		if eventRef != "" && !ctx.ExistEvent(eventRef) {
			structLevel.ReportError(eventRef, "eventRefs", "EventRefs", val.TagExists, "")
		}
	}
}
