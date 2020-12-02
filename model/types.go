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
	"delay":     func(map[string]interface{}) State { return &Delaystate{} },
	"event":     func(map[string]interface{}) State { return &Eventstate{} },
	"operation": func(map[string]interface{}) State { return &Operationstate{} },
	"parallel":  func(map[string]interface{}) State { return &Parallelstate{} },
	"switch": func(s map[string]interface{}) State {
		if _, ok := s["dataConditions"]; ok {
			return &Databasedswitch{}
		}
		return &Eventbasedswitch{}
	},
	"subflow":  func(map[string]interface{}) State { return &Subflowstate{} },
	"inject":   func(map[string]interface{}) State { return &Injectstate{} },
	"foreach":  func(map[string]interface{}) State { return &Foreachstate{} },
	"callback": func(map[string]interface{}) State { return &Callbackstate{} },
}

// WorkflowCommon describes the partial Workflow definition that does not rely on generic interfaces
// to make it easy for custom unmarshalers implementations to unmarshal the common data structure.
type WorkflowCommon struct {
	ID               string   `json:"id"`
	Name             string   `json:"name"`
	Description      string   `json:"description,omitempty"`
	Version          string   `json:"version"`
	SchemaVersion    string   `json:"schemaVersion"`
	DataInputSchema  string   `json:"dataInputSchema,omitempty"`
	DataOutputSchema string   `json:"dataOutputSchema,omitempty"`
	Metadata         Metadata `json:"metadata,omitempty"`
}

// Workflow base definition
type Workflow struct {
	WorkflowCommon
	States    []State    `json:"states"`
	Events    []Eventdef `json:"events,omitempty"`
	Functions []Function `json:"functions,omitempty"`
}

// State definition for a Workflow state
type State interface {
	GetId() string
	GetName() string
	GetType() string
	GetStart() Start
	GetStateDataFilter() Statedatafilter
	GetDataInputSchema() string
	GetDataOutputSchema() string
	GetMetadata() Metadata_1
}
