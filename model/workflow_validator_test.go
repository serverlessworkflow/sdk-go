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
	"testing"

	"github.com/stretchr/testify/assert"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

var workflowStructDefault = Workflow{
	BaseWorkflow: BaseWorkflow{
		ID:          "id",
		SpecVersion: "0.8",
		Auth: AuthArray{
			{
				Name: "auth name",
			},
		},
		Start: &Start{
			StateName: "name state",
		},
	},
	States: []State{
		{
			BaseState: BaseState{
				Name: "name state",
				Type: StateTypeOperation,
				Transition: &Transition{
					NextState: "next name state",
				},
			},
			OperationState: &OperationState{
				ActionMode: "sequential",
				Actions: []Action{
					{},
				},
			},
		},
		{
			BaseState: BaseState{
				Name: "next name state",
				Type: StateTypeOperation,
				End: &End{
					Terminate: true,
				},
			},
			OperationState: &OperationState{
				ActionMode: "sequential",
				Actions: []Action{
					{},
				},
			},
		},
	},
}

var listStateTransition1 = []State{
	{
		BaseState: BaseState{
			Name: "name state",
			Type: StateTypeOperation,
			Transition: &Transition{
				NextState: "next name state",
			},
		},
		OperationState: &OperationState{
			ActionMode: "sequential",
			Actions:    []Action{{}},
		},
	},
	{
		BaseState: BaseState{
			Name: "next name state",
			Type: StateTypeOperation,
			Transition: &Transition{
				NextState: "next name state 2",
			},
		},
		OperationState: &OperationState{
			ActionMode: "sequential",
			Actions:    []Action{{}},
		},
	},
	{
		BaseState: BaseState{
			Name: "next name state 2",
			Type: StateTypeOperation,
			End: &End{
				Terminate: true,
			},
		},
		OperationState: &OperationState{
			ActionMode: "sequential",
			Actions:    []Action{{}},
		},
	},
}

func TestWorkflowStructLevelValidation(t *testing.T) {
	type testCase[T any] struct {
		name     string
		instance T
		err      string
	}
	testCases := []testCase[any]{
		{
			name:     "workflow success",
			instance: workflowStructDefault,
		},
		{
			name: "workflow auth.name repeat",
			instance: func() Workflow {
				w := workflowStructDefault
				w.Auth = append(w.Auth, w.Auth[0])
				return w
			}(),
			err: `Key: 'Workflow.[]Auth.Name' Error:Field validation for '[]Auth.Name' failed on the 'reqnameunique' tag`,
		},
		{
			name: "workflow id exclude key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = ""
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow key exclude id",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = "key"
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow id and key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = "key"
				return w
			}(),
			err: ``,
		},
		{
			name: "workflow without id and key",
			instance: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = ""
				return w
			}(),
			err: `Key: 'Workflow.BaseWorkflow.ID' Error:Field validation for 'ID' failed on the 'required_without' tag
Key: 'Workflow.BaseWorkflow.Key' Error:Field validation for 'Key' failed on the 'required_without' tag`,
		},
		{
			name: "workflow start",
			instance: func() Workflow {
				w := workflowStructDefault
				w.Start = &Start{
					StateName: "start state not found",
				}
				return w
			}(),
			err: `Key: 'Workflow.Start' Error:Field validation for 'Start' failed on the 'startnotexist' tag`,
		},
		{
			name: "workflow states transitions",
			instance: func() Workflow {
				w := workflowStructDefault
				w.States = listStateTransition1
				return w
			}(),
			err: ``,
		},
		{
			name: "valid ContinueAs",
			instance: ContinueAs{
				WorkflowID: "another-test",
				Version:    "2",
				Data:       FromString("${ del(.customerCount) }"),
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration:  "PT1H",
					Interrupt: false,
					RunBefore: "test",
				},
			},
			err: ``,
		},
		{
			name: "invalid WorkflowExecTimeout",
			instance: ContinueAs{
				WorkflowID: "test",
				Version:    "1",
				Data:       FromString("${ del(.customerCount) }"),
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration: "invalid",
				},
			},
			err: `Key: 'ContinueAs.workflowExecTimeout' Error:Field validation for 'workflowExecTimeout' failed on the 'iso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.instance)

			if tc.err != "" {
				assert.Error(t, err)
				if err != nil {
					assert.Equal(t, tc.err, err.Error())
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}
