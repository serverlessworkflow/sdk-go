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

package test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

type ValidationCase[T any] struct {
	Desp  string
	Model T
	Err   string
}

func StructLevelValidation[T any](t *testing.T, testCases []ValidationCase[T]) {
	for _, tc := range testCases {
		t.Run(tc.Desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.Model)
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
