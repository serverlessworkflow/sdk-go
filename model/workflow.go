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
	"bytes"
	"encoding/json"
	"fmt"
)

// InvokeKind defines how the target is invoked.
type InvokeKind string

const (
	// InvokeKindSync meaning that worfklow execution should wait until the target completes.
	InvokeKindSync InvokeKind = "sync"

	// InvokeKindAsync meaning that workflow execution should just invoke the target and should not wait until its completion.
	InvokeKindAsync InvokeKind = "async"
)

// ActionMode specifies how actions are to be performed.
type ActionMode string

const (
	// ActionModeSequential specifies actions should be performed in sequence
	ActionModeSequential ActionMode = "sequential"

	// ActionModeParallel specifies actions should be performed in parallel
	ActionModeParallel ActionMode = "parallel"
)

const (
	// UnlimitedTimeout description for unlimited timeouts
	UnlimitedTimeout = "unlimited"
)

type ExpressionLangType string

const (
	//JqExpressionLang ...
	JqExpressionLang ExpressionLangType = "jq"

	// JsonPathExpressionLang ...
	JsonPathExpressionLang ExpressionLangType = "jsonpath"
)

// BaseWorkflow describes the partial Workflow definition that does not rely on generic interfaces
// to make it easy for custom unmarshalers implementations to unmarshal the common data structure.
type BaseWorkflow struct {
	// Workflow unique identifier
	// +optional
	ID string `json:"id,omitempty" validate:"required_without=Key"`
	// Key Domain-specific workflow identifier
	// +optional
	Key string `json:"key,omitempty" validate:"required_without=ID"`
	// Workflow name
	Name string `json:"name,omitempty"`
	// Workflow description.
	// +optional
	Description string `json:"description,omitempty"`
	// Workflow version.
	// +optional
	Version string `json:"version" validate:"omitempty,min=1"`
	// Workflow start definition.
	// +optional
	Start *Start `json:"start,omitempty"`
	// Annotations List of helpful terms describing the workflows intended purpose, subject areas, or other important
	// qualities.
	// +optional
	Annotations []string `json:"annotations,omitempty"`
	// DataInputSchema URI of the JSON Schema used to validate the workflow data input
	// +optional
	DataInputSchema *DataInputSchema `json:"dataInputSchema,omitempty"`
	// Serverless Workflow schema version
	// +kubebuilder:validation:Required
	// +kubebuilder:default="0.8"
	SpecVersion string `json:"specVersion" validate:"required"`
	// Secrets allow you to access sensitive information, such as passwords, OAuth tokens, ssh keys, etc,
	// inside your Workflow Expressions.
	// +optional
	Secrets Secrets `json:"secrets,omitempty"`
	// Constants Workflow constants are used to define static, and immutable, data which is available to
	// Workflow Expressions.
	// +optional
	Constants *Constants `json:"constants,omitempty"`
	// Identifies the expression language used for workflow expressions. Default is 'jq'.
	// +kubebuilder:validation:Enum=jq;jsonpath
	// +kubebuilder:default=jq
	// +optional
	ExpressionLang ExpressionLangType `json:"expressionLang,omitempty" validate:"omitempty,min=1,oneof=jq jsonpath"`
	// Defines the workflow default timeout settings.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`
	// Defines checked errors that can be explicitly handled during workflow execution.
	// +optional
	Errors []Error `json:"errors,omitempty"`
	// If "true", workflow instances is not terminated when there are no active execution paths.
	// Instance can be terminated with "terminate end definition" or reaching defined "workflowExecTimeout"
	// +optional
	KeepActive bool `json:"keepActive,omitempty"`
	// Metadata custom information shared with the runtime.
	// +optional
	Metadata Metadata `json:"metadata,omitempty"`
	// AutoRetries If set to true, actions should automatically be retried on unchecked errors. Default is false
	// +optional
	AutoRetries bool `json:"autoRetries,omitempty"`
	// Auth definitions can be used to define authentication information that should be applied to resources defined
	// in the operation property of function definitions. It is not used as authentication information for the
	// function invocation, but just to access the resource containing the function invocation information.
	// +optional
	Auth AuthArray `json:"auth,omitempty" validate:"omitempty"`
}

