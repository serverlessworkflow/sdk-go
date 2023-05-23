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
	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/intstr"
	"testing"
)

func TestRetryToString(t *testing.T) {

	multiplier := floatstr.Float32OrString{
		StrVal: "4",
		Type:   1,
	}

	maxAttempt := intstr.IntOrString{
		StrVal: "7",
		Type:   1,
	}
	jitter := floatstr.Float32OrString{
		StrVal: "10",
		Type:   1,
	}

	retry := Retry{
		Name:        "name",
		Increment:   "1",
		Delay:       "2",
		MaxDelay:    "10",
		Multiplier:  &multiplier,
		MaxAttempts: maxAttempt,
		Jitter:      jitter,
	}
	value := retry.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ Name:name, Delay:2, MaxDelay:10, Increment:1,  Multiplier:4, MaxAttempts:{Type:1 IntVal:0 StrVal:7}, Jitter:{Type:1 FloatVal:0 StrVal:10} }", value)
}
