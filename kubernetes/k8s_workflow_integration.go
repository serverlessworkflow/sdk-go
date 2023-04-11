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

// States should be objects that will be in the same array even if it belongs to
// different types. An issue similar to the below will happen when trying to deploy your custom CR:
// strict decoding error: unknown field "spec.states[0].dataConditions"
// To make the CRD is compliant to the specs there are two options,
// a flat struct with all states fields at the same level,
// or use the // +kubebuilder:pruning:PreserveUnknownFields
// kubebuilder validator and delegate the validation  to the sdk-go validator using the admission webhook.
// TODO add a webhook example

// ServerlessWorkflowSpec defines a base API for integration test with operator-sdk
type ServerlessWorkflowSpec struct {
	model.Workflow `json:",inline"`
}

// ServerlessWorkflow ...
// +kubebuilder:object:root=true
// +kubebuilder:object:generate=true
// +kubebuilder:subresource:status
// +k8s:openapi-gen=true
type ServerlessWorkflow struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ServerlessWorkflowSpec `json:"spec,omitempty"`
	Status string                 `json:"status,omitempty"`
}

// ServerlessWorkflowList contains a list of SDKServerlessWorkflow
// +kubebuilder:object:root=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type ServerlessWorkflowList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ServerlessWorkflow `json:"items"`
}

func (S ServerlessWorkflowList) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

func (S ServerlessWorkflow) DeepCopyObject() runtime.Object {
	//TODO implement me
	panic("implement me")
}

func init() {
	SchemeBuilder.Register(&ServerlessWorkflow{}, &ServerlessWorkflowList{})
}
