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

	"github.com/serverlessworkflow/sdk-go/v2/util"
)

// InvokeKind defines how the target is invoked.
type InvokeKind string

func (i InvokeKind) KindValues() []string {
	return []string{
		string(InvokeKindSync),
		string(InvokeKindAsync),
	}
}

func (i InvokeKind) String() string {
	return string(i)
}

const (
	// InvokeKindSync meaning that worfklow execution should wait until the target completes.
	InvokeKindSync InvokeKind = "sync"
	// InvokeKindAsync meaning that workflow execution should just invoke the target and should not wait until its
	// completion.
	InvokeKindAsync InvokeKind = "async"
)

// ActionMode specifies how actions are to be performed.
type ActionMode string

func (i ActionMode) KindValues() []string {
	return []string{
		string(ActionModeSequential),
		string(ActionModeParallel),
	}
}

func (i ActionMode) String() string {
	return string(i)
}

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

func (i ExpressionLangType) KindValues() []string {
	return []string{
		string(JqExpressionLang),
		string(JsonPathExpressionLang),
		string(CELExpressionLang),
	}
}

func (i ExpressionLangType) String() string {
	return string(i)
}

const (
	//JqExpressionLang ...
	JqExpressionLang ExpressionLangType = "jq"

	// JsonPathExpressionLang ...
	JsonPathExpressionLang ExpressionLangType = "jsonpath"

	// CELExpressionLang
	CELExpressionLang ExpressionLangType = "cel"
)

// BaseWorkflow describes the partial Workflow definition that does not rely on generic interfaces
// to make it easy for custom unmarshalers implementations to unmarshal the common data structure.
// +builder-gen:new-call=ApplyDefault
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
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
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
	Secrets Secrets `json:"secrets,omitempty" validate:"unique"`
	// Constants Workflow constants are used to define static, and immutable, data which is available to
	// Workflow Expressions.
	// +optional
	Constants *Constants `json:"constants,omitempty"`
	// Identifies the expression language used for workflow expressions. Default is 'jq'.
	// +kubebuilder:validation:Enum=jq;jsonpath;cel
	// +kubebuilder:default=jq
	// +optional
	ExpressionLang ExpressionLangType `json:"expressionLang,omitempty" validate:"required,oneofkind"`
	// Defines the workflow default timeout settings.
	// +optional
	Timeouts *Timeouts `json:"timeouts,omitempty"`
	// Defines checked errors that can be explicitly handled during workflow execution.
	// +optional
	Errors Errors `json:"errors,omitempty" validate:"unique=Name,dive"`
	// If "true", workflow instances is not terminated when there are no active execution paths.
	// Instance can be terminated with "terminate end definition" or reaching defined "workflowExecTimeout"
	// +optional
	KeepActive bool `json:"keepActive,omitempty"`
	// Metadata custom information shared with the runtime.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Metadata Metadata `json:"metadata,omitempty"`
	// AutoRetries If set to true, actions should automatically be retried on unchecked errors. Default is false
	// +optional
	AutoRetries bool `json:"autoRetries,omitempty"`
	// Auth definitions can be used to define authentication information that should be applied to resources defined
	// in the operation property of function definitions. It is not used as authentication information for the
	// function invocation, but just to access the resource containing the function invocation information.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Auth Auths `json:"auth,omitempty" validate:"unique=Name,dive"`
}

// ApplyDefault set the default values for Workflow
func (w *BaseWorkflow) ApplyDefault() {
	w.SpecVersion = "0.8"
	w.ExpressionLang = JqExpressionLang
}

type Auths []Auth

type authsUnmarshal Auths

// UnmarshalJSON implements json.Unmarshaler
func (r *Auths) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("auth", data, (*authsUnmarshal)(r))
}

type Errors []Error

type errorsUnmarshal Errors

// UnmarshalJSON implements json.Unmarshaler
func (e *Errors) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("errors", data, (*errorsUnmarshal)(e))
}

