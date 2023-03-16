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
	"context"
	"reflect"

	validator "github.com/go-playground/validator/v10"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

type workflowValidator func(mapValues ValidatorContextValue, sl validator.StructLevel)

type ValidatorContextValue struct {
	MapStates    mapValues[State]
	MapFunctions mapValues[Function]
	MapEvents    mapValues[Event]
	MapRetries   mapValues[Retry]
	MapErrors    mapValues[Error]
}

func validationWrap(fn1 validator.StructLevelFunc, fnCtx workflowValidator) validator.StructLevelFuncCtx {
	return func(ctx context.Context, structLevel validator.StructLevel) {
		if fn1 != nil {
			fn1(structLevel)
		}

		if fnCtx != nil {
			if mapValues, ok := ctx.Value("values").(ValidatorContextValue); ok {
				fnCtx(mapValues, structLevel)
			}
		}
	}
}

func NewValidatorContext(workflow *Workflow) context.Context {
	contextValue := ValidatorContextValue{
		MapStates:    newMapValues(workflow.States, "Name"),
		MapFunctions: newMapValues(workflow.Functions, "Name"),
		MapEvents:    newMapValues(workflow.Events, "Name"),
		MapRetries:   newMapValues(workflow.Retries, "Name"),
		MapErrors:    newMapValues(workflow.Errors, "Name"),
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, "values", contextValue)
	return ctx
}

func newMapValues[T any](values []T, field string) mapValues[T] {
	c := mapValues[T]{}
	c.init(values, field)
	return c
}

type mapValues[T any] struct {
	ValuesMap map[string]T
}

func (c *mapValues[T]) init(values []T, field string) {
	c.ValuesMap = make(map[string]T, len(values))
	for _, v := range values {
		name := reflect.ValueOf(v).FieldByName(field).String()
		c.ValuesMap[name] = v
	}
}

func (c *mapValues[T]) contain(name string) bool {
	_, ok := c.ValuesMap[name]
	return ok
}

func init() {
	val.GetValidator().RegisterStructValidationCtx(validationWrap(nil, workflowStructLevelValidation), Workflow{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(onErrorStructLevelValidation, onErrorStructLevelValidationCtx), OnError{})

}

func workflowStructLevelValidation(ctx ValidatorContextValue, structLevel validator.StructLevel) {
	workflow := structLevel.Current().Interface().(Workflow)

	if workflow.Start != nil {
		// if not exists the start transtion stop the states validations
		if !ctx.MapStates.contain(workflow.Start.StateName) {
			structLevel.ReportError(reflect.ValueOf(workflow.Start), "Start", "start", "startnotexist", "")
			return
		}
	}

	if len(workflow.States) == 1 {
		return
	}

	// Naive check if transitions exist
	for _, state := range ctx.MapStates.ValuesMap {
		if state.Transition != nil {
			if !ctx.MapStates.contain(state.Transition.NextState) {
				structLevel.ReportError(reflect.ValueOf(state), "Transition", "transition", "transitionnotexists", state.Transition.NextState)
			}
		}
	}

	// TODO: create states graph to complex check
}

func onErrorStructLevelValidation(structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)

	hasErrorRef := onError.ErrorRef != ""
	hasErrorRefs := len(onError.ErrorRefs) > 0

	if !hasErrorRef && !hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", "required", "")
	} else if hasErrorRef && hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", "exclusive", "")
	}
}

func onErrorStructLevelValidationCtx(ctx ValidatorContextValue, structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)

	if onError.ErrorRef != "" && !ctx.MapErrors.contain(onError.ErrorRef) {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", "exists", "")
	}

	for _, errorRef := range onError.ErrorRefs {
		if !ctx.MapErrors.contain(errorRef) {
			structLevel.ReportError(onError.ErrorRefs, "ErrorRefs", "errorRefs", "exists", "")
		}
	}
}

func validTransitionAndEnd(structLevel validator.StructLevel, field any, transition *Transition, end *End) {
	hasTransition := transition != nil
	isEnd := end != nil && (end.Terminate || end.ContinueAs != nil || len(end.ProduceEvents) > 0) // TODO: check the spec continueAs/produceEvents to see how it influences the end

	if !hasTransition && !isEnd {
		structLevel.ReportError(field, "Transition", "transition", "required", "")
	} else if hasTransition && isEnd {
		structLevel.ReportError(field, "Transition", "transition", "exclusive", "")
	}
}

func validationNotExclusiveParamters(values []bool) bool {
	hasOne := false
	hasTwo := false

	for i, val1 := range values {
		if val1 {
			hasOne = true
			for j, val2 := range values {
				if i != j && val2 {
					hasTwo = true
					break
				}
			}
			break
		}
	}
	return (hasOne && hasTwo) || (!hasOne && !hasTwo)
}
