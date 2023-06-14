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

	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

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

	contextValue := val.ValidatorContext{
		MapStates:    val.NewMapValues(workflow.States, "Name"),
		MapFunctions: val.NewMapValues(workflow.Functions, "Name"),
		MapEvents:    val.NewMapValues(workflow.Events, "Name"),
		MapRetries:   val.NewMapValues(workflow.Retries, "Name"),
		MapErrors:    val.NewMapValues(workflow.Errors, "Name"),
	}

	return context.WithValue(context.Background(), val.ValidatorContextValue, contextValue)
}

func init() {
	// TODO: create states graph to complex check

	// val.GetValidator().RegisterStructValidationCtx(val.ValidationWrap(nil, workflowStructLevelValidation), Workflow{})
	val.GetValidator().RegisterStructValidationCtx(val.ValidationWrap(onErrorStructLevelValidationCtx), OnError{})
	val.GetValidator().RegisterStructValidationCtx(val.ValidationWrap(transitionStructLevelValidationCtx), Transition{})
	val.GetValidator().RegisterStructValidationCtx(val.ValidationWrap(startStructLevelValidationCtx), Start{})
}

func startStructLevelValidationCtx(ctx val.ValidatorContext, structLevel validator.StructLevel) {
	start := structLevel.Current().Interface().(Start)
	if start.StateName != "" && !ctx.MapStates.Contain(start.StateName) {
		structLevel.ReportError(start.StateName, "StateName", "stateName", val.TagExists, "")
		return
	}
}

func onErrorStructLevelValidationCtx(ctx val.ValidatorContext, structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)
	hasErrorRef := onError.ErrorRef != ""
	hasErrorRefs := len(onError.ErrorRefs) > 0

	if !hasErrorRef && !hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagRequired, "")
	} else if hasErrorRef && hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagExclusive, "")
	}

	if onError.ErrorRef != "" && !ctx.MapErrors.Contain(onError.ErrorRef) {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagExists, "")
	}

	for _, errorRef := range onError.ErrorRefs {
		if !ctx.MapErrors.Contain(errorRef) {
			structLevel.ReportError(onError.ErrorRefs, "ErrorRefs", "ErrorRefs", val.TagExists, "")
		}
	}
}

func transitionStructLevelValidationCtx(ctx val.ValidatorContext, structLevel validator.StructLevel) {
	// Naive check if transitions exist
	transition := structLevel.Current().Interface().(Transition)
	if ctx.MapStates.Contain(transition.NextState) {
		if transition.stateParent != nil {
			parentBaseState := transition.stateParent

			if parentBaseState.Name == transition.NextState {
				// TODO: Improve recursive check
				structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagRecursiveState, parentBaseState.Name)
			}

			if parentBaseState.UsedForCompensation && !ctx.MapStates.ValuesMap[transition.NextState].(State).UsedForCompensation {
				structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagTransitionUseForCompensation, "")

			} else if !parentBaseState.UsedForCompensation && ctx.MapStates.ValuesMap[transition.NextState].(State).UsedForCompensation {
				structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagTransitionMainWorkflow, "")
			}
		}

	} else {
		structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagExists, "")
	}
}

func validTransitionAndEnd(structLevel validator.StructLevel, field any, transition *Transition, end *End) {
	hasTransition := transition != nil
	isEnd := end != nil && (end.Terminate || end.ContinueAs != nil || len(end.ProduceEvents) > 0) // TODO: check the spec continueAs/produceEvents to see how it influences the end

	if !hasTransition && !isEnd {
		structLevel.ReportError(field, "Transition", "transition", val.TagRequired, "")
	} else if hasTransition && isEnd {
		structLevel.ReportError(field, "Transition", "transition", val.TagExclusive, "")
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

	return hasOne && hasTwo
}
