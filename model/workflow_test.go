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
	val "github.com/serverlessworkflow/sdk-go/v2/validator"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContinueAsStructLevelValidation(t *testing.T) {
	type testCase struct {
		name       string
		continueAs ContinueAs
		err        string
	}

	testCases := []testCase{
		{
			name: "valid ContinueAs",
			continueAs: ContinueAs{
				WorkflowRef: WorkflowRef{
					WorkflowID:       "another-test",
					Version:          "2",
					Invoke:           "sync",
					OnParentComplete: "terminate",
				},
				Data: "${ del(.customerCount) }",
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
			continueAs: ContinueAs{
				WorkflowRef: WorkflowRef{
					WorkflowID: "test",
					Version:    "1",
				},
				Data: "${ del(.customerCount) }",
				WorkflowExecTimeout: WorkflowExecTimeout{
					Duration: "invalid",
				},
			},
			err: `Key: 'ContinueAs.workflowExecTimeout' Error:Field validation for 'workflowExecTimeout' failed on the 'iso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.continueAs)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
