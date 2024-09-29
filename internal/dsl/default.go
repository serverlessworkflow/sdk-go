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

package dsl

import "github.com/serverlessworkflow/sdk-go/v4/graph"

func ApplyDefault(node *graph.Node) error {
	lookup := node.Lookup("do.*.*.call=http")

	for _, node := range lookup.List() {
		lookupEdged := node.Lookup("with.content")
		for _, nodeEdge := range lookupEdged.List() {
			if !nodeEdge.HasValue() {
				nodeEdge.SetString("content")
			}
		}
	}

	lookup = node.Lookup("do.*.*.then")
	for _, node := range lookup.List() {
		if !node.HasValue() {
			node.SetString("continue")
		}
	}

	lookup = node.Lookup("do.*.*.fork")
	for _, node := range lookup.List() {
		if !node.Edge("compete").HasValue() {
			node.SetBool(false)
		}
	}

	lookup = node.Lookup("do.*.*.run.workflow")
	for _, node := range lookup.List() {
		if !node.Edge("version").HasValue() {
			node.Edge("version").SetString("latest")
		}
	}

	lookup = node.Lookup("do.*.*.catch")
	for _, node := range lookup.List() {
		if !node.Edge("catch").HasValue() {
			node.Edge("catch").SetString("error")
		}
	}

	lookup = node.Lookup("evaluate.language")
	if !lookup.Empty() {
		node.Edge("evaluate").Edge("language").SetString("jq")
	}

	lookup = node.Lookup("evaluate.mode")
	if !lookup.Empty() {
		node.Edge("evaluate").Edge("mode").SetString("strict")
	}

	return nil
}
