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
	"testing"
)

func buildEventRef(workflow *Workflow, action *Action, triggerEvent, resultEvent string) *EventRef {
	produceEvent := Event{
		Name: triggerEvent,
		Type: "event type",
		Kind: EventKindProduced,
	}

	consumeEvent := Event{
		Name: resultEvent,
		Type: "event type",
		Kind: EventKindProduced,
	}

	workflow.Events = append(workflow.Events, produceEvent)
	workflow.Events = append(workflow.Events, consumeEvent)

	eventRef := &EventRef{
		TriggerEventRef: triggerEvent,
		ResultEventRef:  resultEvent,
		Invoke:          InvokeKindSync,
	}

	action.EventRef = eventRef
	return action.EventRef
}

func buildCorrelation(event *Event) *Correlation {
	event.Correlation = append(event.Correlation, Correlation{
		ContextAttributeName: "attribute name",
	})

	return &event.Correlation[len(event.Correlation)-1]
}

func TestEventStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.Events = Events{{
		Name: "event 1",
		Type: "event type",
		Kind: EventKindConsumed,
	}}

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Events = append(model.Events, model.Events[0])
				return *model
			},
			Err: `workflow.events has duplicate "name"`,
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Events[0].Name = ""
				model.Events[0].Type = ""
				model.Events[0].Kind = ""
				return *model
			},
			Err: `workflow.events[0].name is required
workflow.events[0].type is required
workflow.events[0].kind is required`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Events[0].Kind = EventKindConsumed + "invalid"
				return *model
			},
			Err: `workflow.events[0].kind need by one of [consumed produced]`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestCorrelationStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.Events = Events{{
		Name: "event 1",
		Type: "event type",
		Kind: EventKindConsumed,
	}}

	buildCorrelation(&baseWorkflow.Events[0])

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "empty",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Events[0].Correlation = nil
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Events[0].Correlation[0].ContextAttributeName = ""
				return *model
			},
			Err: `workflow.events[0].correlation[0].contextAttributeName is required`,
		},
		//TODO: Add test: correlation only used for `consumed` events
	}

	StructLevelValidationCtx(t, testCases)
}

func TestEventRefStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	eventRef := buildEventRef(baseWorkflow, action1, "event 1", "event 2")
	eventRef.ResultEventTimeout = "PT1H"

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				return *baseWorkflow.DeepCopy()
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].EventRef.TriggerEventRef = ""
				model.States[0].OperationState.Actions[0].EventRef.ResultEventRef = ""
				return *model
			},
			Err: `workflow.states[0].actions[0].eventRef.triggerEventRef is required
workflow.states[0].actions[0].eventRef.resultEventRef is required`,
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].EventRef.TriggerEventRef = "invalid event"
				model.States[0].OperationState.Actions[0].EventRef.ResultEventRef = "invalid event 2"
				return *model
			},
			Err: `workflow.states[0].actions[0].eventRef.triggerEventRef don't exist "invalid event"
workflow.states[0].actions[0].eventRef.triggerEventRef don't exist "invalid event 2"`,
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].EventRef.ResultEventTimeout = "10hs"
				return *model
			},
			Err: `workflow.states[0].actions[0].eventRef.resultEventTimeout invalid iso8601 duration "10hs"`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OperationState.Actions[0].EventRef.Invoke = InvokeKindSync + "invalid"
				return *model
			},
			Err: `workflow.states[0].actions[0].eventRef.invoke need by one of [sync async]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
