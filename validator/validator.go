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

package validator

import (
	"context"

	validator "github.com/go-playground/validator/v10"
	"github.com/senseyeio/duration"
)

// TODO: expose a better validation message. See: https://pkg.go.dev/gopkg.in/go-playground/validator.v8#section-documentation

type Kinds interface {
	AllKinds() []string
	String() string
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	err := validate.RegisterValidationCtx("iso8601duration", validateISO8601TimeDurationFunc)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidation("oneofkind", oneOfKind)
	if err != nil {
		panic(err)
	}

}

// GetValidator gets the default validator.Validate reference
func GetValidator() *validator.Validate {
	return validate
}

// ValidateISO8601TimeDuration validate the string is iso8601 duration format
func ValidateISO8601TimeDuration(s string) error {
	_, err := duration.ParseISO8601(s)
	return err
}

func validateISO8601TimeDurationFunc(_ context.Context, fl validator.FieldLevel) bool {
	err := ValidateISO8601TimeDuration(fl.Field().String())
	return err == nil
}

func oneOfKind(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(Kinds); ok {
		for _, kindValue := range val.AllKinds() {
			if kindValue == val.String() {
				return true
			}
		}
	}

	return false
}
