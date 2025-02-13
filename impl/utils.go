package impl

import (
	"github.com/serverlessworkflow/sdk-go/v3/expr"
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

func traverseAndEvaluate(runtimeExpr *model.ObjectOrRuntimeExpr, input interface{}, taskName string) (output interface{}, err error) {
	if runtimeExpr == nil {
		return input, nil
	}
	output, err = expr.TraverseAndEvaluate(runtimeExpr.AsStringOrMap(), input)
	if err != nil {
		return nil, model.NewErrExpression(err, taskName)
	}
	return output, nil
}
