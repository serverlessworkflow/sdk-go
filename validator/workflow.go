// Copyright 2023 The Serverless Workflow Specification Authors
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

package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	validator "github.com/go-playground/validator/v10"
)

const (
	TagExists    string = "exists"
	TagRequired  string = "required"
	TagExclusive string = "exclusive"

	TagRecursiveState string = "recursivestate"

	// States referenced by compensatedBy (as well as any other states that they transition to) must obey following rules:
	TagTransitionMainWorkflow       string = "transtionmainworkflow"         // They should not have any incoming transitions (should not be part of the main workflow control-flow logic)
	TagCompensatedbyEventState      string = "compensatedbyeventstate"       // They cannot be an event state
	TagRecursiveCompensation        string = "recursivecompensation"         // They cannot themselves set their compensatedBy property to true (compensation is not recursive)
	TagCompensatedby                string = "compensatedby"                 // They must define the usedForCompensation property and set it to true
	TagTransitionUseForCompensation string = "transitionusedforcompensation" // They can transition only to states which also have their usedForCompensation property and set to true
)

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

	var invalidErr *validator.InvalidValidationError
	if errors.As(err, &invalidErr) {
		return err
	}

	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err
	}

	removeNamespace := []string{
		"BaseWorkflow",
		"BaseState",
		"OperationState",
	}

	workflowErrors := []error{}
	for _, err := range validationErrors {
		// normalize namespace
		namespaceList := strings.Split(err.Namespace(), ".")
		normalizedNamespaceList := []string{}
		for i := range namespaceList {
			part := namespaceList[i]
			if !contains(removeNamespace, part) {
				part := strings.ToLower(part[:1]) + part[1:]
				normalizedNamespaceList = append(normalizedNamespaceList, part)
			}
		}
		namespace := strings.Join(normalizedNamespaceList, ".")

		switch err.Tag() {
		case "unique":
			if err.Param() == "" {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s has duplicate value", namespace))
			} else {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s has duplicate %q", namespace, strings.ToLower(err.Param())))
			}
		case "min":
			workflowErrors = append(workflowErrors, fmt.Errorf("%s must have the minimum %s", namespace, err.Param()))
		case "required_without":
			if namespace == "workflow.iD" {
				workflowErrors = append(workflowErrors, errors.New("workflow.id required when \"workflow.key\" is not defined"))
			} else if namespace == "workflow.key" {
				workflowErrors = append(workflowErrors, errors.New("workflow.key required when \"workflow.id\" is not defined"))
			} else if err.StructField() == "FunctionRef" {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s required when \"eventRef\" or \"subFlowRef\" is not defined", namespace))
			} else {
				workflowErrors = append(workflowErrors, err)
			}
		case "oneofkind":
			value := reflect.New(err.Type()).Elem().Interface().(Kind)
			workflowErrors = append(workflowErrors, fmt.Errorf("%s need by one of %s", namespace, value.KindValues()))
		case "gt0":
			workflowErrors = append(workflowErrors, fmt.Errorf("%s must be greater than 0", namespace))
		case TagExists:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s don't exist %q", namespace, err.Value()))
		case TagRequired:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s is required", namespace))
		case TagExclusive:
			if err.StructField() == "ErrorRef" {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s or %s are exclusive", namespace, replaceLastNamespace(namespace, "errorRefs")))
			} else {
				workflowErrors = append(workflowErrors, fmt.Errorf("%s exclusive", namespace))
			}
		case TagCompensatedby:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s = %q is not defined as usedForCompensation", namespace, err.Value()))
		case TagCompensatedbyEventState:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s = %q is defined as usedForCompensation and cannot be an event state", namespace, err.Value()))
		case TagRecursiveCompensation:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s = %q is defined as usedForCompensation (cannot themselves set their compensatedBy)", namespace, err.Value()))
		case TagRecursiveState:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s can't no be recursive %q", namespace, strings.ToLower(err.Param())))
		case TagISO8601Duration:
			workflowErrors = append(workflowErrors, fmt.Errorf("%s invalid iso8601 duration %q", namespace, err.Value()))
		default:
			workflowErrors = append(workflowErrors, err)
		}
	}

	return WorkflowErrors(workflowErrors)
}

func contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func replaceLastNamespace(namespace, replace string) string {
	index := strings.LastIndex(namespace, ".")
	if index == -1 {
		return namespace
	}

	return fmt.Sprintf("%s.%s", namespace[:index], replace)
}
