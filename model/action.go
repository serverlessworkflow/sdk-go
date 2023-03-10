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
	"fmt"
)

// Action specify invocations of services or other workflows during workflow execution.
type Action struct {
	// Defines Unique action identifier.
	// +optional
	ID string `json:"id,omitempty"`
	// Defines Unique action name.
	// +optional
	Name string `json:"name,omitempty"`
	// References a reusable function definition.
	// +optional
	FunctionRef *FunctionRef `json:"functionRef,omitempty"`
	// References a 'trigger' and 'result' reusable event definitions.
	// +optional
	EventRef *EventRef `json:"eventRef,omitempty"`
	// References a workflow to be invoked.
	// +optional
	SubFlowRef *WorkflowRef `json:"subFlowRef,omitempty"`
	// Defines time period workflow execution should sleep before / after function execution.
	// +optional
	Sleep *Sleep `json:"sleep,omitempty"`
	// References a defined workflow retry definition. If not defined uses the default runtime retry definition.
	// +optional
	RetryRef string `json:"retryRef,omitempty"`
	// List of unique references to defined workflow errors for which the action should not be retried.
	// Used only when `autoRetries` is set to `true`
	// +optional
	NonRetryableErrors []string `json:"nonRetryableErrors,omitempty" validate:"omitempty,min=1"`
	// List of unique references to defined workflow errors for which the action should be retried.
	// Used only when `autoRetries` is set to `false`
	// +optional
	RetryableErrors []string `json:"retryableErrors,omitempty" validate:"omitempty,min=1"`
	// Filter the state data to select only the data that can be used within function definition arguments
	// using its fromStateData property. Filter the action results to select only the result data that should
	// be added/merged back into the state data using its results property. Select the part of state data which
	// the action data results should be added/merged to using the toStateData property.
	// +optional
	ActionDataFilter ActionDataFilter `json:"actionDataFilter,omitempty"`
	// Expression, if defined, must evaluate to true for this action to be performed. If false, action is disregarded.
	// +optional
	Condition string `json:"condition,omitempty"`
}

type actionForUnmarshal Action

// UnmarshalJSON implements json.Unmarshaler
func (a *Action) UnmarshalJSON(data []byte) error {
	v := actionForUnmarshal{
		ActionDataFilter: ActionDataFilter{UseResults: true},
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("action value '%s' is not supported, it must be an object or string", string(data))
	}
	*a = Action(v)

	return nil
}

// FunctionRef defines the reference to a reusable function definition
type FunctionRef struct {
	// Name of the referenced function.
	// +kubebuilder:validation:Required
	RefName string `json:"refName" validate:"required"`
	// Arguments (inputs) to be passed to the referenced function
	// +optional
	// TODO: validate it as required if function type is graphql
	Arguments map[string]Object `json:"arguments,omitempty"`
	// Used if function type is graphql. String containing a valid GraphQL selection set.
	// TODO: validate it as required if function type is graphql
	// +optional
	SelectionSet string `json:"selectionSet,omitempty"`
	// Specifies if the function should be invoked sync or async. Default is sync.
	// +kubebuilder:validation:Enum=async;sync
	// +kubebuilder:default=sync
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneof=async sync"`
}

// UnmarshalJSON implements json.Unmarshaler
func (f *FunctionRef) UnmarshalJSON(data []byte) error {
	type functionRefForUnmarshal FunctionRef
	functionRef, refName, err := primitiveOrStruct[string, functionRefForUnmarshal]("functionRef", data)
	if err != nil {
		return err
	}

	if functionRef == nil {
		f.RefName = refName
	} else {
		*f = FunctionRef(*functionRef)
	}

	if f.Invoke == "" {
		f.Invoke = InvokeKindSync
	}

	return nil
}

// Sleep defines time periods workflow execution should sleep before & after function execution
type Sleep struct {
	// Defines amount of time (ISO 8601 duration format) to sleep before function/subflow invocation.
	// Does not apply if 'eventRef' is defined.
	// +optional
	Before string `json:"before,omitempty" validate:"omitempty,iso8601duration"`
	// Defines amount of time (ISO 8601 duration format) to sleep after function/subflow invocation.
	// Does not apply if 'eventRef' is defined.
	// +optional
	After string `json:"after,omitempty" validate:"omitempty,iso8601duration"`
}
