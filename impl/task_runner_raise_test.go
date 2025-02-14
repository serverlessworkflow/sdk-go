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
	"encoding/json"
	"errors"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v3/model"
	"github.com/stretchr/testify/assert"
)

func TestRaiseTaskRunner_WithDefinedError(t *testing.T) {
	input := map[string]interface{}{}

	raiseTask := &model.RaiseTask{
		Raise: model.RaiseTaskConfiguration{
			Error: model.RaiseTaskError{
				Definition: &model.Error{
					Type:   model.NewUriTemplate(model.ErrorTypeValidation),
					Status: 400,
					Title:  model.NewStringOrRuntimeExpr("Validation Error"),
					Detail: model.NewStringOrRuntimeExpr("Invalid input data"),
				},
			},
		},
	}

	runner, err := NewRaiseTaskRunner("task_raise_defined", raiseTask, nil)
	assert.NoError(t, err)

	output, err := runner.Run(input)
	assert.Equal(t, output, input)
	assert.Error(t, err)

	expectedErr := model.NewErrValidation(errors.New("Invalid input data"), "task_raise_defined")

	var modelErr *model.Error
	if errors.As(err, &modelErr) {
		assert.Equal(t, expectedErr.Type.String(), modelErr.Type.String())
		assert.Equal(t, expectedErr.Status, modelErr.Status)
		assert.Equal(t, expectedErr.Title.String(), modelErr.Title.String())
		assert.Equal(t, "Invalid input data", modelErr.Detail.String())
		assert.Equal(t, expectedErr.Instance.String(), modelErr.Instance.String())
	} else {
		t.Errorf("expected error of type *model.Error but got %T", err)
	}
}

func TestRaiseTaskRunner_WithReferencedError(t *testing.T) {
	ref := "someErrorRef"
	raiseTask := &model.RaiseTask{
		Raise: model.RaiseTaskConfiguration{
			Error: model.RaiseTaskError{
				Ref: &ref,
			},
		},
	}

	runner, err := NewRaiseTaskRunner("task_raise_ref", raiseTask, nil)
	assert.Error(t, err)
	assert.Nil(t, runner)
}

func TestRaiseTaskRunner_TimeoutErrorWithExpression(t *testing.T) {
	input := map[string]interface{}{
		"timeoutMessage": "Request took too long",
	}

	raiseTask := &model.RaiseTask{
		Raise: model.RaiseTaskConfiguration{
			Error: model.RaiseTaskError{
				Definition: &model.Error{
					Type:   model.NewUriTemplate(model.ErrorTypeTimeout),
					Status: 408,
					Title:  model.NewStringOrRuntimeExpr("Timeout Error"),
					Detail: model.NewStringOrRuntimeExpr("${ .timeoutMessage }"),
				},
			},
		},
	}

	runner, err := NewRaiseTaskRunner("task_raise_timeout_expr", raiseTask, nil)
	assert.NoError(t, err)

	output, err := runner.Run(input)
	assert.Equal(t, input, output)
	assert.Error(t, err)

	expectedErr := model.NewErrTimeout(errors.New("Request took too long"), "task_raise_timeout_expr")

	var modelErr *model.Error
	if errors.As(err, &modelErr) {
		assert.Equal(t, expectedErr.Type.String(), modelErr.Type.String())
		assert.Equal(t, expectedErr.Status, modelErr.Status)
		assert.Equal(t, expectedErr.Title.String(), modelErr.Title.String())
		assert.Equal(t, "Request took too long", modelErr.Detail.String())
		assert.Equal(t, expectedErr.Instance.String(), modelErr.Instance.String())
	} else {
		t.Errorf("expected error of type *model.Error but got %T", err)
	}
}

func TestRaiseTaskRunner_Serialization(t *testing.T) {
	raiseTask := &model.RaiseTask{
		Raise: model.RaiseTaskConfiguration{
			Error: model.RaiseTaskError{
				Definition: &model.Error{
					Type:     model.NewUriTemplate(model.ErrorTypeRuntime),
					Status:   500,
					Title:    model.NewStringOrRuntimeExpr("Runtime Error"),
					Detail:   model.NewStringOrRuntimeExpr("Unexpected failure"),
					Instance: &model.JsonPointerOrRuntimeExpression{Value: "/task_runtime"},
				},
			},
		},
	}

	data, err := json.Marshal(raiseTask)
	assert.NoError(t, err)

	var deserializedTask model.RaiseTask
	err = json.Unmarshal(data, &deserializedTask)
	assert.NoError(t, err)

	assert.Equal(t, raiseTask.Raise.Error.Definition.Type.String(), deserializedTask.Raise.Error.Definition.Type.String())
	assert.Equal(t, raiseTask.Raise.Error.Definition.Status, deserializedTask.Raise.Error.Definition.Status)
	assert.Equal(t, raiseTask.Raise.Error.Definition.Title.String(), deserializedTask.Raise.Error.Definition.Title.String())
	assert.Equal(t, raiseTask.Raise.Error.Definition.Detail.String(), deserializedTask.Raise.Error.Definition.Detail.String())
	assert.Equal(t, raiseTask.Raise.Error.Definition.Instance.String(), deserializedTask.Raise.Error.Definition.Instance.String())
}

func TestRaiseTaskRunner_ReferenceSerialization(t *testing.T) {
	ref := "errorReference"
	raiseTask := &model.RaiseTask{
		Raise: model.RaiseTaskConfiguration{
			Error: model.RaiseTaskError{
				Ref: &ref,
			},
		},
	}

	data, err := json.Marshal(raiseTask)
	assert.NoError(t, err)

	var deserializedTask model.RaiseTask
	err = json.Unmarshal(data, &deserializedTask)
	assert.NoError(t, err)

	assert.Equal(t, *raiseTask.Raise.Error.Ref, *deserializedTask.Raise.Error.Ref)
	assert.Nil(t, deserializedTask.Raise.Error.Definition)
}
