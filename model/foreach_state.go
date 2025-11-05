// Copyright 2022 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"

	validator "github.com/go-playground/validator/v10"
	"k8s.io/apimachinery/pkg/util/intstr"

	val "github.com/finbox-in/serverlessworkflow-sdk-go/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(ForEachStateStructLevelValidation, ForEachState{})
}

// ForEachModeType Specifies how iterations are to be performed (sequentially or in parallel)
type ForEachModeType string

const (
	// ForEachModeTypeSequential specifies iterations should be done sequentially.
	ForEachModeTypeSequential ForEachModeType = "sequential"
	// ForEachModeTypeParallel specifies iterations should be done parallel.
	ForEachModeTypeParallel ForEachModeType = "parallel"
)

// ForEachState used to execute actions for each element of a data set.
type ForEachState struct {
	BaseState
	// Workflow expression selecting an array element of the states data
	InputCollection string `json:"inputCollection" validate:"required"`
	// Workflow expression specifying an array element of the states data to add the results of each iteration
	OutputCollection string `json:"outputCollection,omitempty"`
	// Name of the iteration parameter that can be referenced in actions/workflow. For each parallel iteration, this param should contain an unique element of the inputCollection array
	IterationParam string `json:"iterationParam,omitempty"`
	// Specifies how upper bound on how many iterations may run in parallel
	BatchSize *intstr.IntOrString `json:"batchSize,omitempty"`
	// Actions to be executed for each of the elements of inputCollection
	Actions []Action `json:"actions,omitempty" validate:"required,min=1,dive"`
	// State specific timeout
	Timeouts *ForEachStateTimeout `json:"timeouts,omitempty"`
	// Mode Specifies how iterations are to be performed (sequentially or in parallel)
	// Defaults to parallel
	Mode ForEachModeType `json:"mode,omitempty"`
}

func (f *ForEachState) DeepCopyState() State {
	return f
}

type forEachStateForUnmarshal ForEachState

func (f *ForEachState) UnmarshalJSON(data []byte) error {
	v := forEachStateForUnmarshal{
		Mode: ForEachModeTypeParallel,
	}
	err := json.Unmarshal(data, &v)
	if err != nil {
		return fmt.Errorf("forEachState value '%s' is not supported, it must be an object or string", string(data))
	}

	*f = ForEachState(v)
	return nil
}

// ForEachStateStructLevelValidation custom validator for ForEachState
func ForEachStateStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	stateObj := structLevel.Current().Interface().(ForEachState)

	if stateObj.Mode != ForEachModeTypeParallel {
		return
	}

	if stateObj.BatchSize == nil {
		return
	}

	switch stateObj.BatchSize.Type {
	case intstr.Int:
		if stateObj.BatchSize.IntVal <= 0 {
			structLevel.ReportError(reflect.ValueOf(stateObj.BatchSize), "BatchSize", "batchSize", "gt0", "")
		}
	case intstr.String:
		v, err := strconv.Atoi(stateObj.BatchSize.StrVal)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(stateObj.BatchSize), "BatchSize", "batchSize", "gt0", err.Error())
			return
		}

		if v <= 0 {
			structLevel.ReportError(reflect.ValueOf(stateObj.BatchSize), "BatchSize", "batchSize", "gt0", "")
		}
	}
}

// ForEachStateTimeout defines timeout settings for foreach state
type ForEachStateTimeout struct {
	StateExecTimeout  *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	ActionExecTimeout string            `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}
