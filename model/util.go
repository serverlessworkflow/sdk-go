// Copyright 2020 The Serverless Workflow Specification Authors
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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"sync/atomic"
	"time"

	"sigs.k8s.io/yaml"

	val "github.com/serverlessworkflow/sdk-go/v2/validator"
)

// Kind ...
// +k8s:deepcopy-gen=false
type Kind interface {
	KindValues() []string
	String() string
}

// TODO: Remove global variable
var httpClient = http.Client{Timeout: time.Duration(1) * time.Second}

// UnmarshalError ...
// +k8s:deepcopy-gen=false
type UnmarshalError struct {
	err           error
	parameterName string
	primitiveType reflect.Kind
	objectType    reflect.Kind
}

func (e *UnmarshalError) Error() string {
	if e.err == nil {
		panic("unmarshalError fail")
	}

	var syntaxErr *json.SyntaxError
	var unmarshalTypeErr *json.UnmarshalTypeError
	if errors.As(e.err, &syntaxErr) {
		return fmt.Sprintf("%s has a syntax error %q", e.parameterName, syntaxErr.Error())

	} else if errors.As(e.err, &unmarshalTypeErr) {
		return e.unmarshalMessageError(unmarshalTypeErr)
	}

	return e.err.Error()
}

func (e *UnmarshalError) unmarshalMessageError(err *json.UnmarshalTypeError) string {
	if err.Struct == "" && err.Field == "" {
		primitiveTypeName := e.primitiveType.String()
		var objectTypeName string
		if e.objectType != reflect.Invalid {
			switch e.objectType {
			case reflect.Struct:
				objectTypeName = "object"
			case reflect.Map:
				objectTypeName = "object"
			case reflect.Slice:
				objectTypeName = "array"
			default:
				objectTypeName = e.objectType.String()
			}
		}
		return fmt.Sprintf("%s must be %s or %s", e.parameterName, primitiveTypeName, objectTypeName)

	} else if err.Struct != "" && err.Field != "" {
		var primitiveTypeName string
		value := reflect.New(err.Type)
		if valKinds, ok := value.Elem().Interface().(val.Kind); ok {
			values := valKinds.KindValues()
			if len(values) <= 2 {
				primitiveTypeName = strings.Join(values, " or ")
			} else {
				primitiveTypeName = fmt.Sprintf("%s, %s", strings.Join(values[:len(values)-2], ", "), strings.Join(values[len(values)-2:], " or "))
			}
		} else {
			primitiveTypeName = err.Type.Name()
		}

		return fmt.Sprintf("%s.%s must be %s", e.parameterName, err.Field, primitiveTypeName)
	}

	return err.Error()
}

func loadExternalResource(url string) (b []byte, err error) {
	index := strings.Index(url, "://")
	if index == -1 {
		b, err = getBytesFromFile(url)
	} else {
		scheme := url[:index]
		switch scheme {
		case "http", "https":
			b, err = getBytesFromHttp(url)
		case "file":
			b, err = getBytesFromFile(url[index+3:])
		default:
			return nil, fmt.Errorf("unsupported scheme: %q", scheme)
		}
	}
	if err != nil {
		return
	}

	// TODO: optimize this
	// NOTE: In specification, we can declare independent definitions with another file format, so
	// we must convert independently yaml source to json format data before unmarshal.
	if !json.Valid(b) {
		b, err = yaml.YAMLToJSON(b)
		if err != nil {
			return nil, err
		}
		return b, nil
	}

	return b, nil
}

func getBytesFromFile(path string) ([]byte, error) {
	if WebAssembly() {
		return nil, fmt.Errorf("unsupported open file")
	}

	// if path is relative, search in include paths
	if !filepath.IsAbs(path) {
		paths := IncludePaths()
		pathFound := false
		for i := 0; i < len(paths) && !pathFound; i++ {
			sn := filepath.Join(paths[i], path)
			_, err := os.Stat(sn)
			if err != nil {
				if !errors.Is(err, os.ErrNotExist) {
					return nil, err
				}
			} else {
				path = sn
				pathFound = true
			}
		}
		if !pathFound {
			return nil, fmt.Errorf("file not found: %q", path)
		}
	}

	return os.ReadFile(filepath.Clean(path))
}

