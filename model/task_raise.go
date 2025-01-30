// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"errors"
)

// RaiseTask represents a task configuration to raise errors.
type RaiseTask struct {
	TaskBase `json:",inline"`       // Inline TaskBase fields
	Raise    RaiseTaskConfiguration `json:"raise" validate:"required"`
}

func (r *RaiseTask) GetBase() *TaskBase {
	return &r.TaskBase
}

type RaiseTaskConfiguration struct {
	Error RaiseTaskError `json:"error" validate:"required"`
}

type RaiseTaskError struct {
	Definition *Error  `json:"-"`
	Ref        *string `json:"-"`
}

// UnmarshalJSON for RaiseTaskError to enforce "oneOf" behavior.
func (rte *RaiseTaskError) UnmarshalJSON(data []byte) error {
	// Try to unmarshal into a string (Ref)
	var ref string
	if err := json.Unmarshal(data, &ref); err == nil {
		rte.Ref = &ref
		rte.Definition = nil
		return nil
	}

	// Try to unmarshal into an Error (Definition)
	var def Error
	if err := json.Unmarshal(data, &def); err == nil {
		rte.Definition = &def
		rte.Ref = nil
		return nil
	}

	// If neither worked, return an error
	return errors.New("invalid RaiseTaskError: data must be either a string (reference) or an object (definition)")
}

// MarshalJSON for RaiseTaskError to ensure proper serialization.
func (rte *RaiseTaskError) MarshalJSON() ([]byte, error) {
	if rte.Definition != nil {
		return json.Marshal(rte.Definition)
	}
	if rte.Ref != nil {
		return json.Marshal(*rte.Ref)
	}
	return nil, errors.New("invalid RaiseTaskError: neither 'definition' nor 'reference' is set")
}
