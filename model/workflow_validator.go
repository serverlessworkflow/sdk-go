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
	"errors"
	"fmt"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

type workflowValidator func(mapValues ValidatorContext, sl validator.StructLevel)

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

type ValidatorContext struct {
	MapStates    mapValues[State]
	MapFunctions mapValues[Function]
	MapEvents    mapValues[Event]
	MapRetries   mapValues[Retry]
	MapErrors    mapValues[Error]
}

func validationWrap(fnCtx workflowValidator) validator.StructLevelFuncCtx {
	return func(ctx context.Context, structLevel validator.StructLevel) {
		if fnCtx != nil {
			if mapValues, ok := ctx.Value(validatorContextValue).(ValidatorContext); ok {
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

	contextValue := ValidatorContext{
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
	val.GetValidator().RegisterStructValidationCtx(validationWrap(onErrorStructLevelValidationCtx), OnError{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(transitionStructLevelValidationCtx), Transition{})
	val.GetValidator().RegisterStructValidationCtx(validationWrap(startStructLevelValidationCtx), Start{})
}

func startStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	start := structLevel.Current().Interface().(Start)
	if !ctx.MapStates.contain(start.StateName) {
		structLevel.ReportError(start.StateName, "StateName", "stateName", TagExists, "")
		return
	}
}

func onErrorStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	onError := structLevel.Current().Interface().(OnError)
	hasErrorRef := onError.ErrorRef != ""
	hasErrorRefs := len(onError.ErrorRefs) > 0

	if !hasErrorRef && !hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagRequired, "")
	} else if hasErrorRef && hasErrorRefs {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagExclusive, "")
	}

	if onError.ErrorRef != "" && !ctx.MapErrors.contain(onError.ErrorRef) {
		structLevel.ReportError(onError.ErrorRef, "ErrorRef", "errorRef", TagExists, "")
	}

	for _, errorRef := range onError.ErrorRefs {
		if !ctx.MapErrors.contain(errorRef) {
			structLevel.ReportError(onError.ErrorRefs, "ErrorRefs", "errorRefs", TagExists, "")
		}
	}
}

func transitionStructLevelValidationCtx(ctx ValidatorContext, structLevel validator.StructLevel) {
	// Naive check if transitions exist
	transition := structLevel.Current().Interface().(Transition)
	if ctx.MapStates.contain(transition.NextState) {
		if transition.stateParent != nil {
			parentBaseState := transition.stateParent

			if parentBaseState.Name == transition.NextState {
				// TODO: Improve recursive check
				structLevel.ReportError(transition.NextState, "NextState", "nextState", TagRecursiveState, parentBaseState.Name)
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

type WorkflowErrors []error

func (e WorkflowErrors) Error() string {
	errors := []string{}
	for _, err := range []error(e) {
		errors = append(errors, err.Error())
	}
	return strings.Join(errors, "\n")
}

func WorkflowError(err error) error {
	if err == nil {
		return nil
	}

	if _, ok := err.(*validator.InvalidValidationError); ok {
		return err
	}

	workflowErrors := []error{}
	for _, err := range err.(validator.ValidationErrors) {
		// fmt.Println("Namespace", err.Namespace())
		// fmt.Println("Field", err.Field())
		// fmt.Println("StructNamespace", err.StructNamespace())
		// fmt.Println("StructField", err.StructField())
		// fmt.Println("Tag", err.Tag())
		// // fmt.Println(err.ActualTag())
		// // fmt.Println(err.Kind())
		// // fmt.Println(err.Type())
		// fmt.Println("value", err.Value())
		// fmt.Println("param", err.Param())
		// fmt.Println()

		// normalize namespace
		namespaceList := strings.Split(err.Namespace(), ".")[1:]
		newNamespaceList := []string{}
		for i := range namespaceList {
			part := namespaceList[i]
			if part != "Workflow" && part != "BaseWorkflow" && part != "BaseState" {
				part := strings.ToLower(part[:1]) + part[1:]
				newNamespaceList = append(newNamespaceList, part)
			}
		}
		namespace := strings.Join(newNamespaceList, ".")

		switch err.Tag() {
		case "exists":
			workflowErrors = append(workflowErrors, fmt.Errorf("%s don't exists %q", namespace, err.Value()))
		case TagCompensatedby:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s compensatedBy don't exists %q", namespace, err.Value()))
		case "unique":
			if err.Param() == "" {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s has duplicate value", namespace))
			} else {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s has duplicate %q", namespace, strings.ToLower(err.Param())))
			}
		case "required_without":
			if err.Param() == "ID" {
				workflowErrors = append(workflowErrors, errors.New("id required when not defined \"key\""))
			} else if err.Param() == "Key" {
				workflowErrors = append(workflowErrors, errors.New("key required when not defined \"id\""))
			} else {
				workflowErrors = append(workflowErrors, err)
			}
		case TagRecursiveState:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s can't no be recursive %q", namespace, strings.ToLower(err.Param())))
		default:
			workflowErrors = append(workflowErrors, err)
		}
	}

	return WorkflowErrors(workflowErrors)
}
