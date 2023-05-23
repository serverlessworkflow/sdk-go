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

import "fmt"

// StateExecTimeout defines workflow state execution timeout
type StateExecTimeout struct {
	// Single state execution timeout, not including retries (ISO 8601 duration format)
	// +optional
	Single string `json:"single,omitempty" validate:"omitempty,iso8601duration"`
	// Total state execution timeout, including retries (ISO 8601 duration format)
	// +kubebuilder:validation:Required
	Total string `json:"total" validate:"required,iso8601duration"`
}

func (s StateExecTimeout) String() string {
	return fmt.Sprintf("{ Single:%s, Total:%s}", s.Single, s.Total)
}

type stateExecTimeoutUnmarshal StateExecTimeout

// UnmarshalJSON unmarshal StateExecTimeout object from json bytes
func (s *StateExecTimeout) UnmarshalJSON(data []byte) error {
	return unmarshalPrimitiveOrObject("stateExecTimeout", data, &s.Total, (*stateExecTimeoutUnmarshal)(s))
}
