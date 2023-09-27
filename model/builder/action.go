// Copyright 2022 The Serverless Workflow Specification Authors
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
	"github.com/serverlessworkflow/sdk-go/v2/model"
)

func NewActionBuilder(name string) *ActionBuilder[any] {
	return newActionBuilder[any](nil, name)
}

func newActionBuilder[T any](parent *T, name string) *ActionBuilder[T] {
	action := model.Action{
		Name: name,
	}
	action.ApplyDefault()

	builder := &ActionBuilder[T]{
		action: action,
		parent: parent,
	}
	return builder
}

type ActionBuilder[T any] struct {
	parent             *T
	functionRefBuilder *FunctionRefBuilder[ActionBuilder[T]]
	eventRefBuilder    *EventRefBuilder
	sleepBuilder       *SleepBuilder
	action             model.Action
}

func (b *ActionBuilder[T]) FunctionRef(refName string) *FunctionRefBuilder[ActionBuilder[T]] {
	b.functionRefBuilder = newFunctionRefBuilder[ActionBuilder[T]](b, refName)
	return b.functionRefBuilder
}

func (b *ActionBuilder[T]) EventRef(triggerEvent, resultEvent string) *EventRefBuilder {
	b.eventRefBuilder = NewEventRefBuilder(triggerEvent, resultEvent)
	return b.eventRefBuilder
}

func (b *ActionBuilder[T]) Sleep() *SleepBuilder {
	b.sleepBuilder = NewSleepBuilder()
	return b.sleepBuilder
}

func (b *ActionBuilder[T]) Parent() *T {
	return b.parent
}

func (b *ActionBuilder[T]) Build() model.Action {
	if b.functionRefBuilder != nil {
		f := b.functionRefBuilder.Build()
		b.action.FunctionRef = &f
	}

	if b.eventRefBuilder != nil {
		e := b.eventRefBuilder.Build()
		b.action.EventRef = &e
	}

	if b.sleepBuilder != nil {
		s := b.sleepBuilder.Build()
		b.action.Sleep = &s
	}

	return b.action
}

func newFunctionRefBuilder[T any](parent *T, refName string) *FunctionRefBuilder[T] {
	functionRef := model.FunctionRef{
		RefName: refName,
	}
	functionRef.ApplyDefault()

	return &FunctionRefBuilder[T]{
		parent:      parent,
		functionRef: functionRef,
	}
}

type FunctionRefBuilder[T any] struct {
	parent      *T
	functionRef model.FunctionRef
}

func (b *FunctionRefBuilder[T]) Sleep() *SleepBuilder {
	return NewSleepBuilder()
}

func (b *FunctionRefBuilder[T]) Build() model.FunctionRef {
	return b.functionRef
}

func NewSleepBuilder() *SleepBuilder {
	sleep := model.Sleep{}
	return &SleepBuilder{
		sleep: sleep,
	}
}

type SleepBuilder struct {
	sleep model.Sleep
}

func (b *SleepBuilder) Before(before string) *SleepBuilder {
	b.sleep.Before = before
	return b
}

func (b *SleepBuilder) After(after string) *SleepBuilder {
	b.sleep.After = after
	return b
}

func (b *SleepBuilder) Build() model.Sleep {
	return b.sleep
}
