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
	"os"
	"path/filepath"
	"testing"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestBasicValidation(t *testing.T) {
	rootPath := "./testdata/workflows"
	files, err := os.ReadDir(rootPath)
	assert.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			workflow, err := FromFile(filepath.Join(rootPath, file.Name()))
			if assert.NoError(t, err, "Test File %s", file.Name()) {
				assert.NotEmpty(t, workflow.ID, "Test File %s", file.Name())
				assert.NotEmpty(t, workflow.States, "Test File %s", file.Name())
			}
		}
	}
}

func TestCustomValidators(t *testing.T) {
	rootPath := "./testdata/workflows/witherrors"
	files, err := os.ReadDir(rootPath)
	assert.NoError(t, err)
	for _, file := range files {
		if !file.IsDir() {
			_, err := FromFile(filepath.Join(rootPath, file.Name()))
			assert.Error(t, err, "Test File %s", file.Name())
		}
	}
}

func TestFromFile(t *testing.T) {
	files := []struct {
		name string
		f    func(*testing.T, *model.Workflow)
	}{
		{
			"./testdata/workflows/greetings.sw.json",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.Equal(t, "greeting", w.ID)
				assert.IsType(t, &model.OperationState{}, w.States[0])
				assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
			},
		}, {
			"./testdata/workflows/greetings.sw.yaml",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.IsType(t, &model.OperationState{}, w.States[0])
				assert.Equal(t, "greeting", w.ID)
				assert.NotEmpty(t, w.States[0].(*model.OperationState).Actions)
				assert.NotNil(t, w.States[0].(*model.OperationState).Actions[0].FunctionRef)
				assert.Equal(t, "greetingFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
			},
		}, {
			"./testdata/workflows/greetings-v08-spec.sw.yaml",
			func(t *testing.T, w *model.Workflow) {
				assert.Empty(t, w.Name)
				assert.Empty(t, w.Start)
				assert.IsType(t, &model.OperationState{}, w.States[0])
				assert.Equal(t, "custom.greeting", w.ID)
				assert.NotEmpty(t, w.States[0].(*model.OperationState).Actions)
				assert.NotEmpty(t, w.Functions[0])
				assert.Equal(t, "greetingCustomFunction", w.Functions[0].Name)
				assert.Equal(t, model.FunctionTypeCustom, w.Functions[0].Type)
				assert.NotNil(t, w.States[0].(*model.OperationState).Actions[0].FunctionRef)
				assert.Equal(t, "greetingCustomFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
				assert.Equal(t, "greetingCustomFunction", w.States[0].(*model.OperationState).Actions[0].Name)
			},
		}, {
			"./testdata/workflows/eventbaseddataandswitch.sw.json",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Switch Transitions", w.Name)
				assert.Equal(t, "Start", w.States[0].GetName())
				assert.Equal(t, "CheckVisaStatus", w.States[1].GetName())
				assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
				assert.IsType(t, &model.EventBasedSwitchState{}, w.States[1])
			},
		}, {
			"./testdata/workflows/conditionbasedstate.yaml", func(t *testing.T, w *model.Workflow) {
				operationState := w.States[0].(*model.OperationState)
				assert.Equal(t, "${ .applicants | .age < 18 }", operationState.Actions[0].Condition)
			},
		}, {
			"./testdata/workflows/eventbasedgreeting.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.IsType(t, &model.EventState{}, w.States[0])
				eventState := w.States[0].(*model.EventState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
				assert.Equal(t, true, eventState.Exclusive)
			},
		}, {
			"./testdata/workflows/eventbasedgreetingexclusive.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.Equal(t, "GreetingEvent2", w.Events[1].Name)
				assert.IsType(t, &model.EventState{}, w.States[0])
				eventState := w.States[0].(*model.EventState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
				assert.Equal(t, "GreetingEvent2", eventState.OnEvents[1].EventRefs[0])
				assert.Equal(t, true, eventState.Exclusive)
			},
		}, {
			"./testdata/workflows/eventbasedgreetingnonexclusive.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.Equal(t, "GreetingEvent2", w.Events[1].Name)
				assert.IsType(t, &model.EventState{}, w.States[0])
				eventState := w.States[0].(*model.EventState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
				assert.Equal(t, "GreetingEvent2", eventState.OnEvents[0].EventRefs[1])
				assert.Equal(t, false, eventState.Exclusive)
			},
		}, {
			"./testdata/workflows/eventbasedgreeting.sw.p.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.IsType(t, &model.EventState{}, w.States[0])
				eventState := w.States[0].(*model.EventState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			},
		}, {
			"./testdata/workflows/eventbasedswitch.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Switch Transitions", w.Name)
				assert.IsType(t, &model.EventBasedSwitchState{}, w.States[0])
				eventState := w.States[0].(*model.EventBasedSwitchState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.EventConditions)
				assert.NotEmpty(t, eventState.Name)
				assert.IsType(t, &model.TransitionEventCondition{}, eventState.EventConditions[0])
			},
		}, {
			"./testdata/workflows/applicationrequest.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
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
				assert.NotNil(t, w.Auth)
				assert.NotNil(t, w.Auth.Defs)
				assert.Equal(t, len(w.Auth.Defs), 1)
				assert.Equal(t, "testAuth", w.Auth.Defs[0].Name)
				assert.Equal(t, model.AuthTypeBearer, w.Auth.Defs[0].Scheme)
				bearerProperties := w.Auth.Defs[0].Properties.(*model.BearerAuthProperties).Token
				assert.Equal(t, "test_token", bearerProperties)
			},
		}, {
			"./testdata/workflows/applicationrequest.multiauth.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
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
				assert.NotNil(t, w.Auth)
				assert.NotNil(t, w.Auth.Defs)
				assert.Equal(t, len(w.Auth.Defs), 2)
				assert.Equal(t, "testAuth", w.Auth.Defs[0].Name)
				assert.Equal(t, model.AuthTypeBearer, w.Auth.Defs[0].Scheme)
				bearerProperties := w.Auth.Defs[0].Properties.(*model.BearerAuthProperties).Token
				assert.Equal(t, "test_token", bearerProperties)
				assert.Equal(t, "testAuth2", w.Auth.Defs[1].Name)
				assert.Equal(t, model.AuthTypeBasic, w.Auth.Defs[1].Scheme)
				basicProperties := w.Auth.Defs[1].Properties.(*model.BasicAuthProperties)
				assert.Equal(t, "test_user", basicProperties.Username)
				assert.Equal(t, "test_pwd", basicProperties.Password)
			},
		}, {
			"./testdata/workflows/applicationrequest.rp.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
				eventState := w.States[0].(*model.DataBasedSwitchState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.DataConditions)
				assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			},
		}, {
			"./testdata/workflows/applicationrequest.url.json", func(t *testing.T, w *model.Workflow) {
				assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
				eventState := w.States[0].(*model.DataBasedSwitchState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.DataConditions)
				assert.IsType(t, &model.TransitionDataCondition{}, eventState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			},
		}, {
			"./testdata/workflows/checkinbox.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Check Inbox Workflow", w.Name)
				assert.IsType(t, &model.OperationState{}, w.States[0])
				operationState := w.States[0].(*model.OperationState)
				assert.NotNil(t, operationState)
				assert.NotEmpty(t, operationState.Actions)
				assert.Len(t, w.States, 2)
			},
		}, {
			// validates: https://github.com/serverlessworkflow/specification/pull/175/
			"./testdata/workflows/provisionorders.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Provision Orders", w.Name)
				assert.IsType(t, &model.OperationState{}, w.States[0])
				operationState := w.States[0].(*model.OperationState)
				assert.NotNil(t, operationState)
				assert.NotEmpty(t, operationState.Actions)
				assert.Len(t, operationState.OnErrors, 3)
				assert.Equal(t, "Missing order id", operationState.OnErrors[0].ErrorRef)
				assert.Equal(t, "MissingId", operationState.OnErrors[0].Transition.NextState)
				assert.Equal(t, "Missing order item", operationState.OnErrors[1].ErrorRef)
				assert.Equal(t, "MissingItem", operationState.OnErrors[1].Transition.NextState)
				assert.Equal(t, "Missing order quantity", operationState.OnErrors[2].ErrorRef)
				assert.Equal(t, "MissingQuantity", operationState.OnErrors[2].Transition.NextState)
			},
		}, {
			"./testdata/workflows/checkinbox.cron-test.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Check Inbox Workflow", w.Name)
				assert.Equal(t, "0 0/15 * * * ?", w.Start.Schedule.Cron.Expression)
				assert.Equal(t, "checkInboxFunction", w.States[0].(*model.OperationState).Actions[0].FunctionRef.RefName)
				assert.Equal(t, "SendTextForHighPriority", w.States[0].GetTransition().NextState)
				assert.False(t, w.States[1].GetEnd().Terminate)
			},
		}, {
			"./testdata/workflows/applicationrequest-issue16.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.IsType(t, &model.DataBasedSwitchState{}, w.States[0])
				dataBaseSwitchState := w.States[0].(*model.DataBasedSwitchState)
				assert.NotNil(t, dataBaseSwitchState)
				assert.NotEmpty(t, dataBaseSwitchState.DataConditions)
				assert.Equal(t, "CheckApplication", w.States[0].GetName())
			},
		}, {
			// validates: https://github.com/serverlessworkflow/sdk-go/issues/36
			"./testdata/workflows/patientonboarding.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Patient Onboarding Workflow", w.Name)
				assert.IsType(t, &model.EventState{}, w.States[0])
				eventState := w.States[0].(*model.EventState)
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, w.Retries)
				assert.Len(t, w.Retries, 1)
				assert.Equal(t, float32(0.0), w.Retries[0].Jitter.FloatVal)
				assert.Equal(t, float32(1.1), w.Retries[0].Multiplier.FloatVal)
			},
		}, {
			"./testdata/workflows/greetings-secret.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.Len(t, w.Secrets, 1)
			},
		}, {
			"./testdata/workflows/greetings-secret-file.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.Len(t, w.Secrets, 3)
			},
		}, {
			"./testdata/workflows/greetings-constants-file.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.NotEmpty(t, w.Constants)
				assert.NotEmpty(t, w.Constants.Data["Translations"])
			},
		}, {
			"./testdata/workflows/roomreadings.timeouts.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Room Temp and Humidity Workflow", w.Name)
				assert.NotNil(t, w.Timeouts)
				assert.Equal(t, "PT1H", w.Timeouts.WorkflowExecTimeout.Duration)
				assert.Equal(t, "GenerateReport", w.Timeouts.WorkflowExecTimeout.RunBefore)
			},
		}, {
			"./testdata/workflows/roomreadings.timeouts.file.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Room Temp and Humidity Workflow", w.Name)
				assert.NotNil(t, w.Timeouts)
				assert.Equal(t, "PT1H", w.Timeouts.WorkflowExecTimeout.Duration)
				assert.Equal(t, "GenerateReport", w.Timeouts.WorkflowExecTimeout.RunBefore)
			},
		}, {
			"./testdata/workflows/purchaseorderworkflow.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Purchase Order Workflow", w.Name)
				assert.NotNil(t, w.Timeouts)
				assert.Equal(t, "PT30D", w.Timeouts.WorkflowExecTimeout.Duration)
				assert.Equal(t, "CancelOrder", w.Timeouts.WorkflowExecTimeout.RunBefore)
			},
		},
	}
	for _, file := range files {
		t.Run(
			file.name, func(t *testing.T) {
				workflow, err := FromFile(file.name)
				if assert.NoError(t, err, "Test File %s", file) {
					assert.NotNil(t, workflow, "Test File %s", file)
					file.f(t, workflow)
				}
			},
		)
	}
}
