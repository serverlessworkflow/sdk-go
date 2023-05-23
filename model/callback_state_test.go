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

func TestCallbackStateStructLevelValidation(t *testing.T) {
	type testCase struct {
		desp             string
		callbackStateObj State
		err              string
	}
	testCases := []testCase{
		{
			desp: "normal",
			callbackStateObj: State{
				BaseState: BaseState{
					Name: "callbackTest",
					Type: StateTypeCallback,
					End: &End{
						Terminate: true,
					},
				},
				CallbackState: &CallbackState{
					Action: Action{
						ID:   "1",
						Name: "action1",
					},
					EventRef: "refExample",
				},
			},
			err: ``,
		},
		{
			desp: "missing required EventRef",
			callbackStateObj: State{
				BaseState: BaseState{
					Name: "callbackTest",
					Type: StateTypeCallback,
				},
				CallbackState: &CallbackState{
					Action: Action{
						ID:   "1",
						Name: "action1",
					},
				},
			},
			err: `Key: 'State.CallbackState.EventRef' Error:Field validation for 'EventRef' failed on the 'required' tag`,
		},
		// TODO need to register custom types - will be fixed by https://github.com/serverlessworkflow/sdk-go/issues/151
		//{
		//	desp: "missing required Action",
		//	callbackStateObj: State{
		//		BaseState: BaseState{
		//			Name: "callbackTest",
		//			Type: StateTypeCallback,
		//		},
		//		CallbackState: &CallbackState{
		//			EventRef: "refExample",
		//		},
		//	},
		//	err: `Key: 'State.CallbackState.Action' Error:Field validation for 'Action' failed on the 'required' tag`,
		//},
	}
	for _, tc := range testCases {
		t.Run(tc.desp, func(t *testing.T) {
			err := val.GetValidator().Struct(&tc.callbackStateObj)

			if tc.err != "" {
				assert.Error(t, err)
				assert.Regexp(t, tc.err, err)
				return
			}

			assert.NoError(t, err)
		})
	}
}

func TestCallbackToString(t *testing.T) {

	dataObj := FromString("data")

	eventRef := EventRef{
		TriggerEventRef:    "triggerEvRef",
		ResultEventRef:     "resEvRef",
		ResultEventTimeout: "0",
		Data:               &dataObj,
	}

	funcRef := FunctionRef{
		RefName:      "funcRefName",
		SelectionSet: "selSet",
	}

	workFlowRef := WorkflowRef{
		WorkflowID:       "58",
		Version:          "0.0.1",
		Invoke:           "invokeKind",
		OnParentComplete: "onPComplete",
	}
	sleep := Sleep{
		Before: "PT10S",
		After:  "PT20S",
	}

	var retryErrs = []string{"errA", "errB"}
	var nonRetryErrs = []string{"nErrA", "nErrB"}

	action := Action{
		ID:                 "46",
		Name:               "ActionName",
		FunctionRef:        &funcRef,
		EventRef:           &eventRef,
		SubFlowRef:         &workFlowRef,
		Sleep:              &sleep,
		RetryableErrors:    retryErrs,
		NonRetryableErrors: nonRetryErrs,
	}

	stateExTimeOut := StateExecTimeout{
		Total:  "30S",
		Single: "10S",
	}

	callbackStateTimeouts := CallbackStateTimeout{
		ActionExecTimeout: "10S",
		EventTimeout:      "20S",
		StateExecTimeout:  &stateExTimeOut,
	}

	evDataFilter := EventDataFilter{
		UseData:     true,
		Data:        "data",
		ToStateData: "Next",
	}

	callback := CallbackState{
		Action:          action,
		EventRef:        "eventRef",
		Timeouts:        &callbackStateTimeouts,
		EventDataFilter: &evDataFilter,
	}
	value := callback.String()
	assert.NotNil(t, value)
	assert.Equal(t, "[[46, ActionName, &{RefName:funcRefName Arguments:map[] SelectionSet:selSet Invoke:}, &{TriggerEventRef:triggerEvRef ResultEventRef:resEvRef ResultEventTimeout:0 Data:[1, 0, data, []] ContextAttributes:map[] Invoke:},"+
		" [58, 0.0.1, invokeKind, onPComplete], &{Before:PT10S After:PT20S}, , [nErrA nErrB], [errA errB], {FromStateData: UseResults:false Results: ToStateData:}, ],"+
		" eventRef, [[10S, 30S], 10S, 20S], [true, data, Next]]", value)
}
