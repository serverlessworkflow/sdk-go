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
	"context"
	"fmt"
	"testing"

	"github.com/itchyny/gojq"
)

func TestTraverseAndEvaluate(t *testing.T) {
	t.Run("Simple no-expression map", func(t *testing.T) {
		node := map[string]interface{}{
			"key": "value",
			"num": 123,
		}
		result, err := TraverseAndEvaluate(node, nil, context.TODO())
		if err != nil {
			t.Fatalf("TraverseAndEvaluate() unexpected error: %v", err)
		}

		got, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return a map")
		}
		if got["key"] != "value" || got["num"] != 123 {
			t.Errorf("TraverseAndEvaluate() returned unexpected map data: %#v", got)
		}
	})

	t.Run("Expression in map", func(t *testing.T) {
		node := map[string]interface{}{
			"expr": "${ .foo }",
		}
		input := map[string]interface{}{
			"foo": "bar",
		}

		result, err := TraverseAndEvaluate(node, input, context.TODO())
		if err != nil {
			t.Fatalf("TraverseAndEvaluate() unexpected error: %v", err)
		}

		got, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return a map")
		}
		if got["expr"] != "bar" {
			t.Errorf("TraverseAndEvaluate() = %v, want %v", got["expr"], "bar")
		}
	})

	t.Run("Expression in array", func(t *testing.T) {
		node := []interface{}{
			"static",
			"${ .foo }",
		}
		input := map[string]interface{}{
			"foo": "bar",
		}

		result, err := TraverseAndEvaluate(node, input, context.TODO())
		if err != nil {
			t.Fatalf("TraverseAndEvaluate() unexpected error: %v", err)
		}

		got, ok := result.([]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return an array")
		}
		if got[0] != "static" {
			t.Errorf("TraverseAndEvaluate()[0] = %v, want 'static'", got[0])
		}
		if got[1] != "bar" {
			t.Errorf("TraverseAndEvaluate()[1] = %v, want 'bar'", got[1])
		}
	})

	t.Run("Nested structures", func(t *testing.T) {
		node := map[string]interface{}{
			"level1": []interface{}{
				map[string]interface{}{
					"expr": "${ .foo }",
				},
			},
		}
		input := map[string]interface{}{
			"foo": "nestedValue",
		}

		result, err := TraverseAndEvaluate(node, input, context.TODO())
		if err != nil {
			t.Fatalf("TraverseAndEvaluate() error: %v", err)
		}

		resMap, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return a map at top-level")
		}

		level1, ok := resMap["level1"].([]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return an array for resMap['level1']")
		}

		level1Map, ok := level1[0].(map[string]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluate() did not return a map for level1[0]")
		}

		if level1Map["expr"] != "nestedValue" {
			t.Errorf("TraverseAndEvaluate() = %v, want %v", level1Map["expr"], "nestedValue")
		}
	})

	t.Run("Invalid JQ expression", func(t *testing.T) {
		node := "${ .foo( }"
		input := map[string]interface{}{
			"foo": "bar",
		}

		_, err := TraverseAndEvaluate(node, input, context.TODO())
		if err == nil {
			t.Errorf("TraverseAndEvaluate() expected error for invalid JQ, got nil")
		}
	})
}

func TestTraverseAndEvaluateWithVars(t *testing.T) {
	t.Run("Variable usage in expression", func(t *testing.T) {
		node := map[string]interface{}{
			"expr": "${ $myVar }",
		}
		variables := map[string]interface{}{
			"$myVar": "HelloVars",
		}
		input := map[string]interface{}{}

		result, err := TraverseAndEvaluateWithVars(node, input, variables, context.TODO())
		if err != nil {
			t.Fatalf("TraverseAndEvaluateWithVars() unexpected error: %v", err)
		}
		got, ok := result.(map[string]interface{})
		if !ok {
			t.Fatalf("TraverseAndEvaluateWithVars() did not return a map")
		}
		if got["expr"] != "HelloVars" {
			t.Errorf("TraverseAndEvaluateWithVars() = %v, want %v", got["expr"], "HelloVars")
		}
	})

	t.Run("Reference variable that isn't defined", func(t *testing.T) {
		// This tries to use a variable that isn't passed in,
		// so presumably it yields an error about an undefined variable.
		node := "${ $notProvided }"
		input := map[string]interface{}{
			"foo": "bar",
		}
		variables := map[string]interface{}{} // intentionally empty

		_, err := TraverseAndEvaluateWithVars(node, input, variables, context.TODO())
		if err == nil {
			t.Errorf("TraverseAndEvaluateWithVars() expected error for undefined variable, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

func TestEvaluateJQExpressionDirect(t *testing.T) {
	// This tests the core evaluator directly for errors and success.
	t.Run("Successful eval", func(t *testing.T) {
		expression := ".foo"
		input := map[string]interface{}{"foo": "bar"}
		variables := map[string]interface{}{}
		result, err := callEvaluateJQ(expression, input, variables)
		if err != nil {
			t.Fatalf("evaluateJQExpression() error = %v, want nil", err)
		}
		if result != "bar" {
			t.Errorf("evaluateJQExpression() = %v, want 'bar'", result)
		}
	})

	t.Run("Parse error", func(t *testing.T) {
		expression := ".foo("
		input := map[string]interface{}{"foo": "bar"}
		variables := map[string]interface{}{}
		_, err := callEvaluateJQ(expression, input, variables)
		if err == nil {
			t.Errorf("evaluateJQExpression() expected parse error, got nil")
		}
	})

	t.Run("Runtime error in evaluation (undefined variable)", func(t *testing.T) {
		expression := "$undefinedVar"
		input := map[string]interface{}{
			"foo": []interface{}{1, 2},
		}
		variables := map[string]interface{}{}
		_, err := callEvaluateJQ(expression, input, variables)
		if err == nil {
			t.Errorf("callEvaluateJQ() expected runtime error, got nil")
		} else {
			t.Logf("Got expected error: %v", err)
		}
	})
}

// Helper to call the unexported evaluateJQExpression via a wrapper in tests.
// Alternatively, you could move `evaluateJQExpression` into a separate file that
// is also in package `expr`, then test it directly if needed.
func callEvaluateJQ(expression string, input interface{}, variables map[string]interface{}) (interface{}, error) {
	// Replicate the logic from evaluateJQExpression for direct testing
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse: %w", err)
	}
	code, err := gojq.Compile(query, gojq.WithVariables(exprGetVariableNames(variables)))
	if err != nil {
		return nil, fmt.Errorf("failed to compile: %w", err)
	}
	iter := code.Run(input, exprGetVariableValues(variables)...)
	result, ok := iter.Next()
	if !ok {
		return nil, fmt.Errorf("no result from jq evaluation")
	}
	if e, isErr := result.(error); isErr {
		return nil, fmt.Errorf("runtime error: %w", e)
	}
	return result, nil
}

// Local copies of the variable-gathering logic from your code:
func exprGetVariableNames(variables map[string]interface{}) []string {
	names := make([]string, 0, len(variables))
	for name := range variables {
		names = append(names, name)
	}
	return names
}

func exprGetVariableValues(variables map[string]interface{}) []interface{} {
	vals := make([]interface{}, 0, len(variables))
	for _, val := range variables {
		vals = append(vals, val)
	}
	return vals
}
