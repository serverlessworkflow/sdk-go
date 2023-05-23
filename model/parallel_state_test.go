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
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
)

func TestParallelStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect *ParallelState
		err    string
	}
	testCases := []testCase{
		{
			desp: "all field set",
			data: `{"completionType": "allOf", "numCompleted": 1}`,
			expect: &ParallelState{
				CompletionType: CompletionTypeAllOf,
				NumCompleted:   intstr.FromInt(1),
			},
			err: ``,
		},
		{
			desp: "all optional field not set",
			data: `{"numCompleted": 1}`,
			expect: &ParallelState{
				CompletionType: CompletionTypeAllOf,
				NumCompleted:   intstr.FromInt(1),
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v ParallelState
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, &v)
		})
	}
}

func TestParallelStateToString(t *testing.T) {
	single := "2023-04-22T10:00:35+02:00"
	total := "2023-04-22T10:01:35+02:00"
	stateExecTimeout := StateExecTimeout{
		Single: single,
		Total:  total,
	}

	timeout := ParallelStateTimeout{
		StateExecTimeout: &stateExecTimeout,
	}

	parallelState := ParallelState{
		CompletionType: CompletionTypeAtLeast,
		NumCompleted:   intstr.FromInt(25),
		Branches:       []Branch{},
		Timeouts:       &timeout,
	}
	value := parallelState.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ Branches:[], CompletionType:atLeast, NumCompleted:{Type:0 IntVal:25 StrVal:}, Timeouts:{ BranchExecTimeout:, StateExecTimeout:{ Single:2023-04-22T10:00:35+02:00, Total:2023-04-22T10:01:35+02:00} } }", value)
}
