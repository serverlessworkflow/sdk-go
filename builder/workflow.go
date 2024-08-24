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
	"github.com/serverlessworkflow/sdk-go/v4/graph"
	"github.com/serverlessworkflow/sdk-go/v4/internal/load"
)

type WorkflowBuilder struct {
	root     *graph.Node
	document *DocumentBuilder
	do       *DoBuilder
	use      *UseBuilder
}

func (b *WorkflowBuilder) Document() *DocumentBuilder {
	if b.document == nil {
		b.document = NewDocumentBuilder(b.root.Edge("document"))
	}
	return b.document
}

func (b *WorkflowBuilder) Do() *DoBuilder {
	if b.do == nil {
		b.do = NewDoBuilder(b.root.Edge("do"))
	}
	return b.do
}

func (b *WorkflowBuilder) Use() *UseBuilder {
	if b.use == nil {
		b.use = NewUseBuilder(b.root.Edge("use"))
	}
	return b.use
}

func (b *WorkflowBuilder) Node() *graph.Node {
	return b.root
}

func NewWorkflowBuilder() *WorkflowBuilder {
	root := graph.NewNode()
	return &WorkflowBuilder{
		root: root,
	}
}

func NewWorkflowBuilderFromFile(path string) (*WorkflowBuilder, error) {
	root, _, err := load.FromFile(path)
	if err != nil {
		return nil, err
	}

	return &WorkflowBuilder{
		root: root,
	}, nil
}

func NewWorkflowBuilderFromYAMLSource(source []byte) (*WorkflowBuilder, error) {
	root, _, err := load.FromYAMLSource(source)
	if err != nil {
		return nil, err
	}

	return &WorkflowBuilder{
		root: root,
	}, nil
}

func NewWorkflowBuilderFromJSONSource(source []byte) (*WorkflowBuilder, error) {
	root, _, err := load.FromJSONSource(source)
	if err != nil {
		return nil, err
	}

	return &WorkflowBuilder{
		root: root,
	}, nil
}
