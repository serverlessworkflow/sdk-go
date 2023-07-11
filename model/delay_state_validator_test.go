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

func buildDelayState(workflow *Workflow, name, timeDelay string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeDelay,
		},
		DelayState: &DelayState{
			TimeDelay: timeDelay,
		},
	}
	workflow.States = append(workflow.States, state)

	return &workflow.States[len(workflow.States)-1]
}

func TestDelayStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	delayState := buildDelayState(baseWorkflow, "start state", "PT5S")
	buildEndByState(delayState, true, false)

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
				model.States[0].DelayState.TimeDelay = ""
				return *model
			},
			Err: `workflow.states[0].delayState.timeDelay is required`,
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].DelayState.TimeDelay = "P5S"
				return *model
			},
			Err: `workflow.states[0].delayState.timeDelay invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
