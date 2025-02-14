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

	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ WorkflowRunner = &workflowRunnerImpl{}

type WorkflowRunner interface {
	GetWorkflowDef() *model.Workflow
	Run(input interface{}) (output interface{}, err error)
	GetContext() *WorkflowContext
}

func NewDefaultRunner(workflow *model.Workflow) WorkflowRunner {
	wfContext := &WorkflowContext{}
	wfContext.SetStatus(PendingStatus)
	// TODO: based on the workflow definition, the context might change.
	ctx := WithWorkflowContext(context.Background(), wfContext)
	return &workflowRunnerImpl{
		Workflow:  workflow,
		Context:   ctx,
		RunnerCtx: wfContext,
	}
}

type workflowRunnerImpl struct {
	Workflow  *model.Workflow
	Context   context.Context
	RunnerCtx *WorkflowContext
}

func (wr *workflowRunnerImpl) GetContext() *WorkflowContext {
	return wr.RunnerCtx
}

func (wr *workflowRunnerImpl) GetTaskContext() TaskContext {
	return wr.RunnerCtx
}

func (wr *workflowRunnerImpl) GetWorkflowDef() *model.Workflow {
	return wr.Workflow
}

// Run executes the workflow synchronously.
func (wr *workflowRunnerImpl) Run(input interface{}) (output interface{}, err error) {
	defer func() {
		if err != nil {
			wr.RunnerCtx.SetStatus(FaultedStatus)
			err = wr.wrapWorkflowError(err, "/")
		}
	}()

	// Process input
	if input, err = wr.processInput(input); err != nil {
		return nil, err
	}

	wr.RunnerCtx.SetInput(input)
	// Run tasks sequentially
	wr.RunnerCtx.SetStatus(RunningStatus)
	doRunner, err := NewDoTaskRunner(wr.Workflow.Do, wr)
	if err != nil {
		return nil, err
	}
	output, err = doRunner.Run(wr.RunnerCtx.GetInput())
	if err != nil {
		return nil, err
	}

	// Process output
	if output, err = wr.processOutput(output); err != nil {
		return nil, err
	}

	wr.RunnerCtx.SetOutput(output)
	wr.RunnerCtx.SetStatus(CompletedStatus)
	return output, nil
}

// wrapWorkflowError ensures workflow errors have a proper instance reference.
func (wr *workflowRunnerImpl) wrapWorkflowError(err error, taskName string) error {
	if knownErr := model.AsError(err); knownErr != nil {
		return knownErr.WithInstanceRef(wr.Workflow, taskName)
	}
	return model.NewErrRuntime(fmt.Errorf("workflow '%s', task '%s': %w", wr.Workflow.Document.Name, taskName, err), taskName)
}

// processInput validates and transforms input if needed.
func (wr *workflowRunnerImpl) processInput(input interface{}) (output interface{}, err error) {
	if wr.Workflow.Input != nil {
		output, err = processIO(input, wr.Workflow.Input.Schema, wr.Workflow.Input.From, "/")
		if err != nil {
			return nil, err
		}
		return output, nil
	}
	return input, nil
}

// processOutput applies output transformations.
func (wr *workflowRunnerImpl) processOutput(output interface{}) (interface{}, error) {
	if wr.Workflow.Output != nil {
		return processIO(output, wr.Workflow.Output.Schema, wr.Workflow.Output.As, "/")
	}
	return output, nil
}
