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

import (
	validator "github.com/go-playground/validator/v10"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

func init() {
	val.GetValidator().RegisterStructValidationCtx(val.ValidationWrap(baseStateStructLevelValidationCtx), BaseState{})
}

func baseStateStructLevelValidationCtx(ctx val.ValidatorContext, structLevel validator.StructLevel) {
	baseState := structLevel.Current().Interface().(BaseState)
	if baseState.Type != StateTypeSwitch {
		validTransitionAndEnd(structLevel, baseState, baseState.Transition, baseState.End)
	}

	if baseState.CompensatedBy != "" {
		if baseState.UsedForCompensation {
			structLevel.ReportError(baseState.CompensatedBy, "CompensatedBy", "compensatedBy", val.TagRecursiveCompensation, "")
		}

		if ctx.MapStates.Contain(baseState.CompensatedBy) {
			value := ctx.MapStates.ValuesMap[baseState.CompensatedBy].(State).BaseState
			if value.UsedForCompensation && value.Type == StateTypeEvent {
				structLevel.ReportError(baseState.CompensatedBy, "CompensatedBy", "compensatedBy", val.TagCompensatedbyEventState, "")

			} else if !value.UsedForCompensation {
				structLevel.ReportError(baseState.CompensatedBy, "CompensatedBy", "compensatedBy", val.TagCompensatedby, "")
			}

		} else {
			structLevel.ReportError(baseState.CompensatedBy, "CompensatedBy", "compensatedBy", val.TagExists, "")
		}
	}
}
