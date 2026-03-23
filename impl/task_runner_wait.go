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
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

func NewWaitTaskRunner(taskName string, task *model.WaitTask) (*WaitTaskRunner, error) {
	if task == nil || task.Wait == nil {
		return nil, model.NewErrValidation(fmt.Errorf("no wait configuration provided for WaitTask %s", taskName), taskName)
	}

	return &WaitTaskRunner{
		Task:     task,
		TaskName: taskName,
	}, nil
}

type WaitTaskRunner struct {
	Task     *model.WaitTask
	TaskName string
}

func (w *WaitTaskRunner) Run(input interface{}, taskSupport TaskSupport) (interface{}, error) {
	if w.Task == nil || w.Task.Wait == nil {
		return nil, model.NewErrValidation(fmt.Errorf("wait configuration is empty for task %s", w.TaskName), w.TaskName)
	}

	waitFor, err := model.DurationToTime(w.Task.Wait, time.Now())
	if err != nil {
		return nil, model.NewErrValidation(err, w.TaskName)
	}

	if waitFor <= 0 {
		return input, nil
	}

	taskSupport.SetTaskStatus(w.TaskName, ctx.WaitingStatus)

	// Do not use time.Sleep(waitFor) here: Sleep blocks unconditionally and
	// ignores context cancellation/deadlines, which can hang task shutdown and
	// make cancellation tests/timeouts flaky. A timer + select lets us stop
	// waiting as soon as the task context is done.
	timer := time.NewTimer(waitFor)
	defer timer.Stop()

	select {
	case <-timer.C:
		return input, nil
	case <-taskSupport.GetContext().Done():
		ctxErr := taskSupport.GetContext().Err()
		if errors.Is(ctxErr, context.DeadlineExceeded) {
			return nil, model.NewErrTimeout(ctxErr, w.TaskName)
		}
		return nil, model.NewErrRuntime(ctxErr, w.TaskName)
	}
}

func (w *WaitTaskRunner) GetTaskName() string {
	return w.TaskName
}
