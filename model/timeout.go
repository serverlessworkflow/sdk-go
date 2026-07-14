// Copyright 2025 The Open Workflow Specification Authors
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

// Timeout specifies a time limit for tasks or workflows.
type Timeout struct {
	// After The duration after which to timeout
	After *Duration `json:"after" validate:"required"`
}

// UnmarshalJSON implements custom unmarshalling for Timeout.
func (t *Timeout) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Check if "after" key exists
	afterData, ok := raw["after"]
	if !ok {
		return errors.New("missing 'after' key in Timeout JSON")
	}

	// Unmarshal "after" using the Duration type
	if err := json.Unmarshal(afterData, &t.After); err != nil {
		return err
	}

	return nil
}

// MarshalJSON implements custom marshalling for Timeout.
func (t *Timeout) MarshalJSON() ([]byte, error) {
	// Check the type of t.After.Value
	switch v := t.After.Value.(type) {
	case DurationInline:
		// Serialize inline duration
		return json.Marshal(map[string]interface{}{
			"after": v,
		})
	case DurationExpression:
		// Serialize expression as a simple string
		return json.Marshal(map[string]string{
			"after": v.Expression,
		})
	case string:
		// Handle direct string values as DurationExpression
		return json.Marshal(map[string]string{
			"after": v,
		})
	default:
		return nil, errors.New("unknown Duration type in Timeout")
	}
}

// TimeoutOrReference handles either a Timeout definition or a reference (string).
type TimeoutOrReference struct {
	Timeout   *Timeout `json:"-" validate:"required_without=Ref"`
	Reference *string  `json:"-" validate:"required_without=Timeout"`
}

// UnmarshalJSON implements custom unmarshalling for TimeoutOrReference.
func (tr *TimeoutOrReference) UnmarshalJSON(data []byte) error {
	// Attempt to unmarshal as a Timeout
	var asTimeout Timeout
	if err := json.Unmarshal(data, &asTimeout); err == nil {
		tr.Timeout = &asTimeout
		tr.Reference = nil
		return nil
	}

	// Attempt to unmarshal as a string (reference)
	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		tr.Reference = &asString
		tr.Timeout = nil
		return nil
	}

	// If neither works, return an error
	return errors.New("invalid TimeoutOrReference: must be a Timeout or a string reference")
}

// MarshalJSON implements custom marshalling for TimeoutOrReference.
func (tr *TimeoutOrReference) MarshalJSON() ([]byte, error) {
	// Marshal as a Timeout if present
	if tr.Timeout != nil {
		return json.Marshal(tr.Timeout)
	}

	// Marshal as a string reference if present
	if tr.Reference != nil {
		return json.Marshal(tr.Reference)
	}

	return nil, errors.New("invalid TimeoutOrReference: neither Timeout nor Ref is set")
}
