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
	"fmt"
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

// Duration represents a flexible duration that can be either inline or an ISO 8601 expression.
type Duration struct {
	Value interface{} `json:"-"`
}

// NewDurationExpr accessor to create a Duration object from a string
func NewDurationExpr(durationExpression string) *Duration {
	return &Duration{DurationExpression{durationExpression}}
}

func (d *Duration) AsExpression() string {
	switch v := d.Value.(type) {
	case string:
		return v
	case DurationExpression:
		return v.String()
	default:
		return ""
	}
}

func (d *Duration) AsInline() *DurationInline {
	switch v := d.Value.(type) {
	case DurationInline:
		return &v
	default:
		return nil
	}
}

// UnmarshalJSON for Duration to handle both inline and expression durations.
func (d *Duration) UnmarshalJSON(data []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err == nil {
		validKeys := map[string]bool{"days": true, "hours": true, "minutes": true, "seconds": true, "milliseconds": true}
		for key := range raw {
			if !validKeys[key] {
				return fmt.Errorf("unexpected key '%s' in duration object", key)
			}
		}

		inline := DurationInline{}
		if err := json.Unmarshal(data, &inline); err != nil {
			return fmt.Errorf("failed to unmarshal DurationInline: %w", err)
		}
		d.Value = inline
		return nil
	}

	var asString string
	if err := json.Unmarshal(data, &asString); err == nil {
		d.Value = DurationExpression{Expression: asString}
		return nil
	}

	return errors.New("data must be a valid duration string or object")
}

// MarshalJSON for Duration to handle both inline and expression durations.
func (d *Duration) MarshalJSON() ([]byte, error) {
	switch v := d.Value.(type) {
	case DurationInline:
		return json.Marshal(v)
	case DurationExpression:
		return json.Marshal(v.Expression)
	case string:
		durationExpression := &DurationExpression{Expression: v}
		return json.Marshal(durationExpression)
	default:
		return nil, errors.New("unknown Duration type")
	}
}

// DurationInline represents the inline definition of a duration.
type DurationInline struct {
	Days         int32 `json:"days,omitempty"`
	Hours        int32 `json:"hours,omitempty"`
	Minutes      int32 `json:"minutes,omitempty"`
	Seconds      int32 `json:"seconds,omitempty"`
	Milliseconds int32 `json:"milliseconds,omitempty"`
}

// MarshalJSON for DurationInline.
func (d *DurationInline) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"days":         d.Days,
		"hours":        d.Hours,
		"minutes":      d.Minutes,
		"seconds":      d.Seconds,
		"milliseconds": d.Milliseconds,
	})
}

// DurationExpression represents the ISO 8601 expression of a duration.
type DurationExpression struct {
	Expression string `json:"-" validate:"required,iso8601_duration"`
}

func (d *DurationExpression) String() string {
	return d.Expression
}

// MarshalJSON for DurationExpression.
func (d *DurationExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(d.Expression)
}

// UnmarshalJSON for DurationExpression to handle ISO 8601 strings.
func (d *DurationExpression) UnmarshalJSON(data []byte) error {
	var asString string
	if err := json.Unmarshal(data, &asString); err != nil {
		return err
	}
	d.Expression = asString
	return nil
}
