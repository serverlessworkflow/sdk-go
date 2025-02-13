package impl

import (
	"context"
	"errors"
	"sync"
)

type ctxKey string

const runnerCtxKey ctxKey = "wfRunnerContext"

// WorkflowContext holds the necessary data for the workflow execution within the instance.
type WorkflowContext struct {
	mu               sync.Mutex
	input            interface{} // input can hold any type
	output           interface{} // output can hold any type
	context          map[string]interface{}
	StatusPhase      []StatusPhaseLog
	TasksStatusPhase map[string][]StatusPhaseLog // Holds `$context` as the key
}

type TaskContext interface {
	SetTaskStatus(task string, status StatusPhase)
}

func (ctx *WorkflowContext) SetStatus(status StatusPhase) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.StatusPhase == nil {
		ctx.StatusPhase = []StatusPhaseLog{}
	}
	ctx.StatusPhase = append(ctx.StatusPhase, NewStatusPhaseLog(status))
}

func (ctx *WorkflowContext) SetTaskStatus(task string, status StatusPhase) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.TasksStatusPhase == nil {
		ctx.TasksStatusPhase = map[string][]StatusPhaseLog{}
	}
	ctx.TasksStatusPhase[task] = append(ctx.TasksStatusPhase[task], NewStatusPhaseLog(status))
}

// SetInstanceCtx safely sets the `$context` value
func (ctx *WorkflowContext) SetInstanceCtx(value interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.context == nil {
		ctx.context = make(map[string]interface{})
	}
	ctx.context["$context"] = value
}

// GetInstanceCtx safely retrieves the `$context` value
func (ctx *WorkflowContext) GetInstanceCtx() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	if ctx.context == nil {
		return nil
	}
	return ctx.context["$context"]
}

// SetInput safely sets the input
func (ctx *WorkflowContext) SetInput(input interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.input = input
}

// GetInput safely retrieves the input
func (ctx *WorkflowContext) GetInput() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.input
}

// SetOutput safely sets the output
func (ctx *WorkflowContext) SetOutput(output interface{}) {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	ctx.output = output
}

// GetOutput safely retrieves the output
func (ctx *WorkflowContext) GetOutput() interface{} {
	ctx.mu.Lock()
	defer ctx.mu.Unlock()
	return ctx.output
}

// GetInputAsMap safely retrieves the input as a map[string]interface{}.
// If input is not a map, it creates a map with an empty string key and the input as the value.
func (ctx *WorkflowContext) GetInputAsMap() map[string]interface{} {
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
func (ctx *WorkflowContext) GetOutputAsMap() map[string]interface{} {
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

// WithWorkflowContext adds the WorkflowContext to a parent context
func WithWorkflowContext(parent context.Context, wfCtx *WorkflowContext) context.Context {
	return context.WithValue(parent, runnerCtxKey, wfCtx)
}

// GetWorkflowContext retrieves the WorkflowContext from a context
func GetWorkflowContext(ctx context.Context) (*WorkflowContext, error) {
	wfCtx, ok := ctx.Value(runnerCtxKey).(*WorkflowContext)
	if !ok {
		return nil, errors.New("workflow context not found")
	}
	return wfCtx, nil
}
