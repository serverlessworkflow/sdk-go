package impl

import (
	"context"
	"errors"
	"sync"
)

type ctxKey string

const runnerCtxKey ctxKey = "wfRunnerContext"

// WorkflowRunnerContext holds the necessary data for the workflow execution within the instance.
type WorkflowRunnerContext struct {
	mu               sync.Mutex
	input            interface{} // input can hold any type
	output           interface{} // output can hold any type
	context          map[string]interface{}
	StatusPhase      []StatusPhaseLog
	TasksStatusPhase map[string][]StatusPhaseLog // Holds `$context` as the key
}

func (runnerCtx *WorkflowRunnerContext) SetStatus(status StatusPhase) {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	if runnerCtx.StatusPhase == nil {
		runnerCtx.StatusPhase = []StatusPhaseLog{}
	}
	runnerCtx.StatusPhase = append(runnerCtx.StatusPhase, NewStatusPhaseLog(status))
}

func (runnerCtx *WorkflowRunnerContext) SetTaskStatus(task string, status StatusPhase) {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	if runnerCtx.TasksStatusPhase == nil {
		runnerCtx.TasksStatusPhase = map[string][]StatusPhaseLog{}
	}
	runnerCtx.TasksStatusPhase[task] = append(runnerCtx.TasksStatusPhase[task], NewStatusPhaseLog(status))
}

// SetWorkflowCtx safely sets the `$context` value
func (runnerCtx *WorkflowRunnerContext) SetWorkflowCtx(value interface{}) {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	if runnerCtx.context == nil {
		runnerCtx.context = make(map[string]interface{})
	}
	runnerCtx.context["$context"] = value
}

// GetWorkflowCtx safely retrieves the `$context` value
func (runnerCtx *WorkflowRunnerContext) GetWorkflowCtx() interface{} {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	if runnerCtx.context == nil {
		return nil
	}
	return runnerCtx.context["$context"]
}

// SetInput safely sets the input
func (runnerCtx *WorkflowRunnerContext) SetInput(input interface{}) {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	runnerCtx.input = input
}

// GetInput safely retrieves the input
func (runnerCtx *WorkflowRunnerContext) GetInput() interface{} {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	return runnerCtx.input
}

// SetOutput safely sets the output
func (runnerCtx *WorkflowRunnerContext) SetOutput(output interface{}) {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	runnerCtx.output = output
}

// GetOutput safely retrieves the output
func (runnerCtx *WorkflowRunnerContext) GetOutput() interface{} {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()
	return runnerCtx.output
}

// GetInputAsMap safely retrieves the input as a map[string]interface{}.
// If input is not a map, it creates a map with an empty string key and the input as the value.
func (runnerCtx *WorkflowRunnerContext) GetInputAsMap() map[string]interface{} {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()

	if inputMap, ok := runnerCtx.input.(map[string]interface{}); ok {
		return inputMap
	}

	// If input is not a map, create a map with an empty key and set input as the value
	return map[string]interface{}{
		"": runnerCtx.input,
	}
}

// GetOutputAsMap safely retrieves the output as a map[string]interface{}.
// If output is not a map, it creates a map with an empty string key and the output as the value.
func (runnerCtx *WorkflowRunnerContext) GetOutputAsMap() map[string]interface{} {
	runnerCtx.mu.Lock()
	defer runnerCtx.mu.Unlock()

	if outputMap, ok := runnerCtx.output.(map[string]interface{}); ok {
		return outputMap
	}

	// If output is not a map, create a map with an empty key and set output as the value
	return map[string]interface{}{
		"": runnerCtx.output,
	}
}

// WithRunnerContext adds the WorkflowRunnerContext to a parent context
func WithRunnerContext(parent context.Context, wfCtx *WorkflowRunnerContext) context.Context {
	return context.WithValue(parent, runnerCtxKey, wfCtx)
}

// GetRunnerContext retrieves the WorkflowRunnerContext from a context
func GetRunnerContext(ctx context.Context) (*WorkflowRunnerContext, error) {
	wfCtx, ok := ctx.Value(runnerCtxKey).(*WorkflowRunnerContext)
	if !ok {
		return nil, errors.New("workflow context not found")
	}
	return wfCtx, nil
}
