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
	"bytes"
	"encoding/json"
	"fmt"
)

// Action specify invocations of services or other workflows during workflow execution.
type Action struct {
	// ID defines Unique action identifier
	ID string `json:"id,omitempty"`
	// Name defines Unique action definition name
	Name string `json:"name,omitempty"`
	// FunctionRef references a reusable function definition
	FunctionRef *FunctionRef `json:"functionRef,omitempty"`
	// EventRef references a 'trigger' and 'result' reusable event definitions
	EventRef *EventRef `json:"eventRef,omitempty"`
	// References a sub-workflow to be executed
	SubFlowRef *WorkflowRef `json:"subFlowRef,omitempty"`
	// Sleep Defines time period workflow execution should sleep before / after function execution
	Sleep Sleep `json:"sleep,omitempty"`
	// RetryRef References a defined workflow retry definition. If not defined the default retry policy is assumed
	RetryRef string `json:"retryRef,omitempty"`
	// List of unique references to defined workflow errors for which the action should not be retried. Used only when `autoRetries` is set to `true`
	NonRetryableErrors []string `json:"nonRetryableErrors,omitempty" validate:"omitempty,min=1"`
	// List of unique references to defined workflow errors for which the action should be retried. Used only when `autoRetries` is set to `false`
	RetryableErrors []string `json:"retryableErrors,omitempty" validate:"omitempty,min=1"`
	// Action data filter
	ActionDataFilter ActionDataFilter `json:"actionDataFilter,omitempty"`
	// Workflow expression evaluated against state data. Must evaluate to true or false
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
		return err
	}
	*a = Action(v)
	return nil
}

// FunctionRef defines the reference to a reusable function definition
type FunctionRef struct {
	// Name of the referenced function
	RefName string `json:"refName" validate:"required"`
	// Function arguments
	// TODO: validate it as required if function type is graphql
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	// String containing a valid GraphQL selection set
	// TODO: validate it as required if function type is graphql
	SelectionSet string `json:"selectionSet,omitempty"`

	// Invoke specifies if the subflow should be invoked sync or async.
	// Defaults to sync.
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneof=async sync"`
}

type functionRefForUnmarshal FunctionRef

// UnmarshalJSON implements json.Unmarshaler
func (f *FunctionRef) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	var err error
	switch data[0] {
	case '"':
		f.RefName, err = unmarshalString(data)
		if err != nil {
			return err
		}
		f.Invoke = InvokeKindSync
		return nil
	case '{':
		v := functionRefForUnmarshal{
			Invoke: InvokeKindSync,
		}
		err = json.Unmarshal(data, &v)
		if err != nil {
			// TODO: replace the error message with correct type's name
			return err
		}
		*f = FunctionRef(v)
		return nil
	}

	return fmt.Errorf("functionRef value '%s' is not supported, it must be an object or string", string(data))
}

// Sleep defines time periods workflow execution should sleep before & after function execution
type Sleep struct {
	// Before defines amount of time (ISO 8601 duration format) to sleep before function/subflow invocation. Does not apply if 'eventRef' is defined.
	Before string `json:"before,omitempty" validate:"omitempty,iso8601duration"`
	// After defines amount of time (ISO 8601 duration format) to sleep after function/subflow invocation. Does not apply if 'eventRef' is defined.
	After string `json:"after,omitempty" validate:"omitempty,iso8601duration"`
}
