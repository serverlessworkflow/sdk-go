// Copyright 2023 The Serverless Workflow Specification Authors
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

package builder

import (
	"fmt"

	"github.com/serverlessworkflow/sdk-go/v2/model"
)

type Builder[T, K any] struct {
	builder build[K]
}

func (b *Builder[T, K]) Build() K {
	return b.builder.Build()
}

type build[T any] interface {
	Build() T
}

// NewWorkflowBuilder create a Workflow Builder
func NewWorkflowBuilder(id, version string) (*WorkflowBuilder, *Builder[WorkflowBuilder, model.Workflow]) {
	workflow := model.Workflow{
		BaseWorkflow: model.BaseWorkflow{
			ID:          id,
			SpecVersion: "0.8",
			Version:     version,
		},
	}
	workflow.ApplyDefault()

	workflowBuilder := &WorkflowBuilder{
		workflow:              workflow,
		functions:             make([]*FunctionBuilder[WorkflowBuilder], 0),
		events:                make([]*EventBuilder[WorkflowBuilder], 0),
		operatorStateBuilders: make([]*OperationStateBuilder[WorkflowBuilder], 0),
		eventStateBuilders:    make([]*EventStateBuilder[WorkflowBuilder], 0),
	}

	builder := &Builder[WorkflowBuilder, model.Workflow]{
		builder: workflowBuilder,
	}

	return workflowBuilder, builder
}

type WorkflowBuilder struct {
	workflow              model.Workflow
	startBuilder          *StartBuidler[WorkflowBuilder]
	functions             []*FunctionBuilder[WorkflowBuilder]
	events                []*EventBuilder[WorkflowBuilder]
	operatorStateBuilders []*OperationStateBuilder[WorkflowBuilder]
	eventStateBuilders    []*EventStateBuilder[WorkflowBuilder]
}

func (b *WorkflowBuilder) Name(name string) *WorkflowBuilder {
	b.workflow.BaseWorkflow.Name = name
	return b
}

func (b *WorkflowBuilder) Start(stateName string) *StartBuidler[WorkflowBuilder] {
	b.startBuilder = newStartBuilder[WorkflowBuilder](b, stateName)
	return b.startBuilder
}

func (b *WorkflowBuilder) Event(name string) *EventBuilder[WorkflowBuilder] {
	builder := newEventBuilder[WorkflowBuilder](b, name)
	b.events = append(b.events, builder)
	return builder
}

func (b *WorkflowBuilder) Function(name, operation string) *FunctionBuilder[WorkflowBuilder] {
	builder := newFunctionBuilder[WorkflowBuilder](b, name, operation)
	b.functions = append(b.functions, builder)
	fmt.Println(len(b.functions))
	fmt.Println(name)
	return builder
}

func (b *WorkflowBuilder) OperationState(name string) *OperationStateBuilder[WorkflowBuilder] {
	builder := newOperationStateBuilder[WorkflowBuilder](b, name)
	b.operatorStateBuilders = append(b.operatorStateBuilders, builder)
	return builder
}

func (b *WorkflowBuilder) EventState(name string) *EventStateBuilder[WorkflowBuilder] {
	builder := newEventStateBuilder(b, name)
	b.eventStateBuilders = append(b.eventStateBuilders, builder)
	return builder
}

func (b *WorkflowBuilder) Build() model.Workflow {
	workflow := b.workflow
	if b.startBuilder != nil {
		workflow.Start = &b.startBuilder.start
	}

	workflow.Functions = make(model.Functions, len(b.functions))
	for i, function := range b.functions {
		workflow.Functions[i] = function.Build()
	}

	workflow.Events = make(model.Events, len(b.events))
	for i, event := range b.events {
		workflow.Events[i] = event.Build()
	}

	workflow.States = make(model.States, len(b.operatorStateBuilders))
	for i, builder := range b.operatorStateBuilders {
		workflow.States[i] = builder.Build()
	}
	return workflow
}

// NewStartBuilder create a Start Builder
func NewStartBuilder(stateName string) (*StartBuidler[any], *Builder[StartBuidler[any], model.Start]) {
	startBuilder := newStartBuilder[any](nil, stateName)
	builder := &Builder[StartBuidler[any], model.Start]{
		builder: startBuilder,
	}
	return startBuilder, builder
}

func newStartBuilder[T any](parent *T, stateName string) *StartBuidler[T] {
	start := model.Start{
		StateName: stateName,
	}
	return &StartBuidler[T]{
		parent: parent,
		start:  start,
	}
}

type StartBuidler[T any] struct {
	parent         *T
	schedulerBuild *ScheduleBuilder[StartBuidler[T]]
	start          model.Start
}

func (b *StartBuidler[T]) Schedule() *ScheduleBuilder[StartBuidler[T]] {
	b.schedulerBuild = newScheduleBuilder(b)
	return b.schedulerBuild
}

func (b *StartBuidler[T]) Parent() *T {
	return b.parent
}

func (b *StartBuidler[T]) Build() model.Start {
	b.start.Schedule = &b.schedulerBuild.schedule
	return b.start
}

// NewScheduleBuilder create a Schedule Builder
func NewScheduleBuilder() *ScheduleBuilder[any] {
	return newScheduleBuilder[any](nil)
}
func newScheduleBuilder[T any](parent *T) *ScheduleBuilder[T] {
	schedule := model.Schedule{}
	return &ScheduleBuilder[T]{
		parent:   parent,
		schedule: schedule,
	}
}

type ScheduleBuilder[T any] struct {
	parent   *T
	schedule model.Schedule
}

func (b *ScheduleBuilder[T]) Parent() *T {
	return b.parent
}

func (b *ScheduleBuilder[T]) Build() model.Schedule {
	return b.schedule
}

// NewEndBuilder create a End Builder
func NewEndBuilder(terminate bool) *EndBuilder[any] {
	return newEndBuilder[any](nil, terminate)
}

func newEndBuilder[T any](parent *T, terminate bool) *EndBuilder[T] {
	return &EndBuilder[T]{
		parent: parent,
		end: model.End{
			Terminate: terminate,
		},
	}
}

type EndBuilder[T any] struct {
	parent *T
	end    model.End
}

func (b *EndBuilder[T]) Parent() *T {
	return b.parent
}

func (b *EndBuilder[T]) Build() model.End {
	return b.end
}
