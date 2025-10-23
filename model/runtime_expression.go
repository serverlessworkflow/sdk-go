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
	"strings"

	"github.com/itchyny/gojq"
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

// IsStrictExpr returns true if the string is enclosed in `${ }`
func IsStrictExpr(expression string) bool {
	return strings.HasPrefix(expression, "${") && strings.HasSuffix(expression, "}")
}

// ContainsExpr returns true if the string contains `${` and `}`
func ContainsExpr(expression string) bool {
	return strings.Contains(expression, "${") && strings.Contains(expression, "}")
}

// SanitizeExpr processes the expression to ensure it's ready for evaluation
// It removes `${}` if present and replaces single quotes with double quotes
func SanitizeExpr(expression string) string {
	// Remove `${}` enclosure if present
	if IsStrictExpr(expression) {
		expression = strings.TrimSpace(expression[2 : len(expression)-1])
	}

	// Replace single quotes with double quotes
	expression = strings.ReplaceAll(expression, "'", "\"")

	return expression
}

func IsValidExpr(expression string) bool {
	expression = SanitizeExpr(expression)
	_, err := gojq.Parse(expression)
	return err == nil
}

// NormalizeExpr adds ${} to the given string
func NormalizeExpr(expr string) string {
	if strings.HasPrefix(expr, "${") {
		return expr
	}
	return fmt.Sprintf("${%s}", expr)
}

// IsValid checks if the RuntimeExpression value is valid, handling both with and without `${}`.
func (r *RuntimeExpression) IsValid() bool {
	return IsValidExpr(r.Value)
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
