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

const (
	DefaultExpressionLang            = "jq"
	ActionModeSequential  ActionMode = "sequential"
	ActionModeParallel    ActionMode = "parallel"
)

type ActionMode string

// BaseWorkflow describes the partial Workflow definition that does not rely on generic interfaces
// to make it easy for custom unmarshalers implementations to unmarshal the common data structure.
type BaseWorkflow struct {
	// Workflow unique identifier
	ID string `json:"id"`
	// Workflow name
	Name string `json:"name"`
	// Workflow description
	Description string `json:"description,omitempty"`
	// Workflow version
	Version string `json:"version"`
	Start   Start  `json:"start"`
	// Serverless Workflow schema version
	SchemaVersion string `json:"schemaVersion"`
	// Identifies the expression language used for workflow expressions. Default is 'jq'
	ExpressionLang string `json:"expressionLang,omitempty"`
	ExecTimeout    ExecTimeout
	// If 'true', workflow instances is not terminated when there are no active execution paths. Instance can be terminated via 'terminate end definition' or reaching defined 'execTimeout'
	KeepActive bool     `json:"keepActive,omitempty"`
	Metadata   Metadata `json:"metadata,omitempty"`
}

type ExecTimeout struct {
	// Timeout duration (ISO 8601 duration format)
	Duration string `json:"duration"`
	// If `false`, workflow instance is allowed to finish current execution. If `true`, current workflow execution is abrupted.
	Interrupt bool `json:"interrupt,omitempty"`
	// Name of a workflow state to be executed before workflow instance is terminated
	RunBefore string `json:"runBefore,omitempty"`
}

// Workflow start definition
type Start struct {
	StateName string   `json:"stateName"`
	Schedule  Schedule `json:"schedule"`
}

// Default definition. Can be either a transition or end definition
type DefaultDef struct {
	Transition Transition `json:"transition,omitempty"`
	End        End        `json:"end,omitempty"`
}

type Schedule struct {
	// Time interval (must be repeating interval) described with ISO 8601 format. Declares when workflow instances will be automatically created.
	Interval string `json:"interval,omitempty"`
	Cron     Cron   `json:"cron,omitempty"`
	// Timezone name used to evaluate the interval & cron-expression. (default: UTC)
	Timezone string `json:"timezone,omitempty"`
}

type Cron struct {
	// Repeating interval (cron expression) describing when the workflow instance should be created
	Expression string `json:"expression"`
	// Specific date and time (ISO 8601 format) when the cron expression invocation is no longer valid
	ValidUntil string `json:"validUntil,omitempty"`
}

type Transition struct {
	// Name of state to transition to
	NextState string `json:"nextState"`
	// Array of events to be produced before the transition happens
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty"`
	// If set to true, triggers workflow compensation when before this transition is taken. Default is false
	Compensate bool `json:"compensate,omitempty"`
}

type Error struct {
	// Domain-specific error name, or '*' to indicate all possible errors
	Error string `json:"error"`
	// Error code. Can be used in addition to the name to help runtimes resolve to technical errors/exceptions. Should not be defined if error is set to '*'
	Code string `json:"code,omitempty"`
	// References a unique name of a retry definition.
	RetryRef string `json:"retryRef,omitempty"`
	// Transition to next state to handle the error. If retryRef is defined, this transition is taken only if retries were unsuccessful.
	Transition Transition `json:"transition,omitempty"`
	// End workflow execution in case of this error. If retryRef is defined, this ends workflow only if retries were unsuccessful.
	End End `json:"end,omitempty"`
}

type OnEvents struct {
	// References one or more unique event names in the defined workflow events
	EventRefs []string `json:"eventRefs"`
	// Specifies how actions are to be performed (in sequence of parallel)
	ActionMode ActionMode `json:"actionMode,omitempty"`
	// Actions to be performed if expression matches
	Actions []Action `json:"actions,omitempty"`
	// Event data filter
	EventDataFilter EventDataFilter `json:"eventDataFilter,omitempty"`
}

type Action struct {
	// Unique action definition name
	Name        string      `json:"name"`
	FunctionRef FunctionRef `json:"functionRef,omitempty"`
	// References a 'trigger' and 'result' reusable event definitions
	EventRef EventRef `json:"eventRef,omitempty"`
	// Time period to wait for function execution to complete
	Timeout string `json:"timeout,omitempty"`
	// Action data filter
	ActionDataFilter ActionDataFilter `json:"actionDataFilter,omitempty"`
}

// State end definition
type End struct {
	// If true, completes all execution flows in the given workflow instance
	Terminate bool `json:"terminate,omitempty"`
	// Defines events that should be produced
	ProduceEvents []ProduceEvent `json:"produceEvents,omitempty"`
	// If set to true, triggers workflow compensation. Default is false
	Compensate bool `json:"compensate,omitempty"`
}

type ProduceEvent struct {
	// References a name of a defined event
	EventRef string `json:"eventRef"`
	// TODO: add object or string data type
	// If String, expression which selects parts of the states data output to become the data of the produced event. If object a custom object to become the data of produced event.
	Data interface{} `json:"data,omitempty"`
	// Add additional event extension context attributes
	ContextAttributes map[string]interface{} `json:"contextAttributes,omitempty"`
}

type StateDataFilter struct {
	// Workflow expression to filter the state data input
	Input string `json:"input,omitempty"`
	// Workflow expression that filters the state data output
	Output string `json:"output,omitempty"`
}

type EventDataFilter struct {
	// Workflow expression that filters of the event data (payload)
	Data string `json:"data,omitempty"`
	// Workflow expression that selects a state data element to which the event payload should be added/merged into. If not specified, denotes, the top-level state data element.
	ToStateData string `json:"toStateData,omitempty"`
}

// Branch Definition
type Branch struct {
	// Branch name
	Name string `json:"name"`
	// Actions to be executed in this branch
	Actions []Action `json:"actions,omitempty"`
	// Unique Id of a workflow to be executed in this branch
	WorkflowID string `json:"workflowId,omitempty"`
}

type ActionDataFilter struct {
	// Workflow expression that selects state data that the state action can use
	FromStateData string `json:"fromStateData,omitempty"`
	// Workflow expression that filters the actions data results
	Results string `json:"results,omitempty"`
	// Workflow expression that selects a state data element to which the action results should be added/merged into. If not specified, denote, the top-level state data element
	ToStateData string `json:"toStateData,omitempty"`
}

type Repeat struct {
	// Expression evaluated against SubFlow state data. SubFlow will repeat execution as long as this expression is true or until the max property count is reached
	Expression string `json:"expression,omitempty"`
	// If true, the expression is evaluated before each repeat execution, if false the expression is evaluated after each repeat execution
	CheckBefore bool `json:"checkBefore,omitempty"`
	// Sets the maximum amount of repeat executions
	Max int `json:"max,omitempty"`
	// If true, repeats executions in a case unhandled errors propagate from the sub-workflow to this state
	ContinueOnError bool `json:"continueOnError,omitempty"`
	// List referencing defined consumed workflow events. SubFlow will repeat execution until one of the defined events is consumed, or until the max property count is reached
	StopOnEvents []string `json:"stopOnEvents,omitempty"`
}
