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
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"
)

const prefix = "file:/"

func getBytesFromFile(s string) (b []byte, err error) {
	// #nosec
	if resp, err := http.Get(s); err == nil {
		defer resp.Body.Close()
		buf := new(bytes.Buffer)
		if _, err = buf.ReadFrom(resp.Body); err != nil {
			return nil, err
		}
		return buf.Bytes(), nil
	}
	if strings.HasPrefix(s, prefix) {
		s = strings.TrimPrefix(s, prefix)
	} else {
		if s, err = filepath.Abs(s); err != nil {
			return nil, err
		}
	}
	if b, err = ioutil.ReadFile(filepath.Clean(s)); err != nil {
		return nil, err
	}
	return b, nil
}

func requiresNotNilOrEmpty(value interface{}) string {
	if value == nil {
		return ""
	}
	return value.(string)
}

func unmarshalString(data []byte) (string, error) {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return "", err
	}
	return value, nil
}

func unmarshalKey(key string, data map[string]json.RawMessage, output interface{}) error {
	if _, found := data[key]; found {
		if err := json.Unmarshal(data[key], output); err != nil {
			return err
		}
	}
	return nil
}
