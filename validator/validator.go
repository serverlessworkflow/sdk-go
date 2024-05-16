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
	"errors"
	"strconv"

	"github.com/relvacode/iso8601"
	"github.com/sosodev/duration"
	"k8s.io/apimachinery/pkg/util/intstr"

	validator "github.com/go-playground/validator/v10"
)

// TODO: expose a better validation message. See: https://pkg.go.dev/gopkg.in/go-playground/validator.v8#section-documentation

type Kind interface {
	KindValues() []string
	String() string
}

var validate *validator.Validate

func init() {
	validate = validator.New()

	err := validate.RegisterValidationCtx("iso8601duration", validateISO8601TimeDurationFunc)
	if err != nil {
		panic(err)
	}

	err = validate.RegisterValidationCtx("iso8601datetime", validateISO8601DatetimeFunc)
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
	if s == "" {
		return errors.New("could not parse duration string")
	}
	_, err := duration.Parse(s)
	if err != nil {
		return errors.New("could not parse duration string")
	}
	return err
}

func validateISO8601TimeDurationFunc(_ context.Context, fl validator.FieldLevel) bool {
	err := ValidateISO8601TimeDuration(fl.Field().String())
	return err == nil
}

// ValidateISO8601Datetime validate the string is iso8601 Datetime format
func ValidateISO8601Datetime(s string) error {
	_, err := iso8601.ParseString(s)
	return err
}

func validateISO8601DatetimeFunc(_ context.Context, fl validator.FieldLevel) bool {
	err := ValidateISO8601Datetime(fl.Field().String())
	return err == nil
}

func oneOfKind(fl validator.FieldLevel) bool {
	if val, ok := fl.Field().Interface().(Kind); ok {
		for _, value := range val.KindValues() {
			if value == val.String() {
				return true
			}
		}
	}

	return false
}

func ValidateGt0IntStr(value *intstr.IntOrString) bool {
	switch value.Type {
	case intstr.Int:
		if value.IntVal <= 0 {
			return false
		}
	case intstr.String:
		v, err := strconv.Atoi(value.StrVal)
		if err != nil {
			return false
		}

		if v <= 0 {
			return false
		}
	}

	return true
}
