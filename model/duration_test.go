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
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDurationToTime_Nil(t *testing.T) {
	_, err := DurationToTime(nil, time.Now())
	assert.Error(t, err)
}

func TestDurationToTime_Inline(t *testing.T) {
	d, err := DurationToTime(&Duration{Value: DurationInline{Seconds: 1, Milliseconds: 500}}, time.Now())
	assert.NoError(t, err)
	assert.Equal(t, 1500*time.Millisecond, d)
}

func TestDurationToTime_NonISOExpressionRejected(t *testing.T) {
	_, err := DurationToTime(NewDurationExpr("250ms"), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_ISO8601Expression(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d, err := DurationToTime(NewDurationExpr("P1DT1H"), start)
	assert.NoError(t, err)
	assert.Equal(t, 25*time.Hour, d)
}

func TestDurationToTime_ISO8601MillisecondsExpression(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	d, err := DurationToTime(NewDurationExpr("PT1S250MS"), start)
	assert.NoError(t, err)
	assert.Equal(t, 1250*time.Millisecond, d)
}

func TestDurationToTime_YearExpressionRejected(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := DurationToTime(NewDurationExpr("P1Y"), start)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_WeekExpressionRejected(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := DurationToTime(NewDurationExpr("P1W"), start)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_FractionalDayExpressionRejected(t *testing.T) {
	start := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	_, err := DurationToTime(NewDurationExpr("P1.5D"), start)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_InvalidExpression(t *testing.T) {
	_, err := DurationToTime(NewDurationExpr("1Y"), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_InvalidBarePTExpression(t *testing.T) {
	_, err := DurationToTime(NewDurationExpr("PT"), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}

func TestDurationToTime_UnsupportedMonthExpressionRejected(t *testing.T) {
	_, err := DurationToTime(NewDurationExpr("P1M"), time.Now())
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid duration expression")
}
