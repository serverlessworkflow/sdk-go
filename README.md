# Go SDK for Serverless Workflow

The Go SDK for Serverless Workflow provides the [specification types](https://github.com/serverlessworkflow/specification/blob/v1.0.0-alpha5/schema/workflow.yaml) defined by the Serverless Workflow DSL in Go, making it easy to parse, validate, and interact with workflows.

---

## Table of Contents

- [Status](#status)
- [Releases](#releases)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Parsing Workflow Files](#parsing-workflow-files)
  - [Programmatic Workflow Creation](#programmatic-workflow-creation)
- [Slack Community](#slack-community)
- [Contributing](#contributing)
  - [Code Style](#code-style)
  - [EditorConfig](#editorconfig)
  - [Known Issues](#known-issues)

---

## Status

The current status of features implemented in the SDK is listed below:

| Feature                                     | Status             |
|-------------------------------------------- | ------------------ |
| Parse workflow JSON and YAML definitions    | :heavy_check_mark: |
| Programmatically build workflow definitions | :heavy_check_mark:    |
| Validate workflow definitions (Schema)      | :heavy_check_mark: |
| Validate workflow definitions (Integrity)   | :no_entry_sign: |
| Generate workflow diagram (SVG)             | :no_entry_sign:    |

---

## Releases

|                              Latest Releases                               |                       Conformance to Spec Version                        |
|:--------------------------------------------------------------------------:|:------------------------------------------------------------------------:|
| [v1.0.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v1.0.0) |  [v0.5](https://github.com/serverlessworkflow/specification/tree/0.5.x)  |
| [v2.0.1](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.0.1) |  [v0.6](https://github.com/serverlessworkflow/specification/tree/0.6.x)  |
| [v2.1.2](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.1.2) |  [v0.7](https://github.com/serverlessworkflow/specification/tree/0.7.x)  |
| [v2.4.3](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.4.1) |  [v0.8](https://github.com/serverlessworkflow/specification/tree/0.8.x)  |
| [v3.0.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v3.0.0) | [v1.0.0](https://github.com/serverlessworkflow/specification/releases/tag/v1.0.0-alpha5) |

---

## Getting Started

### Installation

To use the SDK in your Go project, run the following command:

```shell
$ go get github.com/serverlessworkflow/sdk-go/v3
```

This will update your `go.mod` file to include the Serverless Workflow SDK as a dependency.

Import the SDK in your Go file:

```go
import "github.com/serverlessworkflow/sdk-go/v3/model"
```

You can now use the SDK types and functions, for example:

```go
package main

import (
	"github.com/serverlessworkflow/sdk-go/v3/builder"
    "github.com/serverlessworkflow/sdk-go/v3/model"
)

func main() {
  workflowBuilder := New().
    SetDocument("1.0.0", "examples", "example-workflow", "1.0.0").
          AddTask("task1", &model.CallHTTP{
            TaskBase: model.TaskBase{
              If: &model.RuntimeExpression{Value: "${condition}"},
            },
            Call: "http",
            With: model.HTTPArguments{
              Method:   "GET",
              Endpoint: model.NewEndpoint("http://example.com"),
            },
          })
    workflow, _ := builder.Object(workflowBuilder)
    // use your models
}

```

### Parsing Workflow Files

The Serverless Workflow Specification supports YAML and JSON files. Use the following example to parse a workflow file into a Go data structure:

```go
package main

import (
    "github.com/serverlessworkflow/sdk-go/v3/model"
    "github.com/serverlessworkflow/sdk-go/v3/parser"
)

func ParseWorkflow(filePath string) (*model.Workflow, error) {
    workflow, err := parser.FromFile(filePath)
    if err != nil {
        return nil, err
    }
    return workflow, nil
}
```

This `Workflow` structure can then be used programmatically in your application.

### Programmatic Workflow Creation

Support for building workflows programmatically is planned for future releases. Stay tuned for updates in upcoming versions.

---

## Slack Community

Join the conversation and connect with other contributors on the [CNCF Slack](https://communityinviter.com/apps/cloud-native/cncf). Find us in the `#serverless-workflow-sdk` channel and say hello! 🙋

---

## Contributing

We welcome contributions to improve this SDK. Please refer to the sections below for guidance on maintaining project standards.

### Code Style

- Use `goimports` for import organization.
- Lint your code with:

```bash
make lint
```

To automatically fix lint issues, use:

```bash
make lint params=--fix
```

Example lint error:

```bash
$ make lint
make addheaders
make fmt
./hack/go-lint.sh
util/floatstr/floatstr_test.go:19: File is not `goimports`-ed (goimports)
        "k8s.io/apimachinery/pkg/util/yaml"
make: *** [lint] Error 1
```

### EditorConfig

For IntelliJ users, an example `.editorconfig` file is available [here](contrib/intellij.editorconfig). See the [Jetbrains documentation](https://www.jetbrains.com/help/idea/editorconfig.html) for usage details.

### Known Issues

#### MacOS Issue:

On MacOS, you might encounter the following error:

```
goimports: can't extract issues from gofmt diff output
```

To resolve this, install `diffutils`:

```bash
brew install diffutils
```

