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

package builder

import "github.com/serverlessworkflow/sdk-go/v2/model"

func NewEventBuilder(name string) *EventBuilder[any] {
	return newEventBuilder[any](nil, name)
}
func newEventBuilder[T any](parent *T, name string) *EventBuilder[T] {
	event := model.Event{
		Name: name,
	}
	event.ApplyDefault()

	return &EventBuilder[T]{
		parent: parent,
		Event:  event,
	}
}

type EventBuilder[T any] struct {
	parent *T
	Event  model.Event
}

func (b *EventBuilder[T]) Parent() *T {
	return b.parent
}

func (b *EventBuilder[T]) Build() model.Event {
	return b.Event
}

func NewEventRefBuilder(triggerEvent, resultEvent string) *EventRefBuilder {
	eventRef := model.EventRef{
		TriggerEventRef: triggerEvent,
		ResultEventRef:  resultEvent,
	}
	eventRef.ApplyDefault()

	return &EventRefBuilder{
		eventRef: eventRef,
	}
}

type EventRefBuilder struct {
	eventRef model.EventRef
}

func (b *EventRefBuilder) Build() model.EventRef {
	return b.eventRef
}
