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
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"k8s.io/apimachinery/pkg/util/intstr"

	"github.com/stretchr/testify/assert"

	"github.com/serverlessworkflow/sdk-go/v2/model"
	"github.com/serverlessworkflow/sdk-go/v2/test"
)

func TestBasicValidation(t *testing.T) {
	rootPath := "./testdata/workflows"
	files, err := os.ReadDir(rootPath)
	assert.NoError(t, err)

	model.SetIncludePaths(append(model.IncludePaths(), filepath.Join(test.CurrentProjectPath(), "./parser/testdata")))

	for _, file := range files {
		if !file.IsDir() {
			path := filepath.Join(rootPath, file.Name())
			workflow, err := FromFile(path)

			if assert.NoError(t, err, "Test File %s", path) {
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
				assert.IsType(t, &model.OperationState{}, w.States[0].OperationState)
				assert.Equal(t, "greetingFunction", w.States[0].OperationState.Actions[0].FunctionRef.RefName)
				assert.NotNil(t, w.States[0].End)
				assert.True(t, w.States[0].End.Terminate)
			},
		}, {
			"./testdata/workflows/actiondata-defaultvalue.yaml",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "greeting", w.ID)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].OperationState)
				assert.Equal(t, true, w.States[0].OperationState.Actions[0].ActionDataFilter.UseResults)
				assert.Equal(t, "greeting", w.States[0].OperationState.Actions[0].Name)
				assert.NotNil(t, w.States[0].End)
				assert.True(t, w.States[0].End.Terminate)
			},
		}, {
			"./testdata/workflows/greetings.sw.yaml",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Greeting Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.IsType(t, "idx", w.States[0].ID)
				assert.Equal(t, "greeting", w.ID)
				assert.NotEmpty(t, w.States[0].OperationState.Actions)
				assert.NotNil(t, w.States[0].OperationState.Actions[0].FunctionRef)
				assert.Equal(t, "greetingFunction", w.States[0].OperationState.Actions[0].FunctionRef.RefName)
				assert.True(t, w.States[0].End.Terminate)
			},
		}, {
			"./testdata/workflows/eventbaseddataandswitch.sw.json",
			func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Switch Transitions", w.Name)
				assert.Equal(t, "Start", w.States[0].Name)
				assert.Equal(t, "CheckVisaStatus", w.States[1].Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				assert.NotNil(t, w.States[1])
				assert.NotNil(t, w.States[1].SwitchState)
				assert.Equal(t, "PT1H", w.States[1].SwitchState.Timeouts.EventTimeout)
				assert.Nil(t, w.States[1].End)
				assert.NotNil(t, w.States[2].End)
				assert.True(t, w.States[2].End.Terminate)
			},
		}, {
			"./testdata/workflows/conditionbasedstate.yaml", func(t *testing.T, w *model.Workflow) {
				operationState := w.States[0].OperationState
				assert.Equal(t, "${ .applicants | .age < 18 }", operationState.Actions[0].Condition)
				assert.NotNil(t, w.States[0].End)
				assert.True(t, w.States[0].End.Terminate)
			},
		}, {
			"./testdata/workflows/eventbasedgreeting.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.NotNil(t, w.States[0])
				eventState := w.States[0].EventState
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
				assert.Equal(t, true, eventState.Exclusive)
				assert.NotNil(t, w.States[0].End)
				assert.True(t, w.States[0].End.Terminate)
			},
		}, {
			"./testdata/workflows/eventbasedgreetingexclusive.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Greeting Workflow", w.Name)
				assert.Equal(t, "GreetingEvent", w.Events[0].Name)
				assert.Equal(t, "GreetingEvent2", w.Events[1].Name)
				assert.NotNil(t, w.States[0])
				eventState := w.States[0].EventState
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
				assert.NotNil(t, w.States[0])
				eventState := w.States[0].EventState
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
				assert.NotNil(t, w.States[0])
				eventState := w.States[0].EventState
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.OnEvents)
				assert.Equal(t, "GreetingEvent", eventState.OnEvents[0].EventRefs[0])
			},
		}, {
			"./testdata/workflows/eventbasedswitch.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Event Based Switch Transitions", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				assert.NotEmpty(t, w.States[0].EventConditions)
				assert.Equal(t, "CheckVisaStatus", w.States[0].Name)
				assert.IsType(t, model.EventCondition{}, w.States[0].EventConditions[0])
			},
		}, {
			"./testdata/workflows/applicationrequest.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				switchState := w.States[0].SwitchState
				assert.NotNil(t, switchState)
				assert.NotEmpty(t, switchState.DataConditions)
				assert.IsType(t, model.DataCondition{}, switchState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
				assert.Equal(t, "CheckApplication", w.Start.StateName)
				assert.NotNil(t, w.States[1])
				assert.NotNil(t, w.States[1].OperationState)
				operationState := w.States[1].OperationState
				assert.NotNil(t, operationState)
				assert.NotEmpty(t, operationState.Actions)
				assert.Equal(t, "startApplicationWorkflowId", operationState.Actions[0].SubFlowRef.WorkflowID)
				assert.NotNil(t, w.Auth)
				auth := w.Auth
				assert.Equal(t, len(auth), 1)
				assert.Equal(t, "testAuth", auth[0].Name)
				assert.Equal(t, model.AuthTypeBearer, auth[0].Scheme)
				bearerProperties := auth[0].Properties.Bearer.Token
				assert.Equal(t, "test_token", bearerProperties)
			},
		}, {
			"./testdata/workflows/applicationrequest.multiauth.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				switchState := w.States[0].SwitchState
				assert.NotNil(t, switchState)
				assert.NotEmpty(t, switchState.DataConditions)
				assert.IsType(t, model.DataCondition{}, switchState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
				assert.Equal(t, "CheckApplication", w.Start.StateName)
				assert.NotNil(t, w.States[1])
				assert.NotNil(t, w.States[1].OperationState)
				operationState := w.States[1].OperationState
				assert.NotNil(t, operationState)
				assert.NotEmpty(t, operationState.Actions)
				assert.Equal(t, "startApplicationWorkflowId", operationState.Actions[0].SubFlowRef.WorkflowID)
				assert.NotNil(t, w.Auth)
				auth := w.Auth
				assert.Equal(t, len(auth), 2)
				assert.Equal(t, "testAuth", auth[0].Name)
				assert.Equal(t, model.AuthTypeBearer, auth[0].Scheme)
				bearerProperties := auth[0].Properties.Bearer.Token
				assert.Equal(t, "test_token", bearerProperties)
				assert.Equal(t, "testAuth2", auth[1].Name)
				assert.Equal(t, model.AuthTypeBasic, auth[1].Scheme)
				basicProperties := auth[1].Properties.Basic
				assert.Equal(t, "test_user", basicProperties.Username)
				assert.Equal(t, "test_pwd", basicProperties.Password)
				// metadata
				assert.Equal(t, model.Metadata{"metadata1": model.FromString("metadata1"), "metadata2": model.FromString("metadata2")}, w.Metadata)
				assert.Equal(t, model.Metadata{"auth1": model.FromString("auth1"), "auth2": model.FromString("auth2")}, auth[0].Properties.Bearer.Metadata)
			},
		}, {
			"./testdata/workflows/applicationrequest.rp.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				eventState := w.States[0].SwitchState
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.DataConditions)
				assert.IsType(t, model.DataCondition{}, eventState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			},
		}, {
			"./testdata/workflows/applicationrequest.url.json", func(t *testing.T, w *model.Workflow) {
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				eventState := w.States[0].SwitchState
				assert.NotNil(t, eventState)
				assert.NotEmpty(t, eventState.DataConditions)
				assert.IsType(t, model.DataCondition{}, eventState.DataConditions[0])
				assert.Equal(t, "TimeoutRetryStrategy", w.Retries[0].Name)
			},
		}, {
			"./testdata/workflows/checkinbox.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Check Inbox Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].OperationState)
				operationState := w.States[0].OperationState
				assert.NotNil(t, operationState)
				assert.NotEmpty(t, operationState.Actions)
				assert.Len(t, w.States, 2)
			},
		}, {
			// validates: https://github.com/serverlessworkflow/specification/pull/175/
			"./testdata/workflows/provisionorders.sw.json", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Provision Orders", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].OperationState)
				assert.NotEmpty(t, w.States[0].OperationState.Actions)
				assert.Len(t, w.States[0].OnErrors, 3)
				assert.Equal(t, "Missing order id", w.States[0].OnErrors[0].ErrorRef)
				assert.Equal(t, "MissingId", w.States[0].OnErrors[0].Transition.NextState)
				assert.Equal(t, "Missing order item", w.States[0].OnErrors[1].ErrorRef)
				assert.Equal(t, "MissingItem", w.States[0].OnErrors[1].Transition.NextState)
				assert.Equal(t, "Missing order quantity", w.States[0].OnErrors[2].ErrorRef)
				assert.Equal(t, "MissingQuantity", w.States[0].OnErrors[2].Transition.NextState)
			},
		}, {
			"./testdata/workflows/checkinbox.cron-test.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Check Inbox Workflow", w.Name)
				assert.Equal(t, "0 0/15 * * * ?", w.Start.Schedule.Cron.Expression)
				assert.Equal(t, "checkInboxFunction", w.States[0].OperationState.Actions[0].FunctionRef.RefName)
				assert.Equal(t, "SendTextForHighPriority", w.States[0].Transition.NextState)
				assert.Nil(t, w.States[0].End)
				assert.NotNil(t, w.States[1].End)
				assert.True(t, w.States[1].End.Terminate)
			},
		}, {
			"./testdata/workflows/applicationrequest-issue16.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Applicant Request Decision Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].SwitchState)
				switchState := w.States[0].SwitchState
				assert.NotNil(t, switchState)
				assert.NotEmpty(t, switchState.DataConditions)
				assert.Equal(t, "CheckApplication", w.States[0].Name)
			},
		}, {
			// validates: https://github.com/serverlessworkflow/sdk-go/issues/36
			"./testdata/workflows/patientonboarding.sw.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Patient Onboarding Workflow", w.Name)
				assert.NotNil(t, w.States[0])
				assert.NotNil(t, w.States[0].EventState)
				eventState := w.States[0].EventState
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
		}, {
			"./testdata/workflows/continue-as-example.yaml", func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "Notify Customer", w.Name)
				switchState := w.States[1].SwitchState

				assert.NotNil(t, switchState)
				assert.NotEmpty(t, switchState.DataConditions)
				assert.IsType(t, model.DataCondition{}, switchState.DataConditions[0])

				endDataCondition := switchState.DataConditions[0]
				assert.Equal(t, "notifycustomerworkflow", endDataCondition.End.ContinueAs.WorkflowID)
				assert.Equal(t, "1.0", endDataCondition.End.ContinueAs.Version)
				assert.Equal(t, model.FromString("${ del(.customerCount) }"), endDataCondition.End.ContinueAs.Data)
				assert.Equal(t, "GenerateReport", endDataCondition.End.ContinueAs.WorkflowExecTimeout.RunBefore)
				assert.Equal(t, true, endDataCondition.End.ContinueAs.WorkflowExecTimeout.Interrupt)
				assert.Equal(t, "PT1H", endDataCondition.End.ContinueAs.WorkflowExecTimeout.Duration)
			},
		}, {
			name: "./testdata/workflows/greetings-v08-spec.sw.yaml",
			f: func(t *testing.T, w *model.Workflow) {
				assert.Equal(t, "custom.greeting", w.ID)
				assert.Equal(t, "1.0", w.Version)
				assert.Equal(t, "0.8", w.SpecVersion)

				// Workflow "name" no longer a required property
				assert.Empty(t, w.Name)

				// Functions:
				assert.NotEmpty(t, w.Functions[0])
				assert.Equal(t, "greetingCustomFunction", w.Functions[0].Name)
				assert.Equal(t, model.FunctionTypeCustom, w.Functions[0].Type)
				assert.Equal(t, "/path/to/my/script/greeting.ts#CustomGreeting", w.Functions[0].Operation)

				assert.NotEmpty(t, w.Functions[1])
				assert.Equal(t, "sendTextFunction", w.Functions[1].Name)
				assert.Equal(t, model.FunctionTypeGraphQL, w.Functions[1].Type)
				assert.Equal(t, "http://myapis.org/inboxapi.json#sendText", w.Functions[1].Operation)

				assert.NotEmpty(t, w.Functions[2])
				assert.Equal(t, "greetingFunction", w.Functions[2].Name)
				assert.Empty(t, w.Functions[2].Type)
				assert.Equal(t, "file://myapis/greetingapis.json#greeting", w.Functions[2].Operation)

				// Delay state
				assert.NotEmpty(t, w.States[0].DelayState.TimeDelay)
				assert.Equal(t, "GreetDelay", w.States[0].Name)
				assert.Equal(t, model.StateTypeDelay, w.States[0].Type)
				assert.Equal(t, "StoreCarAuctionBid", w.States[0].Transition.NextState)

				// Event state
				assert.NotEmpty(t, w.States[1].EventState.OnEvents)
				assert.Equal(t, "StoreCarAuctionBid", w.States[1].Name)
				assert.Equal(t, model.StateTypeEvent, w.States[1].Type)
				assert.Equal(t, true, w.States[1].EventState.Exclusive)
				assert.NotEmpty(t, true, w.States[1].EventState.OnEvents[0])
				assert.Equal(t, []string{"CarBidEvent"}, w.States[1].EventState.OnEvents[0].EventRefs)
				assert.Equal(t, true, w.States[1].EventState.OnEvents[0].EventDataFilter.UseData)
				assert.Equal(t, "test", w.States[1].EventState.OnEvents[0].EventDataFilter.Data)
				assert.Equal(t, "testing", w.States[1].EventState.OnEvents[0].EventDataFilter.ToStateData)
				assert.Equal(t, model.ActionModeParallel, w.States[1].EventState.OnEvents[0].ActionMode)

				assert.NotEmpty(t, w.States[1].EventState.OnEvents[0].Actions[0].FunctionRef)
				assert.Equal(t, "StoreBidFunction", w.States[1].EventState.OnEvents[0].Actions[0].FunctionRef.RefName)
				assert.Equal(t, "funcref1", w.States[1].EventState.OnEvents[0].Actions[0].Name)
				assert.Equal(t, map[string]model.Object{"bid": model.FromString("${ .bid }")}, w.States[1].EventState.OnEvents[0].Actions[0].FunctionRef.Arguments)

				assert.NotEmpty(t, w.States[1].EventState.OnEvents[0].Actions[1].EventRef)
				assert.Equal(t, "eventRefName", w.States[1].EventState.OnEvents[0].Actions[1].Name)
				assert.Equal(t, "StoreBidFunction", w.States[1].EventState.OnEvents[0].Actions[1].EventRef.ResultEventRef)

				data := model.FromString("${ .patientInfo }")
				assert.Equal(t, &data, w.States[1].EventState.OnEvents[0].Actions[1].EventRef.Data)
				assert.Equal(t, map[string]model.Object{"customer": model.FromString("${ .customer }"), "time": model.FromInt(48)}, w.States[1].EventState.OnEvents[0].Actions[1].EventRef.ContextAttributes)

				assert.Equal(t, "PT1S", w.States[1].EventState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT2S", w.States[1].EventState.Timeouts.StateExecTimeout.Single)
				assert.Equal(t, "PT1H", w.States[1].EventState.Timeouts.EventTimeout)
				assert.Equal(t, "PT3S", w.States[1].EventState.Timeouts.ActionExecTimeout)

				// Parallel state
				assert.NotEmpty(t, w.States[2].ParallelState.Branches)
				assert.Equal(t, "ShortDelayBranch", w.States[2].ParallelState.Branches[0].Name)
				assert.Equal(t, "shortdelayworkflowid", w.States[2].ParallelState.Branches[0].Actions[0].SubFlowRef.WorkflowID)
				assert.Equal(t, "PT5H", w.States[2].ParallelState.Branches[0].Timeouts.ActionExecTimeout)
				assert.Equal(t, "PT6M", w.States[2].ParallelState.Branches[0].Timeouts.BranchExecTimeout)
				assert.Equal(t, "LongDelayBranch", w.States[2].ParallelState.Branches[1].Name)
				assert.Equal(t, "longdelayworkflowid", w.States[2].ParallelState.Branches[1].Actions[0].SubFlowRef.WorkflowID)
				assert.Equal(t, "ParallelExec", w.States[2].Name)
				assert.Equal(t, model.StateTypeParallel, w.States[2].Type)
				assert.Equal(t, model.CompletionTypeAtLeast, w.States[2].ParallelState.CompletionType)
				assert.Equal(t, "PT6M", w.States[2].ParallelState.Timeouts.BranchExecTimeout)
				assert.Equal(t, "PT1S", w.States[2].ParallelState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT2S", w.States[2].ParallelState.Timeouts.StateExecTimeout.Single)
				assert.Equal(t, intstr.IntOrString{IntVal: 13}, w.States[2].ParallelState.NumCompleted)

				// Switch state
				assert.NotEmpty(t, w.States[3].SwitchState.EventConditions)
				assert.Equal(t, "CheckVisaStatusSwitchEventBased", w.States[3].Name)
				assert.Equal(t, model.StateTypeSwitch, w.States[3].Type)
				assert.Equal(t, "visaApprovedEvent", w.States[3].EventConditions[0].Name)
				assert.Equal(t, "visaApprovedEventRef", w.States[3].EventConditions[0].EventRef)
				assert.Equal(t, "HandleApprovedVisa", w.States[3].EventConditions[0].Transition.NextState)
				assert.Equal(t, model.Metadata{"mastercard": model.Object{Type: 1, IntVal: 0, StrVal: "disallowed", RawValue: json.RawMessage(nil)},
					"visa": model.Object{Type: 1, IntVal: 0, StrVal: "allowed", RawValue: json.RawMessage(nil)}},
					w.States[3].EventConditions[0].Metadata)
				assert.Equal(t, "visaRejectedEvent", w.States[3].EventConditions[1].EventRef)
				assert.Equal(t, "HandleRejectedVisa", w.States[3].EventConditions[1].Transition.NextState)
				assert.Equal(t, model.Metadata{"test": model.Object{Type: 1, IntVal: 0, StrVal: "tested", RawValue: json.RawMessage(nil)}},
					w.States[3].EventConditions[1].Metadata)
				assert.Equal(t, "PT1H", w.States[3].SwitchState.Timeouts.EventTimeout)
				assert.Equal(t, "PT1S", w.States[3].SwitchState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT2S", w.States[3].SwitchState.Timeouts.StateExecTimeout.Single)
				assert.Equal(t, &model.Transition{NextState: "HandleNoVisaDecision"}, w.States[3].SwitchState.DefaultCondition.Transition)

				//  DataBasedSwitchState
				dataBased := w.States[4].SwitchState
				assert.NotEmpty(t, dataBased.DataConditions)
				assert.Equal(t, "CheckApplicationSwitchDataBased", w.States[4].Name)
				dataCondition := dataBased.DataConditions[0]
				assert.Equal(t, "${ .applicants | .age >= 18 }", dataCondition.Condition)
				assert.Equal(t, "StartApplication", dataCondition.Transition.NextState)
				assert.Equal(t, &model.Transition{
					NextState: "RejectApplication",
				}, w.States[4].DefaultCondition.Transition)
				assert.Equal(t, "PT1S", w.States[4].SwitchState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT2S", w.States[4].SwitchState.Timeouts.StateExecTimeout.Single)

				// operation state
				assert.NotEmpty(t, w.States[5].OperationState.Actions)
				assert.Equal(t, "GreetSequential", w.States[5].Name)
				assert.Equal(t, model.StateTypeOperation, w.States[5].Type)
				assert.Equal(t, model.ActionModeSequential, w.States[5].OperationState.ActionMode)
				assert.Equal(t, "greetingCustomFunction", w.States[5].OperationState.Actions[0].Name)
				assert.Equal(t, "greetingCustomFunction", w.States[5].OperationState.Actions[0].Name)
				assert.NotNil(t, w.States[5].OperationState.Actions[0].FunctionRef)
				assert.Equal(t, "greetingCustomFunction", w.States[5].OperationState.Actions[0].FunctionRef.RefName)
				assert.Equal(t, "example", w.States[5].OperationState.Actions[0].EventRef.TriggerEventRef)
				assert.Equal(t, "example", w.States[5].OperationState.Actions[0].EventRef.ResultEventRef)
				assert.Equal(t, "PT1H", w.States[5].OperationState.Actions[0].EventRef.ResultEventTimeout)
				assert.Equal(t, "PT1H", w.States[5].OperationState.Timeouts.ActionExecTimeout)
				assert.Equal(t, "PT1S", w.States[5].OperationState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT2S", w.States[5].OperationState.Timeouts.StateExecTimeout.Single)

				// forEach state
				assert.NotEmpty(t, w.States[6].ForEachState.Actions)
				assert.Equal(t, "SendTextForHighPriority", w.States[6].Name)
				assert.Equal(t, model.ForEachModeTypeSequential, w.States[6].ForEachState.Mode)
				assert.Equal(t, model.StateTypeForEach, w.States[6].Type)
				assert.Equal(t, "${ .messages }", w.States[6].ForEachState.InputCollection)
				assert.Equal(t, "${ .outputMessages }", w.States[6].ForEachState.OutputCollection)
				assert.Equal(t, "${ .this }", w.States[6].ForEachState.IterationParam)

				batchSize := intstr.FromInt(45)
				assert.Equal(t, &batchSize, w.States[6].ForEachState.BatchSize)

				assert.NotNil(t, w.States[6].ForEachState.Actions)
				assert.Equal(t, "test", w.States[6].ForEachState.Actions[0].Name)
				assert.NotNil(t, w.States[6].ForEachState.Actions[0].FunctionRef)
				assert.Equal(t, "sendTextFunction", w.States[6].ForEachState.Actions[0].FunctionRef.RefName)
				assert.Equal(t, map[string]model.Object{"message": model.FromString("${ .singlemessage }")}, w.States[6].ForEachState.Actions[0].FunctionRef.Arguments)

				assert.Equal(t, "example1", w.States[6].ForEachState.Actions[0].EventRef.TriggerEventRef)
				assert.Equal(t, "example2", w.States[6].ForEachState.Actions[0].EventRef.ResultEventRef)
				assert.Equal(t, "PT12H", w.States[6].ForEachState.Actions[0].EventRef.ResultEventTimeout)

				assert.Equal(t, "PT11H", w.States[6].ForEachState.Timeouts.ActionExecTimeout)
				assert.Equal(t, "PT11S", w.States[6].ForEachState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT22S", w.States[6].ForEachState.Timeouts.StateExecTimeout.Single)

				// Inject state
				assert.Equal(t, "HelloInject", w.States[7].Name)
				assert.Equal(t, model.StateTypeInject, w.States[7].Type)
				assert.Equal(t, map[string]model.Object{"result": model.FromString("Hello World, last state!")}, w.States[7].InjectState.Data)
				assert.Equal(t, "PT11M", w.States[7].InjectState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT22M", w.States[7].InjectState.Timeouts.StateExecTimeout.Single)

				// callback state
				assert.NotEmpty(t, w.States[8].CallbackState.Action)
				assert.Equal(t, "CheckCreditCallback", w.States[8].Name)
				assert.Equal(t, model.StateTypeCallback, w.States[8].Type)
				assert.Equal(t, "callCreditCheckMicroservice", w.States[8].CallbackState.Action.FunctionRef.RefName)
				assert.Equal(t, map[string]model.Object{"argsObj": model.FromRaw(map[string]interface{}{"age": 10, "name": "hi"}), "customer": model.FromString("${ .customer }"), "time": model.FromInt(48)},
					w.States[8].CallbackState.Action.FunctionRef.Arguments)
				assert.Equal(t, "PT10S", w.States[8].CallbackState.Action.Sleep.Before)
				assert.Equal(t, "PT20S", w.States[8].CallbackState.Action.Sleep.After)
				assert.Equal(t, "PT150M", w.States[8].CallbackState.Timeouts.ActionExecTimeout)
				assert.Equal(t, "PT34S", w.States[8].CallbackState.Timeouts.EventTimeout)
				assert.Equal(t, "PT115M", w.States[8].CallbackState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT22M", w.States[8].CallbackState.Timeouts.StateExecTimeout.Single)

				assert.Equal(t, true, w.States[8].CallbackState.EventDataFilter.UseData)
				assert.Equal(t, "test data", w.States[8].CallbackState.EventDataFilter.Data)
				assert.Equal(t, "${ .customer }", w.States[8].CallbackState.EventDataFilter.ToStateData)

				// sleepState
				assert.NotEmpty(t, w.States[9].SleepState.Duration)
				assert.Equal(t, "WaitForCompletionSleep", w.States[9].Name)
				assert.Equal(t, model.StateTypeSleep, w.States[9].Type)
				assert.Equal(t, "PT5S", w.States[9].SleepState.Duration)
				assert.NotNil(t, w.States[9].SleepState.Timeouts)
				assert.Equal(t, "PT100S", w.States[9].SleepState.Timeouts.StateExecTimeout.Total)
				assert.Equal(t, "PT200S", w.States[9].SleepState.Timeouts.StateExecTimeout.Single)
				assert.Equal(t, true, w.States[9].End.Terminate)

				// switch state with DefaultCondition as string
				assert.NotEmpty(t, w.States[10].SwitchState)
				assert.Equal(t, "HelloStateWithDefaultConditionString", w.States[10].Name)
				assert.Equal(t, "${ true }", w.States[10].SwitchState.DataConditions[0].Condition)
				assert.Equal(t, "HandleApprovedVisa", w.States[10].SwitchState.DataConditions[0].Transition.NextState)
				assert.Equal(t, "SendTextForHighPriority", w.States[10].SwitchState.DefaultCondition.Transition.NextState)
				assert.Equal(t, true, w.States[10].End.Terminate)
			},
		},
	}
	for _, file := range files {
		t.Run(
			file.name, func(t *testing.T) {
				workflow, err := FromFile(file.name)
				if assert.NoError(t, err, "Test File %s", file.name) {
					assert.NotNil(t, workflow, "Test File %s", file.name)
					file.f(t, workflow)
				}
			},
		)
	}
}

func TestUnmarshalWorkflowBasicTests(t *testing.T) {
	t.Run("BasicWorkflowYamlNoAuthDefs", func(t *testing.T) {
		workflow, err := FromYAMLSource([]byte(`
id: helloworld
version: '1.0.0'
specVersion: '0.8'
name: TestUnmarshalWorkflowBasicTests
description: Inject Hello World
start: Hello State
states:
- name: Hello State
  type: inject
  data:
    result: Hello World!
  end: true
`))
		assert.Nil(t, err)
		assert.NotNil(t, workflow)

		b, err := json.Marshal(workflow)
		assert.Nil(t, err)
		assert.True(t, !strings.Contains(string(b), "auth"))

		workflow = nil
		err = json.Unmarshal(b, &workflow)
		assert.Nil(t, err)
	})

	t.Run("BasicWorkflowBasicAuthJSONSource", func(t *testing.T) {
		workflow, err := FromJSONSource([]byte(`
{
  "id": "applicantrequest",
  "version": "1.0",
  "name": "Applicant Request Decision Workflow",
  "description": "Determine if applicant request is valid",
  "start": "Hello State",
  "specVersion": "0.8",
  "auth": [
    {
      "name": "testAuth",
      "scheme": "bearer",
      "properties": {
        "token": "test_token"
      }
    },
    {
      "name": "testAuth2",
      "scheme": "basic",
      "properties": {
        "username": "test_user",
        "password": "test_pwd"
      }
    }
  ],
  "states": [
    {
	  "name": "Hello State",
	  "type": "inject",
      "data": {
		"result": "Hello World!"
	  },
	  "transition": "Next Hello State"
    },
	{
		"name": "Next Hello State",
		"type": "inject",
		"data": {
		  "result": "Next Hello World!"
		},
		"end": true
	  }
  ]
}
`))
		assert.Nil(t, err)
		assert.NotNil(t, workflow.Auth)

		b, _ := json.Marshal(workflow)
		assert.Equal(t, "{\"id\":\"applicantrequest\",\"name\":\"Applicant Request Decision Workflow\",\"description\":\"Determine if applicant request is valid\",\"version\":\"1.0\",\"start\":{\"stateName\":\"Hello State\"},\"specVersion\":\"0.8\",\"expressionLang\":\"jq\",\"auth\":[{\"name\":\"testAuth\",\"scheme\":\"bearer\",\"properties\":{\"token\":\"test_token\"}},{\"name\":\"testAuth2\",\"scheme\":\"basic\",\"properties\":{\"username\":\"test_user\",\"password\":\"test_pwd\"}}],\"states\":[{\"name\":\"Hello State\",\"type\":\"inject\",\"transition\":{\"nextState\":\"Next Hello State\"},\"data\":{\"result\":\"Hello World!\"}},{\"name\":\"Next Hello State\",\"type\":\"inject\",\"end\":{\"terminate\":true},\"data\":{\"result\":\"Next Hello World!\"}}]}",
			string(b))

	})

	t.Run("BasicWorkflowBasicAuthStringJSONSource", func(t *testing.T) {
		workflow, err := FromJSONSource([]byte(`
{
  "id": "applicantrequest",
  "version": "1.0",
  "name": "Applicant Request Decision Workflow",
  "description": "Determine if applicant request is valid",
  "start": "Hello State",
  "specVersion": "0.8",
  "auth": "testdata/workflows/urifiles/auth.json",
  "states": [
    {
	  "name": "Hello State",
	  "type": "inject",
      "data": {
		"result": "Hello World!"
	  },
	  "end": true
    }
  ]
}
`))
		assert.Nil(t, err)
		assert.NotNil(t, workflow.Auth)

		b, _ := json.Marshal(workflow)
		assert.Equal(t, "{\"id\":\"applicantrequest\",\"name\":\"Applicant Request Decision Workflow\",\"description\":\"Determine if applicant request is valid\",\"version\":\"1.0\",\"start\":{\"stateName\":\"Hello State\"},\"specVersion\":\"0.8\",\"expressionLang\":\"jq\",\"auth\":[{\"name\":\"testAuth\",\"scheme\":\"bearer\",\"properties\":{\"token\":\"test_token\"}},{\"name\":\"testAuth2\",\"scheme\":\"basic\",\"properties\":{\"username\":\"test_user\",\"password\":\"test_pwd\"}}],\"states\":[{\"name\":\"Hello State\",\"type\":\"inject\",\"end\":{\"terminate\":true},\"data\":{\"result\":\"Hello World!\"}}]}",
			string(b))

	})

	t.Run("BasicWorkflowInteger", func(t *testing.T) {
		workflow, err := FromJSONSource([]byte(`
{
  "id": "applicantrequest",
  "version": "1.0",
  "name": "Applicant Request Decision Workflow",
  "description": "Determine if applicant request is valid",
  "start": "Hello State",
  "specVersion": "0.7",
  "auth": 123,
  "states": [
    {
	  "name": "Hello State",
	  "type": "inject",
      "data": {
		"result": "Hello World!"
	  },
	  "end": true
    }
  ]
}
`))

		assert.NotNil(t, err)
		assert.Equal(t, "auth must be string or array", err.Error())
		assert.Nil(t, workflow)
	})
}

func TestUnmarshalWorkflowSwitchState(t *testing.T) {
	t.Run("WorkflowStatesTest", func(t *testing.T) {
		workflow, err := FromYAMLSource([]byte(`
id: helloworld
version: '1.0.0'
specVersion: '0.8'
name: WorkflowStatesTest
description: Inject Hello World
start: GreetDelay
metadata:
  metadata1: metadata1
  metadata2: metadata2
auth:
- name: testAuth
  scheme: bearer
  properties:
    token: test_token
    metadata:
      auth1: auth1
      auth2: auth2
states:
- name: GreetDelay
  type: delay
  timeDelay: PT5S
  transition:
    nextState: StoreCarAuctionBid
- name: StoreCarAuctionBid
  type: event
  exclusive: true
  onEvents:
  - eventRefs:
    - CarBidEvent
    eventDataFilter:
      useData: true
      data: "test"
      toStateData: "testing"
    actionMode: parallel
    actions:
    - functionRef:
        refName: StoreBidFunction
        arguments:
          bid: "${ .bid }"
      name: bidFunctionRef
    - eventRef:
        triggerEventRef: StoreBidFunction
        data: "${ .patientInfo }"
        resultEventRef: StoreBidFunction
        contextAttributes:
          customer: "${ .thatBid }"
          time: 32
      name: bidEventRef
  timeouts:
    eventTimeout: PT1H
    actionExecTimeout: PT3S
    stateExecTimeout:
      total: PT1S
      single: PT2S
  transition: ParallelExec
- name: ParallelExec
  type: parallel
  completionType: atLeast
  branches:
    - name: ShortDelayBranch
      actions:
        - subFlowRef: shortdelayworkflowid
      timeouts:
        actionExecTimeout: "PT5H"
        branchExecTimeout: "PT6M"
    - name: LongDelayBranch
      actions:
        - subFlowRef: longdelayworkflowid
  timeouts:
    branchExecTimeout: "PT6M"
    stateExecTimeout:
      total: PT1S
      single: PT2S
  numCompleted: 13
  transition: CheckVisaStatusSwitchEventBased
- name: CheckVisaStatusSwitchEventBased
  type: switch
  eventConditions:
  - name: visaApprovedEvent
    eventRef: visaApprovedEventRef
    transition:
      nextState: HandleApprovedVisa
    metadata:
      visa: allowed
      mastercard: disallowed
  - eventRef: visaRejectedEvent
    transition:
      nextState: HandleRejectedVisa
    metadata:
      test: tested
  timeouts:
    eventTimeout: PT10H
    stateExecTimeout:
      total: PT10S
      single: PT20S
  defaultCondition:
    transition:
      nextState: HelloStateWithDefaultConditionString
- name: HelloStateWithDefaultConditionString
  type: switch
  dataConditions:
  - condition: ${ true }
    transition:
      nextState: HandleApprovedVisa
  - condition: ${ false }
    transition:
      nextState: HandleRejectedVisa
  defaultCondition: SendTextForHighPriority
- name: SendTextForHighPriority
  type: foreach
  inputCollection: "${ .messages }"
  outputCollection: "${ .outputMessages }"
  iterationParam: "${ .this }"
  batchSize: 45
  mode: sequential
  actions:
    - name: test
      functionRef:
        refName: sendTextFunction
        arguments:
          message: "${ .singlemessage }"
      eventRef:
        triggerEventRef: example1
        resultEventRef: example2
        # Added "resultEventTimeout" for action eventref
        resultEventTimeout: PT12H
  timeouts:
    actionExecTimeout: PT11H
    stateExecTimeout:
      total: PT11S
      single: PT22S
  transition: HelloInject
- name: HelloInject
  type: inject
  data:
    result: Hello World, another state!
  timeouts:
    stateExecTimeout:
      total: PT11M
      single: PT22M
  transition: WaitForCompletionSleep
- name: WaitForCompletionSleep
  type: sleep
  duration: PT5S
  timeouts:
    stateExecTimeout:
      total: PT100S
      single: PT200S
  end: 
    terminate: true
- name: CheckCreditCallback
  type: callback
  action:
    functionRef:
      refName: callCreditCheckMicroservice
      arguments:
        customer: "${ .customer }"
        time: 48
        argsObj: {
          "name" : "hi",
          "age": {
            "initial": 10,
            "final": 32
          }
        }
    sleep:
      before: PT10S
      after: PT20S
  eventRef: CreditCheckCompletedEvent
  eventDataFilter:
    useData: true
    data: "test data"
    toStateData: "${ .customer }"
  timeouts:
    actionExecTimeout: PT199M
    eventTimeout: PT348S
    stateExecTimeout:
      total: PT115M
      single: PT22M
  transition: HandleApprovedVisa
- name: HandleApprovedVisa
  type: operation
  actions:
  - subFlowRef:
      workflowId: handleApprovedVisaWorkflowID
    name: subFlowRefName
  - eventRef:
      triggerEventRef: StoreBidFunction
      data: "${ .patientInfo }"
      resultEventRef: StoreBidFunction
      contextAttributes:
        customer: "${ .customer }"
        time: 50
    name: eventRefName
  timeouts:
    actionExecTimeout: PT777S
    stateExecTimeout:
      total: PT33M
      single: PT123M
  end:
    terminate: true
`))
		assert.Nil(t, err)
		fmt.Println(err)
		assert.NotNil(t, workflow)
		b, err := json.Marshal(workflow)

		assert.Nil(t, err)

		// workflow and auth metadata
		assert.True(t, strings.Contains(string(b), "\"metadata\":{\"metadata1\":\"metadata1\",\"metadata2\":\"metadata2\"}"))
		assert.True(t, strings.Contains(string(b), ":{\"metadata\":{\"auth1\":\"auth1\",\"auth2\":\"auth2\"}"))

		// Callback state
		assert.True(t, strings.Contains(string(b), "{\"name\":\"CheckCreditCallback\",\"type\":\"callback\",\"transition\":{\"nextState\":\"HandleApprovedVisa\"},\"action\":{\"functionRef\":{\"refName\":\"callCreditCheckMicroservice\",\"arguments\":{\"argsObj\":{\"age\":{\"final\":32,\"initial\":10},\"name\":\"hi\"},\"customer\":\"${ .customer }\",\"time\":48},\"invoke\":\"sync\"},\"sleep\":{\"before\":\"PT10S\",\"after\":\"PT20S\"},\"actionDataFilter\":{\"useResults\":true}},\"eventRef\":\"CreditCheckCompletedEvent\",\"eventDataFilter\":{\"useData\":true,\"data\":\"test data\",\"toStateData\":\"${ .customer }\"},\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT22M\",\"total\":\"PT115M\"},\"actionExecTimeout\":\"PT199M\",\"eventTimeout\":\"PT348S\"}}"))

		// Operation State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"HandleApprovedVisa\",\"type\":\"operation\",\"end\":{\"terminate\":true},\"actionMode\":\"sequential\",\"actions\":[{\"name\":\"subFlowRefName\",\"subFlowRef\":{\"workflowId\":\"handleApprovedVisaWorkflowID\",\"invoke\":\"sync\",\"onParentComplete\":\"terminate\"},\"actionDataFilter\":{\"useResults\":true}},{\"name\":\"eventRefName\",\"eventRef\":{\"triggerEventRef\":\"StoreBidFunction\",\"resultEventRef\":\"StoreBidFunction\",\"data\":\"${ .patientInfo }\",\"contextAttributes\":{\"customer\":\"${ .customer }\",\"time\":50},\"invoke\":\"sync\"},\"actionDataFilter\":{\"useResults\":true}}],\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT123M\",\"total\":\"PT33M\"},\"actionExecTimeout\":\"PT777S\"}}"))

		// Delay State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"GreetDelay\",\"type\":\"delay\",\"transition\":{\"nextState\":\"StoreCarAuctionBid\"},\"timeDelay\":\"PT5S\"}"))

		// Event State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"StoreCarAuctionBid\",\"type\":\"event\",\"transition\":{\"nextState\":\"ParallelExec\"},\"exclusive\":true,\"onEvents\":[{\"eventRefs\":[\"CarBidEvent\"],\"actionMode\":\"parallel\",\"actions\":[{\"name\":\"bidFunctionRef\",\"functionRef\":{\"refName\":\"StoreBidFunction\",\"arguments\":{\"bid\":\"${ .bid }\"},\"invoke\":\"sync\"},\"actionDataFilter\":{\"useResults\":true}},{\"name\":\"bidEventRef\",\"eventRef\":{\"triggerEventRef\":\"StoreBidFunction\",\"resultEventRef\":\"StoreBidFunction\",\"data\":\"${ .patientInfo }\",\"contextAttributes\":{\"customer\":\"${ .thatBid }\",\"time\":32},\"invoke\":\"sync\"},\"actionDataFilter\":{\"useResults\":true}}],\"eventDataFilter\":{\"useData\":true,\"data\":\"test\",\"toStateData\":\"testing\"}}],\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT2S\",\"total\":\"PT1S\"},\"actionExecTimeout\":\"PT3S\",\"eventTimeout\":\"PT1H\"}}"))

		// Parallel State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"ParallelExec\",\"type\":\"parallel\",\"transition\":{\"nextState\":\"CheckVisaStatusSwitchEventBased\"},\"branches\":[{\"name\":\"ShortDelayBranch\",\"actions\":[{\"subFlowRef\":{\"workflowId\":\"shortdelayworkflowid\",\"invoke\":\"sync\",\"onParentComplete\":\"terminate\"},\"actionDataFilter\":{\"useResults\":true}}],\"timeouts\":{\"actionExecTimeout\":\"PT5H\",\"branchExecTimeout\":\"PT6M\"}},{\"name\":\"LongDelayBranch\",\"actions\":[{\"subFlowRef\":{\"workflowId\":\"longdelayworkflowid\",\"invoke\":\"sync\",\"onParentComplete\":\"terminate\"},\"actionDataFilter\":{\"useResults\":true}}]}],\"completionType\":\"atLeast\",\"numCompleted\":13,\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT2S\",\"total\":\"PT1S\"},\"branchExecTimeout\":\"PT6M\"}}"))

		// Switch State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"CheckVisaStatusSwitchEventBased\",\"type\":\"switch\",\"defaultCondition\":{\"transition\":{\"nextState\":\"HelloStateWithDefaultConditionString\"}},\"eventConditions\":[{\"name\":\"visaApprovedEvent\",\"eventRef\":\"visaApprovedEventRef\",\"metadata\":{\"mastercard\":\"disallowed\",\"visa\":\"allowed\"},\"transition\":{\"nextState\":\"HandleApprovedVisa\"}},{\"eventRef\":\"visaRejectedEvent\",\"metadata\":{\"test\":\"tested\"},\"transition\":{\"nextState\":\"HandleRejectedVisa\"}}],\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT20S\",\"total\":\"PT10S\"},\"eventTimeout\":\"PT10H\"}}"))

		// Switch State with string DefaultCondition
		assert.True(t, strings.Contains(string(b), "{\"name\":\"HelloStateWithDefaultConditionString\",\"type\":\"switch\",\"defaultCondition\":{\"transition\":{\"nextState\":\"SendTextForHighPriority\"}},\"dataConditions\":[{\"condition\":\"${ true }\",\"transition\":{\"nextState\":\"HandleApprovedVisa\"}},{\"condition\":\"${ false }\",\"transition\":{\"nextState\":\"HandleRejectedVisa\"}}]}"))

		// Foreach State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"SendTextForHighPriority\",\"type\":\"foreach\",\"transition\":{\"nextState\":\"HelloInject\"},\"inputCollection\":\"${ .messages }\",\"outputCollection\":\"${ .outputMessages }\",\"iterationParam\":\"${ .this }\",\"batchSize\":45,\"actions\":[{\"name\":\"test\",\"functionRef\":{\"refName\":\"sendTextFunction\",\"arguments\":{\"message\":\"${ .singlemessage }\"},\"invoke\":\"sync\"},\"eventRef\":{\"triggerEventRef\":\"example1\",\"resultEventRef\":\"example2\",\"resultEventTimeout\":\"PT12H\",\"invoke\":\"sync\"},\"actionDataFilter\":{\"useResults\":true}}],\"mode\":\"sequential\",\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT22S\",\"total\":\"PT11S\"},\"actionExecTimeout\":\"PT11H\"}}"))

		// Inject State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"HelloInject\",\"type\":\"inject\",\"transition\":{\"nextState\":\"WaitForCompletionSleep\"},\"data\":{\"result\":\"Hello World, another state!\"},\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT22M\",\"total\":\"PT11M\"}}}"))

		// Sleep State
		assert.True(t, strings.Contains(string(b), "{\"name\":\"WaitForCompletionSleep\",\"type\":\"sleep\",\"end\":{\"terminate\":true},\"duration\":\"PT5S\",\"timeouts\":{\"stateExecTimeout\":{\"single\":\"PT200S\",\"total\":\"PT100S\"}}}"))

		workflow = nil
		err = json.Unmarshal(b, &workflow)
		assert.Nil(t, err)

	})

	t.Run("WorkflowSwitchStateDataConditions with wrong field name", func(t *testing.T) {
		workflow, err := FromYAMLSource([]byte(`
id: helloworld
version: '1.0.0'
specVersion: '0.8'
name: WorkflowSwitchStateDataConditions with wrong field name
description: Inject Hello World
start: Hello State
states:
- name: Hello State
  type: switch
  dataCondition:
  - condition: ${ true }
    transition:
      nextState: HandleApprovedVisa
  - condition: ${ false }
    transition:
      nextState: HandleRejectedVisa
  defaultCondition:
    transition:
      nextState: HandleApprovedVisa
- name: HandleApprovedVisa
  type: operation
  actions:
  - subFlowRef:
      workflowId: handleApprovedVisaWorkflowID
  end:
    terminate: true
- name: HandleRejectedVisa
  type: operation
  actions:
  - subFlowRef:
      workflowId: handleRejectedVisaWorkflowID
  end:
    terminate: true
- name: HandleNoVisaDecision
  type: operation
  actions:
  - subFlowRef:
      workflowId: handleNoVisaDecisionWorkflowId
  end:
    terminate: true
`))
		assert.Error(t, err)
		assert.Regexp(t, `validation for \'DataConditions\' failed on the \'required\' tag`, err)
		assert.Nil(t, workflow)
	})
}
