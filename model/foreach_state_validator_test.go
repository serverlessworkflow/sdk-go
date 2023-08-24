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

	"k8s.io/apimachinery/pkg/util/intstr"
)

func buildForEachState(workflow *Workflow, name string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeForEach,
		},
		ForEachState: &ForEachState{
			InputCollection: "3",
			Mode:            ForEachModeTypeSequential,
		},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func TestForEachStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	forEachState := buildForEachState(baseWorkflow, "start state")
	buildEndByState(forEachState, true, false)
	action1 := buildActionByForEachState(forEachState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.Mode = ForEachModeTypeParallel
				model.States[0].ForEachState.BatchSize = &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 1,
				}
				return *model
			},
		},
		{
			Desp: "success without batch size",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.Mode = ForEachModeTypeParallel
				model.States[0].ForEachState.BatchSize = nil
				return *model
			},
		},
		{
			Desp: "gt0 int",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.Mode = ForEachModeTypeParallel
				model.States[0].ForEachState.BatchSize = &intstr.IntOrString{
					Type:   intstr.Int,
					IntVal: 0,
				}
				return *model
			},
			Err: `workflow.states[0].forEachState.batchSize must be greater than 0`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.Mode = ForEachModeTypeParallel + "invalid"
				return *model
			},
			Err: `workflow.states[0].forEachState.mode need by one of [sequential parallel]`,
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.InputCollection = ""
				model.States[0].ForEachState.Mode = ""
				model.States[0].ForEachState.Actions = nil
				return *model
			},
			Err: `workflow.states[0].forEachState.inputCollection is required
workflow.states[0].forEachState.actions is required
workflow.states[0].forEachState.mode is required`,
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ForEachState.Actions = []Action{}
				return *model
			},
			Err: ``,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestForEachStateTimeoutStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase{}
	StructLevelValidationCtx(t, testCases)
}
