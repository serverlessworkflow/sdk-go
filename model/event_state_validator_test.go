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

import "testing"

func buildEventState(workflow *Workflow, name string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeEvent,
		},
		EventState: &EventState{},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func buildOnEvents(workflow *Workflow, state *State, name string) *OnEvents {
	event := Event{
		Name: name,
		Type: "type",
		Kind: EventKindProduced,
	}
	workflow.Events = append(workflow.Events, event)

	state.EventState.OnEvents = append(state.EventState.OnEvents, OnEvents{
		EventRefs:  []string{event.Name},
		ActionMode: ActionModeParallel,
	})

	return &state.EventState.OnEvents[len(state.EventState.OnEvents)-1]
}

func buildEventStateTimeout(state *State) *EventStateTimeout {
	state.EventState.Timeouts = &EventStateTimeout{
		ActionExecTimeout: "PT5S",
		EventTimeout:      "PT5S",
	}
	return state.EventState.Timeouts
}

func TestEventStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	eventState := buildEventState(baseWorkflow, "start state")
	buildOnEvents(baseWorkflow, eventState, "event 1")
	buildEndByState(eventState, true, false)

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
				model.States[0].EventState.OnEvents = nil
				return *model
			},
			Err: `workflow.states[0].eventState.onEvents is required`,
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].EventState.OnEvents = []OnEvents{}
				return *model
			},
			Err: `workflow.states[0].eventState.onEvents min > 1`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestOnEventsStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	eventState := buildEventState(baseWorkflow, "start state")
	buildOnEvents(baseWorkflow, eventState, "event 1")
	buildEndByState(eventState, true, false)

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
				model.States[0].EventState.OnEvents[0].EventRefs = nil
				model.States[0].EventState.OnEvents[0].ActionMode = ""
				return *model
			},
			Err: `workflow.states[0].eventState.onEvents[0].eventRefs is required
workflow.states[0].eventState.onEvents[0].actionMode is required`,
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].EventState.OnEvents[0].EventRefs = []string{}
				return *model
			},
			Err: `workflow.states[0].eventState.onEvents[0].eventRefs min > 1`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].EventState.OnEvents[0].ActionMode = ActionModeParallel + "invalid"
				return *model
			},
			Err: `workflow.states[0].eventState.onEvents[0].actionMode need by one of [sequential parallel]`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestEventStateTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	eventState := buildEventState(baseWorkflow, "start state")
	buildEventStateTimeout(eventState)
	buildOnEvents(baseWorkflow, eventState, "event 1")
	buildEndByState(eventState, true, false)

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "omitempty",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].EventState.Timeouts.ActionExecTimeout = ""
				model.States[0].EventState.Timeouts.EventTimeout = ""
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].EventState.Timeouts.ActionExecTimeout = "P5S"
				model.States[0].EventState.Timeouts.EventTimeout = "P5S"
				return *model
			},
			Err: `workflow.states[0].eventState.timeouts.actionExecTimeout invalid iso8601 duration "P5S"
workflow.states[0].eventState.timeouts.eventTimeout invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
