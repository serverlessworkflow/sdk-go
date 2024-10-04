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

	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"

	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

type contextValueKey string

const ValidatorContextValue contextValueKey = "value"

type WorkflowValidator func(mapValues ValidatorContext, sl validator.StructLevel)

func ValidationWrap(fnCtx WorkflowValidator) validator.StructLevelFuncCtx {
	return func(ctx context.Context, structLevel validator.StructLevel) {
		if fnCtx != nil {
			if mapValues, ok := ctx.Value(ValidatorContextValue).(ValidatorContext); ok {
				fnCtx(mapValues, structLevel)
			}
		}
	}
}

// +builder-gen:ignore=true
type ValidatorContext struct {
	States    map[string]State
	Functions map[string]Function
	Events    map[string]Event
	Retries   map[string]Retry
	Errors    map[string]Error
}

func (c *ValidatorContext) init(workflow *Workflow) {
	c.States = make(map[string]State, len(workflow.States))
	for _, state := range workflow.States {
		c.States[state.BaseState.Name] = state
	}

	c.Functions = make(map[string]Function, len(workflow.Functions))
	for _, function := range workflow.Functions {
		c.Functions[function.Name] = function
	}

	c.Events = make(map[string]Event, len(workflow.Events))
	for _, event := range workflow.Events {
		c.Events[event.Name] = event
	}

	c.Retries = make(map[string]Retry, len(workflow.Retries))
	for _, retry := range workflow.Retries {
		c.Retries[retry.Name] = retry
	}

	c.Errors = make(map[string]Error, len(workflow.Errors))
	for _, error := range workflow.Errors {
		c.Errors[error.Name] = error
	}
}

func (c *ValidatorContext) ExistState(name string) bool {
	if c.States == nil {
		return true
	}
	_, ok := c.States[name]
	return ok
}

func (c *ValidatorContext) ExistFunction(name string) bool {
	if c.Functions == nil {
		return true
	}
	_, ok := c.Functions[name]
	return ok
}

func (c *ValidatorContext) ExistEvent(name string) bool {
	if c.Events == nil {
		return true
	}
	_, ok := c.Events[name]
	return ok
}

func (c *ValidatorContext) ExistRetry(name string) bool {
	if c.Retries == nil {
		return true
	}
	_, ok := c.Retries[name]
	return ok
}

func (c *ValidatorContext) ExistError(name string) bool {
	if c.Errors == nil {
		return true
	}
	_, ok := c.Errors[name]
	return ok
}

func NewValidatorContext(object any) context.Context {
	contextValue := ValidatorContext{}

	if workflow, ok := object.(*Workflow); ok {
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
		contextValue.init(workflow)
	}

	return context.WithValue(context.Background(), ValidatorContextValue, contextValue)
}

func init() {
	// TODO: create states graph to complex check

	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(onErrorStructLevelValidationCtx), OnError{})
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(transitionStructLevelValidationCtx), Transition{})
	val.GetValidator().RegisterStructValidationCtx(ValidationWrap(startStructLevelValidationCtx), Start{})

	val.GetValidator().RegisterStructValidation(floatstr.ValidateFloat32OrString, Retry{})
}

func startStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	start := structLevel.Current().Interface().(Start)
	if start.StateName != "" && !ctx.ExistState(start.StateName) {
		structLevel.ReportError(start.StateName, "StateName", "stateName", val.TagExists, "")
		return
	}
}

func onErrorStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)
	hasErrorRef := onError.ErrorRef != ""
	hasErrorRefs := len(onError.ErrorRefs) > 0

	if !hasErrorRef && !hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagRequired, "")
	} else if hasErrorRef && hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagExclusive, "")
		return
	}

	if onError.ErrorRef != "" && !ctx.ExistError(onError.ErrorRef) {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "ErrorRef", val.TagExists, "")
	}

	for _, errorRef := range onError.ErrorRefs {
		if !ctx.ExistError(errorRef) {
			structLevel.ReportError(onError.ErrorRefs, "ErrorRefs", "ErrorRefs", val.TagExists, "")
		}
	}
}

func transitionStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	// Naive check if transitions exist
	transition := structLevel.Current().Interface().(Transition)
	if ctx.ExistState(transition.NextState) {
		if transition.stateParent != nil {
			parentBaseState := transition.stateParent

			if parentBaseState.Name == transition.NextState {
				// TODO: Improve recursive check
				structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagRecursiveState, parentBaseState.Name)
			}

			if parentBaseState.UsedForCompensation && !ctx.States[transition.NextState].BaseState.UsedForCompensation {
				structLevel.ReportError(transition.NextState, "NextState", "NextState", val.TagTransitionUseForCompensation, "")
			}

			if !parentBaseState.UsedForCompensation && ctx.States[transition.NextState].BaseState.UsedForCompensation {
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
