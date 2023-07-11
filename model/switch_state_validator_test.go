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

func buildSwitchState(workflow *Workflow, name string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeSwitch,
		},
		SwitchState: &SwitchState{},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func buildDefaultCondition(state *State) *DefaultCondition {
	state.SwitchState.DefaultCondition = DefaultCondition{}
	return &state.SwitchState.DefaultCondition
}

func buildDataCondition(state *State, name, condition string) *DataCondition {
	if state.SwitchState.DataConditions == nil {
		state.SwitchState.DataConditions = []DataCondition{}
	}

	dataCondition := DataCondition{
		Name:      name,
		Condition: condition,
	}

	state.SwitchState.DataConditions = append(state.SwitchState.DataConditions, dataCondition)
	return &state.SwitchState.DataConditions[len(state.SwitchState.DataConditions)-1]
}

func buildEventCondition(workflow *Workflow, state *State, name, eventRef string) (*Event, *EventCondition) {
	workflow.Events = append(workflow.Events, Event{
		Name: eventRef,
		Type: "event type",
		Kind: EventKindConsumed,
	})

	eventCondition := EventCondition{
		Name:     name,
		EventRef: eventRef,
	}

	state.SwitchState.EventConditions = append(state.SwitchState.EventConditions, eventCondition)
	return &workflow.Events[len(workflow.Events)-1], &state.SwitchState.EventConditions[len(state.SwitchState.EventConditions)-1]
}

func TestSwitchStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	swithState := buildSwitchState(baseWorkflow, "start state")
	defaultCondition := buildDefaultCondition(swithState)
	buildEndByDefaultCondition(defaultCondition, true, false)

	dataCondition := buildDataCondition(swithState, "data condition 1", "1=1")
	buildEndByDataCondition(dataCondition, true, false)

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
				model.States[0].SwitchState.DataConditions = nil
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildEventCondition(model, &model.States[0], "event condition", "event 1")
				buildEndByEventCondition(&model.States[0].SwitchState.EventConditions[0], true, false)
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions exclusive`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestDefaultConditionStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	buildSwitchState(baseWorkflow, "start state")
	buildDefaultCondition(&baseWorkflow.States[0])

	buildDataCondition(&baseWorkflow.States[0], "data condition 1", "1=1")
	buildEndByDataCondition(&baseWorkflow.States[0].SwitchState.DataConditions[0], true, false)
	buildDataCondition(&baseWorkflow.States[0], "data condition 2", "1=1")

	buildOperationState(baseWorkflow, "end state")
	buildEndByState(&baseWorkflow.States[1], true, false)
	buildActionByOperationState(&baseWorkflow.States[1], "action 1")
	buildFunctionRef(baseWorkflow, &baseWorkflow.States[1].OperationState.Actions[0], "function 1")

	buildTransitionByDefaultCondition(&baseWorkflow.States[0].SwitchState.DefaultCondition, &baseWorkflow.States[1])
	buildTransitionByDataCondition(&baseWorkflow.States[0].SwitchState.DataConditions[1], &baseWorkflow.States[1], false)

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
				model.States[0].SwitchState.DataConditions[0].End = nil
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions[0].transition is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildTransitionByDataCondition(&model.States[0].SwitchState.DataConditions[0], &model.States[1], false)
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions[0].transition exclusive`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestSwitchStateTimeoutStructLevelValidation(t *testing.T) {
}

func TestEventConditionStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.States = make(States, 0, 2)

	// switch state
	switchState := buildSwitchState(baseWorkflow, "start state")

	// default condition
	defaultCondition := buildDefaultCondition(switchState)
	buildEndByDefaultCondition(defaultCondition, true, false)

	// event condition 1
	_, eventCondition := buildEventCondition(baseWorkflow, switchState, "data condition 1", "event 1")
	buildEndByEventCondition(eventCondition, true, false)

	// event condition 2
	_, eventCondition2 := buildEventCondition(baseWorkflow, switchState, "data condition 2", "event 2")
	buildEndByEventCondition(eventCondition2, true, false)

	// operation state
	operationState := buildOperationState(baseWorkflow, "end state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	// trasition switch state to operation state
	buildTransitionByEventCondition(eventCondition, operationState, false)

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				return *baseWorkflow.DeepCopy()
			},
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].SwitchState.EventConditions[0].EventRef = "event not found"
				return *model
			},
			Err: `workflow.states[0].switchState.eventConditions[0].eventRef don't exist "event not found"`,
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].SwitchState.EventConditions[0].End = nil
				return *model
			},
			Err: `workflow.states[0].switchState.eventConditions[0].transition is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildTransitionByEventCondition(&model.States[0].SwitchState.EventConditions[0], &model.States[1], false)
				return *model
			},
			Err: `workflow.states[0].switchState.eventConditions[0].transition exclusive`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestDataConditionStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	// switch state
	swithcState := buildSwitchState(baseWorkflow, "start state")

	// default condition
	defaultCondition := buildDefaultCondition(swithcState)
	buildEndByDefaultCondition(defaultCondition, true, false)

	// data condition
	dataCondition := buildDataCondition(swithcState, "data condition 1", "1=1")
	buildEndByDataCondition(dataCondition, true, false)

	// operation state
	operationState := buildOperationState(baseWorkflow, "end state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

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
				model.States[0].SwitchState.DataConditions[0].End = nil
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions[0].transition is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				buildTransitionByDataCondition(&model.States[0].SwitchState.DataConditions[0], &model.States[1], false)
				return *model
			},
			Err: `workflow.states[0].switchState.dataConditions[0].transition exclusive`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
