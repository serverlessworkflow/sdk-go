// Copyright 2020 The Serverless Workflow Specification Authors
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

package parser

import (
	"testing"

	"github.com/serverlessworkflow/sdk-go/model"
	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	files := map[string]func(*testing.T, *model.Workflow){
		"./testdata/greetings.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "greeting", w.ID)
			assert.IsType(t, &model.Operationstate{}, w.States[0])
		},
		"./testdata/greetings.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[0])
			assert.Equal(t, "greeting", w.ID)
			assert.NotEmpty(t, w.States[0].(*model.Operationstate).Actions)
			assert.NotNil(t, w.States[0].(*model.Operationstate).Actions[0].FunctionRef)
		},
		"./testdata/eventbasedgreeting.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", *w.Events[0].Name)
			assert.IsType(t, &model.Eventstate{}, w.States[0])
			eventState := w.States[0].(*model.Eventstate)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			assert.Equal(t, true, eventState.Exclusive)
		},
		"./testdata/eventbasedgreeting.sw.p.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", *w.Events[0].Name)
			assert.IsType(t, &model.Eventstate{}, w.States[0])
			eventState := w.States[0].(*model.Eventstate)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
		},
		"./testdata/eventbasedswitch.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Eventbasedswitch{}, w.States[0])
			eventState := w.States[0].(*model.Eventbasedswitch)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.EventConditions)
			assert.IsType(t, &model.Transitioneventcondition{}, eventState.EventConditions[0])
		},
		"./testdata/applicationrequest.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Databasedswitch{}, w.States[0])
			eventState := w.States[0].(*model.Databasedswitch)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.Transitiondatacondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/applicationrequest.rp.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Databasedswitch{}, w.States[0])
			eventState := w.States[0].(*model.Databasedswitch)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.Transitiondatacondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/applicationrequest.url.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Databasedswitch{}, w.States[0])
			eventState := w.States[0].(*model.Databasedswitch)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.Transitiondatacondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/checkinbox.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[0])
			operationState := w.States[0].(*model.Operationstate)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, w.States, 2)
		},
		"./testdata/testfromissues.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[2])
			switchState := w.States[0].(*model.Databasedswitch)
			eventState := w.States[1].(*model.Subflowstate)
			operationState := w.States[2].(*model.Operationstate)
			assert.NotNil(t, eventState)
			assert.NotNil(t, switchState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, w.States, 3)
			//fmt.Println(switchState)
			assert.Equal(t, "CheckApplication", switchState.GetName())
			assert.IsType(t, &model.Transitiondatacondition{}, switchState.DataConditions[0])
			assert.Equal(t, "StartApplication", switchState.DataConditions[0].(*model.Transitiondatacondition).Transition.NextState)
			//assert.Equal(t, "StartApplication", switchState.DataConditions[0].)
			//fmt.Println("Ch.P. 2")
			assert.Equal(t, "StartApplication", eventState.GetName())
			assert.Equal(t, "startApplicationWorkflowId", eventState.GetId())
			// fmt.Println("Ch.P. 3")
			// assert.Equal(t, "RejectApplication", operationState.GetName())
			// fmt.Println("Ch.P. 4")
		},
		// validates: https://github.com/serverlessworkflow/specification/pull/175/
		"./testdata/provisionorders.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[0])
			operationState := w.States[0].(*model.Operationstate)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, operationState.OnErrors, 3)
			assert.Equal(t, "Missing order id", *operationState.OnErrors[0].Error)
			assert.Equal(t, "Missing order item", *operationState.OnErrors[1].Error)
			assert.Equal(t, "Missing order quantity", *operationState.OnErrors[2].Error)
		},
	}
	for file, f := range files {
		workflow, err := FromFile(file)
		assert.NoError(t, err)
		assert.NotNil(t, workflow)
		f(t, workflow)
	}
}
