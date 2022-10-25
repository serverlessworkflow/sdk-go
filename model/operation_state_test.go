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
)

func TestOperationStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect OperationState
		err    string
	}
	testCases := []testCase{
		{
			desp: "all fields set",
			data: `{"actionMode": "parallel"}`,
			expect: OperationState{
				ActionMode: ActionModeParallel,
			},
			err: ``,
		},
		{
			desp: "actionMode unset",
			data: `{}`,
			expect: OperationState{
				ActionMode: ActionModeSequential,
			},
			err: ``,
		},
		{
			desp: "invalid object format",
			data: `{"actionMode": parallel}`,
			expect: OperationState{
				ActionMode: ActionModeParallel,
			},
			err: `invalid character 'p' looking for beginning of value`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			v := OperationState{}
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, v)
		})
	}
}
