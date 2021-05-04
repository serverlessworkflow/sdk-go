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

package model

var actionsModelMapping = map[string]func(state map[string]interface{}) State{
	"delay":     func(map[string]interface{}) State { return &DelayState{} },
	"event":     func(map[string]interface{}) State { return &EventState{} },
	"operation": func(map[string]interface{}) State { return &OperationState{} },
	"parallel":  func(map[string]interface{}) State { return &ParallelState{} },
	"switch": func(s map[string]interface{}) State {
		if _, ok := s["dataConditions"]; ok {
			return &DataBasedSwitchState{}
		}
		return &EventBasedSwitchState{}
	},
	"subflow":  func(map[string]interface{}) State { return &SubflowState{} },
	"inject":   func(map[string]interface{}) State { return &InjectState{} },
	"foreach":  func(map[string]interface{}) State { return &ForEachState{} },
	"callback": func(map[string]interface{}) State { return &CallbackState{} },
}

// Workflow base definition
type Workflow struct {
	BaseWorkflow
	States    []State    `json:"states"`
	Events    []Event    `json:"events,omitempty"`
	Functions []Function `json:"functions,omitempty"`
	Retries   []Retry    `json:"retries,omitempty"`
}
