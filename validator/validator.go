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

import "gopkg.in/go-playground/validator.v8"

// TODO: expose a better validation message. See: https://pkg.go.dev/gopkg.in/go-playground/validator.v8#section-documentation

var validate *validator.Validate

func init() {
	validate = validator.New(&validator.Config{TagName: "validate"})
}

// GetValidator gets the default validator.Validate reference
func GetValidator() *validator.Validate {
	return validate
}
