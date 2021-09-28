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
	"github.com/serverlessworkflow/sdk-go/v2/util/floatstr"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// Retry ...
type Retry struct {
	// Unique retry strategy name
	Name string `json:"name" validate:"required"`
	// Time delay between retry attempts (ISO 8601 duration format)
	Delay string `json:"delay,omitempty"`
	// Maximum time delay between retry attempts (ISO 8601 duration format)
	MaxDelay string `json:"maxDelay,omitempty"`
	// Static value by which the delay increases during each attempt (ISO 8601 time format)
	Increment string `json:"increment,omitempty"`
	// Numeric value, if specified the delay between retries is multiplied by this value.
	Multiplier floatstr.Float32OrString `json:"multiplier,omitempty" validate:"omitempty,min=0"`
	// Maximum number of retry attempts.
	MaxAttempts intstr.IntOrString `json:"maxAttempts" validate:"required"`
	// If float type, maximum amount of random time added or subtracted from the delay between each retry relative to total delay (between 0 and 1). If string type, absolute maximum amount of random time added or subtracted from the delay between each retry (ISO 8601 duration format)
	Jitter floatstr.Float32OrString `json:"jitter,omitempty" validate:"omitempty,min=0,max=1"`
}
