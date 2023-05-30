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

	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(validationWrap(switchStateStructLevelValidation), SwitchState{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(defaultConditionStructLevelValidation), DefaultCondition{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(eventConditionStructLevelValidationCtx), EventCondition{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(dataConditionStructLevelValidation), DataCondition{})
}

// SwitchStateStructLevelValidation custom validator for SwitchState
func switchStateStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
	switchState := structLevel.Current().Interface().(SwitchState)

	switch {
	case len(switchState.DataConditions) == 0 && len(switchState.EventConditions) == 0:
		structLevel.ReportError(reflect.ValueOf(switchState), "DataConditions", "dataConditions", TagRequired, "")
	case len(switchState.DataConditions) > 0 && len(switchState.EventConditions) > 0:
		structLevel.ReportError(reflect.ValueOf(switchState), "DataConditions", "dataConditions", TagExclusive, "")
	}

}

// DefaultConditionStructLevelValidation custom validator for DefaultCondition
func defaultConditionStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
	defaultCondition := structLevel.Current().Interface().(DefaultCondition)
	validTransitionAndEnd(structLevel, defaultCondition, defaultCondition.Transition, defaultCondition.End)
}

// EventConditionStructLevelValidation custom validator for EventCondition
func eventConditionStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	eventCondition := structLevel.Current().Interface().(EventCondition)
	validTransitionAndEnd(structLevel, eventCondition, eventCondition.Transition, eventCondition.End)

	if eventCondition.EventRef != "" && !ctx.MapEvents.contain(eventCondition.EventRef) {
		structLevel.ReportError(eventCondition, "eventRef", "EventRef", TagExists, "")
	}
}

// DataConditionStructLevelValidation custom validator for DataCondition
func dataConditionStructLevelValidation(ctx ValidatorContext, structLevel validator.StructLevel) {
	dataCondition := structLevel.Current().Interface().(DataCondition)
	validTransitionAndEnd(structLevel, dataCondition, dataCondition.Transition, dataCondition.End)
}
