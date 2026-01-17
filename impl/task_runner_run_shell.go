// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package impl

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/serverlessworkflow/sdk-go/v3/impl/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

// RunTaskShell defines the shell configuration for RunTask.
// It implements the RunTask.shell definition.
type RunTaskShell struct {
}

// NewRunTaskShell creates a new RunTaskShell instance.
func NewRunTaskShell() *RunTaskShell {
	return &RunTaskShell{}
}

func (shellTask *RunTaskShell) RunTask(r *RunTaskRunner, input interface{}, taskSupport TaskSupport) (interface{}, error) {
	await := r.Task.Run.Await
	shell := r.Task.Run.Shell
	var cmdStr string

	if shell == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no shell configuration provided for RunTask %s", r.TaskName), r.TaskName)
	}

	cmdStr = shell.Command

	if cmdStr == "" {
		return nil, model.NewErrValidation(fmt.Errorf("no command provided for RunTask shell: %s ", r.TaskName), r.TaskName)
	}

	if shell.Environment != nil {
		for key, value := range shell.Environment {
			evaluated, evalErr := expr.TraverseAndEvaluate(value, input, taskSupport.GetContext())
			if evalErr != nil {
				return nil, model.NewErrRuntime(fmt.Errorf("error evaluating environment variable value for RunTask shell: %s", r.TaskName), r.TaskName)
			}

			envVal := fmt.Sprint(evaluated)
			if err := os.Setenv(key, envVal); err != nil {
				return nil, model.NewErrRuntime(fmt.Errorf("error setting environment variable for RunTask shell: %s", r.TaskName), r.TaskName)
			}
		}
	}

	evaluated, err := expr.TraverseAndEvaluate(cmdStr, input, taskSupport.GetContext())
	if err != nil {
		return nil, model.NewErrRuntime(fmt.Errorf("error evaluating command for RunTask shell: %s", r.TaskName), r.TaskName)
	}

	cmdEvaluated := fmt.Sprint(evaluated)

	var args []string

	args = append(args, "-c", cmdEvaluated)

	if shell.Arguments != nil {
		for key, value := range shell.Arguments {
			keyEval, evalErr := expr.TraverseAndEvaluate(key, input, taskSupport.GetContext())
			if evalErr != nil {
				return nil, model.NewErrRuntime(fmt.Errorf("error evaluating argument key for RunTask shell: %s", r.TaskName), r.TaskName)
			}

			keyStr := fmt.Sprint(keyEval)

			if value != nil {
				valueEval, evalErr := expr.TraverseAndEvaluate(value, input, taskSupport.GetContext())
				if evalErr != nil {
					return nil, model.NewErrRuntime(fmt.Errorf("error evaluating argument value for RunTask shell: %s", r.TaskName), r.TaskName)
				}
				valueStr := fmt.Sprint(valueEval)
				args = append(args, fmt.Sprintf("%s=%s", keyStr, valueStr))
			} else {
				args = append(args, fmt.Sprintf("%s", keyStr))
			}
		}
	}

	var fullCmd strings.Builder
	fullCmd.WriteString(cmdEvaluated)
	for i := 2; i < len(args); i++ {
		fullCmd.WriteString(" ")
		fullCmd.WriteString(args[i])
	}

	if await != nil && !*await {
		go func() {
			cmd := exec.Command("sh", "-c", fullCmd.String())
			_ = cmd.Start()
			_ = cmd.Wait()
		}()
		return input, nil
	}

	cmd := exec.Command("sh", "-c", fullCmd.String())
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()
	exitCode := 0
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			exitCode = exitErr.ExitCode()
		}
	} else if cmd.ProcessState != nil {
		exitCode = cmd.ProcessState.ExitCode()
	}

	stdoutStr := strings.TrimSpace(stdout.String())
	stderrStr := strings.TrimSpace(stderr.String())

	switch r.Task.Run.Return {
	case "all":
		return NewProcessResult(stdoutStr, stderrStr, exitCode), nil
	case "stderr":
		return stderrStr, nil
	case "code":
		return exitCode, nil
	case "none":
		return nil, nil
	default:
		return stdoutStr, nil
	}
}
