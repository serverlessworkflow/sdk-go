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
	"encoding/json"
	"fmt"
)

const (
	// DefaultExpressionLang ...
	DefaultExpressionLang = "jq"
	// ActionModeSequential ...
	ActionModeSequential ActionMode = "sequential"
	// ActionModeParallel ...
	ActionModeParallel ActionMode = "parallel"
	// UnlimitedTimeout description for unlimited timeouts
	UnlimitedTimeout = "unlimited"
)

var actionsModelMapping = map[string]func(state map[string]interface{}) State{
	StateTypeDelay:     func(map[string]interface{}) State { return &DelayState{} },
	StateTypeEvent:     func(map[string]interface{}) State { return &EventState{} },
	StateTypeOperation: func(map[string]interface{}) State { return &OperationState{} },
	StateTypeParallel:  func(map[string]interface{}) State { return &ParallelState{} },
	StateTypeSwitch: func(s map[string]interface{}) State {
		if _, ok := s["dataConditions"]; ok {
			return &DataBasedSwitchState{}
		}
		return &EventBasedSwitchState{}
	},
	StateTypeInject:   func(map[string]interface{}) State { return &InjectState{} },
	StateTypeForEach:  func(map[string]interface{}) State { return &ForEachState{} },
	StateTypeCallback: func(map[string]interface{}) State { return &CallbackState{} },
	StateTypeSleep:    func(map[string]interface{}) State { return &SleepState{} },
}

// ActionMode ...
type ActionMode string

// BaseWorkflow describes the partial Workflow definition that does not rely on generic interfaces
// to make it easy for custom unmarshalers implementations to unmarshal the common data structure.
type BaseWorkflow struct {
	// Workflow unique identifier
	ID string `json:"id" validate:"omitempty,min=1"`
	// Key Domain-specific workflow identifier
	Key string `json:"key,omitempty" validate:"omitempty,min=1"`
	// Workflow name
	Name string `json:"name" validate:"required"`
	// Workflow description
	Description string `json:"description,omitempty"`
	// Workflow version
	Version string `json:"version" validate:"omitempty,min=1"`
	Start   *Start `json:"start" validate:"required"`
	// Annotations List of helpful terms describing the workflows intended purpose, subject areas, or other important qualities
	Annotations []string `json:"annotations,omitempty"`
	// DataInputSchema URI of the JSON Schema used to validate the workflow data input
	DataInputSchema *DataInputSchema `json:"dataInputSchema,omitempty"`
	// Serverless Workflow schema version
	SpecVersion string `json:"specVersion,omitempty" validate:"required"`
	// Secrets allow you to access sensitive information, such as passwords, OAuth tokens, ssh keys, etc inside your Workflow Expressions.
	Secrets Secrets `json:"secrets,omitempty"`
	// Constants Workflow constants are used to define static, and immutable, data which is available to Workflow Expressions.
	Constants *Constants `json:"constants,omitempty"`
	// Identifies the expression language used for workflow expressions. Default is 'jq'
	ExpressionLang string `json:"expressionLang,omitempty" validate:"omitempty,min=1"`
	// Timeouts definition for Workflow, State, Action, Branch, and Event consumption.
	Timeouts *Timeouts `json:"timeouts,omitempty"`
	// Errors declarations for this Workflow definition
	Errors []Error `json:"errors,omitempty"`
	// If 'true', workflow instances is not terminated when there are no active execution paths. Instance can be terminated via 'terminate end definition' or reaching defined 'execTimeout'
	KeepActive bool `json:"keepActive,omitempty"`
	// Metadata custom information shared with the runtime
	Metadata Metadata `json:"metadata,omitempty"`
	// AutoRetries If set to true, actions should automatically be retried on unchecked errors. Default is false
	AutoRetries bool `json:"autoRetries,omitempty"`
	// Auth definitions can be used to define authentication information that should be applied to resources defined in the operation
	// property of function definitions. It is not used as authentication information for the function invocation,
	// but just to access the resource containing the function invocation information.
	Auth *Auth `json:"auth,omitempty"`
}

