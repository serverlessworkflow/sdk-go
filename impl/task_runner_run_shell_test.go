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
	"path/filepath"
	"strings"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v3/parser"
	"github.com/stretchr/testify/assert"
)

func TestRunShellWithTestData(t *testing.T) {

	t.Run("Simple with echo", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo.yaml"

		input := map[string]interface{}{}
		output, err := runWorkflow(t, workflowPath, input, nil)

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
		output, err := runWorkflow(t, workflowPath, input, nil)
		assert.NoError(t, err)
		// `ls` of a nonexistent directory exits non-zero; the exact code is
		// platform-dependent (2 on GNU/Linux, 1 on macOS), so assert non-zero.
		assert.NotZero(t, output.(int))
	})

	t.Run("JQ expression in command with 'all' return", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_jq.yaml"
		input := map[string]interface{}{
			"user": map[string]interface{}{
				"name": "Matheus Cruz",
			},
		}
		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)
		assert.NoError(t, err)
		assert.Equal(t, "", processResult.Stderr)
		assert.Equal(t, "Hello, Matheus Cruz", processResult.Stdout)
		assert.Equal(t, 0, processResult.Code)
	})

	t.Run("Simple echo with 'none' return", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_none.yaml"
		input := map[string]interface{}{}
		output, err := runWorkflow(t, workflowPath, input, nil)

		assert.NoError(t, err)
		assert.Nil(t, output)
	})

	t.Run("Simple echo with env and await as 'false'", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_env_no_awaiting.yaml"
		input := map[string]interface{}{
			"full_name": "John Doe",
		}
		output, err := runWorkflow(t, workflowPath, input, nil)

		assert.NoError(t, err)
		assert.Equal(t, output, input)
	})

	t.Run("Simple echo not awaiting, function should returns immediately", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_not_awaiting.yaml"
		input := map[string]interface{}{
			"full_name": "John Doe",
		}
		output, err := runWorkflow(t, workflowPath, input, nil)

		assert.NoError(t, err)
		assert.Equal(t, output, input)
	})

	t.Run("Simple 'ls' command getting output as stderr", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_ls_stderr.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflow(t, workflowPath, input, nil)

		assert.NoError(t, err)
		assert.True(t, strings.Contains(output.(string), "ls:"))
	})

	t.Run("Simple echo with args using JQ expression", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_with_args_key_value_jq.yaml"
		input := map[string]interface{}{
			"user":        "Alice",
			"passwordKey": "--password",
		}

		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)

		// Arguments are provided as a map, so iteration order is not guaranteed;
		// assert each evaluated argument independently.
		assert.NoError(t, err)
		assert.True(t, strings.Contains(processResult.Stdout, "--user=Alice"))
		assert.True(t, strings.Contains(processResult.Stdout, "--password=serverless"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with args", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_args.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)

		// Arguments are provided as an ordered array, so the order is preserved.
		assert.NoError(t, err)
		assert.Equal(t, "--user=john --password=doe", processResult.Stdout)
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with args using only key", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_args_only_key.yaml"
		input := map[string]interface{}{
			"firstName": "Mary",
			"lastName":  "Jane",
		}

		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)

		// Arguments are provided as an ordered array, so the order is preserved.
		assert.NoError(t, err)
		assert.Equal(t, "Hello Mary Jane from args!", processResult.Stdout)
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Simple echo with env and JQ", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_echo_with_env.yaml"
		input := map[string]interface{}{
			"lastName": "Doe",
		}

		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)
		assert.True(t, strings.Contains(processResult.Stdout, "Hello John Doe from env!"))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Execute touch and cat command", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_touch_cat.yaml"
		input := map[string]interface{}{}

		output, err := runWorkflow(t, workflowPath, input, nil)

		processResult := output.(*ProcessResult)

		assert.NoError(t, err)
		assert.Equal(t, "hello world", strings.TrimSpace(processResult.Stdout))
		assert.Equal(t, 0, processResult.Code)
		assert.Equal(t, "", processResult.Stderr)
	})

	t.Run("Missing command fails validation", func(t *testing.T) {
		workflowPath := "./testdata/run_shell_missing_command.yaml"
		yamlBytes, err := os.ReadFile(filepath.Clean(workflowPath))
		assert.NoError(t, err)

		_, err = parser.FromYAMLSource(yamlBytes)
		assert.Error(t, err)
	})
}
