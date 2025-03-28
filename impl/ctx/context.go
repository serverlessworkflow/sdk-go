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
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/serverlessworkflow/sdk-go/v3/model"
	"sync"
	"time"
)

var ErrWorkflowContextNotFound = errors.New("workflow context not found")

var _ WorkflowContext = &workflowContext{}

type ctxKey string

const (
	runnerCtxKey ctxKey = "wfRunnerContext"

	varsContext  = "$context"
	varsInput    = "$input"
	varsOutput   = "$output"
	varsWorkflow = "$workflow"
	varsRuntime  = "$runtime"
	varsTask     = "$task"

	// TODO: script during the release to update this value programmatically
	runtimeVersion = "v3.1.0"
	runtimeName    = "CNCF Serverless Workflow Specification Go SDK"
)

type WorkflowContext interface {
	SetStartedAt(t time.Time)
	SetStatus(status StatusPhase)
	SetRawInput(input interface{})
	SetInstanceCtx(value interface{})
	GetInstanceCtx() interface{}
	SetInput(input interface{})
	GetInput() interface{}
	SetOutput(output interface{})
	GetOutput() interface{}
	GetOutputAsMap() map[string]interface{}
	GetVars() map[string]interface{}
	SetTaskStatus(task string, status StatusPhase)
	SetTaskRawInput(input interface{})
	SetTaskRawOutput(output interface{})
	SetTaskDef(task model.Task) error
	SetTaskStartedAt(startedAt time.Time)
	SetTaskName(name string)
	SetTaskReference(ref string)
	GetTaskReference() string
	ClearTaskContext()
	SetLocalExprVars(vars map[string]interface{})
	AddLocalExprVars(vars map[string]interface{})
	RemoveLocalExprVars(keys ...string)
}

// workflowContext holds the necessary data for the workflow execution within the instance.
type workflowContext struct {
	mu                 sync.Mutex
	input              interface{}            // $input can hold any type
	output             interface{}            // $output can hold any type
	context            map[string]interface{} // Holds `$context` as the key
	workflowDescriptor map[string]interface{} // $workflow representation in the context
	taskDescriptor     map[string]interface{} // $task representation in the context
	localExprVars      map[string]interface{} // Local expression variables defined in a given task or private context. E.g. a For task $item.
	StatusPhase        []StatusPhaseLog
	TasksStatusPhase   map[string][]StatusPhaseLog
}

func NewWorkflowContext(workflow *model.Workflow) (WorkflowContext, error) {
	workflowCtx := &workflowContext{}
	workflowDef, err := workflow.AsMap()
	if err != nil {
		return nil, err
	}
	workflowCtx.taskDescriptor = map[string]interface{}{}
	workflowCtx.workflowDescriptor = map[string]interface{}{
		varsWorkflow: map[string]interface{}{
			"id":         uuid.NewString(),
			"definition": workflowDef,
		},
	}
	workflowCtx.SetStatus(PendingStatus)

	return workflowCtx, nil
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

func (ctx *workflowContext) SetStartedAt(t time.Time) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	wf, ok := ctx.workflowDescriptor[varsWorkflow].(map[string]interface{})
	if !ok {
		wf = make(map[string]interface{})
		ctx.workflowDescriptor[varsWorkflow] = wf
	}

	startedAt, ok := wf["startedAt"].(map[string]interface{})
	if !ok {
		startedAt = make(map[string]interface{})
		wf["startedAt"] = startedAt
	}

	startedAt["iso8601"] = t.UTC().Format(time.RFC3339)
}

func (ctx *workflowContext) SetRawInput(input interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	// Ensure the outer "workflow" map
	wf, ok := ctx.workflowDescriptor[varsWorkflow].(map[string]interface{})
	if !ok {
		wf = make(map[string]interface{})
		ctx.workflowDescriptor[varsWorkflow] = wf
	}

	// Store the input
	wf["input"] = input
}

func (ctx *workflowContext) AddLocalExprVars(vars map[string]interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.localExprVars == nil {
		ctx.localExprVars = map[string]interface{}{}
	}
	for k, v := range vars {
		ctx.localExprVars[k] = v
	}
}

func (ctx *workflowContext) RemoveLocalExprVars(keys ...string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if ctx.localExprVars == nil {
		return
	}

	for _, k := range keys {
		delete(ctx.localExprVars, k)
	}
}

func (ctx *workflowContext) SetLocalExprVars(vars map[string]interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.localExprVars = vars
}

func (ctx *workflowContext) GetVars() map[string]interface{} {
	vars := make(map[string]interface{})
	vars[varsInput] = ctx.GetInput()
	vars[varsOutput] = ctx.GetOutput()
	vars[varsContext] = ctx.GetInstanceCtx()
	vars[varsTask] = ctx.taskDescriptor[varsTask]
	vars[varsWorkflow] = ctx.workflowDescriptor[varsWorkflow]
	vars[varsRuntime] = map[string]interface{}{
		"name":    runtimeName,
		"version": runtimeVersion,
	}
	for varName, varValue := range ctx.localExprVars {
		vars[varName] = varValue
	}
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

func (ctx *workflowContext) SetTaskStatus(task string, status StatusPhase) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.TasksStatusPhase == nil {
		ctx.TasksStatusPhase = map[string][]StatusPhaseLog{}
	}
	ctx.TasksStatusPhase[task] = append(ctx.TasksStatusPhase[task], NewStatusPhaseLog(status))
}

func (ctx *workflowContext) SetTaskRawInput(input interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		task = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = task
	}

	task["input"] = input
}

func (ctx *workflowContext) SetTaskRawOutput(output interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		task = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = task
	}

	task["output"] = output
}

func (ctx *workflowContext) SetTaskDef(task model.Task) error {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	if task == nil {
		return errors.New("SetTaskDef called with nil model.Task")
	}

	defBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("failed to marshal task: %w", err)
	}

	var defMap map[string]interface{}
	if err := json.Unmarshal(defBytes, &defMap); err != nil {
		return fmt.Errorf("failed to unmarshal task into map: %w", err)
	}

	taskMap, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		taskMap = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = taskMap
	}

	taskMap["definition"] = defMap

	return nil
}

func (ctx *workflowContext) SetTaskStartedAt(startedAt time.Time) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		task = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = task
	}

	task["startedAt"] = startedAt.UTC().Format(time.RFC3339)
}

func (ctx *workflowContext) SetTaskName(name string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		task = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = task
	}

	task["name"] = name
}

func (ctx *workflowContext) SetTaskReference(ref string) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()

	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		task = make(map[string]interface{})
		ctx.taskDescriptor[varsTask] = task
	}

	task["reference"] = ref
}

func (ctx *workflowContext) GetTaskReference() string {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	task, ok := ctx.taskDescriptor[varsTask].(map[string]interface{})
	if !ok {
		return ""
	}
	return task["reference"].(string)
}

func (ctx *workflowContext) ClearTaskContext() {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.taskDescriptor[varsTask] = make(map[string]interface{})
}
