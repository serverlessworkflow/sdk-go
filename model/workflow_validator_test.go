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
	"testing"

	"github.com/stretchr/testify/assert"

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
				Name: "error 1",
			},
			{
				Name: "error 2",
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
			Kind: EventKindConsumed,
		},
	},
	Functions: []Function{
		{
			Name:      "function 1",
			Operation: "rest",
			Type:      FunctionTypeREST,
		},
		{
			Name:      "function 2",
			Operation: "rest",
			Type:      FunctionTypeREST,
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
				OnErrors: []OnError{
					{
						ErrorRefs: []string{
							"error 1",
						},
					},
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
				OnErrors: []OnError{
					{
						ErrorRef: "error 2",
					},
				},
				CompensatedBy: "compensation state",
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
		{
			BaseState: BaseState{
				Name: "compensation state",
				Type: StateTypeOperation,
				OnErrors: []OnError{
					{
						ErrorRef: "error 2",
					},
				},
				UsedForCompensation: true,
				End: &End{
					Terminate: true,
				},
			},
			OperationState: &OperationState{
				ActionMode: "sequential",
				Actions: []Action{
					{
						FunctionRef: &FunctionRef{
							RefName: "function 2",
							Invoke:  InvokeKindSync,
						},
						RetryRef: "retry 1",
					},
				},
			},
		},
	},
}

type ValidationCase[T any] struct {
	Desp  string
	Model T
	Err   string
}

func StructLevelValidationCtx[T any](t *testing.T, ctx context.Context, testCases []ValidationCase[T]) {
	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			err := val.GetValidator().StructCtx(ctx, tc.Model)
			err = WorkflowError(err)
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

func TestWorkflowStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase[Workflow]{
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
			Err: `states has duplicate "name"`,
		},
		{
			Desp: "workflow event.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Events = append(w.Events, w.Events[0])
				return w
			}(),
			Err: `events has duplicate "name"`,
		},
		{
			Desp: "workflow function.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Functions = append(w.Functions, w.Functions[0])
				return w
			}(),
			Err: `functions has duplicate "name"`,
		},
		{
			Desp: "workflow retrie.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Retries = append(w.Retries, w.Retries[0])
				return w
			}(),
			Err: `retries has duplicate "name"`,
		},
		{
			Desp: "workflow auth.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Auth = append(w.Auth, w.Auth[0])
				return w
			}(),
			Err: `auth has duplicate "name"`,
		},
		{
			Desp: "workflow error.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Errors = append(w.Errors, w.Errors[0])
				return w
			}(),
			Err: `errors has duplicate "name"`,
		},
		{
			Desp: "workflow secrets.name repeat",
			Model: func() Workflow {
				w := workflowStructDefault
				w.Secrets = append(w.Secrets, w.Secrets[0])
				return w
			}(),
			Err: `secrets has duplicate value`,
		},
		{
			Desp: "function not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				f := w.Functions[0]
				f.Name = "function renamed to fail"
				w.Functions = []Function{f, w.Functions[1]}
				return w
			}(),
			Err: `states[0].operationState.actions[0].functionRef.refName don't exists "function 1"`,
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
			Err: `states[1].operationState.actions[0].eventRef.triggerEventRef don't exists "event 1"
states[1].operationState.actions[0].eventRef.triggerEventRef don't exists "event 1"`,
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
			Err: `states[0].operationState.actions[0].retryRef don't exists "retry 1"
states[2].operationState.actions[0].retryRef don't exists "retry 1"`,
		},
		{
			Desp: "error 1 not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				e := w.Errors[0]
				e.Name = "error fail exists"
				w.Errors = []Error{e, w.Errors[1]}
				return w
			}(),
			Err: `states[0].onErrors[0].errorRefs don't exists ["error 1"]`,
		},
		{
			Desp: "error 2 not exists",
			Model: func() Workflow {
				w := workflowStructDefault
				e := w.Errors[1]
				e.Name = "error fail exists"
				w.Errors = []Error{w.Errors[0], e}
				return w
			}(),
			Err: `states[1].onErrors[0].errorRef don't exists "error 2"
states[2].onErrors[0].errorRef don't exists "error 2"`,
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
			Err: `key required when not defined "id"
id required when not defined "key"`,
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
			Err: `start.stateName don't exists "start state not found"`,
		},
		{
			Desp: "workflow transition no exists",
			Model: func() Workflow {
				w := workflowStructDefault
				s := w.States[0]
				t := *s.Transition
				t.NextState = "transtion not exists"
				s.Transition = &t
				w.States = []State{s, w.States[1], w.States[2]}
				return w
			}(),
			Err: `states[0].transition.nextState don't exists "transtion not exists"`,
		},
		{
			Desp: "transition compensation",
			Model: func() Workflow {
				w := workflowStructDefault
				s := w.States[2]
				s.UsedForCompensation = false
				w.States = []State{w.States[0], w.States[1], s}
				return w
			}(),
			Err: `states[1].compensatedBy compensatedBy don't exists "compensation state"`,
		},
		{
			Desp: "state recursive",
			Model: func() Workflow {
				w := workflowStructDefault
				s := w.States[0]
				t := *s.Transition
				t.NextState = s.Name
				s.Transition = &t
				w.States = []State{s}
				return w
			}(),
			Err: `states[0].transition.nextState can't no be recursive "name state"`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			StructLevelValidationCtx(t, NewValidatorContext(&tc.Model), []ValidationCase[Workflow]{tc})
		})
	}
}

func TestContinueAsStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase[ContinueAs]{
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

	StructLevelValidationCtx(t, NewValidatorContext(&Workflow{}), testCases)
}

func TestOnErrorStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase[OnError]{
		{
			Desp: "duplicate ErrorRefs",
			Model: OnError{
				ErrorRefs: []string{"error 1", "error 1"},
			},
			Err: `errorRefs has duplicate value`,
		},
		{
			Desp: "valid OnError",
			Model: OnError{
				ErrorRef: "error 1",
			},
			Err: ``,
		},
		{
			Desp: "valid OnError",
			Model: OnError{
				ErrorRefs: []string{"error 1"},
			},
			Err: ``,
		},
		{
			Desp:  "required errorRef",
			Model: OnError{},
			Err:   `Key: 'OnError.ErrorRef' Error:Field validation for 'ErrorRef' failed on the 'required' tag`,
		},

		{
			Desp: "required errorRef",
			Model: OnError{
				ErrorRef:  "error 1",
				ErrorRefs: []string{"error 1"},
			},
			Err: `Key: 'OnError.ErrorRef' Error:Field validation for 'ErrorRef' failed on the 'exclusive' tag`,
		},
	}

	StructLevelValidationCtx(t, NewValidatorContext(&workflowStructDefault), testCases)
}
