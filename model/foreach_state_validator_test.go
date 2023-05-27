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

	"github.com/serverlessworkflow/sdk-go/v2/model/test"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestForEachStateStructLevelValidation(t *testing.T) {
	test.StructLevelValidation(t, []test.ValidationCase[State]{
		{
			Desp: "normal test & sequential",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "function 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeSequential,
				},
			},
			Err: ``,
		},
		{
			Desp: "normal test & parallel int",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "test 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeParallel,
					BatchSize: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 1,
					},
				},
			},
			Err: ``,
		},
		{
			Desp: "normal test & parallel string",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "test 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeParallel,
					BatchSize: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "1",
					},
				},
			},
			Err: ``,
		},
		{
			Desp: "invalid parallel int",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "test 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeParallel,
					BatchSize: &intstr.IntOrString{
						Type:   intstr.Int,
						IntVal: 0,
					},
				},
			},
			Err: `Key: 'State.ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
		{
			Desp: "invalid parallel string",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "test 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeParallel,
					BatchSize: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "0",
					},
				},
			},
			Err: `Key: 'State.ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
		{
			Desp: "invalid parallel string format",
			Model: State{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeForEach,
					End: &End{
						Terminate: true,
					},
				},
				ForEachState: &ForEachState{
					InputCollection: "3",
					Actions: []Action{{
						FunctionRef: &FunctionRef{
							RefName: "test 1",
							Invoke:  InvokeKindAsync,
						},
					}},
					Mode: ForEachModeTypeParallel,
					BatchSize: &intstr.IntOrString{
						Type:   intstr.String,
						StrVal: "a",
					},
				},
			},
			Err: `Key: 'State.ForEachState.BatchSize' Error:Field validation for 'BatchSize' failed on the 'gt0' tag`,
		},
	})
}
