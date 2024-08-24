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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGraph(t *testing.T) {
	source := []byte(`{
  "test": "val",
  "test2": "val2",
  "list": [
    "test1"
  ],
  "listObject": [
    {
      "test": "val"
    },
	{
      "test": "val1"
    },
    {
      "test2": "va2"
    }
  ],
  "deep": [
    [
      {
        "test": [
	      {
            "test2": "val",
            "test3": "val1",
              "test4": {
              "test5": "val"
			}
          },
          {
            "test2": "val"
          }
        ]
      }
    ]
  ]
}`)

	root, err := UnmarshalJSON(source)
	if !assert.NoError(t, err) {
		return
	}

	t.Run("marshal", func(t *testing.T) {
		_, err := MarshalJSON(root)
		assert.NoError(t, err)
	})

	t.Run("lookup key", func(t *testing.T) {
		lookup := root.Lookup("test")
		assert.NoError(t, err)
		assert.Equal(t, "val", lookup.First().value)
	})

	t.Run("lookup not found", func(t *testing.T) {
		lookup := root.Lookup("list2")
		assert.Equal(t, 0, len(lookup.List()))
	})

	t.Run("lookup list", func(t *testing.T) {
		lookup := root.Lookup("list")
		assert.Nil(t, lookup.First().value)
		assert.Equal(t, 1, len(lookup.First().edges))

		lookup = root.Lookup("listObject.*")
		assert.Nil(t, lookup.Get(0).value)
		assert.Equal(t, 3, len(lookup.List()))
	})

	t.Run("lookup list index", func(t *testing.T) {
		lookup := root.Lookup("list.0")
		assert.Equal(t, "test1", lookup.First().value)
	})

	t.Run("lookup search in a list", func(t *testing.T) {
		lookup := root.Lookup("listObject.*.test")
		assert.Equal(t, 2, len(lookup.List()))
		assert.Equal(t, "val", lookup.Get(0).value)
		assert.Equal(t, "val1", lookup.Get(1).value)
	})

	t.Run("lookup deep", func(t *testing.T) {
		lookup := root.Lookup("deep.*.*.test.*.test2")
		assert.Equal(t, 2, len(lookup.List()))
		assert.Equal(t, "val", lookup.Get(0).value)
		assert.Equal(t, "val", lookup.Get(1).value)

		lookup = root.Lookup("deep.*.*.test.*.test2=val")
		assert.Equal(t, 2, len(lookup.List()))
		assert.Equal(t, 3, len(lookup.Get(0).edges))
		assert.Equal(t, 1, len(lookup.Get(1).edges))
	})

	t.Run("lookup with value equality", func(t *testing.T) {
		lookup := root.Lookup("deep.*.*.test.*.test2=val")
		assert.Equal(t, 2, len(lookup.List()))

		assert.Equal(t, "val", lookup.Get(0).Lookup("test2").First().value)
		assert.Equal(t, "val1", lookup.Get(0).Lookup("test3").First().GetString())

		assert.Equal(t, false, lookup.Get(1).Lookup("test2").Empty())
		assert.Equal(t, true, lookup.Get(1).Lookup("test3").Empty())
		assert.Equal(t, "val", lookup.Get(1).Lookup("test2").First().GetString())
	})
}
