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
	"bytes"
	"encoding/json"
	"log"
	"strings"
)

type Lookup struct {
	nodes []*Node
}

func (l Lookup) Empty() bool {
	return len(l.nodes) == 0
}

func (l Lookup) First() *Node {
	return l.nodes[0]
}

func (l Lookup) Get(index int) *Node {
	return l.nodes[index]
}

func (l Lookup) List() []*Node {
	return l.nodes
}

type Node struct {
	value  interface{}
	order  []string
	parent *Node
	edges  map[string]*Node
	list   bool
}

func (n *Node) List(list bool) {
	n.list = list
}

func (n *Node) IsList() bool {
	return n.list
}

func (n *Node) UnmarshalJSON(data []byte) error {
	return unmarshalNode(n, data)
}

func (n *Node) MarshalJSON() ([]byte, error) {
	return marshalNode(n)
}

func (n *Node) Edge(name string) *Node {
	if n.HasValue() {
		log.Panic("value already defined, execute clear first")
	}
	if _, ok := n.edges[name]; !ok {
		newNode := NewNode()
		newNode.parent = n
		n.edges[name] = newNode
		n.order = append(n.order, name)
	}
	return n.edges[name]
}

func (n *Node) SetString(value string) *Node {
	n.setValue(value)
	return n
}

func (n *Node) SetInt(value int) *Node {
	n.setValue(value)
	return n
}

func (n *Node) SetFloat(value float32) *Node {
	n.setValue(value)
	return n
}

func (n *Node) SetBool(value bool) *Node {
	n.setValue(value)
	return n
}

func (n *Node) setValue(value any) {
	if len(n.edges) > 0 {
		log.Panic("already defined edges, execute clear fist")
	}
	n.value = value
}

func (n *Node) GetString() string {
	return n.value.(string)
}

func (n *Node) GetInt() int {
	return n.value.(int)
}

func (n *Node) GetFloat() float32 {
	return n.value.(float32)
}

func (n *Node) HasValue() bool {
	return n.value != nil
}

func (n *Node) Clear() *Node {
	n.value = nil
	n.edges = map[string]*Node{}
	n.order = []string{}
	return n
}

func (n *Node) Parent() *Node {
	return n.parent
}

func (n *Node) Index(i int) (string, *Node) {
	lookup := n.Lookup(n.order[i])
	if !lookup.Empty() {
		return n.order[i], lookup.First()
	}
	return "", nil
}

func (n *Node) Lookup(path string) Lookup {
	dotIndex := strings.Index(path, ".")
	var key string
	if dotIndex == -1 {
		key = strings.TrimSpace(path)
	} else {
		key = strings.TrimSpace(path[0:dotIndex])
		path = path[dotIndex+1:]
	}

	var currentNode *Node
	if key == "*" {
		nodes := []*Node{}
		if dotIndex == -1 {
			for _, node := range n.edges {
				nodes = append(nodes, node)
			}
			return Lookup{nodes}
		}
		for _, node := range n.edges {
			if nodesLookup := node.Lookup(path); !nodesLookup.Empty() {
				nodes = append(nodes, nodesLookup.List()...)
			}
		}
		return Lookup{nodes}

	}

	equalIndex := strings.Index(key, "=")
	if equalIndex != -1 {
		value := key[equalIndex+1:]
		key := key[:equalIndex]

		lookup := n.Lookup(key)
		if !lookup.Empty() && lookup.First().value != value {
			return Lookup{}
		}

		return Lookup{[]*Node{n}}
	}

	currentNode = n.edges[key]
	if currentNode == nil {
		return Lookup{}
	}
	if dotIndex == -1 {
		return Lookup{[]*Node{currentNode}}
	}

	return currentNode.Lookup(path)
}

func NewNode() *Node {
	return (&Node{}).Clear()
}

func UnmarshalJSON(data []byte) (*Node, error) {
	node := NewNode()
	err := json.Unmarshal(data, &node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func MarshalJSON(n *Node) ([]byte, error) {
	data, err := json.Marshal(n)
	if err != nil {
		return nil, err
	}

	var out bytes.Buffer
	err = json.Indent(&out, data, "", "  ")
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}

func LoadExternalResource(n *Node) error {
	return nil
}
