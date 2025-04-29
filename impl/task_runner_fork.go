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
	"fmt"
	"sync"

	"github.com/serverlessworkflow/sdk-go/v3/model"
)

func NewForkTaskRunner(taskName string, task *model.ForkTask, workflowDef *model.Workflow) (*ForkTaskRunner, error) {
	if task == nil || task.Fork.Branches == nil {
		return nil, model.NewErrValidation(fmt.Errorf("invalid Fork task %s", taskName), taskName)
	}

	var runners []TaskRunner
	for _, branchItem := range *task.Fork.Branches {
		r, err := NewTaskRunner(branchItem.Key, branchItem.Task, workflowDef)
		if err != nil {
			return nil, err
		}
		runners = append(runners, r)
	}

	return &ForkTaskRunner{
		Task:          task,
		TaskName:      taskName,
		BranchRunners: runners,
	}, nil
}

type ForkTaskRunner struct {
	Task          *model.ForkTask
	TaskName      string
	BranchRunners []TaskRunner
}

func (f ForkTaskRunner) GetTaskName() string {
	return f.TaskName
}

func (f ForkTaskRunner) Run(input interface{}, parentSupport TaskSupport) (interface{}, error) {
	cancelCtx, cancel := context.WithCancel(parentSupport.GetContext())
	defer cancel()

	n := len(f.BranchRunners)
	results := make([]interface{}, n)
	errs := make(chan error, n)
	done := make(chan struct{})
	resultCh := make(chan interface{}, 1)

	var (
		wg   sync.WaitGroup
		once sync.Once // <-- declare a Once
	)

	for i, runner := range f.BranchRunners {
		wg.Add(1)
		go func(i int, runner TaskRunner) {
			defer wg.Done()
			// **Isolate context** for each branch!
			branchSupport := parentSupport.CloneWithContext(cancelCtx)

			select {
			case <-cancelCtx.Done():
				return
			default:
			}

			out, err := runner.Run(input, branchSupport)
			if err != nil {
				errs <- err
				return
			}
			results[i] = out

			if f.Task.Fork.Compete {
				select {
				case resultCh <- out:
					once.Do(func() {
						cancel()    // **signal cancellation** to all other branches
						close(done) // signal we have a winner
					})
				default:
				}
			}
		}(i, runner)
	}

	if f.Task.Fork.Compete {
		select {
		case <-done:
			return <-resultCh, nil
		case err := <-errs:
			return nil, err
		}
	}

	wg.Wait()
	select {
	case err := <-errs:
		return nil, err
	default:
	}
	return results, nil
}
