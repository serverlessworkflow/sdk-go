package impl

import (
	"context"
	"errors"
	"sync"
)

type ctxKey string

const executorCtxKey ctxKey = "executorContext"

// ExecutorContext to not confound with Workflow Context as "$context" in the specification.
// This holds the necessary data for the workflow execution within the instance.
type ExecutorContext struct {
	mu     sync.Mutex
	Input  map[string]interface{}
	Output map[string]interface{}
	// Context or `$context` passed through the task executions see https://github.com/serverlessworkflow/specification/blob/main/dsl.md#data-flow
	Context map[string]interface{}
}

// SetWorkflowCtx safely sets the $context
func (execCtx *ExecutorContext) SetWorkflowCtx(wfCtx map[string]interface{}) {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	execCtx.Context = wfCtx
}

// GetWorkflowCtx safely retrieves the $context
func (execCtx *ExecutorContext) GetWorkflowCtx() map[string]interface{} {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	return execCtx.Context
}

// SetInput safely sets the input map
func (execCtx *ExecutorContext) SetInput(input map[string]interface{}) {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	execCtx.Input = input
}

// GetInput safely retrieves the input map
func (execCtx *ExecutorContext) GetInput() map[string]interface{} {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	return execCtx.Input
}

// SetOutput safely sets the output map
func (execCtx *ExecutorContext) SetOutput(output map[string]interface{}) {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	execCtx.Output = output
}

// GetOutput safely retrieves the output map
func (execCtx *ExecutorContext) GetOutput() map[string]interface{} {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	return execCtx.Output
}

// UpdateOutput allows adding or updating a single key-value pair in the output map
func (execCtx *ExecutorContext) UpdateOutput(key string, value interface{}) {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	if execCtx.Output == nil {
		execCtx.Output = make(map[string]interface{})
	}
	execCtx.Output[key] = value
}

// GetOutputValue safely retrieves a single key from the output map
func (execCtx *ExecutorContext) GetOutputValue(key string) (interface{}, bool) {
	execCtx.mu.Lock()
	defer execCtx.mu.Unlock()
	value, exists := execCtx.Output[key]
	return value, exists
}

func WithExecutorContext(parent context.Context, wfCtx *ExecutorContext) context.Context {
	return context.WithValue(parent, executorCtxKey, wfCtx)
}

func GetExecutorContext(ctx context.Context) (*ExecutorContext, error) {
	wfCtx, ok := ctx.Value(executorCtxKey).(*ExecutorContext)
	if !ok {
		return nil, errors.New("workflow context not found")
	}
	return wfCtx, nil
}
