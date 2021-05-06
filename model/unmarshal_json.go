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
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const prefix = "file:/"

func getBytesFromFile(s string) (b []byte, err error) {

	// #nosec
	if resp, err := http.Get(s); err == nil {
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	} else if s, err = filepath.Abs(s); err != nil {
		return nil, err
	}
	if b, err = ioutil.ReadFile(filepath.Clean(s)); err != nil {
		return nil, err
	}
	return b, nil
}

// UnmarshalJSON implementation for json Unmarshal function for the Workflow type
func (w *Workflow) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.WorkflowCommon); err != nil {
		return err
	}

	workflowMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &workflowMap); err != nil {
		return err
	}
	var rawStates []json.RawMessage
	if err := json.Unmarshal(workflowMap["states"], &rawStates); err != nil {
		return err
	}

	w.States = make([]State, len(rawStates))
	var mapState map[string]interface{}
	for i, rawState := range rawStates {
		if err := json.Unmarshal(rawState, &mapState); err != nil {
			return err
		}
		if _, ok := actionsModelMapping[mapState["type"].(string)]; !ok {
			return fmt.Errorf("state %s not supported", mapState["type"])
		}
		state := actionsModelMapping[mapState["type"].(string)](mapState)
		if err := json.Unmarshal(rawState, &state); err != nil {
			return err
		}
		w.States[i] = state
	}
	if _, ok := workflowMap["events"]; ok {
		if err := json.Unmarshal(workflowMap["events"], &w.Events); err != nil {
			var s string
			if err := json.Unmarshal(workflowMap["events"], &s); err != nil {
				return err
			}
			var nestedData []byte
			if nestedData, err = getBytesFromFile(s); err != nil {
				return err
			}
			m := make(map[string][]Eventdef)
			if err := json.Unmarshal(nestedData, &m); err != nil {
				return err
			}
			w.Events = m["events"]
		}
	}
	if _, ok := workflowMap["functions"]; ok {
		if err := json.Unmarshal(workflowMap["functions"], &w.Functions); err != nil {
			var s string
			if err := json.Unmarshal(workflowMap["functions"], &s); err != nil {
				return err
			}
			var nestedData []byte
			if nestedData, err = getBytesFromFile(s); err != nil {
				return err
			}
			m := make(map[string][]Function)
			if err := json.Unmarshal(nestedData, &m); err != nil {
				return err
			}
			w.Functions = m["functions"]
		}
	}
	if _, ok := workflowMap["retries"]; ok {
		if err := json.Unmarshal(workflowMap["retries"], &w.Retries); err != nil {
			var s string
			if err := json.Unmarshal(workflowMap["retries"], &s); err != nil {
				return err
			}
			var nestedData []byte
			if nestedData, err = getBytesFromFile(s); err != nil {
				return err
			}
			m := make(map[string][]Retrydef)
			if err := json.Unmarshal(nestedData, &m); err != nil {
				return err
			}
			w.Retries = m["retries"]
		}
	}
	return nil
}

// UnmarshalJSON implementation for json Unmarshal function for the Eventbasedswitch type
func (j *Eventbasedswitch) UnmarshalJSON(data []byte) error {
	eventBasedSwitch := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &eventBasedSwitch)
	if err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	err = json.Unmarshal(eventBasedSwitch["eventConditions"], &rawConditions)
	if err != nil {
		return err
	}

	j.EventConditions = make([]EventbasedswitchEventConditionsElem, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		err = json.Unmarshal(rawCondition, &mapConditions)
		if err != nil {
			return err
		}
		var condition EventbasedswitchEventConditionsElem
		if _, ok := mapConditions["end"]; ok {
			condition = &Enddeventcondition{}
		} else {
			condition = &Transitioneventcondition{}
		}
		err := json.Unmarshal(rawCondition, condition)
		if err != nil {
			return err
		}
		j.EventConditions[i] = condition
	}
	return nil
}

// UnmarshalJSON implementation for json Unmarshal function for the Databasedswitch type
func (j *Databasedswitch) UnmarshalJSON(data []byte) error {
	dataBasedSwitch := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &dataBasedSwitch)
	if err != nil {
		return err
	}
	var rawConditions []json.RawMessage
	err = json.Unmarshal(dataBasedSwitch["dataConditions"], &rawConditions)
	if err != nil {
		return err
	}

	j.DataConditions = make([]DatabasedswitchDataConditionsElem, len(rawConditions))
	var mapConditions map[string]interface{}
	for i, rawCondition := range rawConditions {
		err = json.Unmarshal(rawCondition, &mapConditions)
		if err != nil {
			return err
		}
		var condition DatabasedswitchDataConditionsElem
		if _, ok := mapConditions["end"]; ok {
			condition = &Enddatacondition{}
		} else {
			condition = &Transitiondatacondition{}
		}
		err := json.Unmarshal(rawCondition, condition)
		if err != nil {
			return err
		}
		j.DataConditions[i] = condition
	}
	return nil
}
