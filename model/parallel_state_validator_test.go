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

func buildParallelState(workflow *Workflow, name string) *State {
	state := State{
		BaseState: BaseState{
			Name: name,
			Type: StateTypeParallel,
		},
		ParallelState: &ParallelState{
			CompletionType: CompletionTypeAllOf,
		},
	}

	workflow.States = append(workflow.States, state)
	return &workflow.States[len(workflow.States)-1]
}

func buildBranch(state *State, name string) *Branch {
	branch := Branch{
		Name: name,
	}

	state.ParallelState.Branches = append(state.ParallelState.Branches, branch)
	return &state.ParallelState.Branches[len(state.ParallelState.Branches)-1]
}

func buildBranchTimeouts(branch *Branch) *BranchTimeouts {
	branch.Timeouts = &BranchTimeouts{}
	return branch.Timeouts
}

func buildParallelStateTimeout(state *State) *ParallelStateTimeout {
	state.ParallelState.Timeouts = &ParallelStateTimeout{
		BranchExecTimeout: "PT5S",
	}
	return state.ParallelState.Timeouts
}

func TestParallelStateStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	parallelState := buildParallelState(baseWorkflow, "start state")
	buildEndByState(parallelState, true, false)
	branch := buildBranch(parallelState, "brach 1")
	action1 := buildActionByBranch(branch, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success completionTypeAllOf",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "success completionTypeAtLeast",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.CompletionType = CompletionTypeAtLeast
				model.States[0].ParallelState.NumCompleted = intstr.FromInt(1)
				return *model
			},
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.CompletionType = CompletionTypeAtLeast + " invalid"
				return *model
			},
			Err: `workflow.states[0].parallelState.completionType need by one of [allOf atLeast]`,
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches = nil
				model.States[0].ParallelState.CompletionType = ""
				return *model
			},
			Err: `workflow.states[0].parallelState.branches is required
workflow.states[0].parallelState.completionType is required`,
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches = []Branch{}
				return *model
			},
			Err: `workflow.states[0].parallelState.branches min > 1`,
		},
		{
			Desp: "required numCompleted",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.CompletionType = CompletionTypeAtLeast
				return *model
			},
			Err: `Key: 'Workflow.States[0].ParallelState.NumCompleted' Error:Field validation for 'NumCompleted' failed on the 'gt0' tag`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestBranchStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	parallelState := buildParallelState(baseWorkflow, "start state")
	buildEndByState(parallelState, true, false)
	branch := buildBranch(parallelState, "brach 1")
	action1 := buildActionByBranch(branch, "action 1")
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
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches[0].Name = ""
				model.States[0].ParallelState.Branches[0].Actions = nil
				return *model
			},
			Err: `workflow.states[0].parallelState.branches[0].name is required
workflow.states[0].parallelState.branches[0].actions is required`,
		},
		{
			Desp: "min",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches[0].Actions = []Action{}
				return *model
			},
			Err: `workflow.states[0].parallelState.branches[0].actions min > 1`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestBranchTimeoutsStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	parallelState := buildParallelState(baseWorkflow, "start state")
	buildEndByState(parallelState, true, false)
	branch := buildBranch(parallelState, "brach 1")
	buildBranchTimeouts(branch)
	action1 := buildActionByBranch(branch, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches[0].Timeouts.ActionExecTimeout = "PT5S"
				model.States[0].ParallelState.Branches[0].Timeouts.BranchExecTimeout = "PT5S"
				return *model
			},
		},
		{
			Desp: "omitempty",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches[0].Timeouts.ActionExecTimeout = ""
				model.States[0].ParallelState.Branches[0].Timeouts.BranchExecTimeout = ""
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Branches[0].Timeouts.ActionExecTimeout = "P5S"
				model.States[0].ParallelState.Branches[0].Timeouts.BranchExecTimeout = "P5S"
				return *model
			},
			Err: `workflow.states[0].parallelState.branches[0].timeouts.actionExecTimeout invalid iso8601 duration "P5S"
workflow.states[0].parallelState.branches[0].timeouts.branchExecTimeout invalid iso8601 duration "P5S"`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestParallelStateTimeoutStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	parallelState := buildParallelState(baseWorkflow, "start state")
	buildParallelStateTimeout(parallelState)
	buildEndByState(parallelState, true, false)
	branch := buildBranch(parallelState, "brach 1")
	action1 := buildActionByBranch(branch, "action 1")
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
			Desp: "omitempty",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Timeouts.BranchExecTimeout = ""
				return *model
			},
		},
		{
			Desp: "iso8601duration",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.States[0].ParallelState.Timeouts.BranchExecTimeout = "P5S"
				return *model
			},
			Err: `workflow.states[0].parallelState.timeouts.branchExecTimeout invalid iso8601 duration "P5S"`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
