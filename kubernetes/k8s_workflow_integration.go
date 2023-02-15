// Copyright 2023 The Serverless Workflow Specification Authors
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

package kubernetes

import (
	"github.com/serverlessworkflow/sdk-go/v2/model"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// This package provides a very simple api for kubernetes operator to test the integration
// of the Serverless SDK-Go with operator-sdk controller-gen and deepcopy-gen tools.
// The purpose of this integration is to stop issues like below beforehand:
// github.com/serverlessworkflow/sdk-go/model/event.go:51:2: encountered struct field "" without JSON tag in type "Event"
// github.com/serverlessworkflow/sdk-go/model/states.go:66:12: unsupported AST kind *ast.InterfaceType

// ServerlessWorkflowSpec defines a base API for integration test with operator-sdk
type ServerlessWorkflowSpec struct {
	BaseWorkflow model.BaseWorkflow `json:"inline"`
	Events       []model.Event      `json:"events,omitempty"`
	Functions    []model.Function   `json:"functions,omitempty"`
	Retries      []model.Retry      `json:"retries,omitempty"`
	States       []model.State      `json:"states"`
}

// SDKServerlessWorkflow ...
// +kubebuilder:object:root=true
// +kubebuilder:object:generate=true
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
type SDKServerlessWorkflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerlessWorkflowSpec `json:"spec,omitempty"`
	Status string                 `json:"status,omitempty"`
}

// SDKServerlessWorkflowList contains a list of SDKServerlessWorkflow
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type SDKServerlessWorkflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SDKServerlessWorkflow `json:"items"`
}

func (S SDKServerlessWorkflowList) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

func (S SDKServerlessWorkflow) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

func init() {
	SchemeBuilder.Register(&SDKServerlessWorkflow{}, &SDKServerlessWorkflowList{})
}
