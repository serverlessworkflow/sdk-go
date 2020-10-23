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

import (
	"encoding/json"
	"fmt"
)

// Workflow base definition
type Workflow struct {
	Id               string     `json:"id"`
	Name             string     `json:"name"`
	Description      string     `json:"description,omitempty"`
	Version          string     `json:"version"`
	SchemaVersion    string     `json:"schemaVersion"`
	DataInputSchema  string     `json:"dataInputSchema,omitempty"`
	DataOutputSchema string     `json:"dataOutputSchema,omitempty"`
	Metadata         Metadata   `json:"metadata,omitempty"`
	Events           []Eventdef `json:"events,omitempty"`
	Functions        []Function `json:"functions,omitempty"`
	States           []State    `json:"states"`
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

var unmarshals = map[string]func(state map[string]interface{}) State{
	"delay":     func(map[string]interface{}) State { return &Delaystate{} },
	"event":     func(map[string]interface{}) State { return &Eventstate{} },
	"operation": func(map[string]interface{}) State { return &Operationstate{} },
	"parallel":  func(map[string]interface{}) State { return &Parallelstate{} },
	"switch":    func(s map[string]interface{}) State {
		if _, ok := s["dataConditions"]; ok {
			return &Databasedswitch{}
		}
		return &Eventbasedswitch{}
	},
	"subflow":   func(map[string]interface{}) State { return &Subflowstate{} },
	"inject":    func(map[string]interface{}) State { return &Injectstate{} },
	"foreach":   func(map[string]interface{}) State { return &Foreachstate{} },
	"callback":  func(map[string]interface{}) State { return &Callbackstate{} },
}

// UnmarshalJSON implementation Unmarshaler interface
// see: http://gregtrowbridge.com/golang-json-serialization-with-interfaces/
func (w *Workflow) UnmarshalJSON(data []byte) error {
	workflowMap := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &workflowMap)
	if err != nil {
		return err
	}
	var rawStates []json.RawMessage
	err = json.Unmarshal(workflowMap["states"], &rawStates)
	if err != nil {
		return err
	}

	w.States = make([]State, len(rawStates))
	var mapState map[string]interface{}
	for i, rawState := range rawStates {
		err = json.Unmarshal(rawState, &mapState)
		if err != nil {
			return err
		}
		if _, ok := mapState["type"]; !ok {
			return fmt.Errorf("state %s not supported", mapState["type"])
		}
		state := unmarshals[mapState["type"].(string)](mapState)
		err := json.Unmarshal(rawState, &state)
		if err != nil {
			return err
		}
		w.States[i] = state
	}
	return nil
}
