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

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestFromFile(t *testing.T) {
	files := map[string]func(*testing.T, *model.Workflow){
		"./testdata/greetings.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "greeting", w.ID)
			assert.IsType(t, &model.OperationState{}, w.States[0])
			assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
		},
		"./testdata/greetings.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			assert.Equal(t, "greeting", w.ID)
			assert.NotEmpty(t, w.States[0].(*model.OperationState).Actions)
			assert.NotNil(t, w.States[0].(*model.OperationState).Actions[0].FunctionRef)
			assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
		},
		"./testdata/eventbasedgreeting.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
		},
		"./testdata/eventbasedgreeting.sw.p.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "GreetingEvent", w.Events[0].Name)
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.OnEvents)
			assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
		},
		"./testdata/eventbasedswitch.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.EventBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.EventBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.EventConditions)
			assert.IsType(t, &model.TransitionEventCondition{}, eventState.EventConditions[0])
		},
		"./testdata/applicationrequest.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			assert.Equal(t, "CheckApplication", w.Start.StateName)
			assert.IsType(t, &model.OperationState{}, w.States[1])
			operationState := w.States[1].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Equal(t, "startApplicationWorkflowId", operationState.Actions[0].SubFlowRef.WorkflowID)
		},
		"./testdata/applicationrequest.rp.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/applicationrequest.url.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			eventState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, eventState.DataConditions)
			assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
			assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
		},
		"./testdata/checkinbox.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			operationState := w.States[0].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, w.States, 2)
		},
		// validates: https://github.com/serverlessworkflow/specification/pull/175/
		"./testdata/provisionorders.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.OperationState{}, w.States[0])
			operationState := w.States[0].(*model.OperationState)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, operationState.OnErrors, 3)
			assert.Equal(t, "Missing order id", operationState.OnErrors[0].Error)
			assert.Equal(t, "Missing order item", operationState.OnErrors[1].Error)
			assert.Equal(t, "Missing order quantity", operationState.OnErrors[2].Error)
		}, "./testdata/checkinbox.cron-test.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "0 0/15 * * * ?", w.Start.Schedule.Cron.Expression)
			assert.Equal(t, "checkInboxFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
			assert.Equal(t, "SendTextForHighPriority", w.States[0].GetTransition().NextState)
			assert.False(t, w.States[1].GetEnd().Terminate)
		}, "./testdata/applicationrequest-issue16.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
			dataBaseSwitchState := w.States[0].(*model.DataBasedSwitchState)
			assert.NotNil(t, dataBaseSwitchState)
			assert.NotEmpty(t, dataBaseSwitchState.DataConditions)
			assert.Equal(t, "CheckApplication", w.States[0].GetName())
		},
		// validates: https://github.com/serverlessworkflow/sdk-go/issues/36
		"./testdata/patientonboarding.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.EventState{}, w.States[0])
			eventState := w.States[0].(*model.EventState)
			assert.NotNil(t, eventState)
			assert.NotEmpty(t, w.Retries)
			assert.Len(t, w.Retries, 1)
			assert.Equal(t, float32(0.0), w.Retries[0].Jitter.FloatVal)
			assert.Equal(t, float32(1.1), w.Retries[0].Multiplier.FloatVal)
		},
	}
	for file, f := range files {
		workflow, err := FromFile(file)
		assert.NoError(t, err, "Test File", file)
		assert.NotNil(t, workflow, "Test File", file)
		f(t, workflow)
	}
}
