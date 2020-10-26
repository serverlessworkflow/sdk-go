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
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"os"
	"strings"
)

const (
	extJson = ".json"
	extYaml = ".yaml"
	extYml  = ".yml"
)

var supportedExt = []string{extYaml, extYml, extJson}

// FromFile parses the given Serverless Workflow file into the Workflow type.
func FromFile(path string) (*model.Workflow, error) {
	if err := checkFilePath(path); err != nil {
		return nil, err
	}
	fileBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	workflow := &model.Workflow{}
	if err := mustUnmarshalerFor(path)(fileBytes, workflow); err != nil {
		return nil, err
	}
	return workflow, nil
}

// mustUnmarshalerFor gets the Unmarshal function for the given file. Does not validate if the file exists.
func mustUnmarshalerFor(path string) func([]byte, interface{}) error {
	if strings.HasSuffix(path, extJson) {
		return json.Unmarshal
	} else if strings.HasSuffix(path, extYaml) || strings.HasSuffix(path, extYml) {
		return yaml.Unmarshal
	}
	// we panic to make it consistent with checkFilePath call
	panic(fmt.Errorf("unmarshal function not found for file '%s'. Supported extensions are %s", path, supportedExt))
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
