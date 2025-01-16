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
	"fmt"
)

// WaitTask represents a task configuration to delay execution for a specified duration.
type WaitTask struct {
	TaskBase `json:",inline"`
	Wait     *Duration `json:"wait" validate:"required"`
}

// MarshalJSON for WaitTask to ensure proper serialization.
func (wt *WaitTask) MarshalJSON() ([]byte, error) {
	type Alias WaitTask
	waitData, err := json.Marshal(wt.Wait)
	if err != nil {
		return nil, err
	}

	alias := struct {
		Alias
		Wait json.RawMessage `json:"wait"`
	}{
		Alias: (Alias)(*wt),
		Wait:  waitData,
	}

	return json.Marshal(alias)
}

// UnmarshalJSON for WaitTask to ensure proper deserialization.
func (wt *WaitTask) UnmarshalJSON(data []byte) error {
	type Alias WaitTask
	alias := struct {
		*Alias
		Wait json.RawMessage `json:"wait"`
	}{
		Alias: (*Alias)(wt),
	}

	// Unmarshal data into alias
	if err := json.Unmarshal(data, &alias); err != nil {
		return fmt.Errorf("failed to unmarshal WaitTask: %w", err)
	}

	// Unmarshal Wait field
	if err := json.Unmarshal(alias.Wait, &wt.Wait); err != nil {
		return fmt.Errorf("failed to unmarshal Wait field: %w", err)
	}

	return nil
}
