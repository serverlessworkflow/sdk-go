package impl

import (
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/serverlessworkflow/sdk-go/v3/parser"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

// runWorkflowTest is a reusable test function for workflows
func runWorkflowTest(t *testing.T, workflowPath string, input, expectedOutput map[string]interface{}) {
	// Run the workflow
	output, err := runWorkflow(t, workflowPath, input, expectedOutput)
	assert.NoError(t, err)

	assertWorkflowRun(t, expectedOutput, output)
}

func runWorkflowWithErr(t *testing.T, workflowPath string, input, expectedOutput map[string]interface{}, assertErr func(error)) {
	output, err := runWorkflow(t, workflowPath, input, expectedOutput)
	assert.Error(t, err)
	assertErr(err)
	assertWorkflowRun(t, expectedOutput, output)
}

func runWorkflow(t *testing.T, workflowPath string, input, expectedOutput map[string]interface{}) (output interface{}, err error) {
	// Read the workflow YAML from the testdata directory
	yamlBytes, err := os.ReadFile(filepath.Clean(workflowPath))
	assert.NoError(t, err, "Failed to read workflow YAML file")

	// Parse the YAML workflow
	workflow, err := parser.FromYAMLSource(yamlBytes)
	assert.NoError(t, err, "Failed to parse workflow YAML")

	// Initialize the workflow runner
	runner := NewDefaultRunner(workflow)

	// Run the workflow
	output, err = runner.Run(input)
	return output, err
}

func assertWorkflowRun(t *testing.T, expectedOutput map[string]interface{}, output interface{}) {
	if expectedOutput == nil {
		assert.Nil(t, output, "Expected nil Workflow run output")
	} else {
		assert.Equal(t, expectedOutput, output, "Workflow output mismatch")
	}
}

// TestWorkflowRunner_Run_YAML validates multiple workflows
func TestWorkflowRunner_Run_YAML(t *testing.T) {
	// Workflow 1: Chained Set Tasks
	t.Run("Chained Set Tasks", func(t *testing.T) {
		workflowPath := "./testdata/chained_set_tasks.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"tripled": float64(60),
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	// Workflow 2: Concatenating Strings
	t.Run("Concatenating Strings", func(t *testing.T) {
		workflowPath := "./testdata/concatenating_strings.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"fullName": "John Doe",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	// Workflow 3: Conditional Logic
	t.Run("Conditional Logic", func(t *testing.T) {
		workflowPath := "./testdata/conditional_logic.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"weather": "hot",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Conditional Logic", func(t *testing.T) {
		workflowPath := "./testdata/sequential_set_colors.yaml"
		// Define the input and expected output
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"resultColors": []interface{}{"red", "green", "blue"},
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})
	t.Run("input From", func(t *testing.T) {
		workflowPath := "./testdata/sequential_set_colors_output_as.yaml"
		// Define the input and expected output
		expectedOutput := map[string]interface{}{
			"result": []interface{}{"red", "green", "blue"},
		}
		runWorkflowTest(t, workflowPath, nil, expectedOutput)
	})
	t.Run("input From", func(t *testing.T) {
		workflowPath := "./testdata/conditional_logic_input_from.yaml"
		// Define the input and expected output
		input := map[string]interface{}{
			"localWeather": map[string]interface{}{
				"temperature": 34,
			},
		}
		expectedOutput := map[string]interface{}{
			"weather": "hot",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})
}

func TestWorkflowRunner_Run_YAML_WithSchemaValidation(t *testing.T) {
	// Workflow 1: Workflow input Schema Validation
	t.Run("Workflow input Schema Validation - Valid input", func(t *testing.T) {
		workflowPath := "./testdata/workflow_input_schema.yaml"
		input := map[string]interface{}{
			"key": "value",
		}
		expectedOutput := map[string]interface{}{
			"outputKey": "value",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Workflow input Schema Validation - Invalid input", func(t *testing.T) {
		workflowPath := "./testdata/workflow_input_schema.yaml"
		input := map[string]interface{}{
			"wrongKey": "value",
		}
		yamlBytes, err := os.ReadFile(filepath.Clean(workflowPath))
		assert.NoError(t, err, "Failed to read workflow YAML file")
		workflow, err := parser.FromYAMLSource(yamlBytes)
		assert.NoError(t, err, "Failed to parse workflow YAML")
		runner := NewDefaultRunner(workflow)
		_, err = runner.Run(input)
		assert.Error(t, err, "Expected validation error for invalid input")
		assert.Contains(t, err.Error(), "JSON schema validation failed")
	})

	// Workflow 2: Task input Schema Validation
	t.Run("Task input Schema Validation", func(t *testing.T) {
		workflowPath := "./testdata/task_input_schema.yaml"
		input := map[string]interface{}{
			"taskInputKey": 42,
		}
		expectedOutput := map[string]interface{}{
			"taskOutputKey": 84,
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Task input Schema Validation - Invalid input", func(t *testing.T) {
		workflowPath := "./testdata/task_input_schema.yaml"
		input := map[string]interface{}{
			"taskInputKey": "invalidValue",
		}
		yamlBytes, err := os.ReadFile(filepath.Clean(workflowPath))
		assert.NoError(t, err, "Failed to read workflow YAML file")
		workflow, err := parser.FromYAMLSource(yamlBytes)
		assert.NoError(t, err, "Failed to parse workflow YAML")
		runner := NewDefaultRunner(workflow)
		_, err = runner.Run(input)
		assert.Error(t, err, "Expected validation error for invalid task input")
		assert.Contains(t, err.Error(), "JSON schema validation failed")
	})

	// Workflow 3: Task output Schema Validation
	t.Run("Task output Schema Validation", func(t *testing.T) {
		workflowPath := "./testdata/task_output_schema.yaml"
		input := map[string]interface{}{
			"taskInputKey": "value",
		}
		expectedOutput := map[string]interface{}{
			"finalOutputKey": "resultValue",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Task output Schema Validation - Invalid output", func(t *testing.T) {
		workflowPath := "./testdata/task_output_schema_with_dynamic_value.yaml"
		input := map[string]interface{}{
			"taskInputKey": 123, // Invalid value (not a string)
		}
		yamlBytes, err := os.ReadFile(filepath.Clean(workflowPath))
		assert.NoError(t, err, "Failed to read workflow YAML file")
		workflow, err := parser.FromYAMLSource(yamlBytes)
		assert.NoError(t, err, "Failed to parse workflow YAML")
		runner := NewDefaultRunner(workflow)
		_, err = runner.Run(input)
		assert.Error(t, err, "Expected validation error for invalid task output")
		assert.Contains(t, err.Error(), "JSON schema validation failed")
	})

	t.Run("Task output Schema Validation - Valid output", func(t *testing.T) {
		workflowPath := "./testdata/task_output_schema_with_dynamic_value.yaml"
		input := map[string]interface{}{
			"taskInputKey": "validValue", // Valid value
		}
		expectedOutput := map[string]interface{}{
			"finalOutputKey": "validValue",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	// Workflow 4: Task Export Schema Validation
	t.Run("Task Export Schema Validation", func(t *testing.T) {
		workflowPath := "./testdata/task_export_schema.yaml"
		input := map[string]interface{}{
			"key": "value",
		}
		expectedOutput := map[string]interface{}{
			"exportedKey": "value",
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})
}

func TestWorkflowRunner_Run_YAML_ControlFlow(t *testing.T) {
	t.Run("Set Tasks with Then Directive", func(t *testing.T) {
		workflowPath := "./testdata/set_tasks_with_then.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"result": float64(90),
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Set Tasks with Termination", func(t *testing.T) {
		workflowPath := "./testdata/set_tasks_with_termination.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"finalValue": float64(20),
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})

	t.Run("Set Tasks with Invalid Then Reference", func(t *testing.T) {
		workflowPath := "./testdata/set_tasks_invalid_then.yaml"
		input := map[string]interface{}{}
		expectedOutput := map[string]interface{}{
			"partialResult": float64(15),
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})
}

func TestWorkflowRunner_Run_YAML_RaiseTasks(t *testing.T) {
	// TODO: add $workflow context to the expr processing
	//t.Run("Raise Inline Error", func(t *testing.T) {
	//	runWorkflowTest(t, "./testdata/raise_inline.yaml", nil, nil)
	//})

	t.Run("Raise Referenced Error", func(t *testing.T) {
		runWorkflowWithErr(t, "./testdata/raise_reusable.yaml", nil, nil,
			func(err error) {
				assert.Equal(t, model.ErrorTypeAuthentication, model.AsError(err).Type.String())
			})
	})

	t.Run("Raise Error with Dynamic Detail", func(t *testing.T) {
		input := map[string]interface{}{
			"reason": "User token expired",
		}
		runWorkflowWithErr(t, "./testdata/raise_error_with_input.yaml", input, nil,
			func(err error) {
				assert.Equal(t, model.ErrorTypeAuthentication, model.AsError(err).Type.String())
				assert.Equal(t, "User authentication failed: User token expired", model.AsError(err).Detail.String())
			})
	})

	t.Run("Raise Undefined Error Reference", func(t *testing.T) {
		runWorkflowWithErr(t, "./testdata/raise_undefined_reference.yaml", nil, nil,
			func(err error) {
				assert.Equal(t, model.ErrorTypeValidation, model.AsError(err).Type.String())
			})
	})
}

func TestWorkflowRunner_Run_YAML_RaiseTasks_ControlFlow(t *testing.T) {
	t.Run("Raise Error with Conditional Logic", func(t *testing.T) {
		input := map[string]interface{}{
			"user": map[string]interface{}{
				"age": 16,
			},
		}
		runWorkflowWithErr(t, "./testdata/raise_conditional.yaml", input, nil,
			func(err error) {
				assert.Equal(t, model.ErrorTypeAuthorization, model.AsError(err).Type.String())
				assert.Equal(t, "User is under the required age", model.AsError(err).Detail.String())
			})
	})

}
