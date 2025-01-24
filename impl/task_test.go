package impl

import (
	"github.com/serverlessworkflow/sdk-go/v3/parser"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

// runWorkflowTest is a reusable test function for workflows
func runWorkflowTest(t *testing.T, workflowPath string, input, expectedOutput map[string]interface{}) {
	// Read the workflow YAML from the testdata directory
	yamlBytes, err := ioutil.ReadFile(filepath.Clean(workflowPath))
	assert.NoError(t, err, "Failed to read workflow YAML file")

	// Parse the YAML workflow
	workflow, err := parser.FromYAMLSource(yamlBytes)
	assert.NoError(t, err, "Failed to parse workflow YAML")

	// Initialize the workflow runner
	runner := NewDefaultRunner(workflow)

	// Run the workflow
	output, err := runner.Run(input)

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedOutput, output, "Workflow output mismatch")
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
			"colors": []interface{}{"red", "green", "blue"},
		}
		runWorkflowTest(t, workflowPath, input, expectedOutput)
	})
}