type AuthArray []Auth

func (r *AuthArray) UnmarshalJSON(data []byte) error {
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	switch data[0] {
	case '"':
		return r.unmarshalFile(data)
	case '[':
		return r.unmarshalMany(data)
	}

	return fmt.Errorf("auth value '%s' is not supported, it must be an array or string", string(data))
}

func (r *AuthArray) unmarshalFile(data []byte) error {
	b, err := unmarshalFile(data)
	if err != nil {
		return fmt.Errorf("authDefinitions value '%s' is not supported, it must be an object or string", string(data))
	}

	return r.unmarshalMany(b)
}

func (r *AuthArray) unmarshalMany(data []byte) error {
	var auths []Auth
	err := json.Unmarshal(data, &auths)
	if err != nil {
		return fmt.Errorf("authDefinitions value '%s' is not supported, it must be an object or string", string(data))
	}

	*r = auths
	return nil
}

// Workflow base definition
type Workflow struct {
	BaseWorkflow
	// +kubebuilder:validation:MinItems=1
	States []State `json:"states" validate:"required,min=1,dive"`
	// +optional
	Events []Event `json:"events,omitempty"`
	// +optional
	Functions []Function `json:"functions,omitempty"`
	// +optional
	Retries []Retry `json:"retries,omitempty" validate:"dive"`
}

// UnmarshalJSON implementation for json Unmarshal function for the Workflow type
func (w *Workflow) UnmarshalJSON(data []byte) error {
	if err := json.Unmarshal(data, &w.BaseWorkflow); err != nil {
		return err
	}

	workflowMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &workflowMap); err != nil {
		return err
	}

	var rawStates []json.RawMessage
	if _, ok := workflowMap["states"]; ok {
		if err := json.Unmarshal(workflowMap["states"], &rawStates); err != nil {
			return err
		}
	}

	w.States = make([]State, len(rawStates))
	for i, rawState := range rawStates {
		if err := json.Unmarshal(rawState, &w.States[i]); err != nil {
			return err
		}
	}

	// if the start is not defined, use the first state
	if w.BaseWorkflow.Start == nil && len(w.States) > 0 {
		w.BaseWorkflow.Start = &Start{
			StateName: w.States[0].Name,
		}
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

			m := make(map[string][]Event)
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
			m := make(map[string][]Retry)
			if err := json.Unmarshal(nestedData, &m); err != nil {
				return err
			}
			w.Retries = m["retries"]
		}
	}
	if _, ok := workflowMap["errors"]; ok {
		if err := json.Unmarshal(workflowMap["errors"], &w.Errors); err != nil {
			nestedData, err := unmarshalFile(workflowMap["errors"])
			if err != nil {
				return err
			}
			m := make(map[string][]Error)
			if err := json.Unmarshal(nestedData, &m); err != nil {
				return err
			}
			w.Errors = m["errors"]
		}
	}
	w.setDefaults()
	return nil
}

func (w *Workflow) setDefaults() {
	if len(w.ExpressionLang) == 0 {
		w.ExpressionLang = JqExpressionLang
	}
}

// Timeouts ...
type Timeouts struct {
	// WorkflowExecTimeout Workflow execution timeout duration (ISO 8601 duration format). If not specified should
	// be 'unlimited'.
	// +optional
	WorkflowExecTimeout *WorkflowExecTimeout `json:"workflowExecTimeout,omitempty"`
	// StateExecTimeout Total state execution timeout (including retries) (ISO 8601 duration format).
	// +optional
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// ActionExecTimeout Single actions definition execution timeout duration (ISO 8601 duration format).
	// +optional
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,min=1"`
	// BranchExecTimeout Single branch execution timeout duration (ISO 8601 duration format).
	// +optional
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,min=1"`
	// EventTimeout Timeout duration to wait for consuming defined events (ISO 8601 duration format).
	// +optional
	EventTimeout string `json:"eventTimeout,omitempty" validate:"omitempty,min=1"`
}

