// Copyright 2024 The Serverless Workflow Specification Authors
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

package load

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"sigs.k8s.io/yaml"

	"github.com/serverlessworkflow/sdk-go/v4/graph"
	"github.com/serverlessworkflow/sdk-go/v4/internal/dsl"
)

const (
	extJSON = ".json"
	extYAML = ".yaml"
	extYML  = ".yml"
)

var supportedExt = []string{extYAML, extYML, extJSON}

func FromFile(path string) (*graph.Node, []byte, error) {
	if err := checkFilePath(path); err != nil {
		return nil, nil, err
	}

	fileBytes, err := os.ReadFile(filepath.Clean(path))
	if err != nil {
		return nil, nil, err
	}

	if strings.HasSuffix(path, extYAML) || strings.HasSuffix(path, extYML) {
		return FromYAMLSource(fileBytes)
	}

	return FromJSONSource(fileBytes)
}

func FromYAMLSource(source []byte) (*graph.Node, []byte, error) {
	jsonBytes, err := yaml.YAMLToJSON(source)
	if err != nil {
		return nil, nil, err
	}
	return FromJSONSource(jsonBytes)
}

func FromJSONSource(fileBytes []byte) (*graph.Node, []byte, error) {
	root, err := graph.UnmarshalJSON(fileBytes)
	if err != nil {
		return nil, nil, err
	}

	err = graph.LoadExternalResource(root)
	if err != nil {
		return nil, nil, err
	}

	err = dsl.ApplyDefault(root)
	if err != nil {
		return nil, nil, err
	}

	return root, fileBytes, nil
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
