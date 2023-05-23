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
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForEachStateUnmarshalJSON(t *testing.T) {
	type testCase struct {
		desp   string
		data   string
		expect *ForEachState
		err    string
	}
	testCases := []testCase{
		{
			desp: "all field",
			data: `{"mode": "sequential"}`,
			expect: &ForEachState{
				Mode: ForEachModeTypeSequential,
			},
			err: ``,
		},
		{
			desp: "mode unset",
			data: `{}`,
			expect: &ForEachState{
				Mode: ForEachModeTypeParallel,
			},
			err: ``,
		},
		{
			desp:   "invalid json format",
			data:   `{"mode": 1}`,
			expect: nil,
			err:    `forEachState.mode must be sequential or parallel`,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			var v ForEachState
			err := json.Unmarshal([]byte(tc.data), &v)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expect, &v)
		})
	}
}

func TestForeachStateToString(t *testing.T) {

	dataObj := FromString("dataF")

	eventRef := EventRef{
		TriggerEventRef:    "triggerEvRef",
		ResultEventRef:     "resEvRef",
		ResultEventTimeout: "3",
		Data:               &dataObj,
	}

	funcRef := FunctionRef{
		RefName:      "funcRefName",
		SelectionSet: "selSet",
	}

	workFlowRef := WorkflowRef{
		WorkflowID:       "76",
		Version:          "0.0.14",
		Invoke:           "invokeKind",
		OnParentComplete: "onPComplete",
	}
	sleep := Sleep{
		Before: "PT13S",
		After:  "PT23S",
	}

	var retryErrs = []string{"errE", "errF"}
	var nonRetryErrs = []string{"nErrE", "nErrF"}

	action := Action{
		ID:                 "13",
		Name:               "ActionName",
		FunctionRef:        &funcRef,
		EventRef:           &eventRef,
		SubFlowRef:         &workFlowRef,
		Sleep:              &sleep,
		RetryableErrors:    retryErrs,
		NonRetryableErrors: nonRetryErrs,
	}
	actions := []Action{action}

	intStr := intstr.IntOrString{
		StrVal: "42",
		IntVal: 42,
	}

	stateExTimeOut := StateExecTimeout{
		Total:  "40S",
		Single: "20S",
	}

	timeouts := ForEachStateTimeout{
		ActionExecTimeout: "10S",
		StateExecTimeout:  &stateExTimeOut,
	}

	state := ForEachState{
		InputCollection:  "inputCollection",
		OutputCollection: "outCollection",
		IterationParam:   "iterationParam",
		Actions:          actions,
		Mode:             "Mode",
		BatchSize:        &intStr,
		Timeouts:         &timeouts,
	}
	value := state.String()
	assert.NotNil(t, value)
	assert.Equal(t, "[inputCollection, outCollection, iterationParam, 42, [[13, ActionName, &{RefName:funcRefName "+
		"Arguments:map[] SelectionSet:selSet Invoke:}, &{TriggerEventRef:triggerEvRef ResultEventRef:resEvRef ResultEventTimeout:3 "+
		"Data:[1, 0, dataF, []] ContextAttributes:map[] Invoke:}, [76, 0.0.14, invokeKind, onPComplete], &{Before:PT13S After:PT23S}, , "+
		"[nErrE nErrF], [errE errF], {FromStateData: UseResults:false Results: ToStateData:}, ]], "+
		"&{StateExecTimeout:[20S, 40S] ActionExecTimeout:10S}, Mode]", value)
}
