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
	"reflect"

	validator "github.com/go-playground/validator/v10"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidation(continueAsStructLevelValidation, ContinueAs{})
	val.GetValidator().RegisterStructValidation(workflowStructLevelValidation, Workflow{})
}

func continueAsStructLevelValidation(structLevel validator.StructLevel) {
	continueAs := structLevel.Current().Interface().(ContinueAs)
	if len(continueAs.WorkflowExecTimeout.Duration) > 0 {
		if err := val.ValidateISO8601TimeDuration(continueAs.WorkflowExecTimeout.Duration); err != nil {
			structLevel.ReportError(reflect.ValueOf(continueAs.WorkflowExecTimeout.Duration),
				"workflowExecTimeout", "duration", "iso8601duration", "")
		}
	}
}

// WorkflowStructLevelValidation custom validator
func workflowStructLevelValidation(structLevel validator.StructLevel) {
	// unique name of the auth methods
	// NOTE: we cannot add the custom validation of auth to Auth
	// because `RegisterStructValidation` only works with struct type
	wf := structLevel.Current().Interface().(Workflow)
	dict := map[string]bool{}

	for _, a := range wf.BaseWorkflow.Auth {
		if !dict[a.Name] {
			dict[a.Name] = true
		} else {
			structLevel.ReportError(reflect.ValueOf(a.Name), "[]Auth.Name", "name", "reqnameunique", "")
		}
	}

	startAndStatesTransitionValidator(structLevel, wf.BaseWorkflow.Start, wf.States)
}

func startAndStatesTransitionValidator(structLevel validator.StructLevel, start *Start, states []State) {
	statesMap := make(map[string]State, len(states))
	for _, state := range states {
		statesMap[state.Name] = state
	}

	if start != nil {
		// if not exists the start transtion stop the states validations
		if _, ok := statesMap[start.StateName]; !ok {
			structLevel.ReportError(reflect.ValueOf(start), "Start", "start", "startnotexist", "")
			return
		}
	}

	if len(states) == 1 {
		return
	}

	// Naive check if transitions exist
	for _, state := range statesMap {
		if state.Transition != nil {
			if _, ok := statesMap[state.Transition.NextState]; !ok {
				structLevel.ReportError(reflect.ValueOf(state), "Transition", "transition", "transitionnotexists", state.Transition.NextState)
			}
		}
	}

	// TODO: create states graph to complex check
}

func validTransitionAndEnd(structLevel validator.StructLevel, field interface{}, transition *Transition, end *End) {
	hasTransition := transition != nil
	isEnd := end != nil && (end.Terminate || end.ContinueAs != nil || len(end.ProduceEvents) > 0) // TODO: check the spec continueAs/produceEvents to see how it influences the end

	if !hasTransition && !isEnd {
		structLevel.ReportError(field, "Transition", "transition", "required", "must have one of transition, end")
	} else if hasTransition && isEnd {
		structLevel.ReportError(field, "Transition", "transition", "exclusive", "must have one of transition, end")
	}
}
