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

package ctx

import (
	"context"
	"errors"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"sync"
)

var ErrWorkflowContextNotFound = errors.New("workflow context not found")

var _ WorkflowContext = &workflowContext{}

type ctxKey string

const (
	runnerCtxKey ctxKey = "wfRunnerContext"
	varsContext         = "$context"
	varsInput           = "$input"
	varsOutput          = "$output"
	varsWorkflow        = "$workflow"
)

type WorkflowContext interface {
	SetStatus(status StatusPhase)
	SetTaskStatus(task string, status StatusPhase)
	SetInstanceCtx(value interface{})
	GetInstanceCtx() interface{}
	SetInput(input interface{})
	GetInput() interface{}
	SetOutput(output interface{})
	GetOutput() interface{}
	GetOutputAsMap() map[string]interface{}
	AsJQVars() map[string]interface{}
}

// workflowContext holds the necessary data for the workflow execution within the instance.
type workflowContext struct {
	mu               sync.Mutex
	input            interface{}            // $input can hold any type
	output           interface{}            // $output can hold any type
	context          map[string]interface{} // Holds `$context` as the key
	definition       map[string]interface{} // $workflow representation in the context
	StatusPhase      []StatusPhaseLog
	TasksStatusPhase map[string][]StatusPhaseLog
}

func NewWorkflowContext(workflow *model.Workflow) (WorkflowContext, error) {
	workflowCtx := &workflowContext{}
	workflowDef, err := workflow.AsMap()
	if err != nil {
		return nil, err
	}

	workflowCtx.definition = workflowDef
	workflowCtx.SetStatus(PendingStatus)

	return workflowCtx, nil
}

func (ctx *workflowContext) AsJQVars() map[string]interface{} {
	vars := make(map[string]interface{})
	vars[varsInput] = ctx.GetInput()
	vars[varsOutput] = ctx.GetOutput()
	vars[varsContext] = ctx.GetInstanceCtx()
	vars[varsOutput] = ctx.definition
	return vars
}

func (ctx *workflowContext) SetStatus(status StatusPhase) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.StatusPhase == nil {
		ctx.StatusPhase = []StatusPhaseLog{}
	}
	ctx.StatusPhase = append(ctx.StatusPhase, NewStatusPhaseLog(status))
}

func (ctx *workflowContext) SetTaskStatus(task string, status StatusPhase) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.TasksStatusPhase == nil {
		ctx.TasksStatusPhase = map[string][]StatusPhaseLog{}
	}
	ctx.TasksStatusPhase[task] = append(ctx.TasksStatusPhase[task], NewStatusPhaseLog(status))
}

// SetInstanceCtx safely sets the `$context` value
func (ctx *workflowContext) SetInstanceCtx(value interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.context == nil {
		ctx.context = make(map[string]interface{})
	}
	ctx.context[varsContext] = value
}

// GetInstanceCtx safely retrieves the `$context` value
func (ctx *workflowContext) GetInstanceCtx() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.context == nil {
		return nil
	}
	return ctx.context[varsContext]
}

// SetInput safely sets the input
func (ctx *workflowContext) SetInput(input interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.input = input
}

// GetInput safely retrieves the input
func (ctx *workflowContext) GetInput() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.input
}

// SetOutput safely sets the output
func (ctx *workflowContext) SetOutput(output interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.output = output
}

// GetOutput safely retrieves the output
func (ctx *workflowContext) GetOutput() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.output
}

// GetInputAsMap safely retrieves the input as a map[string]interface{}.
// If input is not a map, it creates a map with an empty string key and the input as the value.
func (ctx *workflowContext) GetInputAsMap() map[string]interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if inputMap, ok := ctx.input.(map[string]interface{}); ok {
		return inputMap
	}

	// If input is not a map, create a map with an empty key and set input as the value
	return map[string]interface{}{
		"": ctx.input,
	}
}

// GetOutputAsMap safely retrieves the output as a map[string]interface{}.
// If output is not a map, it creates a map with an empty string key and the output as the value.
func (ctx *workflowContext) GetOutputAsMap() map[string]interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if outputMap, ok := ctx.output.(map[string]interface{}); ok {
		return outputMap
	}

	// If output is not a map, create a map with an empty key and set output as the value
	return map[string]interface{}{
		"": ctx.output,
	}
}

// WithWorkflowContext adds the workflowContext to a parent context
func WithWorkflowContext(parent context.Context, wfCtx WorkflowContext) context.Context {
	return context.WithValue(parent, runnerCtxKey, wfCtx)
}

// GetWorkflowContext retrieves the workflowContext from a context
func GetWorkflowContext(ctx context.Context) (WorkflowContext, error) {
	wfCtx, ok := ctx.Value(runnerCtxKey).(*workflowContext)
	if !ok {
		return nil, ErrWorkflowContextNotFound
	}
	return wfCtx, nil
}
