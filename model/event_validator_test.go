// Copyright 2021 The Serverless Workflow Specification Authors
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

	"github.com/serverlessworkflow/sdk-go/v2/model/test"
)

func TestEventRefStructLevelValidation(t *testing.T) {
	test.StructLevelValidation(t, []test.ValidationCase[EventRef]{
		{
			Desp: "valid resultEventTimeout",
			Model: EventRef{
				TriggerEventRef:    "example valid",
				ResultEventRef:     "example valid",
				ResultEventTimeout: "PT1H",
				Invoke:             InvokeKindSync,
			},
		},
		{
			Desp: "invalid resultEventTimeout",
			Model: EventRef{
				TriggerEventRef:    "example invalid",
				ResultEventRef:     "example invalid red",
				ResultEventTimeout: "10hs",
				Invoke:             InvokeKindSync,
			},
			Err: `Key: 'EventRef.ResultEventTimeout' Error:Field validation for 'ResultEventTimeout' failed on the 'iso8601duration' tag`,
		},
	})

}
