// Copyright 2022 The Serverless Workflow Specification Authors
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

package parser

import (
	"embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/go-multierror"
	"github.com/santhosh-tekuri/jsonschema/v5"
	_ "github.com/santhosh-tekuri/jsonschema/v5/httploader"
	"sigs.k8s.io/yaml"
)

func validateSchema(fs embed.FS, path string, schema *jsonschema.Schema) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	workflowFilePaths, err := os.ReadDir(path)
	if err != nil {
		return err
	}

	var result error

	for _, f := range workflowFilePaths {
		if f.IsDir() {
			if err = validateSchema(fs, fmt.Sprintf("%s/%s", path, f.Name()), schema); err != nil {
				result = multierror.Append(result, err)
			}
			continue
		}

		var jsonBytes []byte
		relativeFilePath := fmt.Sprintf("%s/%s", path, f.Name())
		switch filepath.Ext(f.Name()) {
		case extYAML:
			fallthrough
		case extYML:
			fileBytes, err := os.ReadFile(relativeFilePath)
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}
			if jsonBytes, err = yaml.YAMLToJSON(fileBytes); err != nil {
				result = multierror.Append(result, err)
				continue
			}
		case extJSON:
			jsonBytes, err = os.ReadFile(relativeFilePath)
			if err != nil {
				result = multierror.Append(result, err)
				continue
			}
		default:
			result = multierror.Append(result, fmt.Errorf("skipping %s/%s", cwd, relativeFilePath))
			continue
		}
		err = validateYamlAgainstSchema(jsonBytes, schema)
		if err != nil {
			result = multierror.Append(result, fmt.Errorf("%w\n%s/%s:", err, cwd, relativeFilePath))
			continue
		}
	}
	return result
}

func validateYamlAgainstSchema(jsonBytes []byte, schema *jsonschema.Schema) error {
	var m interface{}
	err := json.Unmarshal(jsonBytes, &m)
	if err != nil {
		return err
	}
	if err = schema.Validate(m); err != nil {
		return err
	}
	return nil
}