// UnmarshalJSON ...
func (t *Timeouts) UnmarshalJSON(data []byte) error {
	timeout := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &timeout); err != nil {
		// assumes it's a reference to a file
		file, err := unmarshalFile(data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(file, &t); err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("workflowExecTimeout", timeout, &t.WorkflowExecTimeout); err != nil {
		return err
	}
	if err := unmarshalKey("stateExecTimeout", timeout, &t.StateExecTimeout); err != nil {
		return err
	}
	if err := unmarshalKey("actionExecTimeout", timeout, &t.ActionExecTimeout); err != nil {
		return err
	}
	if err := unmarshalKey("branchExecTimeout", timeout, &t.ActionExecTimeout); err != nil {
		return err
	}
	if err := unmarshalKey("eventTimeout", timeout, &t.ActionExecTimeout); err != nil {
		return err
	}

	return nil
}

// WorkflowExecTimeout  property defines the workflow execution timeout. It is defined using the ISO 8601 duration
// format. If not defined, the workflow execution should be given "unlimited" amount of time to complete.
type WorkflowExecTimeout struct {
	// Workflow execution timeout duration (ISO 8601 duration format). If not specified should be 'unlimited'.
	// +kubebuilder:default=unlimited
	Duration string `json:"duration" validate:"required,min=1"`
	// If false, workflow instance is allowed to finish current execution. If true, current workflow execution
	// is stopped immediately. Default is false.
	// +optional
	Interrupt bool `json:"interrupt,omitempty"`
	// Name of a workflow state to be executed before workflow instance is terminated.
	// +optional
	RunBefore string `json:"runBefore,omitempty" validate:"omitempty,min=1"`
}

// UnmarshalJSON ...
func (w *WorkflowExecTimeout) UnmarshalJSON(data []byte) error {
	execTimeout := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &execTimeout); err != nil {
		w.Duration, err = unmarshalString(data)
		if err != nil {
			return err
		}
	} else {
		if err := unmarshalKey("duration", execTimeout, &w.Duration); err != nil {
			return err
		}
		if err := unmarshalKey("interrupt", execTimeout, &w.Interrupt); err != nil {
			return err
		}
		if err := unmarshalKey("runBefore", execTimeout, &w.RunBefore); err != nil {
			return err
		}
	}
	if len(w.Duration) == 0 {
		w.Duration = UnlimitedTimeout
	}
	return nil
}

// Error declaration for workflow definitions
type Error struct {
	// Name Domain-specific error name.
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// Code OnError code. Can be used in addition to the name to help runtimes resolve to technical errors/exceptions.
	// Should not be defined if error is set to '*'.
	// +optional
	Code string `json:"code,omitempty" validate:"omitempty,min=1"`
	// OnError description.
	// +optional
	Description string `json:"description,omitempty"`
}

// Start definition
type Start struct {
	// Name of the starting workflow state
	// +kubebuilder:validation:Required
	StateName string `json:"stateName" validate:"required"`
	// Define the recurring time intervals or cron expressions at which workflow instances should be automatically
	// started.
	// +optional
	Schedule *Schedule `json:"schedule,omitempty" validate:"omitempty"`
}

// UnmarshalJSON ...
func (s *Start) UnmarshalJSON(data []byte) error {
	startMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &startMap); err != nil {
		s.StateName, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("stateName", startMap, &s.StateName); err != nil {
		return err
	}
	if err := unmarshalKey("schedule", startMap, &s.Schedule); err != nil {
		return err
	}

	return nil
}

// Schedule ...
type Schedule struct {
	// TODO Interval is required if Cron is not set and vice-versa, make a exclusive validation
	// A recurring time interval expressed in the derivative of ISO 8601 format specified below. Declares that
	// workflow instances should be automatically created at the start of each time interval in the series.
	// +optional
	Interval string `json:"interval,omitempty"`
	// Cron expression defining when workflow instances should be automatically created.
	// optional
	Cron *Cron `json:"cron,omitempty"`
	// Timezone name used to evaluate the interval & cron-expression. If the interval specifies a date-time
	// w/ timezone then proper timezone conversion will be applied. (default: UTC).
	// +optional
	Timezone string `json:"timezone,omitempty"`
}

