// Copyright 2024 The Serverless Workflow Specification Authors
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

	"github.com/serverlessworkflow/sdk-go/v4/graph"
)

type DoBuilder struct {
	root  *graph.Node
	tasks []any
}

func (b *DoBuilder) AddCall(name string) (*CallBuilder, int) {
	index := len(b.tasks)
	nodeIndex := b.root.Edge(fmt.Sprintf("%d", index))
	nodeName := nodeIndex.Edge(name)

	callBuilder := NewCallBuilder(nodeName)
	b.tasks = append(b.tasks, callBuilder)
	return callBuilder, index
}

func (b *DoBuilder) AddWait(name string) (*WaitBuilder, int) {
	index := len(b.tasks)
	nodeIndex := b.root.Edge(fmt.Sprintf("%d", index))
	nodeName := nodeIndex.Edge(name)

	waitBuilder := NewWaitBuilder(nodeName)
	b.tasks = append(b.tasks, waitBuilder)
	return waitBuilder, index
}

func (b *DoBuilder) RemoveTask(index int) *DoBuilder {
	b.tasks = append(b.tasks[:index], b.tasks[index+1:]...)
	return b
}

func NewDoBuilder(root *graph.Node) *DoBuilder {
	root.List(true)
	return &DoBuilder{
		root:  root,
		tasks: []any{},
	}
}