func getBytesFromHttp(url string) ([]byte, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf := new(bytes.Buffer)
	if _, err = buf.ReadFrom(resp.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func unmarshalObjectOrFile[U any](parameterName string, data []byte, valObject *U) error {
	var valString string
	err := unmarshalPrimitiveOrObject(parameterName, data, &valString, valObject)
	if err != nil || valString == "" {
		return err
	}

	// Assumes that the value inside `data` is a path to a known location.
	// Returns the content of the file or a not nil error reference.
	data, err = loadExternalResource(valString)
	if err != nil {
		return err
	}

	data = bytes.TrimSpace(data)
	if data[0] != '{' && data[0] != '[' {
		return errors.New("invalid external resource definition")
	}

	if data[0] == '[' && parameterName != "auth" && parameterName != "secrets" {
		return errors.New("invalid external resource definition")
	}

	data = bytes.TrimSpace(data)
	if data[0] == '{' && parameterName != "constants" && parameterName != "timeouts" {
		extractData := map[string]json.RawMessage{}
		err = json.Unmarshal(data, &extractData)
		if err != nil {
			return &UnmarshalError{
				err:           err,
				parameterName: parameterName,
				primitiveType: reflect.TypeOf(*valObject).Kind(),
			}
		}

		var ok bool
		if data, ok = extractData[parameterName]; !ok {
			return fmt.Errorf("external resource parameter not found: %q", parameterName)
		}
	}

	return unmarshalObject(parameterName, data, valObject)
}

func unmarshalPrimitiveOrObject[T string | bool, U any](parameterName string, data []byte, valPrimitive *T, valStruct *U) error {
	data = bytes.TrimSpace(data)
	if len(data) == 0 {
		// TODO: Normalize error messages
		return fmt.Errorf("%s no bytes to unmarshal", parameterName)
	}

	isObject := data[0] == '{' || data[0] == '['
	var err error
	if isObject {
		err = unmarshalObject(parameterName, data, valStruct)
	} else {
		err = unmarshalPrimitive(parameterName, data, valPrimitive)
	}

	var unmarshalError *UnmarshalError
	if errors.As(err, &unmarshalError) {
		unmarshalError.objectType = reflect.TypeOf(*valStruct).Kind()
		unmarshalError.primitiveType = reflect.TypeOf(*valPrimitive).Kind()
	}

	return err
}

func unmarshalPrimitive[T string | bool](parameterName string, data []byte, value *T) error {
	if value == nil {
		return nil
	}

	err := json.Unmarshal(data, value)
	if err != nil {
		return &UnmarshalError{
			err:           err,
			parameterName: parameterName,
			primitiveType: reflect.TypeOf(*value).Kind(),
		}
	}

	return nil
}

func unmarshalObject[U any](parameterName string, data []byte, value *U) error {
	if value == nil {
		return nil
	}

	err := json.Unmarshal(data, value)
	if err != nil {
		return &UnmarshalError{
			err:           err,
			parameterName: parameterName,
			objectType:    reflect.TypeOf(*value).Kind(),
		}
	}

	return nil
}

var defaultIncludePaths atomic.Value

func init() {
	// No execute set include path to suport webassembly
	if WebAssembly() {
		return
	}

	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	SetIncludePaths([]string{wd})
}

// IncludePaths will return the search path for non-absolute import file
func IncludePaths() []string {
	return defaultIncludePaths.Load().([]string)
}

// SetIncludePaths will update the search path for non-absolute import file
func SetIncludePaths(paths []string) {
	for _, path := range paths {
		if !filepath.IsAbs(path) {
			panic(fmt.Errorf("%s must be an absolute file path", path))
		}
	}

	defaultIncludePaths.Store(paths)
}

func WebAssembly() bool {
	return runtime.GOOS == "js" && runtime.GOARCH == "wasm"
}
