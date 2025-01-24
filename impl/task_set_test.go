package impl

import (
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestSetTaskExecutor_Exec(t *testing.T) {
	input := map[string]interface{}{
		"configuration": map[string]interface{}{
			"size": map[string]interface{}{
				"width":  6,
				"height": 6,
			},
			"fill": map[string]interface{}{
				"red":   69,
				"green": 69,
				"blue":  69,
			},
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"shape": "circle",
			"size":  "${ .configuration.size }",
			"fill":  "${ .configuration.fill }",
		},
	}

	executor, err := NewSetTaskExecutor("task1", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"shape": "circle",
		"size": map[string]interface{}{
			"width":  6,
			"height": 6,
		},
		"fill": map[string]interface{}{
			"red":   69,
			"green": 69,
			"blue":  69,
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_StaticValues(t *testing.T) {
	input := map[string]interface{}{}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"status": "completed",
			"count":  10,
		},
	}

	executor, err := NewSetTaskExecutor("task_static", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"status": "completed",
		"count":  10,
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_RuntimeExpressions(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"fullName": "${ \"\\(.user.firstName) \\(.user.lastName)\" }",
		},
	}

	executor, err := NewSetTaskExecutor("task_runtime_expr", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"fullName": "John Doe",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_NestedStructures(t *testing.T) {
	input := map[string]interface{}{
		"order": map[string]interface{}{
			"id":    12345,
			"items": []interface{}{"item1", "item2"},
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"orderDetails": map[string]interface{}{
				"orderId":   "${ .order.id }",
				"itemCount": "${ .order.items | length }",
			},
		},
	}

	executor, err := NewSetTaskExecutor("task_nested_structures", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"orderDetails": map[string]interface{}{
			"orderId":   12345,
			"itemCount": 2,
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_StaticAndDynamicValues(t *testing.T) {
	input := map[string]interface{}{
		"config": map[string]interface{}{
			"threshold": 100,
		},
		"metrics": map[string]interface{}{
			"current": 75,
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"status":    "active",
			"remaining": "${ .config.threshold - .metrics.current }",
		},
	}

	executor, err := NewSetTaskExecutor("task_static_dynamic", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"status":    "active",
		"remaining": 25,
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_MissingInputData(t *testing.T) {
	input := map[string]interface{}{}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"value": "${ .missingField }",
		},
	}

	executor, err := NewSetTaskExecutor("task_missing_input", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)
	assert.Nil(t, output["value"])
}

func TestSetTaskExecutor_ExpressionsWithFunctions(t *testing.T) {
	input := map[string]interface{}{
		"values": []interface{}{1, 2, 3, 4, 5},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"sum": "${ .values | map(.) | add }",
		},
	}

	executor, err := NewSetTaskExecutor("task_expr_functions", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"sum": 15,
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_ConditionalExpressions(t *testing.T) {
	input := map[string]interface{}{
		"temperature": 30,
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"weather": "${ if .temperature > 25 then 'hot' else 'cold' end }",
		},
	}

	executor, err := NewSetTaskExecutor("task_conditional_expr", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"weather": "hot",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_ArrayDynamicIndex(t *testing.T) {
	input := map[string]interface{}{
		"items": []interface{}{"apple", "banana", "cherry"},
		"index": 1,
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"selectedItem": "${ .items[.index] }",
		},
	}

	executor, err := NewSetTaskExecutor("task_array_indexing", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"selectedItem": "banana",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_NestedConditionalLogic(t *testing.T) {
	input := map[string]interface{}{
		"age": 20,
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"status": "${ if .age < 18 then 'minor' else if .age < 65 then 'adult' else 'senior' end end }",
		},
	}

	executor, err := NewSetTaskExecutor("task_nested_condition", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"status": "adult",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_DefaultValues(t *testing.T) {
	input := map[string]interface{}{}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"value": "${ .missingField // 'defaultValue' }",
		},
	}

	executor, err := NewSetTaskExecutor("task_default_values", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"value": "defaultValue",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_ComplexNestedStructures(t *testing.T) {
	input := map[string]interface{}{
		"config": map[string]interface{}{
			"dimensions": map[string]interface{}{
				"width":  10,
				"height": 5,
			},
		},
		"meta": map[string]interface{}{
			"color": "blue",
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"shape": map[string]interface{}{
				"type":   "rectangle",
				"width":  "${ .config.dimensions.width }",
				"height": "${ .config.dimensions.height }",
				"color":  "${ .meta.color }",
				"area":   "${ .config.dimensions.width * .config.dimensions.height }",
			},
		},
	}

	executor, err := NewSetTaskExecutor("task_complex_nested", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"shape": map[string]interface{}{
			"type":   "rectangle",
			"width":  10,
			"height": 5,
			"color":  "blue",
			"area":   50,
		},
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}

func TestSetTaskExecutor_MultipleExpressions(t *testing.T) {
	input := map[string]interface{}{
		"user": map[string]interface{}{
			"name":  "Alice",
			"email": "alice@example.com",
		},
	}

	setTask := &model.SetTask{
		Set: map[string]interface{}{
			"username": "${ .user.name }",
			"contact":  "${ .user.email }",
		},
	}

	executor, err := NewSetTaskExecutor("task_multiple_expr", setTask)
	assert.NoError(t, err)

	output, err := executor.Exec(input)
	assert.NoError(t, err)

	expectedOutput := map[string]interface{}{
		"username": "Alice",
		"contact":  "alice@example.com",
	}

	if !reflect.DeepEqual(output, expectedOutput) {
		t.Errorf("expected %v, got %v", expectedOutput, output)
	}
}
