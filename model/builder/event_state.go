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

import "github.com/serverlessworkflow/sdk-go/v2/model"

func NewEventStateBuilder(name string) *EventStateBuilder[any] {
	return newEventStateBuilder[any](nil, name)
}

func newEventStateBuilder[T any](parent *T, name string) *EventStateBuilder[T] {
	eventState := model.EventState{}
	eventState.ApplyDefault()

	state := model.State{
		BaseState: model.BaseState{
			Name: name,
			Type: model.StateTypeEvent,
		},
		EventState: &eventState,
	}

	return &EventStateBuilder[T]{
		parent: parent,

		State:    state,
		onEvents: make([]*OnEventsBuilder[EventStateBuilder[T]], 0),
	}
}

type EventStateBuilder[T any] struct {
	parent     *T
	State      model.State
	onEvents   []*OnEventsBuilder[EventStateBuilder[T]]
	endBuilder *EndBuilder[EventStateBuilder[T]]
}

func (b *EventStateBuilder[T]) OnEvent(eventRef string) *OnEventsBuilder[EventStateBuilder[T]] {
	onEvents := newOnEventsBuilder[EventStateBuilder[T]](b, eventRef)
	b.onEvents = append(b.onEvents, onEvents)
	return onEvents
}

func (b *EventStateBuilder[T]) End(terminate bool) *EndBuilder[EventStateBuilder[T]] {
	b.endBuilder = newEndBuilder[EventStateBuilder[T]](b, terminate)
	return b.endBuilder
}

func (b *EventStateBuilder[T]) Parent() *T {
	return b.parent
}

func (b *EventStateBuilder[T]) Build() model.State {
	b.State.OnEvents = make([]model.OnEvents, len(b.onEvents))
	for i, onEvents := range b.onEvents {
		b.State.OnEvents[i] = onEvents.Build()
	}
	return b.State
}

func newOnEventsBuilder[T any](parent *T, eventRef string) *OnEventsBuilder[T] {
	onEvents := model.OnEvents{
		EventRefs: []string{
			eventRef,
		},
	}
	onEvents.ApplyDefault()

	return &OnEventsBuilder[T]{
		parent:   parent,
		onEvents: onEvents,
	}
}

type OnEventsBuilder[T any] struct {
	parent        *T
	actionBuilder []*ActionBuilder[OnEventsBuilder[T]]
	onEvents      model.OnEvents
}

func (b *OnEventsBuilder[T]) Action(name string) *ActionBuilder[OnEventsBuilder[T]] {
	builder := newActionBuilder(b, name)
	b.actionBuilder = append(b.actionBuilder, builder)
	return builder
}

func (b *OnEventsBuilder[T]) Parent() *T {
	return b.parent
}

func (b *OnEventsBuilder[T]) Build() model.OnEvents {
	b.onEvents.Actions = make([]model.Action, len(b.actionBuilder))
	for i, action := range b.actionBuilder {
		b.onEvents.Actions[i] = action.Build()
	}
	return b.onEvents
}
