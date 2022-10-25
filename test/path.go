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

package test

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/pkg/errors"
)

// CurrentProjectPath get the project root path
func CurrentProjectPath() string {
	path := currentFilePath()

	ppath, err := filepath.Abs(filepath.Join(filepath.Dir(path), "../"))
	if err != nil {
		panic(errors.Wrapf(err, "Get current project path with %s failed", path))
	}

	f, err := os.Stat(ppath)
	if err != nil {
		panic(errors.Wrapf(err, "Stat project path %v failed", ppath))
	}

	if f.Mode()&os.ModeSymlink != 0 {
		fpath, err := os.Readlink(ppath)
		if err != nil {
			panic(errors.Wrapf(err, "Readlink from path %v failed", fpath))
		}
		ppath = fpath
	}

	return ppath
}

func currentFilePath() string {
	_, file, _, _ := runtime.Caller(1)
	return file
}
