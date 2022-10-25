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

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

// CompletionType define on how to complete branch execution.
type CompletionType string

const (
	// CompletionTypeAllOf defines all branches must complete execution before the state can transition/end.
	CompletionTypeAllOf CompletionType = "allOf"
	// CompletionTypeAtLeast defines state can transition/end once at least the specified number of branches
	// have completed execution.
	CompletionTypeAtLeast CompletionType = "atLeast"
)

// ParallelState Consists of a number of states that are executed in parallel
type ParallelState struct {
	BaseState
	// Branch Definitions
	Branches []Branch `json:"branches" validate:"required,min=1,dive"`
	// Option types on how to complete branch execution.
	// Defaults to `allOf`
	CompletionType CompletionType `json:"completionType,omitempty" validate:"required,oneof=allOf atLeast"`

	// Used when completionType is set to 'atLeast' to specify the minimum number of branches that must complete before the state will transition."
	// TODO: change this field to unmarshal result as int
	NumCompleted intstr.IntOrString `json:"numCompleted,omitempty"`
	// State specific timeouts
	Timeouts *ParallelStateTimeout `json:"timeouts,omitempty"`
}

type parallelStateForUnmarshal ParallelState

// UnmarshalJSON unmarshal ParallelState object from json bytes
func (s *ParallelState) UnmarshalJSON(b []byte) error {
	if len(b) == 0 {
		// TODO: Normalize error messages
		return fmt.Errorf("no bytes to unmarshal")
	}

	v := &parallelStateForUnmarshal{
		CompletionType: CompletionTypeAllOf,
	}
	err := json.Unmarshal(b, v)
	if err != nil {
		return err
	}

	*s = ParallelState(*v)

	return nil
}

// Branch Definition
type Branch struct {
	// Branch name
	Name string `json:"name" validate:"required"`
	// Actions to be executed in this branch
	Actions []Action `json:"actions" validate:"required,min=1,dive"`
	// Timeouts State specific timeouts
	Timeouts *BranchTimeouts `json:"timeouts,omitempty"`
}

// BranchTimeouts defines the specific timeout settings for branch
type BranchTimeouts struct {
	// ActionExecTimeout Single actions definition execution timeout duration (ISO 8601 duration format)
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
	// BranchExecTimeout Single branch execution timeout duration (ISO 8601 duration format)
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}

// ParallelStateTimeout defines the specific timeout settings for parallel state
type ParallelStateTimeout struct {
	StateExecTimeout  *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	BranchExecTimeout string            `json:"branchExecTimeout,omitempty" validate:"omitempty,iso8601duration"`
}

// ParallelStateStructLevelValidation custom validator for ParallelState
func ParallelStateStructLevelValidation(_ context.Context, structLevel validator.StructLevel) {
	parallelStateObj := structLevel.Current().Interface().(ParallelState)

	if parallelStateObj.CompletionType == CompletionTypeAllOf {
		return
	}

	switch parallelStateObj.NumCompleted.Type {
	case intstr.Int:
		if parallelStateObj.NumCompleted.IntVal <= 0 {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", "")
		}
	case intstr.String:
		v, err := strconv.Atoi(parallelStateObj.NumCompleted.StrVal)
		if err != nil {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", err.Error())
			return
		}

		if v <= 0 {
			structLevel.ReportError(reflect.ValueOf(parallelStateObj.NumCompleted), "NumCompleted", "numCompleted", "gt0", "")
		}
	}
}

func init() {
	val.GetValidator().RegisterStructValidationCtx(
		ParallelStateStructLevelValidation,
		ParallelState{},
	)
}
