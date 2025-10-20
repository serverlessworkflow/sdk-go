// Copyright 2025 The Serverless Workflow Specification Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCallHTTP_MarshalJSON(t *testing.T) {
	callHTTP := CallHTTP{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: string(FlowDirectiveContinue)},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Call: "http",
		With: HTTPArguments{
			Method: "GET",
			Endpoint: &Endpoint{
				URITemplate: &LiteralUri{Value: "http://example.com"},
			},
			Headers: map[string]string{
				"Authorization": "Bearer token",
			},
			Query: map[string]interface{}{
				"q": "search",
			},
			Output:   "content",
			Redirect: true,
		},
	}

	data, err := json.Marshal(callHTTP)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
			"if": "${condition}",
			"input": { "from": {"key": "value"} },
			"output": { "as": {"result": "output"} },
			"timeout": { "after": "10s" },
			"then": "continue",
			"metadata": {"meta": "data"},
			"call": "http",
			"with": {
				"method": "GET",
				"endpoint": "http://example.com",
				"headers": {"Authorization": "Bearer token"},
				"query": {"q": "search"},
				"output": "content",
				"redirect": true
			}
		}`, string(data))
}

func TestCallHTTP_UnmarshalJSON(t *testing.T) {
	jsonData := `{
			"if": "${condition}",
			"input": { "from": {"key": "value"} },
			"output": { "as": {"result": "output"} },
			"timeout": { "after": "10s" },
			"then": "continue",
			"metadata": {"meta": "data"},
			"call": "http",
			"with": {
				"method": "GET",
				"endpoint": "http://example.com",
				"headers": {"Authorization": "Bearer token"},
				"query": {"q": "search"},
				"output": "content",
				"redirect": true
			}
		}`

	var callHTTP CallHTTP
	err := json.Unmarshal([]byte(jsonData), &callHTTP)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{"${condition}"}, callHTTP.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, callHTTP.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, callHTTP.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, callHTTP.Timeout)
	assert.Equal(t, &FlowDirective{Value: string(FlowDirectiveContinue)}, callHTTP.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, callHTTP.Metadata)
	assert.Equal(t, "http", callHTTP.Call)
	assert.Equal(t, "GET", callHTTP.With.Method)
	assert.Equal(t, "http://example.com", callHTTP.With.Endpoint.String())
	assert.Equal(t, map[string]string{"Authorization": "Bearer token"}, callHTTP.With.Headers)
	assert.Equal(t, map[string]interface{}{"q": "search"}, callHTTP.With.Query)
	assert.Equal(t, "content", callHTTP.With.Output)
	assert.Equal(t, true, callHTTP.With.Redirect)
}

func TestCallOpenAPI_MarshalJSON(t *testing.T) {
	authPolicy := "my-auth"
	callOpenAPI := CallOpenAPI{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Call: "openapi",
		With: OpenAPIArguments{
			Document: &ExternalResource{
				Name: "MyOpenAPIDoc",
				Endpoint: &Endpoint{
					URITemplate: &LiteralUri{Value: "http://example.com/openapi.json"},
				},
			},
			OperationID: "getUsers",
			Parameters: map[string]interface{}{
				"param1": "value1",
				"param2": "value2",
			},
			Authentication: &ReferenceableAuthenticationPolicy{
				Use: &authPolicy,
			},
			Output:   "content",
			Redirect: true,
		},
	}

	data, err := json.Marshal(callOpenAPI)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "openapi",
		"with": {
			"document": {
				"name": "MyOpenAPIDoc",
				"endpoint": "http://example.com/openapi.json"
			},
			"operationId": "getUsers",
			"parameters": {
				"param1": "value1",
				"param2": "value2"
			},
			"authentication": {
				"use": "my-auth"
			},
			"output": "content",
			"redirect": true
		}
	}`, string(data))
}

