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

func TestEventRefUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect EventRef
		err    string
	}
	testCases := []testCase{
		{
			desp: "all field",
			data: `{"invoke": "async"}`,
			expect: EventRef{
				Invoke: InvokeKindAsync,
			},
			err: ``,
		},
		{
			desp: "invoke unset",
			data: `{}`,
			expect: EventRef{
				Invoke: InvokeKindSync,
			},
			err: ``,
		},
		{
			desp:   "invalid json format",
			data:   `{"invoke": 1}`,
			expect: EventRef{},
			err:    `eventRef.invoke must be an sync,async`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v EventRef
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Equal(t, tc.err, err.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, v)
		})
	}
}

func TestEventUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect Event
		err    string
	}
	testCases := []testCase{
		{
			desp: "all field",
			data: `{"dataOnly": false, "kind": "produced"}`,
			expect: Event{
				DataOnly: false,
				Kind:     EventKindProduced,
			},
			err: ``,
		},
		{
			desp: "optional field dataOnly & kind unset",
			data: `{}`,
			expect: Event{
				DataOnly: true,
				Kind:     EventKindConsumed,
			},
			err: ``,
		},
		{
			desp:   "invalid json format",
			data:   `{"dataOnly": "false", "kind": "produced"}`,
			expect: Event{},
			err:    `event.dataOnly must be an bool`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v Event
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
