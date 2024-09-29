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

type DurationBuilder struct {
	root *graph.Node
}

func (b *DurationBuilder) SetSeconds(seconds int) *DurationBuilder {
	b.root.Edge("seconds").SetInt(seconds)
	return b
}

func (b *DurationBuilder) GetSeconds() int {
	return b.root.Edge("seconds").GetInt()
}

func NewDurationBuilder(root *graph.Node) *DurationBuilder {
	return &DurationBuilder{
		root: root,
	}
}
