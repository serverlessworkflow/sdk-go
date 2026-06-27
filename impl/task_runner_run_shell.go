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

	if shell == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no shell configuration provided for RunTask %s", r.TaskName), r.TaskName)
	}

	if shell.Command == "" {
		return nil, model.NewErrValidation(fmt.Errorf("no command provided for RunTask shell: %s", r.TaskName), r.TaskName)
	}

	// Build the environment for the command without mutating the parent process'
	// global environment.
	env := os.Environ()
	for key, value := range shell.Environment {
		evaluated, evalErr := expr.TraverseAndEvaluate(value, input, taskSupport.GetContext())
		if evalErr != nil {
			return nil, model.NewErrRuntime(fmt.Errorf("error evaluating environment variable value for RunTask shell: %s", r.TaskName), r.TaskName)
		}
		env = append(env, fmt.Sprintf("%s=%s", key, fmt.Sprint(evaluated)))
	}

	evaluated, err := expr.TraverseAndEvaluate(shell.Command, input, taskSupport.GetContext())
	if err != nil {
		return nil, model.NewErrRuntime(fmt.Errorf("error evaluating command for RunTask shell: %s", r.TaskName), r.TaskName)
	}
	cmdEvaluated := fmt.Sprint(evaluated)

	args, err := shellTask.buildArguments(r, shell.Arguments, input, taskSupport)
	if err != nil {
		return nil, err
	}

	var fullCmd strings.Builder
	fullCmd.WriteString(cmdEvaluated)
	for _, arg := range args {
		fullCmd.WriteString(" ")
		fullCmd.WriteString(arg)
	}

	newCmd := func() *exec.Cmd {
		cmd := exec.Command("sh", "-c", fullCmd.String())
		cmd.Env = env
		return cmd
	}

	// When the task is explicitly not awaited, fire and forget.
	if await != nil && !*await {
		go func() {
			cmd := newCmd()
			_ = cmd.Start()
			_ = cmd.Wait()
		}()
		return input, nil
	}

	cmd := newCmd()
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

// buildArguments resolves the shell arguments into an ordered list of strings.
// Arguments may be provided either as an ordered array (preferred when order
// matters) or as a key/value map (order is not guaranteed by Go map iteration).
func (shellTask *RunTaskShell) buildArguments(r *RunTaskRunner, arguments *model.RunArguments, input interface{}, taskSupport TaskSupport) ([]string, error) {
	if arguments == nil {
		return nil, nil
	}

	eval := func(value interface{}) (string, error) {
		evaluated, evalErr := expr.TraverseAndEvaluate(value, input, taskSupport.GetContext())
		if evalErr != nil {
			return "", model.NewErrRuntime(fmt.Errorf("error evaluating argument for RunTask shell: %s", r.TaskName), r.TaskName)
		}
		return fmt.Sprint(evaluated), nil
	}

	if slice := arguments.AsSlice(); slice != nil {
		args := make([]string, 0, len(slice))
		for _, value := range slice {
			arg, err := eval(value)
			if err != nil {
				return nil, err
			}
			args = append(args, arg)
		}
		return args, nil
	}

	if m := arguments.AsMap(); m != nil {
		args := make([]string, 0, len(m))
		for key, value := range m {
			keyStr, err := eval(key)
			if err != nil {
				return nil, err
			}
			if value != nil {
				valueStr, err := eval(value)
				if err != nil {
					return nil, err
				}
				args = append(args, fmt.Sprintf("%s=%s", keyStr, valueStr))
			} else {
				args = append(args, keyStr)
			}
		}
		return args, nil
	}

	return nil, nil
}
