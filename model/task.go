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

package model

import (
	"encoding/json"
	"errors"
	"fmt"
)

type TaskBase struct {
	// A runtime expression, if any, used to determine whether or not the task should be run.
	If *RuntimeExpression `json:"if,omitempty" validate:"omitempty"`
	// Configure the task's input.
	Input *Input `json:"input,omitempty" validate:"omitempty"`
	// Configure the task's output.
	Output *Output `json:"output,omitempty" validate:"omitempty"`
	// Export task output to context.
	Export  *Export             `json:"export,omitempty" validate:"omitempty"`
	Timeout *TimeoutOrReference `json:"timeout,omitempty" validate:"omitempty"`
	// The flow directive to be performed upon completion of the task.
	Then     *FlowDirective         `json:"then,omitempty" validate:"omitempty"`
	Metadata map[string]interface{} `json:"metadata,omitempty"`
}

// Task represents a discrete unit of work in a workflow.
type Task interface {
	GetBase() *TaskBase
}

type NamedTaskMap map[string]Task

// UnmarshalJSON for NamedTaskMap to ensure proper deserialization.
func (ntm *NamedTaskMap) UnmarshalJSON(data []byte) error {
	var rawTasks map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawTasks); err != nil {
		return err
	}

	for name, raw := range rawTasks {
		task, err := unmarshalTask(name, raw)
		if err != nil {
			return err
		}

		if *ntm == nil {
			*ntm = make(map[string]Task)
		}
		(*ntm)[name] = task
	}

	return nil
}

// TaskList represents a list of named tasks to perform.
type TaskList []*TaskItem

// Next gets the next item in the list based on the current index
func (tl *TaskList) Next(currentIdx int) (int, *TaskItem) {
	if currentIdx == -1 || currentIdx >= len(*tl) {
		return -1, nil
	}

	current := (*tl)[currentIdx]
	if current.GetBase() != nil && current.GetBase().Then != nil {
		then := current.GetBase().Then
		if then.IsTermination() {
			return -1, nil
		}
		return tl.KeyAndIndex(then.Value)
	}

	// Proceed sequentially if no 'then' is specified
	if currentIdx+1 < len(*tl) {
		return currentIdx + 1, (*tl)[currentIdx+1]
	}
	return -1, nil
}

// UnmarshalJSON for TaskList to ensure proper deserialization.
func (tl *TaskList) UnmarshalJSON(data []byte) error {
	var rawTasks []json.RawMessage
	if err := json.Unmarshal(data, &rawTasks); err != nil {
		return err
	}

	for _, raw := range rawTasks {
		var taskItemRaw map[string]json.RawMessage
		if err := json.Unmarshal(raw, &taskItemRaw); err != nil {
			return err
		}

		if len(taskItemRaw) != 1 {
			return errors.New("each TaskItem must have exactly one key")
		}

		for key, taskRaw := range taskItemRaw {
			task, err := unmarshalTask(key, taskRaw)
			if err != nil {
				return err
			}
			*tl = append(*tl, &TaskItem{Key: key, Task: task})
		}
	}

	return nil
}

var taskTypeRegistry = map[string]func() Task{
	"call_http":     func() Task { return &CallHTTP{} },
	"call_openapi":  func() Task { return &CallOpenAPI{} },
	"call_grpc":     func() Task { return &CallGRPC{} },
	"call_asyncapi": func() Task { return &CallAsyncAPI{} },
	"call":          func() Task { return &CallFunction{} },
	"do":            func() Task { return &DoTask{} },
	"fork":          func() Task { return &ForkTask{} },
	"emit":          func() Task { return &EmitTask{} },
	"for":           func() Task { return &ForTask{} },
	"listen":        func() Task { return &ListenTask{} },
	"raise":         func() Task { return &RaiseTask{} },
	"run":           func() Task { return &RunTask{} },
	"set":           func() Task { return &SetTask{} },
	"switch":        func() Task { return &SwitchTask{} },
	"try":           func() Task { return &TryTask{} },
	"wait":          func() Task { return &WaitTask{} },
}

func unmarshalTask(key string, taskRaw json.RawMessage) (Task, error) {
	var taskType map[string]interface{}
	if err := json.Unmarshal(taskRaw, &taskType); err != nil {
		return nil, fmt.Errorf("failed to parse task type for key '%s': %w", key, err)
	}

	// Determine task type
	var task Task
	if callValue, hasCall := taskType["call"].(string); hasCall {
		// Form composite key and check if it's in the registry
		registryKey := fmt.Sprintf("call_%s", callValue)
		if constructor, exists := taskTypeRegistry[registryKey]; exists {
			task = constructor()
		} else {
			// Default to CallFunction for unrecognized call values
			task = &CallFunction{}
		}
	} else {
		// Handle non-call tasks (e.g., "do", "fork")
		for typeKey := range taskType {
			if constructor, exists := taskTypeRegistry[typeKey]; exists {
				task = constructor()
				break
			}
		}
	}

	if task == nil {
		return nil, fmt.Errorf("unknown task type for key '%s'", key)
	}

	// Populate the task with raw data
	if err := json.Unmarshal(taskRaw, task); err != nil {
		return nil, fmt.Errorf("failed to unmarshal task '%s': %w", key, err)
	}

	return task, nil
}

