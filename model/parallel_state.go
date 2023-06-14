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
	"encoding/json"

	"github.com/serverlessworkflow/sdk-go/v2/util"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// CompletionType define on how to complete branch execution.
type CompletionType string

func (i CompletionType) KindValues() []string {
	return []string{
		string(CompletionTypeAllOf),
		string(CompletionTypeAtLeast),
	}
}

func (i CompletionType) String() string {
	return string(i)
}

const (
	// CompletionTypeAllOf defines all branches must complete execution before the state can transition/end.
	CompletionTypeAllOf CompletionType = "allOf"
	// CompletionTypeAtLeast defines state can transition/end once at least the specified number of branches
	// have completed execution.
	CompletionTypeAtLeast CompletionType = "atLeast"
)

// ParallelState Consists of a number of states that are executed in parallel
type ParallelState struct {
	// List of branches for this parallel state.
	// +kubebuilder:validation:MinItems=1
	Branches []Branch `json:"branches" validate:"required,min=1,dive"`
	// Option types on how to complete branch execution. Defaults to `allOf`.
	// +kubebuilder:validation:Enum=allOf;atLeast
	// +kubebuilder:default=allOf
	CompletionType CompletionType `json:"completionType,omitempty" validate:"required,oneofkind"`
	// Used when branchCompletionType is set to atLeast to specify the least number of branches that must complete
	// in order for the state to transition/end.
	// +optional
	// TODO: change this field to unmarshal result as int
	NumCompleted intstr.IntOrString `json:"numCompleted,omitempty"`
	// State specific timeouts
	// +optional
	Timeouts *ParallelStateTimeout `json:"timeouts,omitempty"`
}

func (p *ParallelState) MarshalJSON() ([]byte, error) {
	type Alias ParallelState
	custom, err := json.Marshal(&struct {
		*Alias
		Timeouts *ParallelStateTimeout `json:"timeouts,omitempty"`
	}{
		Alias:    (*Alias)(p),
		Timeouts: p.Timeouts,
	})
	return custom, err
}

type parallelStateUnmarshal ParallelState

// UnmarshalJSON unmarshal ParallelState object from json bytes
func (ps *ParallelState) UnmarshalJSON(data []byte) error {
	ps.ApplyDefault()
	return util.UnmarshalObject("parallelState", data, (*parallelStateUnmarshal)(ps))
}

// ApplyDefault set the default values for Parallel State
func (ps *ParallelState) ApplyDefault() {
	ps.CompletionType = CompletionTypeAllOf
}

// Branch Definition
type Branch struct {
	// Branch name
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// Actions to be executed in this branch
	// +kubebuilder:validation:MinItems=1
	Actions []Action `json:"actions" validate:"min=1,dive"`
	// Branch specific timeout settings
	// +optional
	Timeouts *BranchTimeouts `json:"timeouts,omitempty"`
}

// BranchTimeouts defines the specific timeout settings for branch
type BranchTimeouts struct {
	// Single actions definition execution timeout duration (ISO 8601 duration format)
	// +optional
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
	// Single branch execution timeout duration (ISO 8601 duration format)
	// +optional
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}

// ParallelStateTimeout defines the specific timeout settings for parallel state
type ParallelStateTimeout struct {
	// Default workflow state execution timeout (ISO 8601 duration format)
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// Default single branch execution timeout (ISO 8601 duration format)
	// +optional
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
