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

// This file contains comprehensive tests for the task registration system.
// The tests cover:
// - Thread-safe registration and access
// - Error handling and validation
// - Integration with the unmarshaling process
// - Concurrent access patterns
// - Performance benchmarks
// - Global registry functions

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"
	"strings"
	"sync"
	"testing"
)

// TestTask is a simple test implementation of the Task interface
type TestTask struct {
	TaskBase `json:",inline"`
	TestData string `json:"test_data,omitempty"`
}

func (t *TestTask) GetBase() *TaskBase {
	return &t.TaskBase
}

// AnotherTestTask is another test implementation
type AnotherTestTask struct {
	TaskBase   `json:",inline"`
	OtherField int `json:"other_field,omitempty"`
}

func (t *AnotherTestTask) GetBase() *TaskBase {
	return &t.TaskBase
}

func TestNewTaskRegistry(t *testing.T) {
	registry := NewTaskRegistry()
	if registry == nil {
		t.Fatal("NewTaskRegistry returned nil")
	}
	if registry.constructors == nil {
		t.Fatal("constructors map not initialized")
	}
	if len(registry.constructors) != 0 {
		t.Fatal("new registry should be empty")
	}
}

func TestTaskRegistry_RegisterTask(t *testing.T) {
	tests := []struct {
		name          string
		taskType      string
		constructor   TaskConstructor
		expectError   bool
		expectedError string
	}{
		{
			name:        "Valid registration",
			taskType:    "test_task",
			constructor: func() Task { return &TestTask{} },
			expectError: false,
		},
		{
			name:          "Empty task type",
			taskType:      "",
			constructor:   func() Task { return &TestTask{} },
			expectError:   true,
			expectedError: "task type cannot be empty",
		},
		{
			name:          "Nil constructor",
			taskType:      "test_task",
			constructor:   nil,
			expectError:   true,
			expectedError: "constructor function cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewTaskRegistry()
			err := registry.RegisterTask(tt.taskType, tt.constructor)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got nil")
				} else if err.Error() != tt.expectedError {
					t.Errorf("expected error '%s', got '%s'", tt.expectedError, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				// Verify the constructor was registered
				constructor, exists := registry.GetConstructor(tt.taskType)
				if !exists {
					t.Errorf("task type '%s' was not registered", tt.taskType)
				}
				if constructor == nil {
					t.Errorf("registered constructor is nil")
				}
			}
		})
	}
}

func TestTaskRegistry_RegisterTask_Duplicate(t *testing.T) {
	registry := NewTaskRegistry()
	taskType := "duplicate_task"
	constructor := func() Task { return &TestTask{} }

	// First registration should succeed
	err := registry.RegisterTask(taskType, constructor)
	if err != nil {
		t.Fatalf("first registration failed: %v", err)
	}

	// Second registration should fail
	err = registry.RegisterTask(taskType, constructor)
	if err == nil {
		t.Fatal("expected error for duplicate registration")
	}

	expectedError := "task type 'duplicate_task' is already registered"
	if err.Error() != expectedError {
		t.Errorf("expected error '%s', got '%s'", expectedError, err.Error())
	}
}

func TestTaskRegistry_GetConstructor(t *testing.T) {
	registry := NewTaskRegistry()
	taskType := "get_test_task"
	constructor := func() Task { return &TestTask{} }

	// Test getting non-existent constructor
	_, exists := registry.GetConstructor(taskType)
	if exists {
		t.Errorf("constructor should not exist yet")
	}

	// Register and test getting existing constructor
	err := registry.RegisterTask(taskType, constructor)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	gotConstructor, exists := registry.GetConstructor(taskType)
	if !exists {
		t.Errorf("constructor should exist after registration")
	}
	if gotConstructor == nil {
		t.Errorf("returned constructor should not be nil")
	}

	// Test that the constructor actually works
	task := gotConstructor()
	if task == nil {
		t.Errorf("constructor should return a task")
	}
	if _, ok := task.(*TestTask); !ok {
		t.Errorf("constructor should return a TestTask")
	}
}

func TestTaskRegistry_UnregisterTask(t *testing.T) {
	registry := NewTaskRegistry()
	taskType := "unregister_test_task"
	constructor := func() Task { return &TestTask{} }

	// Register task
	err := registry.RegisterTask(taskType, constructor)
	if err != nil {
		t.Fatalf("registration failed: %v", err)
	}

	// Verify it exists
	_, exists := registry.GetConstructor(taskType)
	if !exists {
		t.Fatal("task should exist after registration")
	}

	// Unregister task
	registry.UnregisterTask(taskType)

	// Verify it no longer exists
	_, exists = registry.GetConstructor(taskType)
	if exists {
		t.Error("task should not exist after unregistration")
	}

	// Unregistering non-existent task should not panic
	registry.UnregisterTask("non_existent_task")
}