// UnmarshalJSON ...
func (s *Schedule) UnmarshalJSON(data []byte) error {
	scheduleMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &scheduleMap); err != nil {
		s.Interval, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}

	if err := unmarshalKey("interval", scheduleMap, &s.Interval); err != nil {
		return err
	}
	if err := unmarshalKey("cron", scheduleMap, &s.Cron); err != nil {
		return err
	}
	if err := unmarshalKey("timezone", scheduleMap, &s.Timezone); err != nil {
		return err
	}

	return nil
}

// Cron ...
type Cron struct {
	// Cron expression describing when the workflow instance should be created (automatically).
	// +kubebuilder:validation:Required
	Expression string `json:"expression" validate:"required"`
	// Specific date and time (ISO 8601 format) when the cron expression is no longer valid.
	// +optional
	ValidUntil string `json:"validUntil,omitempty" validate:"omitempty,iso8601duration"`
}

// UnmarshalJSON custom unmarshal function for Cron
func (c *Cron) UnmarshalJSON(data []byte) error {
	cron := make(map[string]interface{})
	if err := json.Unmarshal(data, &cron); err != nil {
		c.Expression, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}

	c.Expression = requiresNotNilOrEmpty(cron["expression"])
	c.ValidUntil = requiresNotNilOrEmpty(cron["validUntil"])

	return nil
}

// Transition Serverless workflow states can have one or more incoming and outgoing transitions (from/to other states).
// Each state can define a transition definition that is used to determine which state to transition to next.
type Transition struct {
	// Name of the state to transition to next.
	// +kubebuilder:validation:Required
	NextState string `json:"nextState" validate:"required,min=1"`
	// Array of producedEvent definitions. Events to be produced before the transition takes place.
	// +optional
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty" validate:"omitempty,dive"`
	// If set to true, triggers workflow compensation before this transition is taken. Default is false.
	// +kubebuilder:default=false
	// +optional
	Compensate bool `json:"compensate,omitempty"`
}

// UnmarshalJSON ...
func (e *Transition) UnmarshalJSON(data []byte) error {
	type defTransitionUnmarshal Transition

	obj, str, err := primitiveOrStruct[string, defTransitionUnmarshal](data)
	if err != nil {
		return err
	}

	if obj == nil {
		e.NextState = str
	} else {
		*e = Transition(*obj)
	}
	return nil
}

// OnError ...
type OnError struct {
	// ErrorRef Reference to a unique workflow error definition. Used of errorRefs is not used
	ErrorRef string `json:"errorRef,omitempty"`
	// ErrorRefs References one or more workflow error definitions. Used if errorRef is not used
	ErrorRefs []string `json:"errorRefs,omitempty"`
	// Transition to next state to handle the error. If retryRef is defined, this transition is taken only if retries were unsuccessful.
	Transition *Transition `json:"transition,omitempty"`
	// End workflow execution in case of this error. If retryRef is defined, this ends workflow only if retries were unsuccessful.
	End *End `json:"end,omitempty"`
}

// End definition
type End struct {
	// If true, completes all execution flows in the given workflow instance.
	// +optional
	Terminate bool `json:"terminate,omitempty"`
	// Array of producedEvent definitions. Defines events that should be produced.
	// +optional
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty"`
	// If set to true, triggers workflow compensation before workflow execution completes. Default is false.
	// +optional
	Compensate bool `json:"compensate,omitempty"`
	// Defines that current workflow execution should stop, and execution should continue as a new workflow
	// instance of the provided id
	// +optional
	ContinueAs *ContinueAs `json:"continueAs,omitempty"`
}

// UnmarshalJSON ...
func (e *End) UnmarshalJSON(data []byte) error {
	type endUnmarshal End
	end, endBool, err := primitiveOrStruct[bool, endUnmarshal](data)
	if err != nil {
		return err
	}

	if end == nil {
		e.Terminate = endBool
		e.Compensate = false
	} else {
		*e = End(*end)
	}

	return nil
}

