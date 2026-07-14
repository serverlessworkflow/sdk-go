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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTimeout_UnmarshalJSON(t *testing.T) {
	// Test cases for Timeout unmarshalling
	tests := []struct {
		name    string
		jsonStr string
		expect  *Timeout
		err     bool
	}{
		{
			name:    "Valid inline duration",
			jsonStr: `{"after": {"days": 1, "hours": 2}}`,
			expect: &Timeout{
				After: &Duration{DurationInline{
					Days:  1,
					Hours: 2,
				}},
			},
			err: false,
		},
		{
			name:    "Valid ISO 8601 duration",
			jsonStr: `{"after": "P3DT4H5M6S250MS"}`,
			expect: &Timeout{
				After: NewDurationExpr("P3DT4H5M6S250MS"),
			},
			err: false,
		},
		{
			name:    "Invalid duration type",
			jsonStr: `{"after": {"unknown": "value"}}`,
			expect:  nil,
			err:     true,
		},
		{
			name:    "Invalid empty inline duration",
			jsonStr: `{"after": {}}`,
			expect:  nil,
			err:     true,
		},
		{
			name:    "Invalid non-ISO expression",
			jsonStr: `{"after": "10s"}`,
			expect:  nil,
			err:     true,
		},
		{
			name:    "Unsupported ISO unit",
			jsonStr: `{"after": "P1Y"}`,
			expect:  nil,
			err:     true,
		},
		{
			name:    "Missing after key",
			jsonStr: `{}`,
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var timeout Timeout
			err := json.Unmarshal([]byte(test.jsonStr), &timeout)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expect, &timeout)
			}
		})
	}
}

func TestTimeout_MarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		input    *Timeout
		expected string
		wantErr  bool
	}{
		{
			name: "ISO 8601 Duration",
			input: &Timeout{
				After: &Duration{
					Value: DurationExpression{Expression: "PT1H"},
				},
			},
			expected: `{"after":"PT1H"}`,
			wantErr:  false,
		},
		{
			name: "Inline Duration",
			input: &Timeout{
				After: &Duration{
					Value: DurationInline{
						Days:    1,
						Hours:   2,
						Minutes: 30,
					},
				},
			},
			expected: `{"after":{"days":1,"hours":2,"minutes":30}}`,
			wantErr:  false,
		},
		{
			name:     "Invalid Duration",
			input:    &Timeout{After: &Duration{Value: 123}},
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, tt.expected, string(data))
			}
		})
	}
}

func TestTimeoutOrReference_UnmarshalJSON(t *testing.T) {
	// Test cases for TimeoutOrReference unmarshalling
	tests := []struct {
		name    string
		jsonStr string
		expect  *TimeoutOrReference
		err     bool
	}{
		{
			name:    "Valid Timeout",
			jsonStr: `{"after": {"days": 1, "hours": 2}}`,
			expect: &TimeoutOrReference{
				Timeout: &Timeout{
					After: &Duration{DurationInline{
						Days:  1,
						Hours: 2,
					}},
				},
			},
			err: false,
		},
		{
			name:    "Valid Ref",
			jsonStr: `"some-timeout-reference"`,
			expect: &TimeoutOrReference{
				Reference: ptrString("some-timeout-reference"),
			},
			err: false,
		},
		{
			name:    "Invalid JSON",
			jsonStr: `{"invalid": }`,
			expect:  nil,
			err:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var tor TimeoutOrReference
			err := json.Unmarshal([]byte(test.jsonStr), &tor)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, test.expect, &tor)
			}
		})
	}
}

func ptrString(s string) *string {
	return &s
}

func TestTimeoutOrReference_MarshalJSON(t *testing.T) {
	// Test cases for TimeoutOrReference marshalling
	tests := []struct {
		name   string
		input  *TimeoutOrReference
		expect string
		err    bool
	}{
		{
			name: "Valid Timeout",
			input: &TimeoutOrReference{
				Timeout: &Timeout{
					After: &Duration{DurationInline{
						Days:  1,
						Hours: 2,
					}},
				},
			},
			expect: `{"after":{"days":1,"hours":2}}`,
			err:    false,
		},
		{
			name: "Valid Ref",
			input: &TimeoutOrReference{
				Reference: ptrString("some-timeout-reference"),
			},
			expect: `"some-timeout-reference"`,
			err:    false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			data, err := json.Marshal(test.input)
			if test.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.JSONEq(t, test.expect, string(data))
			}
		})
	}
}

func TestDuration_AsTimeDuration(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)

	t.Run("Inline", func(t *testing.T) {
		d := &Duration{Value: DurationInline{Seconds: 2, Milliseconds: 500}}
		result, err := d.AsTimeDuration(start)
		assert.NoError(t, err)
		assert.Equal(t, 2500*time.Millisecond, result)
	})

	t.Run("NonISOExpressionRejected", func(t *testing.T) {
		d := NewDurationExpr("150ms")
		_, err := d.AsTimeDuration(start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration expression")
	})

	t.Run("ISO8601Expression", func(t *testing.T) {
		d := NewDurationExpr("P1DT1H")
		result, err := d.AsTimeDuration(start)
		assert.NoError(t, err)
		assert.Equal(t, 25*time.Hour, result)
	})

	t.Run("ISO8601MillisecondsExpression", func(t *testing.T) {
		d := NewDurationExpr("PT2S500MS")
		result, err := d.AsTimeDuration(start)
		assert.NoError(t, err)
		assert.Equal(t, 2500*time.Millisecond, result)
	})

	t.Run("ISO8601WeekExpressionRejected", func(t *testing.T) {
		d := NewDurationExpr("P1W")
		_, err := d.AsTimeDuration(start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration expression")
	})

	t.Run("ISO8601FractionalDayExpressionRejected", func(t *testing.T) {
		d := NewDurationExpr("P1.5D")
		_, err := d.AsTimeDuration(start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration expression")
	})

	t.Run("InvalidExpression", func(t *testing.T) {
		d := NewDurationExpr("1Y")
		_, err := d.AsTimeDuration(start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration expression")
	})

	t.Run("UnsupportedMonthExpressionRejected", func(t *testing.T) {
		d := NewDurationExpr("P1M")
		_, err := d.AsTimeDuration(start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid duration expression")
	})
}