func TestCallOpenAPI_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "openapi",
		"with": {
			"document": {
				"name": "MyOpenAPIDoc",
				"endpoint": { "uri": "http://example.com/openapi.json" }
			},
			"operationId": "getUsers",
			"parameters": {
				"param1": "value1",
				"param2": "value2"
			},
			"authentication": {
				"use": "my-auth"
			},
			"output": "content",
			"redirect": true
		}
	}`

	var callOpenAPI CallOpenAPI
	err := json.Unmarshal([]byte(jsonData), &callOpenAPI)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, callOpenAPI.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, callOpenAPI.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, callOpenAPI.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, callOpenAPI.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, callOpenAPI.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, callOpenAPI.Metadata)
	assert.Equal(t, "openapi", callOpenAPI.Call)
	assert.Equal(t, "MyOpenAPIDoc", callOpenAPI.With.Document.Name)
	assert.Equal(t, "http://example.com/openapi.json", callOpenAPI.With.Document.Endpoint.EndpointConfig.URI.String())
	assert.Equal(t, "getUsers", callOpenAPI.With.OperationID)
	assert.Equal(t, map[string]interface{}{"param1": "value1", "param2": "value2"}, callOpenAPI.With.Parameters)
	assert.Equal(t, "my-auth", *callOpenAPI.With.Authentication.Use)
	assert.Equal(t, "content", callOpenAPI.With.Output)
	assert.Equal(t, true, callOpenAPI.With.Redirect)
}

func TestCallGRPC_MarshalJSON(t *testing.T) {
	callGRPC := CallGRPC{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Call: "grpc",
		With: GRPCArguments{
			Proto: &ExternalResource{
				Name: "MyProtoFile",
				Endpoint: &Endpoint{
					URITemplate: &LiteralUri{Value: "http://example.com/protofile"},
				},
			},
			Service: GRPCService{
				Name: "UserService",
				Host: "example.com",
				Port: 50051,
			},
			Method:    "GetUser",
			Arguments: map[string]interface{}{"userId": "12345"},
		},
	}

	data, err := json.Marshal(callGRPC)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "grpc",
		"with": {
			"proto": {
				"name": "MyProtoFile",
				"endpoint": "http://example.com/protofile"
			},
			"service": {
				"name": "UserService",
				"host": "example.com",
				"port": 50051
			},
			"method": "GetUser",
			"arguments": {
				"userId": "12345"
			}
		}
	}`, string(data))
}

func TestCallGRPC_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "grpc",
		"with": {
			"proto": {
				"name": "MyProtoFile",
				"endpoint": "http://example.com/protofile"
			},
			"service": {
				"name": "UserService",
				"host": "example.com",
				"port": 50051
			},
			"method": "GetUser",
			"arguments": {
				"userId": "12345"
			}
		}
	}`

	var callGRPC CallGRPC
	err := json.Unmarshal([]byte(jsonData), &callGRPC)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, callGRPC.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, callGRPC.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, callGRPC.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, callGRPC.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, callGRPC.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, callGRPC.Metadata)
	assert.Equal(t, "grpc", callGRPC.Call)
	assert.Equal(t, "MyProtoFile", callGRPC.With.Proto.Name)
	assert.Equal(t, "http://example.com/protofile", callGRPC.With.Proto.Endpoint.String())
	assert.Equal(t, "UserService", callGRPC.With.Service.Name)
	assert.Equal(t, "example.com", callGRPC.With.Service.Host)
	assert.Equal(t, 50051, callGRPC.With.Service.Port)
	assert.Equal(t, "GetUser", callGRPC.With.Method)
	assert.Equal(t, map[string]interface{}{"userId": "12345"}, callGRPC.With.Arguments)
}

func TestCallAsyncAPI_MarshalJSON(t *testing.T) {
	callAsyncAPI := CallAsyncAPI{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Call: "asyncapi",
		With: AsyncAPIArguments{
			Document: &ExternalResource{
				Name: "MyAsyncAPIDoc",
				Endpoint: &Endpoint{
					URITemplate: &LiteralUri{Value: "http://example.com/asyncapi.json"},
				},
			},
			Operation: "user.signup",
			Server:    &AsyncAPIServer{Name: "default-server"},
			Message:   &AsyncAPIOutboundMessage{Payload: map[string]interface{}{"userId": "12345"}},
			Protocol:  "http",
		},
	}

	data, err := json.Marshal(callAsyncAPI)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "asyncapi",
		"with": {
			"document": {
				"name": "MyAsyncAPIDoc",
				"endpoint": "http://example.com/asyncapi.json"
			},
			"operation": "user.signup",
			"server": { "name": "default-server" },
			"protocol": "http",
			"message": {
				"payload": { "userId": "12345" }
			}
		}
	}`, string(data))
}

