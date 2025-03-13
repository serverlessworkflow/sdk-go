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

package expr

import (
	"errors"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

// IsStrictExpr returns true if the string is enclosed in `${ }`
func IsStrictExpr(expression string) bool {
	return strings.HasPrefix(expression, "${") && strings.HasSuffix(expression, "}")
}

// Sanitize processes the expression to ensure it's ready for evaluation
// It removes `${}` if present and replaces single quotes with double quotes
func Sanitize(expression string) string {
	// Remove `${}` enclosure if present
	if IsStrictExpr(expression) {
		expression = strings.TrimSpace(expression[2 : len(expression)-1])
	}

	// Replace single quotes with double quotes
	expression = strings.ReplaceAll(expression, "'", "\"")

	return expression
}

// IsValid tries to parse and check if the given value is a valid expression
func IsValid(expression string) bool {
	expression = Sanitize(expression)
	_, err := gojq.Parse(expression)
	return err == nil
}

// TraverseAndEvaluate recursively processes and evaluates all expressions in a JSON-like structure
func TraverseAndEvaluate(node interface{}, input interface{}) (interface{}, error) {
	switch v := node.(type) {
	case map[string]interface{}:
		// Traverse map
		for key, value := range v {
			evaluatedValue, err := TraverseAndEvaluate(value, input)
			if err != nil {
				return nil, err
			}
			v[key] = evaluatedValue
		}
		return v, nil

	case []interface{}:
		// Traverse array
		for i, value := range v {
			evaluatedValue, err := TraverseAndEvaluate(value, input)
			if err != nil {
				return nil, err
			}
			v[i] = evaluatedValue
		}
		return v, nil

	case string:
		// Check if the string is a runtime expression (e.g., ${ .some.path })
		if IsStrictExpr(v) {
			return evaluateJQExpression(Sanitize(v), input)
		}
		return v, nil

	default:
		// Return other types as-is
		return v, nil
	}
}

// TODO: add support to variables see https://github.com/itchyny/gojq/blob/main/option_variables_test.go

// evaluateJQExpression evaluates a jq expression against a given JSON input
func evaluateJQExpression(expression string, input interface{}) (interface{}, error) {
	// Parse the sanitized jq expression
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jq expression: %s, error: %w", expression, err)
	}

	// Compile and evaluate the expression
	iter := query.Run(input)
	result, ok := iter.Next()
	if !ok {
		return nil, errors.New("no result from jq evaluation")
	}

	// Check if an error occurred during evaluation
	if err, isErr := result.(error); isErr {
		return nil, fmt.Errorf("jq evaluation error: %w", err)
	}

	return result, nil
}
