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
	"errors"
	"fmt"

	"github.com/itchyny/gojq"
	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

func TraverseAndEvaluateWithVars(node interface{}, input interface{}, variables map[string]interface{}, nodeContext context.Context) (interface{}, error) {
	if err := mergeContextInVars(nodeContext, variables); err != nil {
		return nil, err
	}
	return traverseAndEvaluate(node, input, variables)
}

// TraverseAndEvaluate recursively processes and evaluates all expressions in a JSON-like structure
func TraverseAndEvaluate(node interface{}, input interface{}, nodeContext context.Context) (interface{}, error) {
	return TraverseAndEvaluateWithVars(node, input, map[string]interface{}{}, nodeContext)
}

func traverseAndEvaluate(node interface{}, input interface{}, variables map[string]interface{}) (interface{}, error) {
	switch v := node.(type) {
	case map[string]interface{}:
		// Traverse map
		for key, value := range v {
			evaluatedValue, err := traverseAndEvaluate(value, input, variables)
			if err != nil {
				return nil, err
			}
			v[key] = evaluatedValue
		}
		return v, nil

	case []interface{}:
		// Traverse array
		for i, value := range v {
			evaluatedValue, err := traverseAndEvaluate(value, input, variables)
			if err != nil {
				return nil, err
			}
			v[i] = evaluatedValue
		}
		return v, nil

	case string:
		// Check if the string is a runtime expression (e.g., ${ .some.path })
		if model.IsStrictExpr(v) {
			return evaluateJQExpression(model.SanitizeExpr(v), input, variables)
		}
		return v, nil

	default:
		// Return other types as-is
		return v, nil
	}
}

// evaluateJQExpression evaluates a jq expression against a given JSON input
func evaluateJQExpression(expression string, input interface{}, variables map[string]interface{}) (interface{}, error) {
	query, err := gojq.Parse(expression)
	if err != nil {
		return nil, fmt.Errorf("failed to parse jq expression: %s, error: %w", expression, err)
	}

	// Get the variable names & values in a single pass:
	names, values := getVariableNamesAndValues(variables)

	code, err := gojq.Compile(query, gojq.WithVariables(names))
	if err != nil {
		return nil, fmt.Errorf("failed to compile jq expression: %s, error: %w", expression, err)
	}

	iter := code.Run(input, values...)
	result, ok := iter.Next()
	if !ok {
		return nil, errors.New("no result from jq evaluation")
	}

	// If there's an error from the jq engine, report it
	if errVal, isErr := result.(error); isErr {
		return nil, fmt.Errorf("jq evaluation error: %w", errVal)
	}

	return result, nil
}

// getVariableNamesAndValues constructs two slices, where 'names[i]' matches 'values[i]'.
func getVariableNamesAndValues(vars map[string]interface{}) ([]string, []interface{}) {
	names := make([]string, 0, len(vars))
	values := make([]interface{}, 0, len(vars))

	for k, v := range vars {
		names = append(names, k)
		values = append(values, v)
	}
	return names, values
}

func mergeContextInVars(nodeCtx context.Context, variables map[string]interface{}) error {
	if variables == nil {
		variables = make(map[string]interface{})
	}
	wfCtx, err := ctx.GetWorkflowContext(nodeCtx)
	if err != nil {
		if errors.Is(err, ctx.ErrWorkflowContextNotFound) {
			return nil
		}
		return err
	}
	// merge
	for k, val := range wfCtx.GetVars() {
		variables[k] = val
	}

	return nil
}

func TraverseAndEvaluateObj(runtimeExpr *model.ObjectOrRuntimeExpr, input interface{}, taskName string, wfCtx context.Context) (output interface{}, err error) {
	if runtimeExpr == nil {
		return input, nil
	}
	output, err = TraverseAndEvaluate(runtimeExpr.AsStringOrMap(), input, wfCtx)
	if err != nil {
		return nil, model.NewErrExpression(err, taskName)
	}
	return output, nil
}

func TraverseAndEvaluateBool(runtimeExpr string, input interface{}, wfCtx context.Context) (bool, error) {
	if len(runtimeExpr) == 0 {
		return false, nil
	}
	output, err := TraverseAndEvaluate(runtimeExpr, input, wfCtx)
	if err != nil {
		return false, nil
	}
	if result, ok := output.(bool); ok {
		return result, nil
	}
	return false, nil
}
