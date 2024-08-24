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

package graph

import (
	"encoding/json"
)

func marshalNode(n *Node) ([]byte, error) {
	if n.value != nil {
		return json.Marshal(n.value)
	}

	var out []byte
	if n.list {
		out = append(out, '[')
	} else {
		out = append(out, '{')
	}

	nEdge := len(n.order) - 1
	for i, edge := range n.order {
		node := n.edges[edge]
		val, err := json.Marshal(node)
		if err != nil {
			return nil, err
		}

		if n.list {
			out = append(out, val...)
		} else {
			out = append(out, []byte("\""+edge+"\":")...)
			out = append(out, val...)
		}

		if nEdge != i {
			out = append(out, byte(','))
		}
	}

	if n.list {
		out = append(out, ']')
	} else {
		out = append(out, '}')
	}

	return out, nil
}