func TestTaskRegistry_ListRegisteredTypes(t *testing.T) {
	registry := NewTaskRegistry()

	// Test empty registry
	types := registry.ListRegisteredTypes()
	if len(types) != 0 {
		t.Errorf("empty registry should return empty slice, got %v", types)
	}

	// Register some tasks
	tasks := map[string]TaskConstructor{
		"task_a": func() Task { return &TestTask{} },
		"task_b": func() Task { return &AnotherTestTask{} },
		"task_c": func() Task { return &TestTask{} },
	}

	for taskType, constructor := range tasks {
		err := registry.RegisterTask(taskType, constructor)
		if err != nil {
			t.Fatalf("failed to register task '%s': %v", taskType, err)
		}
	}

	// Test listing registered types
	types = registry.ListRegisteredTypes()
	if len(types) != len(tasks) {
		t.Errorf("expected %d types, got %d", len(tasks), len(types))
	}

	// Sort both slices for comparison
	expectedTypes := make([]string, 0, len(tasks))
	for taskType := range tasks {
		expectedTypes = append(expectedTypes, taskType)
	}
	sort.Strings(expectedTypes)
	sort.Strings(types)

	for i, expected := range expectedTypes {
		if types[i] != expected {
			t.Errorf("expected type '%s' at index %d, got '%s'", expected, i, types[i])
		}
	}
}

func TestTaskRegistry_ConcurrentAccess(t *testing.T) {
	registry := NewTaskRegistry()
	numGoroutines := 10
	numOperationsPerGoroutine := 100

	var wg sync.WaitGroup
	wg.Add(numGoroutines * 3) // 3 types of operations

	// Concurrent registrations
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range numOperationsPerGoroutine {
				taskType := fmt.Sprintf("concurrent_task_%d_%d", id, j)
				err := registry.RegisterTask(taskType, func() Task { return &TestTask{} })
				if err != nil {
					t.Errorf("registration failed: %v", err)
				}
			}
		}(i)
	}

	// Concurrent reads
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range numOperationsPerGoroutine {
				taskType := fmt.Sprintf("concurrent_task_%d_%d", id%2, j%10) // Read some existing tasks
				registry.GetConstructor(taskType)
				registry.ListRegisteredTypes()
			}
		}(i)
	}

	// Concurrent unregistrations
	for i := range numGoroutines {
		go func(id int) {
			defer wg.Done()
			for j := range numOperationsPerGoroutine {
				taskType := fmt.Sprintf("concurrent_task_%d_%d", id, j)
				registry.UnregisterTask(taskType)
			}
		}(i)
	}

	wg.Wait()
}

func TestGlobalRegistry_RegisterTask(t *testing.T) {
	// Save original registered types
	originalTypes := ListRegisteredTaskTypes()

	taskType := "global_test_task"
	constructor := func() Task { return &TestTask{} }

	// Test global registration
	err := RegisterTask(taskType, constructor)
	if err != nil {
		t.Fatalf("global registration failed: %v", err)
	}

	// Verify it was registered
	gotConstructor, exists := GetTaskConstructor(taskType)
	if !exists {
		t.Fatal("task type should exist in global registry")
	}
	if gotConstructor == nil {
		t.Fatal("returned constructor should not be nil")
	}

	// Verify it appears in the list
	newTypes := ListRegisteredTaskTypes()
	found := slices.Contains(newTypes, taskType)
	if !found {
		t.Error("task type should appear in registered types list")
	}

	// Cleanup
	defaultRegistry.UnregisterTask(taskType)

	// Verify cleanup
	finalTypes := ListRegisteredTaskTypes()
	if len(finalTypes) != len(originalTypes) {
		t.Errorf("cleanup failed, expected %d types, got %d", len(originalTypes), len(finalTypes))
	}
}

func TestUnmarshalTask_WithCustomTask(t *testing.T) {
	// Register a custom task type
	taskType := "test_task"
	constructor := func() Task { return &TestTask{} }

	err := RegisterTask(taskType, constructor)
	if err != nil {
		t.Fatalf("failed to register custom task: %v", err)
	}
	defer defaultRegistry.UnregisterTask(taskType) // Cleanup

	// Test unmarshaling JSON with custom task type
	// The JSON should have the task-specific field at the top level along with TaskBase fields
	taskJSON := `{
		"test_task": {},
		"test_data": "hello world",
		"output": {
			"as": "${ .result }"
		}
	}`

	task, err := unmarshalTask("test_key", json.RawMessage(taskJSON))
	if err != nil {
		t.Fatalf("failed to unmarshal custom task: %v", err)
	}

	// Verify the task is of the correct type
	testTask, ok := task.(*TestTask)
	if !ok {
		t.Fatalf("expected *TestTask, got %T", task)
	}

	if testTask.TestData != "hello world" {
		t.Errorf("expected test_data 'hello world', got '%s'", testTask.TestData)
	}

	// Verify TaskBase fields are also populated
	if testTask.Output == nil {
		t.Error("expected Output to be populated")
	}
}

