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

package main

import (
	"fmt"
	"github.com/serverlessworkflow/sdk-go/v3/parser"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <examples-directory>")
		os.Exit(1)
	}

	baseDir := os.Args[1]
	supportedExt := []string{".json", ".yaml", ".yml"}
	errCount := 0

	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			for _, ext := range supportedExt {
				if filepath.Ext(path) == ext {
					fmt.Printf("Validating: %s\n", path)
					_, err := parser.FromFile(path)
					if err != nil {
						fmt.Printf("Validation failed for %s: %v\n", path, err)
						errCount++
					} else {
						fmt.Printf("Validation succeeded for %s\n", path)
					}
					break
				}
			}
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error walking the path %s: %v\n", baseDir, err)
		os.Exit(1)
	}

	if errCount > 0 {
		fmt.Printf("Validation failed for %d file(s).\n", errCount)
		os.Exit(1)
	}

	fmt.Println("All workflows validated successfully.")
}
