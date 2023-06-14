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

func TestAuthStructLevelValidation(t *testing.T) {
	baseWorkflow := buildWorkflow()

	baseWorkflow.Auth = Auths{{
		Name: "auth 1",
	}}

	operationState := buildOperationState(baseWorkflow, "start state")
	buildEndByState(operationState, true, false)
	action1 := buildActionByOperationState(operationState, "action 1")
	buildFunctionRef(baseWorkflow, action1, "function 1")

	testCases := []ValidationCase{
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
	testCases := []ValidationCase{}
	StructLevelValidationCtx(t, testCases)
}

func TestBearerAuthPropertiesStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase{}
	StructLevelValidationCtx(t, testCases)
}

func TestOAuth2AuthPropertiesPropertiesStructLevelValidation(t *testing.T) {
	testCases := []ValidationCase{}
	StructLevelValidationCtx(t, testCases)
}
