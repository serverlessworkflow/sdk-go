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

func TestEventDataFilterUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect EventDataFilter
		err    string
	}
	testCases := []testCase{
		{
			desp: "normal test",
			data: `{"data": "1", "toStateData": "2"}`,
			expect: EventDataFilter{
				UseData:     true,
				Data:        "1",
				ToStateData: "2",
			},
			err: ``,
		},
		{
			desp: "add UseData to false",
			data: `{"UseData": false, "data": "1", "toStateData": "2"}`,
			expect: EventDataFilter{
				UseData:     false,
				Data:        "1",
				ToStateData: "2",
			},
			err: ``,
		},
		{
			desp:   "empty data",
			data:   ` `,
			expect: EventDataFilter{},
			err:    `unexpected end of JSON input`,
		},
		{
			desp:   "invalid json format",
			data:   `{"data": 1, "toStateData": "2"}`,
			expect: EventDataFilter{},
			err:    `eventDataFilter.data must be an string`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v EventDataFilter
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
