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

package parser

import (
	"encoding/json"
	"fmt"
	"github.com/serverlessworkflow/sdk-go/model"
	"io/ioutil"
	"os"
	"sigs.k8s.io/yaml"
	"strings"
)

const (
	extJSON = ".json"
	extYAML = ".yaml"
	extYML  = ".yml"
)

var supportedExt = []string{extYAML, extYML, extJSON}

// FromYAMLSource parses the given Serverless Workflow YAML source into the Workflow type.
func FromYAMLSource(source []byte) (workflow *model.Workflow, err error) {
	var jsonBytes []byte
	if jsonBytes, err = yaml.YAMLToJSON(source); err != nil {
		return nil, err
	}
	return FromJSONSource(jsonBytes)
}

// FromJSONSource parses the given Serverless Workflow JSON source into the Workflow type.
func FromJSONSource(source []byte) (workflow *model.Workflow, err error) {
	workflow = &model.Workflow{}
	if err := json.Unmarshal(source, workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

// FromFile parses the given Serverless Workflow file into the Workflow type.
func FromFile(path string) (*model.Workflow, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if strings.HasSuffix(path, extYAML) || strings.HasSuffix(path, extYML) {
		return FromYAMLSource(fileBytes)
	}
	return FromJSONSource(fileBytes)
}

// checkFilePath verifies if the file exists in the given path and if it's supported by the parser package
func checkFilePath(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return err
	}
	if info.IsDir() {
		return fmt.Errorf("file path '%s' must stand to a file", path)
	}
	for _, ext := range supportedExt {
		if strings.HasSuffix(path, ext) {
			return nil
		}
	}
	return fmt.Errorf("file extension not supported for '%s'. supported formats are %s", path, supportedExt)
}
