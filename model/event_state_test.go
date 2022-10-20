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

func TestEventStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect EventState
		err    string
	}
	testCases := []testCase{
		{
			desp: "all fields set",
			data: `{"name": "1", "Type": "event", "exclusive": false, "onEvents": [{"eventRefs": ["E1", "E2"], "actionMode": "parallel"}], "timeouts": {"actionExecTimeout": "PT5M", "eventTimeout": "PT5M", "stateExecTimeout": "PT5M"}}`,
			expect: EventState{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeEvent,
				},
				Exclusive: false,
				OnEvents: []OnEvents{
					{
						EventRefs:  []string{"E1", "E2"},
						ActionMode: "parallel",
					},
				},
				Timeouts: &EventStateTimeout{
					EventTimeout:      "PT5M",
					ActionExecTimeout: "PT5M",
					StateExecTimeout: &StateExecTimeout{
						Total: "PT5M",
					},
				},
			},
			err: ``,
		},
		{
			desp: "default exclusive",
			data: `{"name": "1", "Type": "event", "onEvents": [{"eventRefs": ["E1", "E2"], "actionMode": "parallel"}], "timeouts": {"actionExecTimeout": "PT5M", "eventTimeout": "PT5M", "stateExecTimeout": "PT5M"}}`,
			expect: EventState{
				BaseState: BaseState{
					Name: "1",
					Type: StateTypeEvent,
				},
				Exclusive: true,
				OnEvents: []OnEvents{
					{
						EventRefs:  []string{"E1", "E2"},
						ActionMode: "parallel",
					},
				},
				Timeouts: &EventStateTimeout{
					EventTimeout:      "PT5M",
					ActionExecTimeout: "PT5M",
					StateExecTimeout: &StateExecTimeout{
						Total: "PT5M",
					},
				},
			},
			err: ``,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			v := EventState{}
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

func TestOnEventsUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect OnEvents
		err    string
	}
	testCases := []testCase{
		{
			desp: "all fields set",
			data: `{"eventRefs": ["E1", "E2"], "actionMode": "parallel"}`,
			expect: OnEvents{
				EventRefs:  []string{"E1", "E2"},
				ActionMode: ActionModeParallel,
			},
			err: ``,
		},
		{
			desp: "default action mode",
			data: `{"eventRefs": ["E1", "E2"]}`,
			expect: OnEvents{
				EventRefs:  []string{"E1", "E2"},
				ActionMode: ActionModeSequential,
			},
			err: ``,
		},
		{
			desp:   "invalid object format",
			data:   `"eventRefs": ["E1", "E2"], "actionMode": "parallel"}`,
			expect: OnEvents{},
			err:    `invalid character ':' after top-level value`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			v := OnEvents{}
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