// Workflow base definition
// +builder-gen:embedded-ignore-method=BaseWorkflow
type Workflow struct {
	BaseWorkflow `json:",inline"`
	// +kubebuilder:pruning:PreserveUnknownFields
	States States `json:"states" validate:"min=1,unique=Name,dive"`
	// +optional
	Events Events `json:"events,omitempty" validate:"unique=Name,dive"`
	// +optional
	Functions Functions `json:"functions,omitempty" validate:"unique=Name,dive"`
	// +optional
	Retries Retries `json:"retries,omitempty" validate:"unique=Name,dive"`
}

type workflowUnmarshal Workflow

// UnmarshalJSON implementation for json Unmarshal function for the Workflow type
func (w *Workflow) UnmarshalJSON(data []byte) error {
	w.ApplyDefault()
	err := util.UnmarshalObject("workflow", data, (*workflowUnmarshal)(w))
	if err != nil {
		return err
	}

	if w.Start == nil && len(w.States) > 0 {
		w.Start = &Start{
			StateName: w.States[0].Name,
		}
	}

	return nil
}

// +kubebuilder:validation:MinItems=1
type States []State

type statesUnmarshal States

// UnmarshalJSON implements json.Unmarshaler
func (s *States) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObject("states", data, (*statesUnmarshal)(s))
}

type Events []Event

type eventsUnmarshal Events

// UnmarshalJSON implements json.Unmarshaler
func (e *Events) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("events", data, (*eventsUnmarshal)(e))
}

type Functions []Function

type functionsUnmarshal Functions

// UnmarshalJSON implements json.Unmarshaler
func (f *Functions) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("functions", data, (*functionsUnmarshal)(f))
}

type Retries []Retry

type retriesUnmarshal Retries

// UnmarshalJSON implements json.Unmarshaler
func (r *Retries) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("retries", data, (*retriesUnmarshal)(r))
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

type timeoutsUnmarshal Timeouts

// UnmarshalJSON implements json.Unmarshaler
func (t *Timeouts) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("timeouts", data, (*timeoutsUnmarshal)(t))
}

// WorkflowExecTimeout  property defines the workflow execution timeout. It is defined using the ISO 8601 duration
// format. If not defined, the workflow execution should be given "unlimited" amount of time to complete.
// +builder-gen:new-call=ApplyDefault
type WorkflowExecTimeout struct {
	// Workflow execution timeout duration (ISO 8601 duration format). If not specified should be 'unlimited'.
	// +kubebuilder:default=unlimited
	Duration string `json:"duration" validate:"required,min=1,iso8601duration"`
	// If false, workflow instance is allowed to finish current execution. If true, current workflow execution
	// is stopped immediately. Default is false.
	// +optional
	Interrupt bool `json:"interrupt,omitempty"`
	// Name of a workflow state to be executed before workflow instance is terminated.
	// +optional
	RunBefore string `json:"runBefore,omitempty" validate:"omitempty,min=1"`
}

type workflowExecTimeoutUnmarshal WorkflowExecTimeout

// UnmarshalJSON implements json.Unmarshaler
func (w *WorkflowExecTimeout) UnmarshalJSON(data []byte) error {
	w.ApplyDefault()
	return util.UnmarshalPrimitiveOrObject("workflowExecTimeout", data, &w.Duration, (*workflowExecTimeoutUnmarshal)(w))
}

// ApplyDefault set the default values for Workflow Exec Timeout
func (w *WorkflowExecTimeout) ApplyDefault() {
	w.Duration = UnlimitedTimeout
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

type startUnmarshal Start

// UnmarshalJSON implements json.Unmarshaler
func (s *Start) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("start", data, &s.StateName, (*startUnmarshal)(s))
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

type scheduleUnmarshal Schedule

// UnmarshalJSON implements json.Unmarshaler
func (s *Schedule) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("schedule", data, &s.Interval, (*scheduleUnmarshal)(s))
}