// ContinueAs can be used to stop the current workflow execution and start another one (of the same or a different type)
type ContinueAs struct {
	// Unique id of the workflow to continue execution as.
	// +kubebuilder:validation:Required
	WorkflowID string `json:"workflowId" validate:"required"`
	// Version of the workflow to continue execution as.
	// +optional
	Version string `json:"version,omitempty"`
	// If string type, an expression which selects parts of the states data output to become the workflow data input of
	// continued execution. If object type, a custom object to become the workflow data input of the continued execution
	// +optional
	Data Object `json:"data,omitempty"`
	// WorkflowExecTimeout Workflow execution timeout to be used by the workflow continuing execution.
	// Overwrites any specific settings set by that workflow
	// +optional
	WorkflowExecTimeout WorkflowExecTimeout `json:"workflowExecTimeout,omitempty"`
}

type continueAsForUnmarshal ContinueAs

func (c *ContinueAs) UnmarshalJSON(data []byte) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		return fmt.Errorf("no bytes to unmarshal")
	}

	var err error
	switch data[0] {
	case '"':
		c.WorkflowID, err = unmarshalString(data)
		return err
	case '{':
		v := continueAsForUnmarshal{}
		err = json.Unmarshal(data, &v)
		if err != nil {
			return err
		}

		*c = ContinueAs(v)
		return nil
	}

	return fmt.Errorf("continueAs value '%s' is not supported, it must be an object or string", string(data))
}

// ProduceEvent Defines the event (CloudEvent format) to be produced when workflow execution completes or during a
// workflow transitions. The eventRef property must match the name of one of the defined produced events in the
// events definition.
type ProduceEvent struct {
	// Reference to a defined unique event name in the events definition
	// +kubebuilder:validation:Required
	EventRef string `json:"eventRef" validate:"required"`
	// If String, expression which selects parts of the states data output to become the data of the produced event.
	// If object a custom object to become the data of produced event.
	// +optional
	Data Object `json:"data,omitempty"`
	// Add additional event extension context attributes.
	// +optional
	ContextAttributes map[string]string `json:"contextAttributes,omitempty"`
}

// StateDataFilter ...
type StateDataFilter struct {
	// Workflow expression to filter the state data input
	Input string `json:"input,omitempty"`
	// Workflow expression that filters the state data output
	Output string `json:"output,omitempty"`
}

// DataInputSchema Used to validate the workflow data input against a defined JSON Schema
type DataInputSchema struct {
	// +kubebuilder:validation:Required
	Schema string `json:"schema" validate:"required"`
	// +kubebuilder:validation:Required
	FailOnValidationErrors *bool `json:"failOnValidationErrors" validate:"required"`
}

// UnmarshalJSON ...
func (d *DataInputSchema) UnmarshalJSON(data []byte) error {
	dataInSchema := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &dataInSchema); err != nil {
		d.Schema, err = unmarshalString(data)
		if err != nil {
			return err
		}
		d.FailOnValidationErrors = &TRUE
		return nil
	}
	if err := unmarshalKey("schema", dataInSchema, &d.Schema); err != nil {
		return err
	}
	if err := unmarshalKey("failOnValidationErrors", dataInSchema, &d.FailOnValidationErrors); err != nil {
		return err
	}

	return nil
}

// Secrets allow you to access sensitive information, such as passwords, OAuth tokens, ssh keys, etc inside your
// Workflow Expressions.
type Secrets []string

// UnmarshalJSON ...
func (s *Secrets) UnmarshalJSON(data []byte) error {
	var secretArray []string
	if err := json.Unmarshal(data, &secretArray); err != nil {
		file, err := unmarshalFile(data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(file, &secretArray); err != nil {
			return err
		}
	}
	*s = secretArray
	return nil
}

// Constants Workflow constants are used to define static, and immutable, data which is available to Workflow Expressions.
type Constants struct {
	// Data represents the generic structure of the constants value
	// +optional
	Data map[string]json.RawMessage `json:",omitempty"`
}

// UnmarshalJSON ...
func (c *Constants) UnmarshalJSON(data []byte) error {
	constantData := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &constantData); err != nil {
		// assumes it's a reference to a file
		file, err := unmarshalFile(data)
		if err != nil {
			return err
		}
		if err := json.Unmarshal(file, &constantData); err != nil {
			return err
		}
	}
	c.Data = constantData
	return nil
}
