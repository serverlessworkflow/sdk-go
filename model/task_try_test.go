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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetryPolicy_MarshalJSON(t *testing.T) {
	retryPolicy := RetryPolicy{
		When:       &RuntimeExpression{"${someCondition}"},
		ExceptWhen: &RuntimeExpression{"${someOtherCondition}"},
		Delay:      NewDurationExpr("PT5S"),
		Backoff: &RetryBackoff{
			Exponential: &BackoffDefinition{
				Definition: map[string]interface{}{"factor": 2},
			},
		},
		Limit: RetryLimit{
			Attempt: &RetryLimitAttempt{
				Count:    3,
				Duration: NewDurationExpr("PT1M"),
			},
			Duration: NewDurationExpr("PT10M"),
		},
		Jitter: &RetryPolicyJitter{
			From: NewDurationExpr("PT1S"),
			To:   NewDurationExpr("PT3S"),
		},
	}

	data, err := json.Marshal(retryPolicy)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"when": "${someCondition}",
		"exceptWhen": "${someOtherCondition}",
		"delay": "PT5S",
		"backoff": {"exponential": {"factor": 2}},
		"limit": {
			"attempt": {"count": 3, "duration": "PT1M"},
			"duration": "PT10M"
		},
		"jitter": {"from": "PT1S", "to": "PT3S"}
	}`, string(data))
}

func TestRetryPolicy_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"when": "${someCondition}",
		"exceptWhen": "${someOtherCondition}",
		"delay": "PT5S",
		"backoff": {"exponential": {"factor": 2}},
		"limit": {
			"attempt": {"count": 3, "duration": "PT1M"},
			"duration": "PT10M"
		},
		"jitter": {"from": "PT1S", "to": "PT3S"}
	}`

	var retryPolicy RetryPolicy
	err := json.Unmarshal([]byte(jsonData), &retryPolicy)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{"${someCondition}"}, retryPolicy.When)
	assert.Equal(t, &RuntimeExpression{"${someOtherCondition}"}, retryPolicy.ExceptWhen)
	assert.Equal(t, NewDurationExpr("PT5S"), retryPolicy.Delay)
	assert.NotNil(t, retryPolicy.Backoff.Exponential)
	assert.Equal(t, map[string]interface{}{"factor": float64(2)}, retryPolicy.Backoff.Exponential.Definition)
	assert.Equal(t, 3, retryPolicy.Limit.Attempt.Count)
	assert.Equal(t, NewDurationExpr("PT1M"), retryPolicy.Limit.Attempt.Duration)
	assert.Equal(t, NewDurationExpr("PT10M"), retryPolicy.Limit.Duration)
	assert.Equal(t, NewDurationExpr("PT1S"), retryPolicy.Jitter.From)
	assert.Equal(t, NewDurationExpr("PT3S"), retryPolicy.Jitter.To)
}

func TestRetryPolicy_Validation(t *testing.T) {
	// Valid RetryPolicy
	retryPolicy := RetryPolicy{
		When:       &RuntimeExpression{"${someCondition}"},
		ExceptWhen: &RuntimeExpression{"${someOtherCondition}"},
		Delay:      NewDurationExpr("PT5S"),
		Backoff: &RetryBackoff{
			Constant: &BackoffDefinition{
				Definition: map[string]interface{}{"delay": 5},
			},
		},
		Limit: RetryLimit{
			Attempt: &RetryLimitAttempt{
				Count:    3,
				Duration: NewDurationExpr("PT1M"),
			},
			Duration: NewDurationExpr("PT10M"),
		},
		Jitter: &RetryPolicyJitter{
			From: NewDurationExpr("PT1S"),
			To:   NewDurationExpr("PT3S"),
		},
	}
	assert.NoError(t, validate.Struct(retryPolicy))

	// Invalid RetryPolicy (missing required fields in Jitter)
	invalidRetryPolicy := RetryPolicy{
		Jitter: &RetryPolicyJitter{
			From: NewDurationExpr("PT1S"),
		},
	}
	assert.Error(t, validate.Struct(invalidRetryPolicy))
}

func TestRetryPolicy_UnmarshalJSON_WithReference(t *testing.T) {
	retries := map[string]*RetryPolicy{
		"default": {
			Delay: &Duration{DurationInline{Seconds: 3}},
			Backoff: &RetryBackoff{
				Exponential: &BackoffDefinition{},
			},
			Limit: RetryLimit{
				Attempt: &RetryLimitAttempt{Count: 5},
			},
		},
	}

	jsonData := `{
		"retry": "default"
	}`

	var task TryTaskCatch
	err := json.Unmarshal([]byte(jsonData), &task)
	assert.NoError(t, err)

	// Resolve the reference
	err = task.Retry.ResolveReference(retries)
	assert.NoError(t, err)

	assert.Equal(t, retries["default"].Delay, task.Retry.Delay)
	assert.Equal(t, retries["default"].Backoff, task.Retry.Backoff)
	assert.Equal(t, retries["default"].Limit, task.Retry.Limit)
}

func TestRetryPolicy_UnmarshalJSON_Inline(t *testing.T) {
	jsonData := `{
		"retry": {
			"delay": { "seconds": 3 },
			"backoff": { "exponential": {} },
			"limit": { "attempt": { "count": 5 } }
		}
	}`

	var task TryTaskCatch
	err := json.Unmarshal([]byte(jsonData), &task)
	assert.NoError(t, err)

	assert.NotNil(t, task.Retry)
	assert.Equal(t, int32(3), task.Retry.Delay.AsInline().Seconds)
	assert.NotNil(t, task.Retry.Backoff.Exponential)
	assert.Equal(t, 5, task.Retry.Limit.Attempt.Count)
}
