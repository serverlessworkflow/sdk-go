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

func buildSleepState(workflow *Workflow, name, duration string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeSleep,
		},
		SleepState: &SleepState{
			Duration: duration,
		},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func buildSleepStateTimeout(state *State) *SleepStateTimeout {
	state.SleepState.Timeouts = &SleepStateTimeout{}
	return state.SleepState.Timeouts
}

func TestSleepStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	sleepState := buildSleepState(baseWorkflow, "start state", "PT5S")
	buildEndByState(sleepState, true, false)

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
				model.States[0].SleepState.Duration = ""
				return *model
			},
			Err: `workflow.states[0].sleepState.duration is required`,
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].SleepState.Duration = "P5S"
				return *model
			},
			Err: `workflow.states[0].sleepState.duration invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestSleepStateTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	sleepState := buildSleepState(baseWorkflow, "start state", "PT5S")
	buildEndByState(sleepState, true, false)
	sleepStateTimeout := buildSleepStateTimeout(sleepState)
	buildStateExecTimeoutBySleepStateTimeout(sleepStateTimeout)

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
	}

	StructLevelValidationCtx(t, testCases)
}
