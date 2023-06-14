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
	"testing"
)

func buildCallbackState(workflow *Workflow, name, eventRef string) *State {
	consumeEvent := Event{
		Name: eventRef,
		Type: "event type",
		Kind: EventKindProduced,
	}
	workflow.Events = append(workflow.Events, consumeEvent)

	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeCallback,
		},
		CallbackState: &CallbackState{
			EventRef: eventRef,
		},
	}
	workflow.States = append(workflow.States, state)

	return &workflow.States[len(workflow.States)-1]
}

func buildCallbackStateTimeout(callbackState *CallbackState) *CallbackStateTimeout {
	callbackState.Timeouts = &CallbackStateTimeout{}
	return callbackState.Timeouts
}

func TestCallbackStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	callbackState := buildCallbackState(baseWorkflow, "start state", "event 1")
	buildEndByState(callbackState, true, false)
	buildFunctionRef(baseWorkflow, &callbackState.Action, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].CallbackState.EventRef = ""
				return *model
			},
			Err: `workflow.states[0].callbackState.eventRef is required`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestCallbackStateTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	callbackState := buildCallbackState(baseWorkflow, "start state", "event 1")
	buildEndByState(callbackState, true, false)
	buildCallbackStateTimeout(callbackState.CallbackState)
	buildFunctionRef(baseWorkflow, &callbackState.Action, "function 1")

	testCases := []ValidationCase{
		{
			Desp: `success`,
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].CallbackState.Timeouts.ActionExecTimeout = "P5S"
				model.States[0].CallbackState.Timeouts.EventTimeout = "P5S"
				return *model
			},
			Err: `workflow.states[0].callbackState.timeouts.actionExecTimeout invalid iso8601 duration "P5S"
workflow.states[0].callbackState.timeouts.eventTimeout invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
