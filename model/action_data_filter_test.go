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

func TestActionDataFilterUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect ActionDataFilter
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			data: `{"fromStateData": "1", "results": "2", "toStateData": "3"}`,
			expect: ActionDataFilter{
				FromStateData: "1",
				Results:       "2",
				ToStateData:   "3",
				UseResults:    true,
			},
			err: ``,
		},
		{
			desp: "add UseData to false",
			data: `{"fromStateData": "1", "results": "2", "toStateData": "3", "useResults": false}`,
			expect: ActionDataFilter{
				FromStateData: "1",
				Results:       "2",
				ToStateData:   "3",
				UseResults:    false,
			},
			err: ``,
		},
		{
			desp:   "empty data",
			data:   ` `,
			expect: ActionDataFilter{},
			err:    `unexpected end of JSON input`,
		},
		{
			desp:   "invalid json format",
			data:   `{"fromStateData": 1, "results": "2", "toStateData": "3"}`,
			expect: ActionDataFilter{},
			err:    `json: cannot unmarshal number into Go struct field actionDataFilterForUnmarshal.fromStateData of type string`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v ActionDataFilter
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
