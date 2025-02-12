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
	"github.com/serverlessworkflow/sdk-go/v3/expr"
)

// RuntimeExpression represents a runtime expression.
type RuntimeExpression struct {
	Value string `json:"-" validate:"required"`
}

// NewRuntimeExpression is an alias for NewExpr
var NewRuntimeExpression = NewExpr

// NewExpr creates a new RuntimeExpression instance
func NewExpr(runtimeExpression string) *RuntimeExpression {
	return &RuntimeExpression{Value: runtimeExpression}
}

// IsValid checks if the RuntimeExpression value is valid, handling both with and without `${}`.
func (r *RuntimeExpression) IsValid() bool {
	return expr.IsValid(r.Value)
}

// UnmarshalJSON implements custom unmarshalling for RuntimeExpression.
func (r *RuntimeExpression) UnmarshalJSON(data []byte) error {
	// Decode the input as a string
	var raw string
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("failed to unmarshal RuntimeExpression: %w", err)
	}

	// Assign the value
	r.Value = raw

	// Validate the runtime expression
	if !r.IsValid() {
		return fmt.Errorf("invalid runtime expression format: %s", raw)
	}

	return nil
}

// MarshalJSON implements custom marshalling for RuntimeExpression.
func (r *RuntimeExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.Value)
}

func (r *RuntimeExpression) String() string {
	return r.Value
}

func (r *RuntimeExpression) GetValue() interface{} {
	return r.Value
}
