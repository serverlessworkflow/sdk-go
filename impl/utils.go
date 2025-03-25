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

package impl

import (
	"context"
	"github.com/serverlessworkflow/sdk-go/v3/impl/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

// Deep clone a map to avoid modifying the original object
func deepClone(obj map[string]interface{}) map[string]interface{} {
	clone := make(map[string]interface{})
	for key, value := range obj {
		clone[key] = deepCloneValue(value)
	}
	return clone
}

func deepCloneValue(value interface{}) interface{} {
	if m, ok := value.(map[string]interface{}); ok {
		return deepClone(m)
	}
	if s, ok := value.([]interface{}); ok {
		clonedSlice := make([]interface{}, len(s))
		for i, v := range s {
			clonedSlice[i] = deepCloneValue(v)
		}
		return clonedSlice
	}
	return value
}

func validateSchema(data interface{}, schema *model.Schema, taskName string) error {
	if schema != nil {
		if err := ValidateJSONSchema(data, schema); err != nil {
			return model.NewErrValidation(err, taskName)
		}
	}
	return nil
}

func traverseAndEvaluate(runtimeExpr *model.ObjectOrRuntimeExpr, input interface{}, taskName string, wfCtx context.Context) (output interface{}, err error) {
	if runtimeExpr == nil {
		return input, nil
	}
	output, err = expr.TraverseAndEvaluate(runtimeExpr.AsStringOrMap(), input, wfCtx)
	if err != nil {
		return nil, model.NewErrExpression(err, taskName)
	}
	return output, nil
}
