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

type contextValueKey string

const validatorContextValue contextValueKey = "value"

const (
	TagExists    string = "exists"
	TagRequired  string = "required"
	TagExclusive string = "exclusive"

	TagRecursiveState string = "recursivestate"

	// States referenced by compensatedBy (as well as any other states that they transition to) must obey following rules:
	TagTransitionMainWorkflow       string = "transtionmainworkflow"         // They should not have any incoming transitions (should not be part of the main workflow control-flow logic)
	TagEventState                   string = "eventstate"                    // They cannot be an event state
	TagRecursiveCompensation        string = "recursivecompensation"         // They cannot themselves set their compensatedBy property to true (compensation is not recursive)
	TagCompensatedby                string = "compensatedby"                 // They must define the usedForCompensation property and set it to true
	TagTransitionUseForCompensation string = "transitionusedforcompensation" // They can transition only to states which also have their usedForCompensation property and set to true
)

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
			if mapValues, ok := ctx.Value(validatorContextValue).(ValidatorContextValue); ok {
				fnCtx(mapValues, structLevel)
			}
		}
	}
}

func NewValidatorContext(workflow *Workflow) context.Context {
	for i := range workflow.States {
		s := &workflow.States[i]
		if s.BaseState.Transition != nil {
			s.BaseState.Transition.stateParent = s
		}
		for _, onError := range s.BaseState.OnErrors {
			if onError.Transition != nil {
				onError.Transition.stateParent = s
			}
		}
		if s.Type == StateTypeSwitch {
			if s.SwitchState.DefaultCondition.Transition != nil {
				s.SwitchState.DefaultCondition.Transition.stateParent = s
			}
			for _, e := range s.SwitchState.EventConditions {
				if e.Transition != nil {
					e.Transition.stateParent = s
				}
			}
			for _, d := range s.SwitchState.DataConditions {
				if d.Transition != nil {
					d.Transition.stateParent = s
				}
			}
		}
	}

	contextValue := ValidatorContextValue{
		MapStates:    newMapValues(workflow.States, "Name"),
		MapFunctions: newMapValues(workflow.Functions, "Name"),
		MapEvents:    newMapValues(workflow.Events, "Name"),
		MapRetries:   newMapValues(workflow.Retries, "Name"),
		MapErrors:    newMapValues(workflow.Errors, "Name"),
	}

	return context.WithValue(context.Background(), validatorContextValue, contextValue)
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
	// TODO: create states graph to complex check

	// val.GetValidator().RegisterStructValidationCtx(validationWrap(nil, workflowStructLevelValidation), Workflow{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(onErrorStructLevelValidation, onErrorStructLevelValidationCtx), OnError{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(nil, transitionStructLevelValidationCtx), Transition{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(nil, startStructLevelValidationCtx), Start{})
}

func startStructLevelValidationCtx(ctx ValidatorContextValue, structLevel validator.StructLevel) {
	start := structLevel.Current().Interface().(Start)
	if !ctx.MapStates.contain(start.StateName) {
		structLevel.ReportError(start.StateName, "StateName", "stateName", TagExists, "")
		return
	}
}

func onErrorStructLevelValidation(structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)

	hasErrorRef := onError.ErrorRef != ""
	hasErrorRefs := len(onError.ErrorRefs) > 0

	if !hasErrorRef && !hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagRequired, "")
	} else if hasErrorRef && hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagExclusive, "")
	}
}

func onErrorStructLevelValidationCtx(ctx ValidatorContextValue, structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)

	if onError.ErrorRef != "" && !ctx.MapErrors.contain(onError.ErrorRef) {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagExists, "")
	}

	for _, errorRef := range onError.ErrorRefs {
		if !ctx.MapErrors.contain(errorRef) {
			structLevel.ReportError(onError.ErrorRefs, "ErrorRefs", "errorRefs", TagExists, "")
		}
	}
}

func transitionStructLevelValidationCtx(ctx ValidatorContextValue, structLevel validator.StructLevel) {
	// Naive check if transitions exist
	transition := structLevel.Current().Interface().(Transition)
	if ctx.MapStates.contain(transition.NextState) {
		if transition.stateParent != nil {
			parentBaseState := transition.stateParent

			if parentBaseState.Name == transition.NextState {
				structLevel.ReportError(transition.NextState, "NextState", "nextState", TagRecursiveState, "")
			}

			if parentBaseState.UsedForCompensation && !ctx.MapStates.ValuesMap[transition.NextState].UsedForCompensation {
				structLevel.ReportError(transition.NextState, "NextState", "nextState", TagTransitionUseForCompensation, "")

			} else if !parentBaseState.UsedForCompensation && ctx.MapStates.ValuesMap[transition.NextState].UsedForCompensation {
				structLevel.ReportError(transition.NextState, "NextState", "nextState", TagTransitionMainWorkflow, "")
			}
		}

	} else {
		structLevel.ReportError(transition.NextState, "NextState", "nextState", TagExists, "")
	}
}

func validTransitionAndEnd(structLevel validator.StructLevel, field any, transition *Transition, end *End) {
	hasTransition := transition != nil
	isEnd := end != nil && (end.Terminate || end.ContinueAs != nil || len(end.ProduceEvents) > 0) // TODO: check the spec continueAs/produceEvents to see how it influences the end

	if !hasTransition && !isEnd {
		structLevel.ReportError(field, "Transition", "transition", TagRequired, "")
	} else if hasTransition && isEnd {
		structLevel.ReportError(field, "Transition", "transition", TagExclusive, "")
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