// MarshalJSON for TaskList to ensure proper serialization.
func (tl *TaskList) MarshalJSON() ([]byte, error) {
	return json.Marshal([]*TaskItem(*tl))
}

// Key retrieves a TaskItem by its key.
func (tl *TaskList) Key(key string) *TaskItem {
	_, keyItem := tl.KeyAndIndex(key)
	return keyItem
}

func (tl *TaskList) KeyAndIndex(key string) (int, *TaskItem) {
	for i, item := range *tl {
		if item.Key == key {
			return i, item
		}
	}
	// TODO: Add logging here for missing task references
	return -1, nil
}

// TaskItem represents a named task and its associated definition.
type TaskItem struct {
	Key  string `json:"-" validate:"required"`
	Task Task   `json:"-" validate:"required"`
}

// MarshalJSON for TaskItem to ensure proper serialization as a key-value pair.
func (ti *TaskItem) MarshalJSON() ([]byte, error) {
	if ti == nil {
		return nil, fmt.Errorf("cannot marshal a nil TaskItem")
	}

	// Serialize the Task
	taskJSON, err := json.Marshal(ti.Task)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal task: %w", err)
	}

	// Create a map with the Key and Task
	taskEntry := map[string]json.RawMessage{
		ti.Key: taskJSON,
	}

	// Marshal the map into JSON
	return json.Marshal(taskEntry)
}

func (ti *TaskItem) GetBase() *TaskBase {
	return ti.Task.GetBase()
}

// AsCallHTTPTask casts the Task to a CallTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsCallHTTPTask() *CallHTTP {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*CallHTTP); ok {
		return task
	}
	return nil
}

// AsCallOpenAPITask casts the Task to a CallOpenAPI task if possible, returning nil if the cast fails.
func (ti *TaskItem) AsCallOpenAPITask() *CallOpenAPI {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*CallOpenAPI); ok {
		return task
	}
	return nil
}

// AsCallGRPCTask casts the Task to a CallGRPC task if possible, returning nil if the cast fails.
func (ti *TaskItem) AsCallGRPCTask() *CallGRPC {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*CallGRPC); ok {
		return task
	}
	return nil
}

// AsCallAsyncAPITask casts the Task to a CallAsyncAPI task if possible, returning nil if the cast fails.
func (ti *TaskItem) AsCallAsyncAPITask() *CallAsyncAPI {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*CallAsyncAPI); ok {
		return task
	}
	return nil
}

// AsCallFunctionTask casts the Task to a CallFunction task if possible, returning nil if the cast fails.
func (ti *TaskItem) AsCallFunctionTask() *CallFunction {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*CallFunction); ok {
		return task
	}
	return nil
}

// AsDoTask casts the Task to a DoTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsDoTask() *DoTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*DoTask); ok {
		return task
	}
	return nil
}

// AsForkTask casts the Task to a ForkTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsForkTask() *ForkTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*ForkTask); ok {
		return task
	}
	return nil
}

// AsEmitTask casts the Task to an EmitTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsEmitTask() *EmitTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*EmitTask); ok {
		return task
	}
	return nil
}

// AsForTask casts the Task to a ForTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsForTask() *ForTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*ForTask); ok {
		return task
	}
	return nil
}

// AsListenTask casts the Task to a ListenTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsListenTask() *ListenTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*ListenTask); ok {
		return task
	}
	return nil
}

// AsRaiseTask casts the Task to a RaiseTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsRaiseTask() *RaiseTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*RaiseTask); ok {
		return task
	}
	return nil
}

// AsRunTask casts the Task to a RunTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsRunTask() *RunTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*RunTask); ok {
		return task
	}
	return nil
}

// AsSetTask casts the Task to a SetTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsSetTask() *SetTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*SetTask); ok {
		return task
	}
	return nil
}

// AsSwitchTask casts the Task to a SwitchTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsSwitchTask() *SwitchTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*SwitchTask); ok {
		return task
	}
	return nil
}

// AsTryTask casts the Task to a TryTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsTryTask() *TryTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*TryTask); ok {
		return task
	}
	return nil
}

// AsWaitTask casts the Task to a WaitTask if possible, returning nil if the cast fails.
func (ti *TaskItem) AsWaitTask() *WaitTask {
	if ti == nil {
		return nil
	}
	if task, ok := ti.Task.(*WaitTask); ok {
		return task
	}
	return nil
}
