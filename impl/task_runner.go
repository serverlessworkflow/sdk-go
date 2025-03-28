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
	"context"
	"time"

	"github.com/serverlessworkflow/sdk-go/v3/impl/ctx"
	"github.com/serverlessworkflow/sdk-go/v3/model"
)

var _ TaskRunner = &SetTaskRunner{}
var _ TaskRunner = &RaiseTaskRunner{}
var _ TaskRunner = &ForTaskRunner{}
var _ TaskRunner = &DoTaskRunner{}

type TaskRunner interface {
	Run(input interface{}) (interface{}, error)
	GetTaskName() string
}

type TaskSupport interface {
	SetTaskStatus(task string, status ctx.StatusPhase)
	GetWorkflowDef() *model.Workflow
	// SetWorkflowInstanceCtx is the `$context` variable accessible in JQ expressions and set in `export.as`
	SetWorkflowInstanceCtx(value interface{})
	// GetContext gets the sharable Workflow context. Accessible via ctx.GetWorkflowContext.
	GetContext() context.Context
	SetTaskRawInput(value interface{})
	SetTaskRawOutput(value interface{})
	SetTaskDef(task model.Task) error
	SetTaskStartedAt(value time.Time)
	SetTaskName(name string)
	// SetTaskReferenceFromName based on the taskName and the model.Workflow definition, set the JSON Pointer reference to the context
	SetTaskReferenceFromName(taskName string) error
	GetTaskReference() string
	// SetLocalExprVars overrides local variables in expression processing
	SetLocalExprVars(vars map[string]interface{})
	// AddLocalExprVars adds to the local variables in expression processing. Won't override previous entries.
	AddLocalExprVars(vars map[string]interface{})
	// RemoveLocalExprVars removes local variables added in AddLocalExprVars or SetLocalExprVars
	RemoveLocalExprVars(keys ...string)
}
