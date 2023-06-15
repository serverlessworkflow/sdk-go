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

func buildStateExecTimeoutByTimeouts(timeouts *Timeouts, total string) *StateExecTimeout {
	stateExecTimeout := StateExecTimeout{
		Total: total,
	}
	timeouts.StateExecTimeout = &stateExecTimeout
	return timeouts.StateExecTimeout
}

func buildStateExecTimeoutBySleepStateTimeout(timeouts *SleepStateTimeout, total string) *StateExecTimeout {
	stateExecTimeout := StateExecTimeout{
		Total: total,
	}
	timeouts.StateExecTimeout = &stateExecTimeout
	return timeouts.StateExecTimeout
}

func buildStateExecTimeoutByOperationStateTimeout(timeouts *OperationStateTimeout, total string) *StateExecTimeout {
	stateExecTimeout := StateExecTimeout{
		Total: total,
	}
	timeouts.StateExecTimeout = &stateExecTimeout
	return timeouts.StateExecTimeout
}

func TestStateExecTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	timeouts := buildTimeouts(baseWorkflow)
	buildStateExecTimeoutByTimeouts(timeouts, "PT5S")

	callbackState := buildCallbackState(baseWorkflow, "start state", "event 1")
	buildEndByState(callbackState, true, false)
	buildCallbackStateTimeout(callbackState.CallbackState)
	buildFunctionRef(baseWorkflow, &callbackState.Action, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.BaseWorkflow.Timeouts.StateExecTimeout.Single = "PT5S"
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.BaseWorkflow.Timeouts.StateExecTimeout.Total = ""
				return *model
			},
			Err: `workflow.timeouts.stateExecTimeout.total is required`,
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.BaseWorkflow.Timeouts.StateExecTimeout.Single = "P5S"
				model.BaseWorkflow.Timeouts.StateExecTimeout.Total = "P5S"
				return *model
			},
			Err: `workflow.timeouts.stateExecTimeout.single invalid iso8601 duration "P5S"
workflow.timeouts.stateExecTimeout.total invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
