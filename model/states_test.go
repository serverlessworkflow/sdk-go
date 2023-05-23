// Copyright 2023 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStateToString(t *testing.T) {
	end := End{
		Terminate:     true,
		Compensate:    false,
		ProduceEvents: []ProduceEvent{},
		ContinueAs:    &ContinueAs{},
	}
	baseState := BaseState{
		ID:                  "46",
		Name:                "name",
		End:                 &end,
		Type:                StateTypeOperation,
		Transition:          &Transition{},
		Metadata:            &Metadata{},
		UsedForCompensation: false,
		StateDataFilter:     &StateDataFilter{},
		CompensatedBy:       "compensatedBy",
		OnErrors:            []OnError{},
	}

	state := State{
		BaseState: baseState,
	}
	value := state.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ BaseState:{ ID:46, Name:name, Type:operation, OnErrors:[], Transition:[, [], false], StateDataFilter:{ Input:, Output: }, CompensatedBy:compensatedBy, UsedForCompensation:false, End:{ Terminate:true, ProduceEvents:[], Compensate:false, ContinueAs:{ WorkflowID:, Version:, WorkflowExecTimeout:[, false, ], Data:{ Type:0, IntVal:0, StrVal:, RawValue:[] } } }, Metadata:&map[] } }", value)
}
