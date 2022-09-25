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

package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateISO8601TimeDuration(t *testing.T) {
	type testCase struct {
		desp string
		s    string
		err  string
	}
	testCases := []testCase{
		{
			desp: "normal_all_designator",
			s:    "P3Y6M4DT12H30M5S",
			err:  ``,
		},
		{
			desp: "normal_second_designator",
			s:    "PT5S",
			err:  ``,
		},
		{
			desp: "empty value",
			s:    "",
			err:  `could not parse duration string`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := ValidateISO8601TimeDuration(tc.s)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}
