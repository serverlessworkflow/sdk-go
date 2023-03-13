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

import "encoding/json"

// DelayState Causes the workflow execution to delay for a specified duration
type DelayState struct {
	// Amount of time (ISO 8601 format) to delay
	// +kubebuilder:validation:Required
	TimeDelay string `json:"timeDelay" validate:"required,iso8601duration"`
}

func (a *DelayState) MarshalJSON() ([]byte, error) {
	custom, err := json.Marshal(&struct {
		TimeDelay string `json:"timeDelay" validate:"required,iso8601duration"`
	}{
		TimeDelay: a.TimeDelay,
	})
	return custom, err
}
