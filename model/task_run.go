// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"errors"
)

// RunTask represents a task configuration to execute external processes.
type RunTask struct {
	TaskBase `json:",inline"`     // Inline TaskBase fields
	Run      RunTaskConfiguration `json:"run" validate:"required"`
}

type RunTaskConfiguration struct {
	Await     *bool        `json:"await,omitempty"`
	Container *Container   `json:"container,omitempty"`
	Script    *Script      `json:"script,omitempty"`
	Shell     *Shell       `json:"shell,omitempty"`
	Workflow  *RunWorkflow `json:"workflow,omitempty"`
}

type Container struct {
	Image       string                 `json:"image" validate:"required"`
	Command     string                 `json:"command,omitempty"`
	Ports       map[string]interface{} `json:"ports,omitempty"`
	Volumes     map[string]interface{} `json:"volumes,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
}

type Script struct {
	Language    string                 `json:"language" validate:"required"`
	Arguments   map[string]interface{} `json:"arguments,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
	InlineCode  *string                `json:"code,omitempty"`
	External    *ExternalResource      `json:"source,omitempty"`
}

type Shell struct {
	Command     string                 `json:"command" validate:"required"`
	Arguments   map[string]interface{} `json:"arguments,omitempty"`
	Environment map[string]string      `json:"environment,omitempty"`
}

type RunWorkflow struct {
	Namespace string                 `json:"namespace" validate:"required,hostname_rfc1123"`
	Name      string                 `json:"name" validate:"required,hostname_rfc1123"`
	Version   string                 `json:"version" validate:"required,semver_pattern"`
	Input     map[string]interface{} `json:"input,omitempty"`
}

// UnmarshalJSON for RunTaskConfiguration to enforce "oneOf" behavior.
func (rtc *RunTaskConfiguration) UnmarshalJSON(data []byte) error {
	temp := struct {
		Await     *bool        `json:"await"`
		Container *Container   `json:"container"`
		Script    *Script      `json:"script"`
		Shell     *Shell       `json:"shell"`
		Workflow  *RunWorkflow `json:"workflow"`
	}{}

	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}

	// Count non-nil fields
	count := 0
	if temp.Container != nil {
		count++
		rtc.Container = temp.Container
	}
	if temp.Script != nil {
		count++
		rtc.Script = temp.Script
	}
	if temp.Shell != nil {
		count++
		rtc.Shell = temp.Shell
	}
	if temp.Workflow != nil {
		count++
		rtc.Workflow = temp.Workflow
	}

	// Ensure only one of the options is set
	if count != 1 {
		return errors.New("invalid RunTaskConfiguration: only one of 'container', 'script', 'shell', or 'workflow' must be specified")
	}

	rtc.Await = temp.Await
	return nil
}

// MarshalJSON for RunTaskConfiguration to ensure proper serialization.
func (rtc *RunTaskConfiguration) MarshalJSON() ([]byte, error) {
	temp := struct {
		Await     *bool        `json:"await,omitempty"`
		Container *Container   `json:"container,omitempty"`
		Script    *Script      `json:"script,omitempty"`
		Shell     *Shell       `json:"shell,omitempty"`
		Workflow  *RunWorkflow `json:"workflow,omitempty"`
	}{
		Await:     rtc.Await,
		Container: rtc.Container,
		Script:    rtc.Script,
		Shell:     rtc.Shell,
		Workflow:  rtc.Workflow,
	}

	return json.Marshal(temp)
}
