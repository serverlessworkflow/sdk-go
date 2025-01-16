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
	"errors"
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	iso8601DurationPattern = regexp.MustCompile(`^P(\d+Y)?(\d+M)?(\d+W)?(\d+D)?(T(\d+H)?(\d+M)?(\d+S)?)?$`)
	semanticVersionPattern = regexp.MustCompile(`^(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)(?:-((?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*)(?:\.(?:0|[1-9]\d*|\d*[a-zA-Z-][0-9a-zA-Z-]*))*))?(?:\+([0-9a-zA-Z-]+(?:\.[0-9a-zA-Z-]+)*))?$`)
	hostnameRFC1123Pattern = regexp.MustCompile(`^[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?$`)
)

var validate *validator.Validate

func registerValidator(tag string, fn validator.Func) {

	if err := validate.RegisterValidation(tag, fn); err != nil {
		panic(fmt.Sprintf("Failed to register validator '%s': %v", tag, err))
	}
}

func init() {
	validate = validator.New()

	registerValidator("basic_policy", validateBasicPolicy)
	registerValidator("bearer_policy", validateBearerPolicy)
	registerValidator("digest_policy", validateDigestPolicy)
	registerValidator("oauth2_policy", validateOAuth2Policy)
	registerValidator("client_auth_type", validateOptionalOAuthClientAuthentication)
	registerValidator("encoding_type", validateOptionalOAuth2TokenRequestEncoding)

	registerValidator("semver_pattern", func(fl validator.FieldLevel) bool {
		return semanticVersionPattern.MatchString(fl.Field().String())
	})
	registerValidator("hostname_rfc1123", func(fl validator.FieldLevel) bool {
		return hostnameRFC1123Pattern.MatchString(fl.Field().String())
	})
	registerValidator("uri_pattern", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return LiteralUriPattern.MatchString(value)
	})
	registerValidator("uri_template_pattern", func(fl validator.FieldLevel) bool {
		value, ok := fl.Field().Interface().(string)
		if !ok {
			return false
		}
		return LiteralUriTemplatePattern.MatchString(value)
	})
	registerValidator("iso8601_duration", validateISO8601Duration)

	registerValidator("object_or_string", validateObjectOrString)
	registerValidator("object_or_runtime_expr", validateObjectOrRuntimeExpr)
	registerValidator("string_or_runtime_expr", validateStringOrRuntimeExpr)
	registerValidator("uri_template_or_runtime_expr", validateURITemplateOrRuntimeExpr)
	registerValidator("json_pointer_or_runtime_expr", validateJsonPointerOrRuntimeExpr)

	registerValidator("switch_item", validateSwitchItem)
	validate.RegisterStructValidation(validateTaskItem, TaskItem{})
}

func GetValidator() *validator.Validate {
	return validate
}

// validateTaskItem is a struct-level validation function for TaskItem.
func validateTaskItem(sl validator.StructLevel) {
	taskItem := sl.Current().Interface().(TaskItem)

	// Validate Key
	if taskItem.Key == "" {
		sl.ReportError(taskItem.Key, "Key", "Key", "required", "")
		return
	}

	// Validate Task is not nil
	if taskItem.Task == nil {
		sl.ReportError(taskItem.Task, "Task", "Task", "required", "")
		return
	}

	// Validate the concrete type of Task and capture nested errors
	switch t := taskItem.Task.(type) {
	case *CallHTTP:
		validateConcreteTask(sl, t, "Task")
	case *CallOpenAPI:
		validateConcreteTask(sl, t, "Task")
	case *CallGRPC:
		validateConcreteTask(sl, t, "Task")
	case *CallAsyncAPI:
		validateConcreteTask(sl, t, "Task")
	case *CallFunction:
		validateConcreteTask(sl, t, "Task")
	case *DoTask:
		validateConcreteTask(sl, t, "Task")
	case *ForkTask:
		validateConcreteTask(sl, t, "Task")
	case *EmitTask:
		validateConcreteTask(sl, t, "Task")
	case *ForTask:
		validateConcreteTask(sl, t, "Task")
	case *ListenTask:
		validateConcreteTask(sl, t, "Task")
	case *RaiseTask:
		validateConcreteTask(sl, t, "Task")
	case *RunTask:
		validateConcreteTask(sl, t, "Task")
	case *SetTask:
		validateConcreteTask(sl, t, "Task")
	case *SwitchTask:
		validateConcreteTask(sl, t, "Task")
	case *TryTask:
		validateConcreteTask(sl, t, "Task")
	case *WaitTask:
		validateConcreteTask(sl, t, "Task")
	default:
		sl.ReportError(taskItem.Task, "Task", "Task", "unknown_task", "unrecognized task type")
	}
}

// validateConcreteTask validates a concrete Task type and reports nested errors.
func validateConcreteTask(sl validator.StructLevel, task interface{}, fieldName string) {
	err := validate.Struct(task)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, ve := range validationErrors {
				// Report only nested fields to avoid duplicates
				if ve.Namespace() != fieldName {
					sl.ReportError(ve.Value(), fieldName+"."+ve.StructNamespace(), ve.StructField(), ve.Tag(), ve.Param())
				}
			}
		}
	}
}

// func validateSwitchItem(fl validator.FieldLevel) bool { is a custom validation function for SwitchItem.
func validateSwitchItem(fl validator.FieldLevel) bool {
	switchItem, ok := fl.Field().Interface().(SwitchItem)
	if !ok {
		return false
	}
	return len(switchItem) == 1
}

