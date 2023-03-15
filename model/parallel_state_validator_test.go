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

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestParallelStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp  string
		state *State
		err   string
	}
	testCases := []testCase{
		{
			desp: "normal",
			state: &State{
				BaseState: BaseState{
					Name: "1",
					Type: "parallel",
					End: &End{
						Terminate: true,
					},
				},
				ParallelState: &ParallelState{
					Branches: []Branch{
						{
							Name: "b1",
							Actions: []Action{{
								FunctionRef: &FunctionRef{
									RefName: "test",
									Invoke:  InvokeKindAsync,
								},
							}},
						},
					},
					CompletionType: CompletionTypeAllOf,
					NumCompleted:   intstr.FromInt(1),
				},
			},
			err: ``,
		},
		{
			desp: "invalid completeType",
			state: &State{
				BaseState: BaseState{
					Name: "1",
					Type: "parallel",
					End: &End{
						Terminate: true,
					},
				},
				ParallelState: &ParallelState{
					Branches: []Branch{
						{
							Name: "b1",
							Actions: []Action{{
								FunctionRef: &FunctionRef{
									RefName: "test",
									Invoke:  InvokeKindAsync,
								},
							}},
						},
					},
					CompletionType: CompletionTypeAllOf + "1",
				},
			},
			err: `Key: 'State.ParallelState.CompletionType' Error:Field validation for 'CompletionType' failed on the 'oneof' tag`,
		},
		{
			desp: "invalid numCompleted `int`",
			state: &State{
				BaseState: BaseState{
					Name: "1",
					Type: "parallel",
					End: &End{
						Terminate: true,
					},
				},
				ParallelState: &ParallelState{
					Branches: []Branch{
						{
							Name: "b1",
							Actions: []Action{{
								FunctionRef: &FunctionRef{
									RefName: "test",
									Invoke:  InvokeKindAsync,
								},
							}},
						},
					},
					CompletionType: CompletionTypeAtLeast,
					NumCompleted:   intstr.FromInt(0),
				},
			},
			err: `Key: 'State.ParallelState.NumCompleted' Error:Field validation for 'NumCompleted' failed on the 'gt0' tag`,
		},
		{
			desp: "invalid numCompleted string format",
			state: &State{
				BaseState: BaseState{
					Name: "1",
					Type: "parallel",
					End: &End{
						Terminate: true,
					},
				},
				ParallelState: &ParallelState{
					Branches: []Branch{
						{
							Name: "b1",
							Actions: []Action{{
								FunctionRef: &FunctionRef{
									RefName: "test",
									Invoke:  InvokeKindAsync,
								},
							}},
						},
					},
					CompletionType: CompletionTypeAtLeast,
					NumCompleted:   intstr.FromString("a"),
				},
			},
			err: `Key: 'State.ParallelState.NumCompleted' Error:Field validation for 'NumCompleted' failed on the 'gt0' tag`,
		},
		{
			desp: "normal",
			state: &State{
				BaseState: BaseState{
					Name: "1",
					Type: "parallel",
					End: &End{
						Terminate: true,
					},
				},
				ParallelState: &ParallelState{
					Branches: []Branch{
						{
							Name: "b1",
							Actions: []Action{{
								FunctionRef: &FunctionRef{
									RefName: "test",
									Invoke:  InvokeKindAsync,
								},
							}},
						},
					},
					CompletionType: CompletionTypeAtLeast,
					NumCompleted:   intstr.FromString("0"),
				},
			},
			err: `Key: 'State.ParallelState.NumCompleted' Error:Field validation for 'NumCompleted' failed on the 'gt0' tag`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.state)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
