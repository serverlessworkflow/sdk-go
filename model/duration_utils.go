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
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"
)

type parsedISO8601Duration struct {
	Days            float64
	HasDays         bool
	Hours           float64
	HasHours        bool
	Minutes         float64
	HasMinutes      bool
	Seconds         float64
	HasSeconds      bool
	Milliseconds    float64
	HasMilliseconds bool
}

// DurationToTime converts a model Duration to time.Duration.
// It supports inline duration objects and ISO 8601 expressions.
func DurationToTime(duration *Duration, now time.Time) (time.Duration, error) {
	if duration == nil {
		return 0, fmt.Errorf("duration is nil")
	}

	if inline := duration.AsInline(); inline != nil {
		return inlineDurationToTime(inline)
	}

	expression := strings.TrimSpace(duration.AsExpression())
	if expression == "" {
		return 0, fmt.Errorf("duration expression is empty")
	}

	return iso8601DurationToTime(expression, now)
}

func inlineDurationToTime(inline *DurationInline) (time.Duration, error) {
	if inline.Days < 0 || inline.Hours < 0 || inline.Minutes < 0 || inline.Seconds < 0 || inline.Milliseconds < 0 {
		return 0, fmt.Errorf("duration fields must be non-negative")
	}

	duration := (time.Duration(inline.Days) * 24 * time.Hour) +
		(time.Duration(inline.Hours) * time.Hour) +
		(time.Duration(inline.Minutes) * time.Minute) +
		(time.Duration(inline.Seconds) * time.Second) +
		(time.Duration(inline.Milliseconds) * time.Millisecond)

	return duration, nil
}

func iso8601DurationToTime(expression string, now time.Time) (time.Duration, error) {
	_ = now // Reserved for API compatibility.

	parsed, err := parseAndValidateISO8601Duration(expression)
	if err != nil {
		return 0, err
	}

	fixedDuration := time.Duration(0)
	if err := addFixedUnitDuration(&fixedDuration, parsed.Days, 24*time.Hour, "days"); err != nil {
		return 0, err
	}
	if err := addFixedUnitDuration(&fixedDuration, parsed.Hours, time.Hour, "hours"); err != nil {
		return 0, err
	}
	if err := addFixedUnitDuration(&fixedDuration, parsed.Minutes, time.Minute, "minutes"); err != nil {
		return 0, err
	}
	if err := addFixedUnitDuration(&fixedDuration, parsed.Seconds, time.Second, "seconds"); err != nil {
		return 0, err
	}
	if err := addFixedUnitDuration(&fixedDuration, parsed.Milliseconds, time.Millisecond, "milliseconds"); err != nil {
		return 0, err
	}

	return fixedDuration, nil
}

func addFixedUnitDuration(total *time.Duration, value float64, unit time.Duration, fieldName string) error {
	if value == 0 {
		return nil
	}

	delta := value * float64(unit)
	if math.IsNaN(delta) || math.IsInf(delta, 0) {
		return fmt.Errorf("invalid duration %s", fieldName)
	}

	*total += time.Duration(math.Round(delta))
	return nil
}

func parseISO8601DurationExpression(expression string) (parsedISO8601Duration, error) {
	invalidDuration := func(reason string) error {
		return fmt.Errorf("invalid duration expression '%s': %s", expression, reason)
	}

	if !iso8601DurationPattern.MatchString(expression) {
		return parsedISO8601Duration{}, invalidDuration("must match supported ISO-8601 duration format")
	}

	if expression == "P" || expression == "PT" {
		return parsedISO8601Duration{}, invalidDuration("must include at least one duration component")
	}

	matches := iso8601DurationPattern.FindStringSubmatch(expression)
	if len(matches) == 0 {
		return parsedISO8601Duration{}, invalidDuration("parser could not extract duration components")
	}
	if matches[1] == "" && matches[2] == "" && matches[3] == "" && matches[4] == "" && matches[5] == "" && matches[6] == "" {
		return parsedISO8601Duration{}, invalidDuration("must include at least one duration component")
	}

	if strings.Contains(expression, "T") && matches[3] == "" && matches[4] == "" && matches[5] == "" && matches[6] == "" {
		return parsedISO8601Duration{}, invalidDuration("time designator 'T' must be followed by H, M, S, or MS")
	}

	// Keep parsing aligned with iso8601DurationPattern capture groups:
	// 1=D, 2=full T-block, 3=H, 4=M, 5=S, 6=MS.
	days, hasDays, err := parseISO8601Part(matches[1], "D", "days")
	if err != nil {
		return parsedISO8601Duration{}, err
	}
	hours, hasHours, err := parseISO8601Part(matches[3], "H", "hours")
	if err != nil {
		return parsedISO8601Duration{}, err
	}
	minutes, hasMinutes, err := parseISO8601Part(matches[4], "M", "minutes")
	if err != nil {
		return parsedISO8601Duration{}, err
	}
	seconds, hasSeconds, err := parseISO8601Part(matches[5], "S", "seconds")
	if err != nil {
		return parsedISO8601Duration{}, err
	}
	milliseconds, hasMilliseconds, err := parseISO8601Part(matches[6], "MS", "milliseconds")
	if err != nil {
		return parsedISO8601Duration{}, err
	}

	return parsedISO8601Duration{
		Days:            days,
		HasDays:         hasDays,
		Hours:           hours,
		HasHours:        hasHours,
		Minutes:         minutes,
		HasMinutes:      hasMinutes,
		Seconds:         seconds,
		HasSeconds:      hasSeconds,
		Milliseconds:    milliseconds,
		HasMilliseconds: hasMilliseconds,
	}, nil
}

func parseISO8601Part(partWithSuffix string, suffix string, fieldName string) (float64, bool, error) {
	if partWithSuffix == "" {
		return 0, false, nil
	}

	part := partWithSuffix[:len(partWithSuffix)-len(suffix)]
	value, err := strconv.ParseFloat(part, 64)
	if err != nil {
		return 0, false, fmt.Errorf("invalid duration %s: %w", fieldName, err)
	}
	if value < 0 {
		return 0, false, fmt.Errorf("duration %s cannot be negative", fieldName)
	}
	return value, true, nil
}
