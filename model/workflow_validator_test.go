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

	"github.com/stretchr/testify/assert"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func buildWorkflow() *Workflow {
	return &Workflow{
		BaseWorkflow: BaseWorkflow{
			ID:             "id",
			Key:            "key",
			Name:           "name",
			SpecVersion:    "0.8",
			Version:        "0.1",
			ExpressionLang: JqExpressionLang,
		},
	}
}

func buildEndByState(state *State, terminate, compensate bool) *End {
	end := &End{
		Terminate:  terminate,
		Compensate: compensate,
	}
	state.BaseState.End = end
	return end
}

func buildEndByDefaultCondition(defaultCondition *DefaultCondition, terminate, compensate bool) *End {
	end := &End{
		Terminate:  terminate,
		Compensate: compensate,
	}
	defaultCondition.End = end
	return end
}

func buildEndByDataCondition(dataCondition *DataCondition, terminate, compensate bool) *End {
	end := &End{
		Terminate:  terminate,
		Compensate: compensate,
	}
	dataCondition.End = end
	return end
}

func buildEndByEventCondition(eventCondition *EventCondition, terminate, compensate bool) *End {
	end := &End{
		Terminate:  terminate,
		Compensate: compensate,
	}
	eventCondition.End = end
	return end
}

func buildStart(workflow *Workflow, state *State) {
	start := &Start{
		StateName: state.BaseState.Name,
	}
	workflow.BaseWorkflow.Start = start
}

func buildTransitionByState(state, nextState *State, compensate bool) {
	state.BaseState.Transition = &Transition{
		NextState:  nextState.BaseState.Name,
		Compensate: compensate,
	}
}

func buildTransitionByDataCondition(dataCondition *DataCondition, state *State, compensate bool) {
	dataCondition.Transition = &Transition{
		NextState:  state.BaseState.Name,
		Compensate: compensate,
	}
}

func buildTransitionByEventCondition(eventCondition *EventCondition, state *State, compensate bool) {
	eventCondition.Transition = &Transition{
		NextState:  state.BaseState.Name,
		Compensate: compensate,
	}
}

func buildTransitionByDefaultCondition(defaultCondition *DefaultCondition, state *State) {
	defaultCondition.Transition = &Transition{
		NextState: state.BaseState.Name,
	}
}

func buildTimeouts(workflow *Workflow) *Timeouts {
	timeouts := Timeouts{}
	workflow.BaseWorkflow.Timeouts = &timeouts
	return workflow.BaseWorkflow.Timeouts
}

func TestBaseWorkflowStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
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
			Desp: "id exclude key",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.ID = "id"
				model.Key = ""
				return *model
			},
		},
		{
			Desp: "key exclude id",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.ID = ""
				model.Key = "key"
				return *model
			},
		},
		{
			Desp: "without id and key",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.ID = ""
				model.Key = ""
				return *model
			},
			Err: `workflow.id required when "workflow.key" is not defined
workflow.key required when "workflow.id" is not defined`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.BaseWorkflow.ExpressionLang = JqExpressionLang + "invalid"
				return *model
			},
			Err: `workflow.expressionLang need by one of [jq jsonpath cel]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestContinueAsStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	baseWorkflow.States[0].BaseState.End.ContinueAs = &ContinueAs{
		WorkflowID: "sub workflow",
		WorkflowExecTimeout: WorkflowExecTimeout{
			Duration: "P1M",
		},
	}

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
				model.States[0].BaseState.End.ContinueAs.WorkflowID = ""
				return *model
			},
			Err: `workflow.states[0].end.continueAs.workflowID is required`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestOnErrorStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	baseWorkflow.BaseWorkflow.Errors = Errors{{
		Name: "error 1",
	}, {
		Name: "error 2",
	}}
	baseWorkflow.States[0].BaseState.OnErrors = []OnError{{
		ErrorRef: "error 1",
	}, {
		ErrorRefs: []string{"error 1", "error 2"},
	}}

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
				model.States[0].BaseState.OnErrors[0].ErrorRef = ""
				return *model
			},
			Err: `workflow.states[0].onErrors[0].errorRef is required`,
		},
		{
			Desp: "exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OnErrors[0].ErrorRef = "error 1"
				model.States[0].OnErrors[0].ErrorRefs = []string{"error 2"}
				return *model
			},
			Err: `workflow.states[0].onErrors[0].errorRef or workflow.states[0].onErrors[0].errorRefs are exclusive`,
		},
		{
			Desp: "exists and exclusive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.OnErrors[0].ErrorRef = "invalid error name"
				model.States[0].BaseState.OnErrors[0].ErrorRefs = []string{"invalid error name"}
				return *model
			},
			Err: `workflow.states[0].onErrors[0].errorRef or workflow.states[0].onErrors[0].errorRefs are exclusive`,
		},
		{
			Desp: "exists errorRef",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.OnErrors[0].ErrorRef = "invalid error name"
				return *model
			},
			Err: `workflow.states[0].onErrors[0].errorRef don't exist "invalid error name"`,
		},
		{
			Desp: "exists errorRefs",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.OnErrors[0].ErrorRef = ""
				model.States[0].BaseState.OnErrors[0].ErrorRefs = []string{"invalid error name"}
				return *model
			},
			Err: `workflow.states[0].onErrors[0].errorRefs don't exist ["invalid error name"]`,
		},
		{
			Desp: "duplicate",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].OnErrors[1].ErrorRefs = []string{"error 1", "error 1"}
				return *model
			},
			Err: `workflow.states[0].onErrors[1].errorRefs has duplicate value`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestStartStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildStart(baseWorkflow, operationState)
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
				model.Start.StateName = ""
				return *model
			},
			Err: `workflow.start.stateName is required`,
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Start.StateName = "start state not found"
				return *model
			},
			Err: `workflow.start.stateName don't exist "start state not found"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestTransitionStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	baseWorkflow.States = make(States, 0, 5)

	operationState := buildOperationState(baseWorkflow, "start state")
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	operationState2 := buildOperationState(baseWorkflow, "next state")
	buildEndByState(operationState2, true, false)
	operationState2.BaseState.CompensatedBy = "compensation next state 1"
	action2 := buildActionByOperationState(operationState2, "action 1")
	buildFunctionRef(baseWorkflow, action2, "function 2")

	buildTransitionByState(operationState, operationState2, false)

	operationState3 := buildOperationState(baseWorkflow, "compensation next state 1")
	operationState3.BaseState.UsedForCompensation = true
	action3 := buildActionByOperationState(operationState3, "action 1")
	buildFunctionRef(baseWorkflow, action3, "function 3")

	operationState4 := buildOperationState(baseWorkflow, "compensation next state 2")
	operationState4.BaseState.UsedForCompensation = true
	action4 := buildActionByOperationState(operationState4, "action 1")
	buildFunctionRef(baseWorkflow, action4, "function 4")

	buildTransitionByState(operationState3, operationState4, false)

	operationState5 := buildOperationState(baseWorkflow, "compensation next state 3")
	buildEndByState(operationState5, true, false)
	operationState5.BaseState.UsedForCompensation = true
	action5 := buildActionByOperationState(operationState5, "action 5")
	buildFunctionRef(baseWorkflow, action5, "function 5")

	buildTransitionByState(operationState4, operationState5, false)

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				return *baseWorkflow.DeepCopy()
			},
		},
		{
			Desp: "state recursive",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.Transition.NextState = model.States[0].BaseState.Name
				return *model
			},
			Err: `workflow.states[0].transition.nextState can't no be recursive "start state"`,
		},
		{
			Desp: "exists",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.Transition.NextState = "invalid next state"
				return *model
			},
			Err: `workflow.states[0].transition.nextState don't exist "invalid next state"`,
		},
		{
			Desp: "transitionusedforcompensation",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[3].BaseState.UsedForCompensation = false
				return *model
			},
			Err: `Key: 'Workflow.States[2].BaseState.Transition.NextState' Error:Field validation for 'NextState' failed on the 'transitionusedforcompensation' tag
Key: 'Workflow.States[3].BaseState.Transition.NextState' Error:Field validation for 'NextState' failed on the 'transtionmainworkflow' tag`,
		},
		{
			Desp: "transtionmainworkflow",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].BaseState.Transition.NextState = model.States[3].BaseState.Name
				return *model
			},
			Err: `Key: 'Workflow.States[0].BaseState.Transition.NextState' Error:Field validation for 'NextState' failed on the 'transtionmainworkflow' tag`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestDataInputSchemaStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		// TODO Empty DataInoputSchema will have this instead nil:
		// 	&{Schema:{Type:0 StringValue: IntValue:0 FloatValue:0 MapValue:map[] SliceValue:[] BoolValue:false}
		// We can, make Schema pointer, or, find a way to make all fields from Object as pointer.
		// Using Schema: FromNull does have the same effect than just not set it.
		//{
		//	Desp: "empty DataInputSchema",
		//	Model: func() Workflow {
		//		model := baseWorkflow.DeepCopy()
		//		model.DataInputSchema = &DataInputSchema{}
		//		return *model
		//	},
		//	Err: `workflow.dataInputSchema.schema is required`,
		//},
		{
			Desp: "filled Schema, default failOnValidationErrors",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.DataInputSchema = &DataInputSchema{
					Schema: FromString("sample schema"),
				}
				return *model
			},
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestSecretsStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "workflow secrets.name repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Secrets = []string{"secret 1", "secret 1"}
				return *model
			},
			Err: `workflow.secrets has duplicate value`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestErrorStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	baseWorkflow.BaseWorkflow.Errors = Errors{{
		Name: "error 1",
	}, {
		Name: "error 2",
	}}

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
				model.Errors[0].Name = ""
				return *model
			},
			Err: `workflow.errors[0].name is required`,
		},
		{
			Desp: "repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Errors = Errors{model.Errors[0], model.Errors[0]}
				return *model
			},
			Err: `workflow.errors has duplicate "name"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

type ValidationCase struct {
	Desp  string
	Model func() Workflow
	Err   string
}

func StructLevelValidationCtx(t *testing.T, testCases []ValidationCase) {
	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			model := tc.Model()
			err := val.GetValidator().StructCtx(NewValidatorContext(&model), model)
			err = val.WorkflowError(err)
			if tc.Err != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.Err, err.Error())
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
