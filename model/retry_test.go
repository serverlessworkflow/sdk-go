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

	"github.com/finbox-in/serverlessworkflow-sdk-go/util/floatstr"
	val "github.com/finbox-in/serverlessworkflow-sdk-go/validator"
)

func TestRetryStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp     string
		retryObj Retry
		err      string
	}
	testCases := []testCase{
		{
			desp: "normal",
			retryObj: Retry{
				Name:      "1",
				Delay:     "PT5S",
				MaxDelay:  "PT5S",
				Increment: "PT5S",
				Jitter:    floatstr.FromString("PT5S"),
			},
			err: ``,
		},
		{
			desp: "normal with all optinal",
			retryObj: Retry{
				Name: "1",
			},
			err: ``,
		},
		{
			desp: "missing required name",
			retryObj: Retry{
				Name:      "",
				Delay:     "PT5S",
				MaxDelay:  "PT5S",
				Increment: "PT5S",
				Jitter:    floatstr.FromString("PT5S"),
			},
			err: `Key: 'Retry.Name' Error:Field validation for 'Name' failed on the 'required' tag`,
		},
		{
			desp: "invalid delay duration",
			retryObj: Retry{
				Name:      "1",
				Delay:     "P5S",
				MaxDelay:  "PT5S",
				Increment: "PT5S",
				Jitter:    floatstr.FromString("PT5S"),
			},
			err: `Key: 'Retry.Delay' Error:Field validation for 'Delay' failed on the 'iso8601duration' tag`,
		},
		{
			desp: "invdalid max delay duration",
			retryObj: Retry{
				Name:      "1",
				Delay:     "PT5S",
				MaxDelay:  "P5S",
				Increment: "PT5S",
				Jitter:    floatstr.FromString("PT5S"),
			},
			err: `Key: 'Retry.MaxDelay' Error:Field validation for 'MaxDelay' failed on the 'iso8601duration' tag`,
		},
		{
			desp: "invalid increment duration",
			retryObj: Retry{
				Name:      "1",
				Delay:     "PT5S",
				MaxDelay:  "PT5S",
				Increment: "P5S",
				Jitter:    floatstr.FromString("PT5S"),
			},
			err: `Key: 'Retry.Increment' Error:Field validation for 'Increment' failed on the 'iso8601duration' tag`,
		},
		{
			desp: "invalid jitter duration",
			retryObj: Retry{
				Name:      "1",
				Delay:     "PT5S",
				MaxDelay:  "PT5S",
				Increment: "PT5S",
				Jitter:    floatstr.FromString("P5S"),
			},
			err: `Key: 'Retry.Jitter' Error:Field validation for 'Jitter' failed on the 'iso8601duration' tag`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(tc.retryObj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
