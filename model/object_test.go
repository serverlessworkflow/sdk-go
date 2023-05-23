// Copyright 2023 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package model

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestObjectToString(t *testing.T) {

	bytez := []byte("hello")
	object := Object{
		Type:     1,
		IntVal:   13,
		StrVal:   "26",
		RawValue: bytez,
	}
	value := object.String()
	assert.NotNil(t, value)
	assert.Equal(t, "{ Type:1, IntVal:13, StrVal:26, RawValue:[104 101 108 108 111] }", value)
}