// validateBasicPolicy ensures BasicAuthenticationPolicy has mutually exclusive fields set.
func validateBasicPolicy(fl validator.FieldLevel) bool {
	policy, ok := fl.Parent().Interface().(BasicAuthenticationPolicy)
	if !ok {
		return false
	}
	if (policy.Username != "" || policy.Password != "") && policy.Use != "" {
		return false
	}
	return true
}

// validateBearerPolicy ensures BearerAuthenticationPolicy has mutually exclusive fields set.
func validateBearerPolicy(fl validator.FieldLevel) bool {
	policy, ok := fl.Parent().Interface().(BearerAuthenticationPolicy)
	if !ok {
		return false
	}
	if policy.Token != "" && policy.Use != "" {
		return false
	}
	return true
}

// validateDigestPolicy ensures DigestAuthenticationPolicy has mutually exclusive fields set.
func validateDigestPolicy(fl validator.FieldLevel) bool {
	policy, ok := fl.Parent().Interface().(DigestAuthenticationPolicy)
	if !ok {
		return false
	}
	if (policy.Username != "" || policy.Password != "") && policy.Use != "" {
		return false
	}
	return true
}

func validateOAuth2Policy(fl validator.FieldLevel) bool {
	policy, ok := fl.Parent().Interface().(OAuth2AuthenticationPolicy)
	if !ok {
		return false
	}

	if (policy.Properties != nil || policy.Endpoints != nil) && policy.Use != "" {
		return false // Both fields are set, invalid
	}
	if policy.Properties == nil && policy.Use == "" {
		return false // Neither field is set, invalid
	}
	return true
}

// validateOptionalOAuthClientAuthentication checks if the given value is a valid OAuthClientAuthenticationType.
func validateOptionalOAuthClientAuthentication(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	if len(value) == 0 {
		return true
	}
	switch OAuthClientAuthenticationType(value) {
	case
		OAuthClientAuthClientSecretBasic,
		OAuthClientAuthClientSecretPost,
		OAuthClientAuthClientSecretJWT,
		OAuthClientAuthPrivateKeyJWT,
		OAuthClientAuthNone:
		return true
	default:
		return false
	}
}

func validateOptionalOAuth2TokenRequestEncoding(fl validator.FieldLevel) bool {
	value := fl.Field().String()

	// Allow empty fields (optional case)
	if value == "" {
		return true
	}

	// Validate against allowed constants
	switch OAuth2TokenRequestEncodingType(value) {
	case
		EncodingTypeFormUrlEncoded,
		EncodingTypeApplicationJson:
		return true
	default:
		return false
	}
}

func validateObjectOrString(fl validator.FieldLevel) bool {
	// Access the "Value" field
	value := fl.Field().Interface()

	// Validate based on the type of "Value"
	switch v := value.(type) {
	case string:
		return v != "" // Validate non-empty strings.
	case map[string]interface{}:
		return len(v) > 0 // Validate non-empty objects.
	default:
		return false // Reject unsupported types.
	}
}

func validateObjectOrRuntimeExpr(fl validator.FieldLevel) bool {
	// Retrieve the field value using reflection
	value := fl.Field().Interface()

	// Validate based on the type
	switch v := value.(type) {
	case RuntimeExpression:
		return v.IsValid() // Validate runtime expression format.
	case map[string]interface{}:
		return len(v) > 0 // Validate non-empty objects.
	default:
		return false // Unsupported types.
	}
}

func validateStringOrRuntimeExpr(fl validator.FieldLevel) bool {
	// Retrieve the field value using reflection
	value := fl.Field().Interface()

	// Validate based on the type
	switch v := value.(type) {
	case RuntimeExpression:
		return v.IsValid() // Validate runtime expression format.
	case string:
		return v != "" // Validate non-empty strings.
	default:
		return false // Unsupported types.
	}
}

func validateURITemplateOrRuntimeExpr(fl validator.FieldLevel) bool {
	value := fl.Field().Interface()

	// Handle nil or empty values when 'omitempty' is used
	if value == nil {
		return true
	}

	switch v := value.(type) {
	case LiteralUri:
		return LiteralUriPattern.MatchString(v.String())
	case LiteralUriTemplate:
		return LiteralUriTemplatePattern.MatchString(v.String())
	case RuntimeExpression:
		return v.IsValid()
	case string:
		// Check if the string is a valid URI
		if LiteralUriPattern.MatchString(v) {
			return true
		}

		// Check if the string is a valid URI Template
		if LiteralUriTemplatePattern.MatchString(v) {
			return true
		}

		// Check if the string is a valid RuntimeExpression
		expression := RuntimeExpression{Value: v}
		return expression.IsValid()
	default:
		fmt.Printf("Unsupported type in URITemplateOrRuntimeExpr.Value: %T\n", v)
		return false
	}
}

func validateJsonPointerOrRuntimeExpr(fl validator.FieldLevel) bool {
	// Retrieve the field value using reflection
	value := fl.Field().Interface()

	// Validate based on the type
	switch v := value.(type) {
	case string: // JSON Pointer
		return JSONPointerPattern.MatchString(v)
	case RuntimeExpression:
		return v.IsValid()
	default:
		return false // Unsupported types.
	}
}

func validateISO8601Duration(fl validator.FieldLevel) bool {
	value, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	return iso8601DurationPattern.MatchString(value)
}
