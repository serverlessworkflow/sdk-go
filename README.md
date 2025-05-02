# Go SDK for Serverless Workflow

The Go SDK for Serverless Workflow provides strongly-typed structures for the [Serverless Workflow specification](https://github.com/serverlessworkflow/specification/blob/v1.0.0/schema/workflow.yaml). It simplifies parsing, validating, and interacting with workflows in Go. Starting from version `v3.1.0`, the SDK also includes a partial reference implementation, allowing users to execute workflows directly within their Go applications.

---

## Table of Contents

- [Status](#status)
- [Releases](#releases)
- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Basic Usage](#basic-usage)
  - [Parsing Workflow Files](#parsing-workflow-files)
  - [Programmatic Workflow Creation](#programmatic-workflow-creation)
- [Reference Implementation](#reference-implementation)
  - [Example: Running a Workflow](#example-running-a-workflow)
- [Slack Community](#slack-community)
- [Contributing](#contributing)
  - [Code Style](#code-style)
  - [EditorConfig](#editorconfig)
  - [Known Issues](#known-issues)

---

## Status

This table indicates the current state of implementation of various SDK features:

| Feature                                     | Status              |
|-------------------------------------------- |---------------------|
| Parse workflow JSON and YAML definitions    | :heavy_check_mark:  |
| Programmatically build workflow definitions | :heavy_check_mark:  |
| Validate workflow definitions (Schema)      | :heavy_check_mark:  |
| Specification Implementation                | :heavy_check_mark:* |
| Validate workflow definitions (Integrity)   | :no_entry_sign:     |
| Generate workflow diagram (SVG)             | :no_entry_sign:     |

> **Note**: *Implementation is partial; contributions are encouraged.

---

## Releases

|                              Latest Releases                               |                            Conformance to Spec Version                            |
|:--------------------------------------------------------------------------:|:---------------------------------------------------------------------------------:|
| [v1.0.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v1.0.0) |      [v0.5](https://github.com/serverlessworkflow/specification/tree/0.5.x)       |
| [v2.0.1](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.0.1) |      [v0.6](https://github.com/serverlessworkflow/specification/tree/0.6.x)       |
| [v2.1.2](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.1.2) |      [v0.7](https://github.com/serverlessworkflow/specification/tree/0.7.x)       |
| [v2.5.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.5.0) |      [v0.8](https://github.com/serverlessworkflow/specification/tree/0.8.x)       |
| [v3.0.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v3.0.0) | [v1.0.0](https://github.com/serverlessworkflow/specification/releases/tag/v1.0.0) |

---

## Reference Implementation

The SDK provides a partial reference runner to execute your workflows:

### Example: Running a Workflow

Below is a simple YAML workflow that sets a message and then prints it:

```yaml
document:
  dsl: "1.0.0"
  namespace: "examples"
  name: "simple-workflow"
  version: "1.0.0"
do:
  - set:
      message: "Hello from the Serverless Workflow SDK in Go!"
```

You can execute this workflow using the following Go program:

Example of executing a workflow defined in YAML:

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"

    "github.com/serverlessworkflow/sdk-go/v3/impl"
    "github.com/serverlessworkflow/sdk-go/v3/parser"
)

func RunWorkflow(workflowFilePath string, input map[string]interface{}) (interface{}, error) {
    data, err := os.ReadFile(filepath.Clean(workflowFilePath))
    if err != nil {
        return nil, err
    }
    workflow, err := parser.FromYAMLSource(data)
    if err != nil {
        return nil, err
    }

    runner := impl.NewDefaultRunner(workflow)
    output, err := runner.Run(input)
    if err != nil {
        return nil, err
    }
    return output, nil
}

func main() {
    output, err := RunWorkflow("./myworkflow.yaml", map[string]interface{}{"shouldCall": true})
    if err != nil {
        panic(err)
    }
    fmt.Printf("Workflow completed with output: %v\n", output)
}
```

### Implementation Roadmap

The table below lists the current state of this implementation. This table is a roadmap for the project based on the [DSL Reference doc](https://github.com/serverlessworkflow/specification/blob/v1.0.0/dsl-reference.md).

| Feature | State |
| ----------- | --------------- |
| Workflow Document | ✅  |
| Workflow Use | 🟡 |
| Workflow Schedule | ❌ | 
| Task Call | ❌ |
| Task Do | ✅ |
| Task Emit | ❌ | 
| Task For | ✅ |
| Task Fork | ❌ | 
| Task Listen | ❌ | 
| Task Raise | ✅ |
| Task Run | ❌ |
| Task Set | ✅ | 
| Task Switch | ✅ | 
| Task Try | ❌ | 
| Task Wait | ❌ |
| Lifecycle Events | 🟡 |
| External Resource | ❌ |
| Authentication | ❌ |
| Catalog | ❌ |
| Extension | ❌ |
| Error | ✅ | 
| Event Consumption Strategies | ❌ | 
| Retry | ❌ |
| Input | ✅ |
| Output | ✅ |
| Export | ✅ |
| Timeout | ❌ |
| Duration | ❌ |
| Endpoint | ✅ |
| HTTP Response | ❌ |
| HTTP Request | ❌ |
| URI Template | ✅ | 
| Container Lifetime | ❌ |
| Process Result | ❌ |
| AsyncAPI Server | ❌ |
| AsyncAPI Outbound Message | ❌ |
| AsyncAPI Subscription | ❌ |
| Workflow Definition Reference | ✅ |
| Subscription Iterator | ❌ |

We love contributions! Our aim is to have a complete implementation to serve as a reference or to become a project on its own to favor the CNCF Ecosystem.

If you are willing to help, please [file a sub-task](https://github.com/serverlessworkflow/sdk-go/issues/221) in this EPIC describing what you are planning to work on first.

---

## Slack Community

Join our community on the CNCF Slack to collaborate, ask questions, and contribute:

[CNCF Slack Invite](https://communityinviter.com/apps/cloud-native/cncf)

Find us in the `#serverless-workflow-sdk` channel.

---

## Contributing

Your contributions are very welcome!

### Code Style

- Format imports with `goimports`.
- Run static analysis using:

```shell
make lint
```

Automatically fix lint issues:

```shell
make lint params=--fix
```

### EditorConfig

A sample `.editorconfig` for IntelliJ or GoLand users can be found [here](contrib/intellij.editorconfig).

### Known Issues

- **MacOS Issue**: If you encounter `goimports: can't extract issues from gofmt diff output`, resolve it with:

```shell
brew install diffutils
```

---

Contributions are greatly appreciated! Check [this EPIC](https://github.com/serverlessworkflow/sdk-go/issues/221) and contribute to completing more features.

Happy coding!
