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

func NewFunctionBuilder(name, operation string) *FunctionBuilder[any] {
	return newFunctionBuilder[any](nil, name, operation)
}

func newFunctionBuilder[T any](parent *T, name, operation string) *FunctionBuilder[T] {
	function := model.Function{
		Name:      name,
		Operation: operation,
	}
	function.ApplyDefault()

	return &FunctionBuilder[T]{
		parent:   parent,
		Function: function,
	}
}

type FunctionBuilder[T any] struct {
	parent   *T
	Function model.Function
}

func (b *FunctionBuilder[T]) Parent() *T {
	return b.parent
}

func (b *FunctionBuilder[T]) Build() model.Function {
	return b.Function
}
