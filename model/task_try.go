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

type TryTask struct {
	TaskBase `json:",inline"`
	Try      *TaskList     `json:"try" validate:"required,dive"`
	Catch    *TryTaskCatch `json:"catch" validate:"required"`
}

type TryTaskCatch struct {
	Errors struct {
		With *ErrorFilter `json:"with,omitempty"`
	} `json:"errors,omitempty"`
	As         string             `json:"as,omitempty"`
	When       *RuntimeExpression `json:"when,omitempty"`
	ExceptWhen *RuntimeExpression `json:"exceptWhen,omitempty"`
	Retry      *RetryPolicy       `json:"retry,omitempty"`
	Do         *TaskList          `json:"do,omitempty" validate:"omitempty,dive"`
}

// RetryPolicy defines a retry policy.
type RetryPolicy struct {
	When       *RuntimeExpression `json:"when,omitempty"`
	ExceptWhen *RuntimeExpression `json:"exceptWhen,omitempty"`
	Delay      *Duration          `json:"delay,omitempty"`
	Backoff    *RetryBackoff      `json:"backoff,omitempty"`
	Limit      RetryLimit         `json:"limit,omitempty"`
	Jitter     *RetryPolicyJitter `json:"jitter,omitempty"`
	Ref        string             `json:"-"` // Reference to a reusable retry policy
}

// MarshalJSON for RetryPolicy to ensure proper serialization.
func (rp *RetryPolicy) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		When       *RuntimeExpression `json:"when,omitempty"`
		ExceptWhen *RuntimeExpression `json:"exceptWhen,omitempty"`
		Delay      *Duration          `json:"delay,omitempty"`
		Backoff    *RetryBackoff      `json:"backoff,omitempty"`
		Limit      RetryLimit         `json:"limit,omitempty"`
		Jitter     *RetryPolicyJitter `json:"jitter,omitempty"`
	}{
		When:       rp.When,
		ExceptWhen: rp.ExceptWhen,
		Delay:      rp.Delay,
		Backoff:    rp.Backoff,
		Limit:      rp.Limit,
		Jitter:     rp.Jitter,
	})
}

// UnmarshalJSON for RetryPolicy to ensure proper deserialization.
func (rp *RetryPolicy) UnmarshalJSON(data []byte) error {
	var raw interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal RetryPolicy: %w", err)
	}

	switch v := raw.(type) {
	case string:
		// If it's a string, treat it as a reference
		rp.Ref = v
	case map[string]interface{}:
		// If it's an object, unmarshal into the struct
		type Alias RetryPolicy
		alias := &struct {
			*Alias
		}{
			Alias: (*Alias)(rp),
		}
		if err := json.Unmarshal(data, alias); err != nil {
			return fmt.Errorf("failed to unmarshal RetryPolicy object: %w", err)
		}
	default:
		return fmt.Errorf("invalid RetryPolicy type: %T", v)
	}

	return nil
}

func (rp *RetryPolicy) ResolveReference(retries map[string]*RetryPolicy) error {
	if rp.Ref == "" {
		// No reference to resolve
		return nil
	}

	resolved, exists := retries[rp.Ref]
	if !exists {
		return fmt.Errorf("retry policy reference %q not found", rp.Ref)
	}

	// Copy resolved policy fields into the current RetryPolicy
	*rp = *resolved
	rp.Ref = "" // Clear the reference to avoid confusion

	return nil
}

func ResolveRetryPolicies(tasks []TryTaskCatch, retries map[string]*RetryPolicy) error {
	for i := range tasks {
		if tasks[i].Retry != nil {
			if err := tasks[i].Retry.ResolveReference(retries); err != nil {
				return fmt.Errorf("failed to resolve retry policy for task %q: %w", tasks[i].As, err)
			}
		}
	}
	return nil
}

// RetryBackoff defines the retry backoff strategies.
type RetryBackoff struct {
	Constant    *BackoffDefinition `json:"constant,omitempty"`
	Exponential *BackoffDefinition `json:"exponential,omitempty"`
	Linear      *BackoffDefinition `json:"linear,omitempty"`
}

// MarshalJSON for RetryBackoff to ensure oneOf behavior.
func (rb *RetryBackoff) MarshalJSON() ([]byte, error) {
	switch {
	case rb.Constant != nil:
		return json.Marshal(map[string]interface{}{"constant": rb.Constant.Definition})
	case rb.Exponential != nil:
		return json.Marshal(map[string]interface{}{"exponential": rb.Exponential.Definition})
	case rb.Linear != nil:
		return json.Marshal(map[string]interface{}{"linear": rb.Linear.Definition})
	default:
		return nil, errors.New("RetryBackoff must have one of 'constant', 'exponential', or 'linear' defined")
	}
}

func (rb *RetryBackoff) UnmarshalJSON(data []byte) error {
	var raw map[string]json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal RetryBackoff: %w", err)
	}

	if rawConstant, ok := raw["constant"]; ok {
		rb.Constant = &BackoffDefinition{}
		if err := json.Unmarshal(rawConstant, &rb.Constant.Definition); err != nil {
			return fmt.Errorf("failed to unmarshal constant backoff: %w", err)
		}
		return nil
	}

	if rawExponential, ok := raw["exponential"]; ok {
		rb.Exponential = &BackoffDefinition{}
		if err := json.Unmarshal(rawExponential, &rb.Exponential.Definition); err != nil {
			return fmt.Errorf("failed to unmarshal exponential backoff: %w", err)
		}
		return nil
	}

	if rawLinear, ok := raw["linear"]; ok {
		rb.Linear = &BackoffDefinition{}
		if err := json.Unmarshal(rawLinear, &rb.Linear.Definition); err != nil {
			return fmt.Errorf("failed to unmarshal linear backoff: %w", err)
		}
		return nil
	}

	return errors.New("RetryBackoff must have one of 'constant', 'exponential', or 'linear' defined")
}

type BackoffDefinition struct {
	Definition map[string]interface{} `json:"definition,omitempty"`
}

// RetryLimit defines the retry limit configurations.
type RetryLimit struct {
	Attempt  *RetryLimitAttempt `json:"attempt,omitempty"`
	Duration *Duration          `json:"duration,omitempty"`
}

// RetryLimitAttempt defines the limit for each retry attempt.
type RetryLimitAttempt struct {
	Count    int       `json:"count,omitempty"`
	Duration *Duration `json:"duration,omitempty"`
}

// RetryPolicyJitter defines the randomness or variability of retry delays.
type RetryPolicyJitter struct {
	From *Duration `json:"from" validate:"required"`
	To   *Duration `json:"to" validate:"required"`
}
