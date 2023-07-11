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

import "github.com/serverlessworkflow/sdk-go/v2/util"

// CompletionType define on how to complete branch execution.
type OnParentCompleteType string

func (i OnParentCompleteType) KindValues() []string {
	return []string{
		string(OnParentCompleteTypeTerminate),
		string(OnParentCompleteTypeContinue),
	}
}

func (i OnParentCompleteType) String() string {
	return string(i)
}

const (
	OnParentCompleteTypeTerminate OnParentCompleteType = "terminate"
	OnParentCompleteTypeContinue  OnParentCompleteType = "continue"
)

// WorkflowRef holds a reference for a workflow definition
type WorkflowRef struct {
	// Sub-workflow unique id
	// +kubebuilder:validation:Required
	WorkflowID string `json:"workflowId" validate:"required"`
	// Sub-workflow version
	// +optional
	Version string `json:"version,omitempty"`
	// Specifies if the subflow should be invoked sync or async.
	// Defaults to sync.
	// +kubebuilder:validation:Enum=async;sync
	// +kubebuilder:default=sync
	// +optional
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneofkind"`
	// onParentComplete specifies how subflow execution should behave when parent workflow completes if invoke
	// is 'async'. Defaults to terminate.
	// +kubebuilder:validation:Enum=terminate;continue
	// +kubebuilder:default=terminate
	OnParentComplete OnParentCompleteType `json:"onParentComplete,omitempty" validate:"required,oneofkind"`
}

type workflowRefUnmarshal WorkflowRef

// UnmarshalJSON implements json.Unmarshaler
func (s *WorkflowRef) UnmarshalJSON(data []byte) error {
	s.ApplyDefault()
	return util.UnmarshalPrimitiveOrObject("subFlowRef", data, &s.WorkflowID, (*workflowRefUnmarshal)(s))
}

// ApplyDefault set the default values for Workflow Ref
func (s *WorkflowRef) ApplyDefault() {
	s.Invoke = InvokeKindSync
	s.OnParentComplete = "terminate"
}
