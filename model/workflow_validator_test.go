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

	"github.com/serverlessworkflow/sdk-go/v2/model/test"
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

var workflowStructDefault = Workflow{
	BaseWorkflow: BaseWorkflow{
		ID:          "id",
		SpecVersion: "0.8",
		Auth: Auths{
			{
				Name: "auth name",
			},
		},
		Errors: []Error{
			{
				Name: "error1",
			},
		},
		Secrets: Secrets{
			"Secret1",
		},
		Start: &Start{
			StateName: "name state",
		},
	},
	Events: []Event{
		{
			Name: "event 1",
			Type: "consumer",
		},
	},
	Functions: []Function{
		{
			Name:      "function 1",
			Operation: "rest",
		},
	},
	Retries: []Retry{
		{
			Name: "retry 1",
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
					{
						FunctionRef: &FunctionRef{
							RefName: "function 1",
							Invoke:  InvokeKindSync,
						},
						RetryRef: "retry 1",
					},
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
					{
						EventRef: &EventRef{
							TriggerEventRef: "event 1",
							ResultEventRef:  "event 1",
							Invoke:          InvokeKindSync,
						},
					},
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
			Actions: []Action{
				{
					FunctionRef: &FunctionRef{
						RefName: "function 1",
						Invoke:  InvokeKindSync,
					},
				},
			},
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
			Actions: []Action{
				{
					FunctionRef: &FunctionRef{
						RefName: "function 1",
						Invoke:  InvokeKindSync,
					},
				},
			},
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
			Actions: []Action{
				{
					FunctionRef: &FunctionRef{
						RefName: "function 1",
						Invoke:  InvokeKindSync,
					},
				},
			},
		},
	},
}

func TestWorkflowStructLevelValidation(t *testing.T) {
	testCases := []test.ValidationCase[Workflow]{
		{
			Desp:  "workflow success",
			Model: workflowStructDefault,
		},

		{
			Desp: "workflow state.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.States = append(w.States, w.States[0])
				return w
			}(),
			Err: `Key: 'Workflow.States' Error:Field validation for 'States' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow event.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Events = append(w.Events, w.Events[0])
				return w
			}(),
			Err: `Key: 'Workflow.Events' Error:Field validation for 'Events' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow function.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Functions = append(w.Functions, w.Functions[0])
				return w
			}(),
			Err: `Key: 'Workflow.Functions' Error:Field validation for 'Functions' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow retrie.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Retries = append(w.Retries, w.Retries[0])
				return w
			}(),
			Err: `Key: 'Workflow.Retries' Error:Field validation for 'Retries' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow auth.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Auth = append(w.Auth, w.Auth[0])
				return w
			}(),
			Err: `Key: 'Workflow.BaseWorkflow.Auth' Error:Field validation for 'Auth' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow error.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Errors = append(w.Errors, w.Errors[0])
				return w
			}(),
			Err: `Key: 'Workflow.BaseWorkflow.Errors' Error:Field validation for 'Errors' failed on the 'unique_struct' tag`,
		},
		{
			Desp: "workflow secrets.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Secrets = append(w.Secrets, w.Secrets[0])
				return w
			}(),
			Err: `Key: 'Workflow.BaseWorkflow.Secrets' Error:Field validation for 'Secrets' failed on the 'unique' tag`,
		},
		{
			Desp: "function not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				f := w.Functions[0]
				f.Name = "function renamed to fail"
				w.Functions = []Function{f}
				return w
			}(),
			Err: `Key: 'Workflow.States[0].OperationState.Actions[0].FunctionRef.refName' Error:Field validation for 'refName' failed on the 'exists' tag`,
		},
		{
			Desp: "event not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				e := w.Events[0]
				e.Name = "event renamed to fail"
				w.Events = []Event{e}
				return w
			}(),
			Err: `Key: 'Workflow.States[1].OperationState.Actions[0].EventRef.triggerEventRef' Error:Field validation for 'triggerEventRef' failed on the 'exists' tag
Key: 'Workflow.States[1].OperationState.Actions[0].EventRef.triggerEventRef' Error:Field validation for 'triggerEventRef' failed on the 'exists' tag`,
		},
		{
			Desp: "retry not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				r := w.Retries[0]
				r.Name = "retry renamed to fail"
				w.Retries = []Retry{r}
				return w
			}(),
			Err: `Key: 'Workflow.States[0].OperationState.Actions[0].retryRef' Error:Field validation for 'retryRef' failed on the 'exists' tag`,
		},
		{
			Desp: "workflow id exclude key",
			Model: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = ""
				return w
			}(),
			Err: ``,
		},
		{
			Desp: "workflow key exclude id",
			Model: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = "key"
				return w
			}(),
			Err: ``,
		},
		{
			Desp: "workflow id and key",
			Model: func() Workflow {
				w := workflowStructDefault
				w.ID = "id"
				w.Key = "key"
				return w
			}(),
			Err: ``,
		},
		{
			Desp: "workflow without id and key",
			Model: func() Workflow {
				w := workflowStructDefault
				w.ID = ""
				w.Key = ""
				return w
			}(),
			Err: `Key: 'Workflow.BaseWorkflow.ID' Error:Field validation for 'ID' failed on the 'required_without' tag
Key: 'Workflow.BaseWorkflow.Key' Error:Field validation for 'Key' failed on the 'required_without' tag`,
		},
		{
			Desp: "workflow start",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Start = &Start{
					StateName: "start state not found",
				}
				return w
			}(),
			Err: `Key: 'Workflow.Start' Error:Field validation for 'Start' failed on the 'startnotexist' tag`,
		},
		{
			Desp: "workflow states transitions",
			Model: func() Workflow {
				w := workflowStructDefault
				w.States = listStateTransition1
				return w
			}(),
			Err: ``,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			ctx := NewValidatorContext(&tc.Model)
			err := val.GetValidator().StructCtx(ctx, tc.Model)
			if tc.Err != "" {
				if assert.Error(t, err) {
					assert.Equal(t, tc.Err, err.Error())
				}
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestContinueAsStructLevelValidation(t *testing.T) {
	testCases := []test.ValidationCase[ContinueAs]{
		{
			Desp: "valid ContinueAs",
			Model: ContinueAs{
				WorkflowID: "another-test",
				Version:    "2",
				Data:       FromString("${ del(.customerCount) }"),
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration:  "PT1H",
					Interrupt: false,
					RunBefore: "test",
				},
			},
			Err: ``,
		},
		{
			Desp: "invalid WorkflowExecTimeout",
			Model: ContinueAs{
				WorkflowID: "test",
				Version:    "1",
				Data:       FromString("${ del(.customerCount) }"),
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration: "invalid",
				},
			},
			Err: `Key: 'ContinueAs.WorkflowExecTimeout.Duration' Error:Field validation for 'Duration' failed on the 'iso8601duration' tag`,
		},
	}

	test.StructLevelValidation(t, testCases)
}
