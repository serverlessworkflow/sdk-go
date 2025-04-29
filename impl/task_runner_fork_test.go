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

// dummyRunner simulates a TaskRunner that returns its name after an optional delay.
type dummyRunner struct {
	name  string
	delay time.Duration
}

func (d *dummyRunner) GetTaskName() string {
	return d.name
}

func (d *dummyRunner) Run(input interface{}, ts TaskSupport) (interface{}, error) {
	select {
	case <-ts.GetContext().Done():
		// canceled
		return nil, ts.GetContext().Err()
	case <-time.After(d.delay):
		// complete after delay
		return d.name, nil
	}
}

func TestForkTaskRunner_NonCompete(t *testing.T) {
	// Prepare a TaskSupport with a background context
	ts := newTaskSupport(withContext(context.Background()))

	// Two branches that complete immediately
	branches := []TaskRunner{
		&dummyRunner{name: "r1", delay: 0},
		&dummyRunner{name: "r2", delay: 0},
	}
	fork := ForkTaskRunner{
		Task: &model.ForkTask{
			Fork: model.ForkTaskConfiguration{
				Compete: false,
			},
		},
		TaskName:      "fork",
		BranchRunners: branches,
	}

	output, err := fork.Run("in", ts)
	assert.NoError(t, err)

	results, ok := output.([]interface{})
	assert.True(t, ok, "expected output to be []interface{}")
	assert.Equal(t, []interface{}{"r1", "r2"}, results)
}

func TestForkTaskRunner_Compete(t *testing.T) {
	// Prepare a TaskSupport with a background context
	ts := newTaskSupport(withContext(context.Background()))

	// One fast branch and one slow branch
	branches := []TaskRunner{
		&dummyRunner{name: "fast", delay: 10 * time.Millisecond},
		&dummyRunner{name: "slow", delay: 50 * time.Millisecond},
	}
	fork := ForkTaskRunner{
		Task: &model.ForkTask{
			Fork: model.ForkTaskConfiguration{
				Compete: true,
			},
		},
		TaskName:      "fork",
		BranchRunners: branches,
	}

	start := time.Now()
	output, err := fork.Run("in", ts)
	elapsed := time.Since(start)

	assert.NoError(t, err)
	assert.Equal(t, "fast", output)
	// ensure compete returns before the slow branch would finish
	assert.Less(t, elapsed, 50*time.Millisecond, "compete should cancel the slow branch")
}
