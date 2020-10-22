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

// Workflow base definition
type Workflow struct {
	Id               string                 `json:"id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description,omitempty"`
	Version          string                 `json:"version"`
	SchemaVersion    string                 `json:"schemaVersion"`
	DataInputSchema  string                 `json:"dataInputSchema,omitempty"`
	DataOutputSchema string                 `json:"dataOutputSchema,omitempty"`
	Metadata         Metadata               `json:"metadata,omitempty"`
	Events           []Eventdef             `json:"events,omitempty"`
	Functions        []Function             `json:"functions,omitempty"`
	States           interface{} `json:"states"`
}
