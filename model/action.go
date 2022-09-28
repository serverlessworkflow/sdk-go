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
)

// Action ...
type Action struct {
	// Unique action definition name
	Name        string       `json:"name,omitempty"`
	FunctionRef *FunctionRef `json:"functionRef,omitempty"`
	// References a 'trigger' and 'result' reusable event definitions
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
}

// FunctionRef ...
type FunctionRef struct {
	// Name of the referenced function
	RefName string `json:"refName" validate:"required"`
	// Function arguments
	Arguments map[string]interface{} `json:"arguments,omitempty"`
	// String containing a valid GraphQL selection set
	SelectionSet string `json:"selectionSet,omitempty"`
}

// UnmarshalJSON ...
func (f *FunctionRef) UnmarshalJSON(data []byte) error {
	funcRef := make(map[string]interface{})
	if err := json.Unmarshal(data, &funcRef); err != nil {
		f.RefName, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}

	f.RefName = requiresNotNilOrEmpty(funcRef["refName"])
	if _, found := funcRef["arguments"]; found {
		f.Arguments = funcRef["arguments"].(map[string]interface{})
	}
	f.SelectionSet = requiresNotNilOrEmpty(funcRef["selectionSet"])

	return nil
}

// WorkflowRef holds a reference for a workflow definition
type WorkflowRef struct {
	// Sub-workflow unique id
	WorkflowID string `json:"workflowId" validate:"required"`
	// Sub-workflow version
	Version string `json:"version,omitempty"`
}

// UnmarshalJSON ...
func (s *WorkflowRef) UnmarshalJSON(data []byte) error {
	subflowRef := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &subflowRef); err != nil {
		s.WorkflowID, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("version", subflowRef, &s.Version); err != nil {
		return err
	}
	if err := unmarshalKey("workflowId", subflowRef, &s.WorkflowID); err != nil {
		return err
	}

	return nil
}

// Sleep ...
type Sleep struct {
	// Before Amount of time (ISO 8601 duration format) to sleep before function/subflow invocation. Does not apply if 'eventRef' is defined.
	Before string `json:"before,omitempty"`
	// After Amount of time (ISO 8601 duration format) to sleep after function/subflow invocation. Does not apply if 'eventRef' is defined.
	After string `json:"after,omitempty"`
}
