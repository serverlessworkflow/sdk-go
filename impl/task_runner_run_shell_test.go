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
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRunShellWithTestData(t *testing.T) {

	t.Run("Simple with echo", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo.yaml"

		input := map[string]interface{}{}
		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		assert.NotNilf(t, output, "output should not be nil")
		assert.Equal(t, "Hello, anonymous", processResult.Stdout)
		assert.Equal(t, "", processResult.Stderr)
		assert.Equal(t, 0, processResult.Code)
		assert.NoError(t, err)
	})

	t.Run("Simple echo looking exit code", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_exitcode.yaml"
		input := map[string]interface{}{}
		expectedOutput := 2
		output, err := runWorkflowExpectString(t, workflowPath, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output.(int))
	})

	t.Run("JQ expression in command with 'all' return", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_jq.yaml"
		input := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "Matheus Cruz",
			},
		}
		output, err := runWorkflowExpectString(t, workflowPath, input)

		processResult := output.(*ProcessResult)
		assert.NoError(t, err)
		assert.Equal(t, "", processResult.Stderr)
		assert.Equal(t, "Hello, Matheus Cruz", processResult.Stdout)
		assert.Equal(t, 0, processResult.Code)
	})

	t.Run("Simple echo with 'none' return", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_none.yaml"
		input := map[string]interface{}{}
		output, err := runWorkflowExpectString(t, workflowPath, input)

		assert.NoError(t, err)
		assert.Nil(t, output)
	})

	t.Run("Simple echo with env and await as 'false'", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_env_no_awaiting.yaml"
		input := map[string]interface{}{
			"full_name": "John Doe",
		}
		output, err := runWorkflowExpectString(t, workflowPath, input)

		assert.NoError(t, err)
		assert.Equal(t, output, input)
		file, err := os.ReadFile("/tmp/hello-world.txt")
		assert.Equal(t, "hello world not awaiting (John Doe)", strings.TrimSpace(string(file)))
	})

	t.Run("Simple echo not awaiting, function should returns immediately", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_not_awaiting.yaml"
		input := map[string]interface{}{
			"full_name": "John Doe",
		}
		output, err := runWorkflow(t, workflowPath, input)

		assert.NoError(t, err)
		assert.Equal(t, output, input)
	})

	t.Run("Simple 'ls' command getting output as stderr", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_ls_stderr.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflowExpectString(t, workflowPath, input)

		assert.NoError(t, err)
		assert.True(t, strings.Contains(output.(string), "ls:"))
	})

	t.Run("Simple echo with args using JQ expression", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_with_args_key_value_jq.yaml"
		input := map[string]interface{}{
			"user":        "Alice",
			"passwordKey": "--password",
		}

		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)
		assert.True(t, strings.Contains(processResult.Stdout, "--user=Alice"))
		assert.True(t, strings.Contains(processResult.Stdout, "--password=serverless"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with args", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_args.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		// Go does not keep the order of map iteration
		// TODO: improve the UnMarshal of args to keep the order

		assert.NoError(t, err)
		assert.True(t, strings.Contains(processResult.Stdout, "--user=john"))
		assert.True(t, strings.Contains(processResult.Stdout, "--password=doe"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with args using only key", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_args_only_key.yaml"
		input := map[string]interface{}{
			"firstName": "Mary",
			"lastName":  "Jane",
		}

		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)

		// Go does not keep the order of map iteration
		// TODO: improve the UnMarshal of args to keep the order
		assert.True(t, strings.Contains(processResult.Stdout, "Mary"))
		assert.True(t, strings.Contains(processResult.Stdout, "Jane"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with env and JQ", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_env.yaml"
		input := map[string]interface{}{
			"lastName": "Doe",
		}

		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)
		assert.True(t, strings.Contains(processResult.Stdout, "Hello John Doe from env!"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Execute touch and cat command", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_touch_cat.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflow(t, workflowPath, input)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)
		assert.Equal(t, "hello world", strings.TrimSpace(processResult.Stdout))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})
}
