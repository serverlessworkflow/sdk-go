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
	"github.com/serverlessworkflow/sdk-go/model"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFromFile(t *testing.T) {
	files := map[string]func(*testing.T, *model.Workflow){
		"./testdata/greetings.sw.json": func(t *testing.T, w *model.Workflow) {
			assert.Equal(t, "greeting", w.Id)
			assert.IsType(t, &model.Operationstate{}, w.States[0])
		},
		"./testdata/greetings.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[0])
			assert.Equal(t, "greeting", w.Id)
			assert.NotEmpty(t, w.States[0].(*model.Operationstate).Actions)
			assert.NotNil(t, w.States[0].(*model.Operationstate).Actions[0].FunctionRef)
		},
		"./testdata/eventbasedgreeting.sw.json": func(t *testing.T, w *model.Workflow) {
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
		"./testdata/checkinbox.sw.yaml": func(t *testing.T, w *model.Workflow) {
			assert.IsType(t, &model.Operationstate{}, w.States[0])
			operationState := w.States[0].(*model.Operationstate)
			assert.NotNil(t, operationState)
			assert.NotEmpty(t, operationState.Actions)
			assert.Len(t, w.States, 2)
		},
	}
	for file, f := range files {
		workflow, err := FromFile(file)
		assert.NoError(t, err)
		assert.NotNil(t, workflow)
		f(t, workflow)
	}
}