func TestCallAsyncAPI_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "asyncapi",
		"with": {
			"document": {
				"name": "MyAsyncAPIDoc",
				"endpoint": "http://example.com/asyncapi.json"
			},
			"operation": "user.signup",
			"server": { "name": "default-server"},
			"protocol": "http",
			"message": {
				"payload": { "userId": "12345" }
			},
			"authentication": {
				"use": "asyncapi-auth-policy"
			}
		}
	}`

	var callAsyncAPI CallAsyncAPI
	err := json.Unmarshal([]byte(jsonData), &callAsyncAPI)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, callAsyncAPI.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, callAsyncAPI.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, callAsyncAPI.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, callAsyncAPI.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, callAsyncAPI.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, callAsyncAPI.Metadata)
	assert.Equal(t, "asyncapi", callAsyncAPI.Call)
	assert.Equal(t, "MyAsyncAPIDoc", callAsyncAPI.With.Document.Name)
	assert.Equal(t, "http://example.com/asyncapi.json", callAsyncAPI.With.Document.Endpoint.String())
	assert.Equal(t, "user.signup", callAsyncAPI.With.Operation)
	assert.Equal(t, "default-server", callAsyncAPI.With.Server.Name)
	assert.Equal(t, "http", callAsyncAPI.With.Protocol)
	assert.Equal(t, map[string]interface{}{"userId": "12345"}, callAsyncAPI.With.Message.Payload)
	assert.Equal(t, "asyncapi-auth-policy", *callAsyncAPI.With.Authentication.Use)
}

func TestCallFunction_MarshalJSON(t *testing.T) {
	callFunction := CallFunction{
		TaskBase: TaskBase{
			If:      &RuntimeExpression{Value: "${condition}"},
			Input:   &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}},
			Output:  &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}},
			Timeout: &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}},
			Then:    &FlowDirective{Value: "continue"},
			Metadata: map[string]interface{}{
				"meta": "data",
			},
		},
		Call: "myFunction",
		With: map[string]interface{}{
			"param1": "value1",
			"param2": 42,
		},
	}

	data, err := json.Marshal(callFunction)
	assert.NoError(t, err)
	assert.JSONEq(t, `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "myFunction",
		"with": {
			"param1": "value1",
			"param2": 42
		}
	}`, string(data))
}

func TestCallFunction_UnmarshalJSON(t *testing.T) {
	jsonData := `{
		"if": "${condition}",
		"input": { "from": {"key": "value"} },
		"output": { "as": {"result": "output"} },
		"timeout": { "after": "10s" },
		"then": "continue",
		"metadata": {"meta": "data"},
		"call": "myFunction",
		"with": {
			"param1": "value1",
			"param2": 42
		}
	}`

	var callFunction CallFunction
	err := json.Unmarshal([]byte(jsonData), &callFunction)
	assert.NoError(t, err)
	assert.Equal(t, &RuntimeExpression{Value: "${condition}"}, callFunction.If)
	assert.Equal(t, &Input{From: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"key": "value"}}}, callFunction.Input)
	assert.Equal(t, &Output{As: &ObjectOrRuntimeExpr{Value: map[string]interface{}{"result": "output"}}}, callFunction.Output)
	assert.Equal(t, &TimeoutOrReference{Timeout: &Timeout{After: NewDurationExpr("10s")}}, callFunction.Timeout)
	assert.Equal(t, &FlowDirective{Value: "continue"}, callFunction.Then)
	assert.Equal(t, map[string]interface{}{"meta": "data"}, callFunction.Metadata)
	assert.Equal(t, "myFunction", callFunction.Call)

	// Adjust numeric values for comparison
	expectedWith := map[string]interface{}{
		"param1": "value1",
		"param2": float64(42), // Match JSON unmarshaling behavior
	}
	assert.Equal(t, expectedWith, callFunction.With)
}
