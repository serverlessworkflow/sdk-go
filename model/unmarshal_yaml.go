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

import "gopkg.in/yaml.v3"

type workflowForYAML struct {
	WorkflowMeta
	States []yaml.Node `yaml:"states"`
}

type actionForYAML struct {
	ActionDataFilter Actiondatafilter `yaml:"actionDataFilter,omitempty"`
	EventRef         Eventref         `yaml:"eventRef,omitempty"`
	FunctionRef      Functionref      `yaml:"functionRef,omitempty"`
	Name             string           `yaml:"name,omitempty"`
	Timeout          string           `yaml:"timeout,omitempty"`
}

func (a *Action) UnmarshalYAML(value *yaml.Node) error {
	var tmpAction actionForYAML
	if err := value.Decode(&tmpAction); err != nil {
		return err
	}
	return nil
}

func (w *Workflow) UnmarshalYAML(value *yaml.Node) error {
	var tmpWorkflow workflowForYAML
	if err := value.Decode(&tmpWorkflow); err != nil {
		return err
	}
	states := make([]State, len(tmpWorkflow.States))
	for i, state := range tmpWorkflow.States {
		for _, node := range state.Content {
			if _, ok := actionsModelMapping[node.Value]; !ok {
				continue
			}
			var stateMap map[string]interface{}
			if err := state.Decode(&stateMap); err != nil {
				return err
			}

			stateModel := actionsModelMapping[node.Value](stateMap)
			if err := state.Decode(stateModel); err != nil {
				return err
			}

			states[i] = stateModel
		}
	}
	w.States = states

	if err := value.Decode(&tmpWorkflow.WorkflowMeta); err != nil {
		return err
	}
	w.WorkflowMeta = tmpWorkflow.WorkflowMeta
	return nil
}
