// Copyright 2023 The Serverless Workflow Specification Authors
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

package builder

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestNewWorkflowBuilder(t *testing.T) {
	workflowBuilder, builder := NewWorkflowBuilder("test", "1.0")
	workflowBuilder.
		Function("function0", "http://...").
		Parent().
		Event("event0").
		Parent().
		OperationState("state0").Action("action0").FunctionRef("function0")

	workflow := builder.Build()
	assert.Equal(t, 1, len(workflow.Functions))
	assert.Equal(t, "function0", workflow.Functions[0].Name)
	assert.Equal(t, 1, len(workflow.Events))
	assert.Equal(t, "event0", workflow.Events[0].Name)
	assert.Equal(t, 1, len(workflow.States))
	assert.Equal(t, model.StateTypeOperation, workflow.States[0].Type)
	assert.Equal(t, "state0", workflow.States[0].Name)
	assert.NotNil(t, workflow.States[0].OperationState)
	assert.Equal(t, 1, len(workflow.States[0].OperationState.Actions))
	assert.Equal(t, "action0", workflow.States[0].OperationState.Actions[0].Name)
	assert.Equal(t, "function0", workflow.States[0].OperationState.Actions[0].FunctionRef.RefName)
}
