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

import "testing"

func buildAuth(workflow *Workflow, name string) *Auth {
	auth := Auth{
		Name:   name,
		Scheme: AuthTypeBasic,
	}
	workflow.Auth = append(workflow.Auth, auth)
	return &workflow.Auth[len(workflow.Auth)-1]
}

func buildBasicAuthProperties(auth *Auth) *BasicAuthProperties {
	auth.Scheme = AuthTypeBasic
	auth.Properties = AuthProperties{
		Basic: &BasicAuthProperties{
			Username: "username",
			Password: "password",
		},
	}

	return auth.Properties.Basic
}

func buildOAuth2AuthProperties(auth *Auth) *OAuth2AuthProperties {
	auth.Scheme = AuthTypeOAuth2
	auth.Properties = AuthProperties{
		OAuth2: &OAuth2AuthProperties{
			ClientID:  "clientId",
			GrantType: GrantTypePassword,
		},
	}

	return auth.Properties.OAuth2
}

func buildBearerAuthProperties(auth *Auth) *BearerAuthProperties {
	auth.Scheme = AuthTypeBearer
	auth.Properties = AuthProperties{
		Bearer: &BearerAuthProperties{
			Token: "token",
		},
	}

	return auth.Properties.Bearer
}

func TestAuthStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	auth := buildAuth(baseWorkflow, "auth 1")
	buildBasicAuthProperties(auth)

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth[0].Name = ""
				return *model
			},
			Err: `workflow.auth[0].name is required`,
		},
		{
			Desp: "repeat",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth = append(model.Auth, model.Auth[0])
				return *model
			},
			Err: `workflow.auth has duplicate "name"`,
		},
	}
	StructLevelValidationCtx(t, testCases)
}

func TestBasicAuthPropertiesStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	auth := buildAuth(baseWorkflow, "auth 1")
	buildBasicAuthProperties(auth)

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth[0].Properties.Basic.Username = ""
				model.Auth[0].Properties.Basic.Password = ""
				return *model
			},
			Err: `workflow.auth[0].properties.basic.username is required
workflow.auth[0].properties.basic.password is required`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestBearerAuthPropertiesStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	auth := buildAuth(baseWorkflow, "auth 1")
	buildBearerAuthProperties(auth)

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth[0].Properties.Bearer.Token = ""
				return *model
			},
			Err: `workflow.auth[0].properties.bearer.token is required`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}

func TestOAuth2AuthPropertiesPropertiesStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()
	auth := buildAuth(baseWorkflow, "auth 1")
	buildOAuth2AuthProperties(auth)

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
		{
			Desp: "success",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				return *model
			},
		},
		{
			Desp: "required",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth[0].Properties.OAuth2.GrantType = ""
				model.Auth[0].Properties.OAuth2.ClientID = ""
				return *model
			},
			Err: `workflow.auth[0].properties.oAuth2.grantType is required
workflow.auth[0].properties.oAuth2.clientID is required`,
		},
		{
			Desp: "oneofkind",
			Model: func() Workflow {
				model := baseWorkflow.DeepCopy()
				model.Auth[0].Properties.OAuth2.GrantType = GrantTypePassword + "invalid"
				return *model
			},
			Err: `workflow.auth[0].properties.oAuth2.grantType need by one of [password clientCredentials tokenExchange]`,
		},
	}

	StructLevelValidationCtx(t, testCases)
}
