package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

// List of Standard Errors based on the Serverless Workflow specification.
// See: https://github.com/serverlessworkflow/specification/blob/main/dsl-reference.md#standard-error-types
const (
	ErrorTypeConfiguration  = "https://serverlessworkflow.io/spec/1.0.0/errors/configuration"
	ErrorTypeValidation     = "https://serverlessworkflow.io/spec/1.0.0/errors/validation"
	ErrorTypeExpression     = "https://serverlessworkflow.io/spec/1.0.0/errors/expression"
	ErrorTypeAuthentication = "https://serverlessworkflow.io/spec/1.0.0/errors/authentication"
	ErrorTypeAuthorization  = "https://serverlessworkflow.io/spec/1.0.0/errors/authorization"
	ErrorTypeTimeout        = "https://serverlessworkflow.io/spec/1.0.0/errors/timeout"
	ErrorTypeCommunication  = "https://serverlessworkflow.io/spec/1.0.0/errors/communication"
	ErrorTypeRuntime        = "https://serverlessworkflow.io/spec/1.0.0/errors/runtime"
)

type Error struct {
	Type     *URITemplateOrRuntimeExpr       `json:"type" validate:"required"`
	Status   int                             `json:"status" validate:"required"`
	Title    string                          `json:"title,omitempty"`
	Detail   string                          `json:"detail,omitempty"`
	Instance *JsonPointerOrRuntimeExpression `json:"instance,omitempty" validate:"omitempty"`
}

