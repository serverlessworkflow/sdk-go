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

package graph

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"sigs.k8s.io/yaml"

	"github.com/serverlessworkflow/sdk-go/v4/internal/util"
)

// TODO: Remove global variable
var HttpClient = http.Client{Timeout: time.Duration(1) * time.Second}

func unmarshalNode(n *Node, data []byte) error {
	data = bytes.TrimSpace(data)
	if data[0] == '{' {
		return unmarshalObject(n, data)
	} else if data[0] == '[' {
		return unmarshalList(n, data)
	}

	return json.Unmarshal(data, &n.value)
}

func unmarshalObject(n *Node, data []byte) error {
	dataMap := map[string]json.RawMessage{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	for key, val := range dataMap {
		node := n.Edge(key)
		err := json.Unmarshal(val, &node)
		if err != nil {
			return err
		}

	}

	return nil
}

func unmarshalList(n *Node, data []byte) error {
	dataMap := []json.RawMessage{}
	err := json.Unmarshal(data, &dataMap)
	if err != nil {
		return err
	}

	n.List(true)

	for i, val := range dataMap {
		key := fmt.Sprintf("%d", i)
		node := n.Edge(key)
		err := json.Unmarshal(val, &node)
		if err != nil {
			return err
		}
	}

	return nil
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
	if util.WebAssembly() {
		return nil, fmt.Errorf("unsupported open file")
	}

	// if path is relative, search in include paths
	if !filepath.IsAbs(path) {
		paths := util.IncludePaths()
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

	resp, err := HttpClient.Do(req)
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