// Cron ...
type Cron struct {
	// Cron expression describing when the workflow instance should be created (automatically).
	// +kubebuilder:validation:Required
	Expression string `json:"expression" validate:"required"`
	// Specific date and time (ISO 8601 format) when the cron expression is no longer valid.
	// +optional
	ValidUntil string `json:"validUntil,omitempty" validate:"omitempty,iso8601datetime"`
}

type cronUnmarshal Cron

// UnmarshalJSON custom unmarshal function for Cron
func (c *Cron) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("cron", data, &c.Expression, (*cronUnmarshal)(c))
}

// Transition Serverless workflow states can have one or more incoming and outgoing transitions (from/to other states).
// Each state can define a transition definition that is used to determine which state to transition to next.
type Transition struct {
	stateParent *State `json:"-"` // used in validation
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

type transitionUnmarshal Transition

// UnmarshalJSON implements json.Unmarshaler
func (t *Transition) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("transition", data, &t.NextState, (*transitionUnmarshal)(t))
}

// OnError ...
type OnError struct {
	// ErrorRef Reference to a unique workflow error definition. Used of errorRefs is not used
	ErrorRef string `json:"errorRef,omitempty"`
	// ErrorRefs References one or more workflow error definitions. Used if errorRef is not used
	ErrorRefs []string `json:"errorRefs,omitempty" validate:"omitempty,unique"`
	// Transition to next state to handle the error. If retryRef is defined, this transition is taken only if
	// retries were unsuccessful.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
	Transition *Transition `json:"transition,omitempty"`
	// End workflow execution in case of this error. If retryRef is defined, this ends workflow only if
	// retries were unsuccessful.
	// +kubebuilder:validation:Schemaless
	// +kubebuilder:pruning:PreserveUnknownFields
	// +optional
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

type endUnmarshal End

// UnmarshalJSON implements json.Unmarshaler
func (e *End) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("end", data, &e.Terminate, (*endUnmarshal)(e))
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

type continueAsUnmarshal ContinueAs

// UnmarshalJSON implements json.Unmarshaler
func (c *ContinueAs) UnmarshalJSON(data []byte) error {
	return util.UnmarshalPrimitiveOrObject("continueAs", data, &c.WorkflowID, (*continueAsUnmarshal)(c))
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
// +builder-gen:new-call=ApplyDefault
type DataInputSchema struct {
	// +kubebuilder:validation:Required
	Schema string `json:"schema" validate:"required"`
	// +kubebuilder:validation:Required
	FailOnValidationErrors bool `json:"failOnValidationErrors" validate:"required"`
}

type dataInputSchemaUnmarshal DataInputSchema

// UnmarshalJSON implements json.Unmarshaler
func (d *DataInputSchema) UnmarshalJSON(data []byte) error {
	d.ApplyDefault()
	return util.UnmarshalPrimitiveOrObject("dataInputSchema", data, &d.Schema, (*dataInputSchemaUnmarshal)(d))
}

// ApplyDefault set the default values for Data Input Schema
func (d *DataInputSchema) ApplyDefault() {
	d.FailOnValidationErrors = true
}

// Secrets allow you to access sensitive information, such as passwords, OAuth tokens, ssh keys, etc inside your
// Workflow Expressions.
type Secrets []string

type secretsUnmarshal Secrets

// UnmarshalJSON implements json.Unmarshaler
func (s *Secrets) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("secrets", data, (*secretsUnmarshal)(s))
}

// Constants Workflow constants are used to define static, and immutable, data which is available to Workflow Expressions.
type Constants struct {
	// Data represents the generic structure of the constants value
	// +optional
	Data ConstantsData `json:",omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler
func (c *Constants) UnmarshalJSON(data []byte) error {
	return util.UnmarshalObjectOrFile("constants", data, &c.Data)
}

type ConstantsData map[string]json.RawMessage
