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

import "github.com/serverlessworkflow/sdk-go/v4/graph"

type WaitBuilder struct {
	root     *graph.Node
	duration *DurationBuilder
}

func (b *WaitBuilder) SetWait(wait string) {
	b.root.Edge("wait").Clear().SetString(string(wait))
}

func (b *WaitBuilder) GetWait() string {
	return b.root.Edge("wait").GetString()
}

func (b *WaitBuilder) Duration() *DurationBuilder {
	if b.duration == nil {
		node := b.root.Edge("wait").Clear()
		b.duration = NewDurationBuilder(node)
	}
	return b.duration
}

func NewWaitBuilder(root *graph.Node) *WaitBuilder {
	return &WaitBuilder{
		root: root,
	}
}
