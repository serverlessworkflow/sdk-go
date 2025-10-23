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
	"fmt"
	"strings"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/stretchr/testify/assert"
)

func testingRunShell(t *testing.T, task model.RunTask, expected interface{}, input map[string]interface{}) {

	wfCtx, err := ctx.NewWorkflowContext(&model.Workflow{
		Input: &model.Input{
			From: &model.ObjectOrRuntimeExpr{Value: input},
		},
	})
	assert.NoError(t, err)
	wfCtx.SetTaskReference("task_run_defined")
	wfCtx.SetInput(input)

	runner, err := NewRunTaskRunner("runShell", &task)
	assert.NoError(t, err)

	taskSupport := newTaskSupport(withRunnerCtx(wfCtx))

	if input == nil {
		input = map[string]interface{}{}
	}

	output, err := runner.Run(input, taskSupport)

	assert.NoError(t, err)

	switch exp := expected.(type) {

	case int:
		// expected an exit code
		codeOut, ok := output.(int)
		assert.True(t, ok, "output should be int (exit code), got %T", output)
		assert.Equal(t, exp, codeOut)
	case string:
		var outStr string
		switch v := output.(type) {
		case string:
			outStr = v
		case []byte:
			outStr = string(v)
		case int:
			outStr = fmt.Sprintf("%d", v)
		default:
			t.Fatalf("unexpected output type %T", output)
		}
		outStr = strings.TrimSpace(outStr)
		assert.Equal(t, exp, outStr)
	case ProcessResult:
		resultOut, ok := output.(*ProcessResult)
		assert.True(t, ok, "output should be ProcessResult, got %T", output)
		assert.Equal(t, exp.Stdout, strings.TrimSpace(resultOut.Stdout))
		assert.Equal(t, exp.Stderr, strings.TrimSpace(resultOut.Stderr))
		assert.Equal(t, exp.Code, resultOut.Code)
	default:
		t.Fatalf("unsupported expected type %T", expected)
	}
}

func TestWithTestData(t *testing.T) {

	t.Run("Simple with echo", func(t *testing.T) {
		workflowPath := "./testdata/runshell_echo.yaml"

		input := map[string]interface{}{}
		expectedOutput := "Hello, anonymous"
		output, err := runWorkflowExpectString(t, workflowPath, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output)
	})

	t.Run("Simple echo looking exit code", func(t *testing.T) {
		workflowPath := "./testdata/runshell_exitcode.yaml"
		input := map[string]interface{}{}
		expectedOutput := 2
		output, err := runWorkflowExpectString(t, workflowPath, input)
		assert.NoError(t, err)
		assert.Equal(t, expectedOutput, output.(int))
	})
}

func TestRunTaskRunner(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		ret      string
		expected interface{}
		input    map[string]interface{}
	}{
		{
			name:     "echoLookCode",
			cmd:      "echo 'hello world'",
			ret:      "code",
			expected: 0,
		},
		{
			name:     "echoLookStdout",
			cmd:      "echo 'hello world'",
			ret:      "stdout",
			expected: "hello world",
		},
		{
			name: "echoLookAll",
			cmd:  "echo 'hello world'",
			ret:  "all",
			expected: *NewProcessResult(
				"hello world",
				"",
				0,
			),
		},
		{
			name:     "echoJqExpression",
			cmd:      `${ "echo Hello, I love \(.project)" }`,
			ret:      "stdout",
			expected: "Hello, I love ServerlessWorkflow",
			input: map[string]interface{}{
				"project": "ServerlessWorkflow",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			task := model.RunTask{
				Run: model.RunTaskConfiguration{
					Shell: &model.Shell{
						Command: tc.cmd,
					},
					Return: tc.ret,
				},
			}
			testingRunShell(t, task, tc.expected, tc.input)
		})
	}
}
