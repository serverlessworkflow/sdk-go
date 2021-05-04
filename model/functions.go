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

package model

const (
	FunctionTypeREST       FunctionType = "rest"
	FunctionTypeRPC        FunctionType = "rpc"
	FunctionTypeExpression FunctionType = "expression"
)

type FunctionType string

type Function struct {
	Common
	// Unique function name
	Name string `json:"name"`
	// If type is `rest`, <path_to_openapi_definition>#<operation_id>. If type is `rpc`, <path_to_grpc_proto_file>#<service_name>#<service_method>. If type is `expression`, defines the workflow expression.
	Operation string `json:"operation"`
	// Defines the function type. Is either `rest`, `rpc` or `expression`. Default is `rest`
	Type FunctionType `json:"type,omitempty"`
}

// Function Reference
type FunctionRef struct {
	// Name of the referenced function
	RefName string `json:"refName"`
	// Function arguments
	Arguments map[string]interface{} `json:"arguments,omitempty"`
}