func TestUnmarshalTask_UnknownTaskType(t *testing.T) {
	// Test unmarshaling with unknown task type
	taskJSON := `{
		"unknown_task_type": {
			"some_field": "some_value"
		}
	}`

	_, err := unmarshalTask("test_key", json.RawMessage(taskJSON))
	if err == nil {
		t.Fatal("expected error for unknown task type")
	}

	// Verify error message includes available types
	errorMsg := err.Error()
	if !strings.Contains(errorMsg, "unknown task type for key 'test_key'") {
		t.Errorf("error should mention unknown task type, got: %s", errorMsg)
	}
	if !strings.Contains(errorMsg, "Available types:") {
		t.Errorf("error should list available types, got: %s", errorMsg)
	}
}

func TestBuiltInTasksRegistered(t *testing.T) {
	expectedBuiltInTasks := []string{
		"call_http", "call_openapi", "call_grpc", "call_asyncapi", "call",
		"do", "fork", "emit", "for", "listen", "raise", "run", "set", "switch", "try", "wait",
	}

	registeredTypes := ListRegisteredTaskTypes()

	for _, expectedType := range expectedBuiltInTasks {
		found := slices.Contains(registeredTypes, expectedType)
		if !found {
			t.Errorf("built-in task type '%s' should be registered", expectedType)
		}

		// Verify we can get the constructor
		constructor, exists := GetTaskConstructor(expectedType)
		if !exists {
			t.Errorf("constructor for built-in task '%s' should exist", expectedType)
		}
		if constructor == nil {
			t.Errorf("constructor for built-in task '%s' should not be nil", expectedType)
		}
	}
}

func TestTaskRegistry_ErrorMessages(t *testing.T) {
	registry := NewTaskRegistry()

	// Test empty task type error
	err := registry.RegisterTask("", func() Task { return &TestTask{} })
	if err == nil || err.Error() != "task type cannot be empty" {
		t.Errorf("expected 'task type cannot be empty' error, got: %v", err)
	}

	// Test nil constructor error
	err = registry.RegisterTask("test", nil)
	if err == nil || err.Error() != "constructor function cannot be nil" {
		t.Errorf("expected 'constructor function cannot be nil' error, got: %v", err)
	}

	// Test duplicate registration error
	taskType := "duplicate"
	err = registry.RegisterTask(taskType, func() Task { return &TestTask{} })
	if err != nil {
		t.Fatalf("first registration should succeed: %v", err)
	}

	err = registry.RegisterTask(taskType, func() Task { return &TestTask{} })
	expectedError := "task type 'duplicate' is already registered"
	if err == nil || err.Error() != expectedError {
		t.Errorf("expected '%s' error, got: %v", expectedError, err)
	}
}

// Benchmark tests
func BenchmarkTaskRegistry_RegisterTask(b *testing.B) {
	registry := NewTaskRegistry()
	constructor := func() Task { return &TestTask{} }

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskType := fmt.Sprintf("benchmark_task_%d", i)
		_ = registry.RegisterTask(taskType, constructor)
	}
}

func BenchmarkTaskRegistry_GetConstructor(b *testing.B) {
	registry := NewTaskRegistry()
	constructor := func() Task { return &TestTask{} }

	// Pre-register some tasks
	for i := range 1000 {
		taskType := fmt.Sprintf("benchmark_task_%d", i)
		_ = registry.RegisterTask(taskType, constructor)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		taskType := fmt.Sprintf("benchmark_task_%d", i%1000)
		registry.GetConstructor(taskType)
	}
}

func BenchmarkTaskRegistry_ConcurrentAccess(b *testing.B) {
	registry := NewTaskRegistry()
	constructor := func() Task { return &TestTask{} }

	// Pre-register some tasks
	for i := range 100 {
		taskType := fmt.Sprintf("concurrent_benchmark_task_%d", i)
		_ = registry.RegisterTask(taskType, constructor)
	}

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		i := 0
		for pb.Next() {
			taskType := fmt.Sprintf("concurrent_benchmark_task_%d", i%100)
			registry.GetConstructor(taskType)
			i++
		}
	})
}
