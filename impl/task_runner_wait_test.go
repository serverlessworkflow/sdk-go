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
	"testing"
	"time"

	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/stretchr/testify/assert"
)

func TestNewWaitTaskRunner_InvalidTask(t *testing.T) {
	runner, err := NewWaitTaskRunner("wait_task", nil)
	assert.Error(t, err)
	assert.Nil(t, runner)
}

func TestNewWaitTaskRunner_MissingWaitConfiguration(t *testing.T) {
	runner, err := NewWaitTaskRunner("wait_task", &model.WaitTask{})
	assert.Error(t, err)
	assert.Nil(t, runner)
}

func TestNewWaitTaskRunner_ValidTask(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("PT1S"),
	}

	runner, err := NewWaitTaskRunner("wait_task", task)
	assert.NoError(t, err)
	assert.NotNil(t, runner)
	assert.Equal(t, "wait_task", runner.GetTaskName())
	assert.Equal(t, task, runner.Task)
}

func TestWaitTaskRunner_Run_WithInlineDuration(t *testing.T) {
	task := &model.WaitTask{
		Wait: &model.Duration{
			Value: model.DurationInline{
				Milliseconds: 25,
			},
		},
	}

	runner, err := NewWaitTaskRunner("wait_inline", task)
	assert.NoError(t, err)

	input := map[string]interface{}{"value": "test"}
	start := time.Now()
	output, runErr := runner.Run(input, newTaskSupport(withContext(context.Background())))
	elapsed := time.Since(start)

	assert.NoError(t, runErr)
	assert.Equal(t, input, output)
	assert.GreaterOrEqual(t, elapsed, 20*time.Millisecond)
}

func TestWaitTaskRunner_Run_WithISO8601Duration(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("PT1S"),
	}

	runner, err := NewWaitTaskRunner("wait_iso", task)
	assert.NoError(t, err)

	input := "test-input"
	start := time.Now()
	output, runErr := runner.Run(input, newTaskSupport(withContext(context.Background())))
	elapsed := time.Since(start)

	assert.NoError(t, runErr)
	assert.Equal(t, input, output)
	assert.GreaterOrEqual(t, elapsed, 900*time.Millisecond)
}

func TestWaitTaskRunner_Run_WithInvalidDuration(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("1Y"),
	}

	runner, err := NewWaitTaskRunner("wait_invalid", task)
	assert.NoError(t, err)

	output, runErr := runner.Run("input", newTaskSupport())
	assert.Error(t, runErr)
	assert.Nil(t, output)
	assert.True(t, model.IsErrValidation(runErr))
}

func TestWaitTaskRunner_Run_ContextCancelled(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("PT1S"),
	}

	runner, err := NewWaitTaskRunner("wait_cancelled", task)
	assert.NoError(t, err)

	cancelCtx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(20 * time.Millisecond)
		cancel()
	}()

	start := time.Now()
	output, runErr := runner.Run("input", newTaskSupport(withContext(cancelCtx)))
	elapsed := time.Since(start)

	assert.Error(t, runErr)
	assert.Nil(t, output)
	assert.True(t, model.IsErrRuntime(runErr))
	assert.Less(t, elapsed, 500*time.Millisecond)
}

func TestNewTaskRunner_WithWaitTask(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("PT1S"),
	}

	runner, err := NewTaskRunner("wait_task", task, nil)
	assert.NoError(t, err)
	assert.NotNil(t, runner)

	_, ok := runner.(*WaitTaskRunner)
	assert.True(t, ok)
}

func TestWaitTaskRunner_Run_WithValidLongDuration_EntersWaitPath(t *testing.T) {
	task := &model.WaitTask{
		Wait: model.NewDurationExpr("P365D"),
	}
	runner, err := NewWaitTaskRunner("wait_long_valid", task)
	assert.NoError(t, err)

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancel()

	output, runErr := runner.Run("input", newTaskSupport(withContext(cancelCtx)))
	assert.Error(t, runErr)
	assert.Nil(t, output)
	assert.True(t, model.IsErrRuntime(runErr))
	assert.False(t, model.IsErrValidation(runErr))
}
