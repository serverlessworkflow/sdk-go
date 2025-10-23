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
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/serverlessworkflow/sdk-go/v3/impl/expr"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

type RunTaskRunner struct {
	Task     *model.RunTask
	TaskName string
}

func (d *RunTaskRunner) GetTaskName() string {
	return d.TaskName
}

// RunTaskRunnable defines the interface for running a subtask for RunTask.
type RunTaskRunnable interface {
	RunTask(taskConfiguration *model.RunTaskConfiguration, support *TaskSupport, input interface{}) (output interface{}, err error)
}

func NewRunTaskRunner(taskName string, task *model.RunTask) (*RunTaskRunner, error) {

	if task == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no set configuration provided for RunTask %s", taskName), taskName)
	}

	return &RunTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

func (d *RunTaskRunner) Run(input interface{}, taskSupport TaskSupport) (output interface{}, err error) {

	if d.Task.Run.Shell != nil {
		shellTask := NewRunTaskShell()
		return shellTask.RunTask(d, input, taskSupport)
	}

	return nil, fmt.Errorf("no set configuration provided for RunTask %s", d.TaskName)

}

// ProcessResult Describes the result of a process.
type ProcessResult struct {
	Stdout string
	Stderr string
	Code   int
}

// NewProcessResult creates a new ProcessResult instance.
func NewProcessResult(stdout, stderr string, code int) *ProcessResult {
	return &ProcessResult{
		Stdout: stdout,
		Stderr: stderr,
		Code:   code,
	}
}

// RunTaskShell defines the shell configuration for RunTask.
// It implements the RunTask.shell definition.
type RunTaskShell struct {
}

// NewRunTaskShell creates a new RunTaskShell instance.
func NewRunTaskShell() *RunTaskShell {
	return &RunTaskShell{}
}

func (shellTask *RunTaskShell) RunTask(r *RunTaskRunner, input interface{}, taskSupport TaskSupport) (interface{}, error) {

	shell := r.Task.Run.Shell
	var cmdStr string

	if shell != nil {
		cmdStr = shell.Command
	}

	if cmdStr == "" {
		return nil, model.NewErrValidation(fmt.Errorf("no command provided for RunTask %shellTask", r.TaskName), r.TaskName)
	}

	evaluated, err := expr.TraverseAndEvaluate(cmdStr, input, taskSupport.GetContext())
	if err != nil {
		return nil, err
	}

	cmdEvaluated, ok := evaluated.(string)
	if !ok {
		return nil, model.NewErrRuntime(fmt.Errorf("expected evaluated command to be a string, but got a different type. Got: %v", evaluated), r.TaskName)
	}

	cmd := exec.Command("sh", "-c", cmdEvaluated)
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
		return stdoutStr, nil
	case "code":
		return exitCode, nil
	default:
		return stdoutStr, nil
	}
}
