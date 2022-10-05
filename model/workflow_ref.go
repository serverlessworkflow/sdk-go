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

// WorkflowRef holds a reference for a workflow definition
type WorkflowRef struct {
	// Sub-workflow unique id
	WorkflowID string `json:"workflowId" validate:"required"`
	// Sub-workflow version
	Version string `json:"version,omitempty"`

	// Invoke specifies if the subflow should be invoked sync or async.
	// Defaults to sync.
	Invoke InvokeKind `json:"invoke,omitempty" validate:"required,oneof=async sync"`

	// OnParantComplete specifies how subflow execution should behave when parent workflow completes if invoke is 'async'。
	// Defaults to terminate.
	OnParentComplete string `json:"onParentComplete,omitempty" validate:"required,oneof=terminate continue"`
}

type workflowRefForUnmarshal WorkflowRef

// UnmarshalJSON implements json.Unmarshaler
func (s *WorkflowRef) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	var err error
	switch data[0] {
	case '"':
		s.WorkflowID, err = unmarshalString(data)
		if err != nil {
			return err
		}
		s.Invoke, s.OnParentComplete = InvokeKindSync, "terminate"
		return nil
	case '{':
		v := workflowRefForUnmarshal{
			Invoke:           InvokeKindSync,
			OnParentComplete: "terminate",
		}
		err = json.Unmarshal(data, &v)
		if err != nil {
			// TODO: replace the error message with correct type's name
			return err
		}
		*s = WorkflowRef(v)
		return nil
	}

	return fmt.Errorf("subFlowRef value '%s' is not supported, it must be an object or string", string(data))
}
