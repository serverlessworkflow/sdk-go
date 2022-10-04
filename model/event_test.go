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

func TestEventRefStructLevelValidation(t *testing.T) {
	type testCase struct {
		name     string
		eventRef EventRef
		err      string
	}

	testCases := []testCase{
		{
			name: "valid resultEventTimeout",
			eventRef: EventRef{
				TriggerEventRef:    "example valid",
				ResultEventRef:     "example valid",
				ResultEventTimeout: "PT1H",
				Invoke:             "sync",
			},
			err: ``,
		},
		{
			name: "invalid resultEventTimeout",
			eventRef: EventRef{
				TriggerEventRef:    "example invalid",
				ResultEventRef:     "example invalid red",
				ResultEventTimeout: "10hs",
				Invoke:             "sync",
			},
			err: `Key: 'EventRef.ResultEventTimeout' Error:Field validation for 'ResultEventTimeout' failed on the 'reqiso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.eventRef)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
