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

import "github.com/serverlessworkflow/sdk-go/v2/util"

const (
	// FunctionTypeREST a combination of the function/service OpenAPI definition document URI and the particular service
	// operation that needs to be invoked, separated by a '#'.
	FunctionTypeREST FunctionType = "rest"
	// FunctionTypeRPC a combination of the gRPC proto document URI and the particular service name and service method
	// name that needs to be invoked, separated by a '#'.
	FunctionTypeRPC FunctionType = "rpc"
	// FunctionTypeExpression defines the expression syntax.
	FunctionTypeExpression FunctionType = "expression"
	// FunctionTypeGraphQL a combination of the GraphQL schema definition URI and the particular service name and
	// service method name that needs to be invoked, separated by a '#'
	FunctionTypeGraphQL FunctionType = "graphql"
	// FunctionTypeAsyncAPI a combination of the AsyncApi definition document URI and the particular service operation
	// that needs to be invoked, separated by a '#'
	FunctionTypeAsyncAPI FunctionType = "asyncapi"
	// FunctionTypeOData a combination of the GraphQL schema definition URI and the particular service name and service
	// method name that needs to be invoked, separated by a '#'
	FunctionTypeOData FunctionType = "odata"
	// FunctionTypeCustom property defines a list of function types that are set by the specification. Some runtime
	// implementations might support additional function types that extend the ones defined in the specification
	FunctionTypeCustom FunctionType = "custom"
)

// FunctionType ...
type FunctionType string

func (i FunctionType) KindValues() []string {
	return []string{
		string(FunctionTypeREST),
		string(FunctionTypeRPC),
		string(FunctionTypeExpression),
		string(FunctionTypeGraphQL),
		string(FunctionTypeAsyncAPI),
		string(FunctionTypeOData),
		string(FunctionTypeCustom),
	}
}

func (i FunctionType) String() string {
	return string(i)
}

// Function ...
type Function struct {
	Common `json:",inline"`
	// Unique function name
	// +kubebuilder:validation:Required
	Name string `json:"name" validate:"required"`
	// If type is `rest`, <path_to_openapi_definition>#<operation_id>.
	// If type is `rpc`, <path_to_grpc_proto_file>#<service_name>#<service_method>.
	// If type is `expression`, defines the workflow expression. If the type is `custom`,
	// <path_to_custom_script>#<custom_service_method>.
	// +kubebuilder:validation:Required
	Operation string `json:"operation" validate:"required"`
	// Defines the function type. Is either `custom`, `rest`, `rpc`, `expression`, `graphql`, `odata` or `asyncapi`.
	// Default is `rest`.
	// +kubebuilder:validation:Enum=rest;rpc;expression;graphql;odata;asyncapi;custom
	// +kubebuilder:default=rest
	Type FunctionType `json:"type,omitempty" validate:"required,oneofkind"`
	// References an auth definition name to be used to access to resource defined in the operation parameter.
	// +optional
	AuthRef string `json:"authRef,omitempty"`
}

type functionUnmarshal Function

// UnmarshalJSON implements json unmarshaler interface
func (f *Function) UnmarshalJSON(data []byte) error {
	f.ApplyDefault()
	return util.UnmarshalObject("function", data, (*functionUnmarshal)(f))
}

// ApplyDefault set the default values for Function
func (f *Function) ApplyDefault() {
	f.Type = FunctionTypeREST
}