type ErrorFilter struct {
	Type     string `json:"type,omitempty"`
	Status   int    `json:"status,omitempty"`
	Instance string `json:"instance,omitempty"`
	Title    string `json:"title,omitempty"`
	Details  string `json:"details,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("[%d] %s: %s (%s). Origin: '%s'", e.Status, e.Title, e.Detail, e.Type, e.Instance)
}

// WithInstanceRef ensures the error has a valid JSON Pointer reference
func (e *Error) WithInstanceRef(workflow *Workflow, taskName string) *Error {
	if e == nil {
		return nil
	}

	// Check if the instance is already set
	if e.Instance.IsValid() {
		return e
	}

	// Generate a JSON pointer reference for the task within the workflow
	instance, pointerErr := GenerateJSONPointer(workflow, taskName)
	if pointerErr == nil {
		e.Instance = &JsonPointerOrRuntimeExpression{Value: instance}
	}
	// TODO: log the pointer error

	return e
}

// newError creates a new structured error
func newError(errType string, status int, title string, detail error, instance string) *Error {
	if detail != nil {
		return &Error{
			Type:   NewUriTemplate(errType),
			Status: status,
			Title:  title,
			Detail: detail.Error(),
			Instance: &JsonPointerOrRuntimeExpression{
				Value: instance,
			},
		}
	}

	return &Error{
		Type:   NewUriTemplate(errType),
		Status: status,
		Title:  title,
		Instance: &JsonPointerOrRuntimeExpression{
			Value: instance,
		},
	}
}

// Convenience Functions for Standard Errors

func NewErrConfiguration(detail error, instance string) *Error {
	return newError(
		ErrorTypeConfiguration,
		400,
		"Configuration Error",
		detail,
		instance,
	)
}

func NewErrValidation(detail error, instance string) *Error {
	return newError(
		ErrorTypeValidation,
		400,
		"Validation Error",
		detail,
		instance,
	)
}

func NewErrExpression(detail error, instance string) *Error {
	return newError(
		ErrorTypeExpression,
		400,
		"Expression Error",
		detail,
		instance,
	)
}

func NewErrAuthentication(detail error, instance string) *Error {
	return newError(
		ErrorTypeAuthentication,
		401,
		"Authentication Error",
		detail,
		instance,
	)
}

func NewErrAuthorization(detail error, instance string) *Error {
	return newError(
		ErrorTypeAuthorization,
		403,
		"Authorization Error",
		detail,
		instance,
	)
}

func NewErrTimeout(detail error, instance string) *Error {
	return newError(
		ErrorTypeTimeout,
		408,
		"Timeout Error",
		detail,
		instance,
	)
}

func NewErrCommunication(detail error, instance string) *Error {
	return newError(
		ErrorTypeCommunication,
		500,
		"Communication Error",
		detail,
		instance,
	)
}

func NewErrRuntime(detail error, instance string) *Error {
	return newError(
		ErrorTypeRuntime,
		500,
		"Runtime Error",
		detail,
		instance,
	)
}

// Error Classification Functions

func IsErrConfiguration(err error) bool {
	return isErrorType(err, ErrorTypeConfiguration)
}

func IsErrValidation(err error) bool {
	return isErrorType(err, ErrorTypeValidation)
}

func IsErrExpression(err error) bool {
	return isErrorType(err, ErrorTypeExpression)
}

func IsErrAuthentication(err error) bool {
	return isErrorType(err, ErrorTypeAuthentication)
}

func IsErrAuthorization(err error) bool {
	return isErrorType(err, ErrorTypeAuthorization)
}

func IsErrTimeout(err error) bool {
	return isErrorType(err, ErrorTypeTimeout)
}

func IsErrCommunication(err error) bool {
	return isErrorType(err, ErrorTypeCommunication)
}

func IsErrRuntime(err error) bool {
	return isErrorType(err, ErrorTypeRuntime)
}

// Helper function to check error type
func isErrorType(err error, errorType string) bool {
	var e *Error
	if ok := errors.As(err, &e); ok && strings.EqualFold(e.Type.String(), errorType) {
		return true
	}
	return false
}

// AsError attempts to extract a known error type from the given error.
// If the error is one of the predefined structured errors, it returns the *Error.
// Otherwise, it returns nil.
func AsError(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e // Successfully extracted as a known error type
	}
	return nil // Not a known error
}

// Serialization and Deserialization Functions

func ErrorToJSON(err *Error) (string, error) {
	if err == nil {
		return "", fmt.Errorf("error is nil")
	}
	jsonBytes, marshalErr := json.Marshal(err)
	if marshalErr != nil {
		return "", fmt.Errorf("failed to marshal error: %w", marshalErr)
	}
	return string(jsonBytes), nil
}

func ErrorFromJSON(jsonStr string) (*Error, error) {
	var errObj Error
	if err := json.Unmarshal([]byte(jsonStr), &errObj); err != nil {
		return nil, fmt.Errorf("failed to unmarshal error JSON: %w", err)
	}
	return &errObj, nil
}

// JsonPointer functions

func findJsonPointer(data interface{}, target string, path string) (string, bool) {
	switch node := data.(type) {
	case map[string]interface{}:
		for key, value := range node {
			newPath := fmt.Sprintf("%s/%s", path, key)
			if key == target {
				return newPath, true
			}
			if result, found := findJsonPointer(value, target, newPath); found {
				return result, true
			}
		}
	case []interface{}:
		for i, item := range node {
			newPath := fmt.Sprintf("%s/%d", path, i)
			if result, found := findJsonPointer(item, target, newPath); found {
				return result, true
			}
		}
	}
	return "", false
}

// GenerateJSONPointer Function to generate JSON Pointer from a Workflow reference
func GenerateJSONPointer(workflow *Workflow, targetNode interface{}) (string, error) {
	// Convert struct to JSON
	jsonData, err := json.Marshal(workflow)
	if err != nil {
		return "", fmt.Errorf("error marshalling to JSON: %v", err)
	}

	// Convert JSON to a generic map for traversal
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return "", fmt.Errorf("error unmarshalling JSON: %v", err)
	}

	transformedNode := ""
	switch node := targetNode.(type) {
	case string:
		transformedNode = node
	default:
		transformedNode = strings.ToLower(reflect.TypeOf(targetNode).Name())
	}

	// Search for the target node
	jsonPointer, found := findJsonPointer(jsonMap, transformedNode, "")
	if !found {
		return "", fmt.Errorf("node '%s' not found", targetNode)
	}

	return jsonPointer, nil
}
