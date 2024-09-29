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
	"github.com/serverlessworkflow/sdk-go/v4/internal/dsl"
)

type DocumentBuilder struct {
	root *graph.Node
}

func (b *DocumentBuilder) SetDSL(dsl string) *DocumentBuilder {
	node := b.root.Edge("dsl")
	node.SetString(dsl)
	return b
}

func (b *DocumentBuilder) GetDSL() string {
	node := b.root.Edge("dsl")
	return node.GetString()
}

func (b *DocumentBuilder) SetNamespace(dsl string) *DocumentBuilder {
	node := b.root.Edge("namespace")
	node.SetString(dsl)
	return b
}

func (b *DocumentBuilder) GetNamespace() string {
	node := b.root.Edge("namespace")
	return node.GetString()
}

func (b *DocumentBuilder) SetName(dsl string) *DocumentBuilder {
	node := b.root.Edge("name")
	node.SetString(dsl)
	return b
}

func (b *DocumentBuilder) GetName() string {
	node := b.root.Edge("name")
	return node.GetString()
}

func (b *DocumentBuilder) SetVersion(dsl string) *DocumentBuilder {
	node := b.root.Edge("version")
	node.SetString(dsl)
	return b
}

func (b *DocumentBuilder) GetVersion() string {
	node := b.root.Edge("version")
	return node.GetString()
}

func NewDocumentBuilder(root *graph.Node) *DocumentBuilder {
	documentBuilder := &DocumentBuilder{
		root: root,
	}
	documentBuilder.SetDSL(dsl.DSLVersion)
	return documentBuilder
}