// Workflow base definition
type Workflow struct {
	BaseWorkflow
	States    []State    `json:"states" validate:"required,min=1"`
	Events    []Event    `json:"events,omitempty"`
	Functions []Function `json:"functions,omitempty"`
	Retries   []Retry    `json:"retries,omitempty"`
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
		w.ExpressionLang = DefaultExpressionLang
	}
}

// WorkflowRef holds a reference for a workflow definition
type WorkflowRef struct {
	// Sub-workflow unique id
	WorkflowID string `json:"workflowId" validate:"required"`
	// Sub-workflow version
	Version string `json:"version,omitempty"`
}

// UnmarshalJSON ...
func (s *WorkflowRef) UnmarshalJSON(data []byte) error {
	subflowRef := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &subflowRef); err != nil {
		s.WorkflowID, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("version", subflowRef, &s.Version); err != nil {
		return err
	}
	if err := unmarshalKey("workflowId", subflowRef, &s.WorkflowID); err != nil {
		return err
	}

	return nil
}

// Timeouts ...
type Timeouts struct {
	// WorkflowExecTimeout Workflow execution timeout duration (ISO 8601 duration format). If not specified should be 'unlimited'
	WorkflowExecTimeout *WorkflowExecTimeout `json:"workflowExecTimeout,omitempty"`
	// StateExecTimeout Total state execution timeout (including retries) (ISO 8601 duration format)
	StateExecTimeout *StateExecTimeout `json:"stateExecTimeout,omitempty"`
	// ActionExecTimeout Single actions definition execution timeout duration (ISO 8601 duration format)
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,min=1"`
	// BranchExecTimeout Single branch execution timeout duration (ISO 8601 duration format)
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,min=1"`
	// EventTimeout Timeout duration to wait for consuming defined events (ISO 8601 duration format)
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

// WorkflowExecTimeout ...
type WorkflowExecTimeout struct {
	// Duration Workflow execution timeout duration (ISO 8601 duration format). If not specified should be 'unlimited'
	Duration string `json:"duration,omitempty" validate:"omitempty,min=1"`
	// If `false`, workflow instance is allowed to finish current execution. If `true`, current workflow execution is abrupted.
	Interrupt bool `json:"interrupt,omitempty"`
	// Name of a workflow state to be executed before workflow instance is terminated
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

// StateExecTimeout ...
type StateExecTimeout struct {
	// Single state execution timeout, not including retries (ISO 8601 duration format)
	Single string `json:"single,omitempty" validate:"omitempty,min=1"`
	// Total state execution timeout, including retries (ISO 8601 duration format)
	Total string `json:"total" validate:"required"`
}

// UnmarshalJSON ...
func (s *StateExecTimeout) UnmarshalJSON(data []byte) error {
	stateTimeout := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &stateTimeout); err != nil {
		s.Total, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}
	if err := unmarshalKey("total", stateTimeout, &s.Total); err != nil {
		return err
	}
	if err := unmarshalKey("single", stateTimeout, &s.Single); err != nil {
		return err
	}
	return nil
}

// Error declaration for workflow definitions
type Error struct {
	// Name Domain-specific error name
	Name string `json:"name" validate:"required"`
	// Code OnError code. Can be used in addition to the name to help runtimes resolve to technical errors/exceptions. Should not be defined if error is set to '*'
	Code string `json:"code,omitempty" validate:"omitempty,min=1"`
	// OnError description
	Description string `json:"description,omitempty"`
}

// Start definition
type Start struct {
	StateName string    `json:"stateName" validate:"required"`
	Schedule  *Schedule `json:"schedule,omitempty" validate:"omitempty"`
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

// DefaultCondition Can be either a transition or end definition
type DefaultCondition struct {
	Transition Transition `json:"transition,omitempty"`
	End        End        `json:"end,omitempty"`
}

// Schedule ...
type Schedule struct {
	// Time interval (must be repeating interval) described with ISO 8601 format. Declares when workflow instances will be automatically created.
	Interval string `json:"interval,omitempty"`
	Cron     *Cron  `json:"cron,omitempty"`
	// Timezone name used to evaluate the interval & cron-expression. (default: UTC)
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
	// Repeating interval (cron expression) describing when the workflow instance should be created
	Expression string `json:"expression" validate:"required"`
	// Specific date and time (ISO 8601 format) when the cron expression invocation is no longer valid
	ValidUntil string `json:"validUntil,omitempty"`
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

// Transition ...
type Transition struct {
	// Name of state to transition to
	NextState string `json:"nextState" validate:"required,min=1"`
	// Array of events to be produced before the transition happens
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty" validate:"omitempty,dive"`
	// If set to true, triggers workflow compensation when before this transition is taken. Default is false
	Compensate bool `json:"compensate,omitempty"`
}

// UnmarshalJSON ...
func (t *Transition) UnmarshalJSON(data []byte) error {
	transitionMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &transitionMap); err != nil {
		t.NextState, err = unmarshalString(data)
		if err != nil {
			return err
		}
		return nil
	}

	if err := unmarshalKey("compensate", transitionMap, &t.Compensate); err != nil {
		return err
	}
	if err := unmarshalKey("produceEvents", transitionMap, &t.ProduceEvents); err != nil {
		return err
	}
	if err := unmarshalKey("nextState", transitionMap, &t.NextState); err != nil {
		return err
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

// OnEvents ...
type OnEvents struct {
	// References one or more unique event names in the defined workflow events
	EventRefs []string `json:"eventRefs" validate:"required,min=1"`
	// Specifies how actions are to be performed (in sequence of parallel)
	ActionMode ActionMode `json:"actionMode,omitempty"`
	// Actions to be performed if expression matches
	Actions []Action `json:"actions,omitempty" validate:"omitempty,dive"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

// Action ...
type Action struct {
	// Unique action definition name
	Name        string      `json:"name,omitempty"`
	FunctionRef FunctionRef `json:"functionRef,omitempty"`
	// References a 'trigger' and 'result' reusable event definitions
	EventRef EventRef `json:"eventRef,omitempty"`
	// References a sub-workflow to be executed
	SubFlowRef WorkflowRef `json:"subFlowRef,omitempty"`
	// Sleep Defines time period workflow execution should sleep before / after function execution
	Sleep Sleep `json:"sleep,omitempty"`
	// RetryRef References a defined workflow retry definition. If not defined the default retry policy is assumed
	RetryRef string `json:"retryRef,omitempty"`
	// List of unique references to defined workflow errors for which the action should not be retried. Used only when `autoRetries` is set to `true`
	NonRetryableErrors []string `json:"nonRetryableErrors,omitempty" validate:"omitempty,min=1"`
	// List of unique references to defined workflow errors for which the action should be retried. Used only when `autoRetries` is set to `false`
	RetryableErrors []string `json:"retryableErrors,omitempty" validate:"omitempty,min=1"`
	// Action data filter
	ActionDataFilter ActionDataFilter `json:"actionDataFilter,omitempty"`
}

// End definition
type End struct {
	// If true, completes all execution flows in the given workflow instance
	Terminate bool `json:"terminate,omitempty"`
	// Defines events that should be produced
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty"`
	// If set to true, triggers workflow compensation. Default is false
	Compensate bool       `json:"compensate,omitempty"`
	ContinueAs ContinueAs `json:"continueAs,omitempty"`
}

// UnmarshalJSON ...
func (e *End) UnmarshalJSON(data []byte) error {
	endMap := make(map[string]json.RawMessage)
	if err := json.Unmarshal(data, &endMap); err != nil {
		e.Terminate = false
		e.Compensate = false
		return nil
	}

	if err := unmarshalKey("compensate", endMap, &e.Compensate); err != nil {
		return err
	}
	if err := unmarshalKey("terminate", endMap, &e.Terminate); err != nil {
		return err
	}
	if err := unmarshalKey("produceEvents", endMap, &e.ProduceEvents); err != nil {
		return err
	}
	if err := unmarshalKey("continueAs", endMap, &e.ContinueAs); err != nil {
		return err
	}

	return nil
}

// ContinueAs ...
type ContinueAs struct {
	WorkflowRef
	// TODO: add object or string data type
	// If string type, an expression which selects parts of the states data output to become the workflow data input of continued execution. If object type, a custom object to become the workflow data input of the continued execution
	Data interface{} `json:"data,omitempty"`
	// WorkflowExecTimeout Workflow execution timeout to be used by the workflow continuing execution. Overwrites any specific settings set by that workflow
	WorkflowExecTimeout WorkflowExecTimeout `json:"workflowExecTimeout,omitempty"`
}

// ProduceEvent ...
type ProduceEvent struct {
	// References a name of a defined event
	EventRef string `json:"eventRef" validate:"required"`
	// TODO: add object or string data type
	// If String, expression which selects parts of the states data output to become the data of the produced event. If object a custom object to become the data of produced event.
	Data interface{} `json:"data,omitempty"`
	// Add additional event extension context attributes
	ContextAttributes map[string]interface{} `json:"contextAttributes,omitempty"`
}

// StateDataFilter ...
type StateDataFilter struct {
	// Workflow expression to filter the state data input
	Input string `json:"input,omitempty"`
	// Workflow expression that filters the state data output
	Output string `json:"output,omitempty"`
}

// EventDataFilter ...
type EventDataFilter struct {
	// Workflow expression that filters of the event data (payload)
	Data string `json:"data,omitempty"`
	// Workflow expression that selects a state data element to which the event payload should be added/merged into. If not specified, denotes, the top-level state data element.
	ToStateData string `json:"toStateData,omitempty"`
}

// Branch Definition
type Branch struct {
	// Branch name
	Name string `json:"name" validate:"required"`
	// Actions to be executed in this branch
	Actions []Action `json:"actions" validate:"required,min=1"`
	// Timeouts State specific timeouts
	Timeouts BranchTimeouts `json:"timeouts,omitempty"`
}

// BranchTimeouts ...
type BranchTimeouts struct {
	// ActionExecTimeout Single actions definition execution timeout duration (ISO 8601 duration format)
	ActionExecTimeout string `json:"actionExecTimeout,omitempty" validate:"omitempty,min=1"`
	// BranchExecTimeout Single branch execution timeout duration (ISO 8601 duration format)
	BranchExecTimeout string `json:"branchExecTimeout,omitempty" validate:"omitempty,min=1"`
}

// ActionDataFilter ...
type ActionDataFilter struct {
	// Workflow expression that selects state data that the state action can use
	FromStateData string `json:"fromStateData,omitempty"`
	// Workflow expression that filters the actions' data results
	Results string `json:"results,omitempty"`
	// Workflow expression that selects a state data element to which the action results should be added/merged into. If not specified, denote, the top-level state data element
	ToStateData string `json:"toStateData,omitempty"`
}

// DataInputSchema ...
type DataInputSchema struct {
	Schema                 string `json:"schema" validate:"required"`
	FailOnValidationErrors *bool  `json:"failOnValidationErrors" validate:"required"`
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

// Secrets allow you to access sensitive information, such as passwords, OAuth tokens, ssh keys, etc inside your Workflow Expressions.
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

// Sleep ...
type Sleep struct {
	// Before Amount of time (ISO 8601 duration format) to sleep before function/subflow invocation. Does not apply if 'eventRef' is defined.
	Before string `json:"before,omitempty"`
	// After Amount of time (ISO 8601 duration format) to sleep after function/subflow invocation. Does not apply if 'eventRef' is defined.
	After string `json:"after,omitempty"`
}
