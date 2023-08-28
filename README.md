# Go SDK for Serverless Workflow
Here you will find all the [specification types](https://github.com/serverlessworkflow/specification/blob/main/schema/workflow.json) defined by our Json Schemas, in Go.

Table of Contents
=================

- [Status](#status)
- [Releases](#releases)
- [How to Use](#how-to-use)
  - [Parsing Serverless Workflow files](#parsing-serverless-workflow-files)
- [Slack Channel](#slack-channel)
- [Contributors Guide](#contributors-guide)
  - [Code Style](#code-style)
  - [EditorConfig](#editorconfig)
  - [Known Issues](#known-issues)


## Status
Current status of features implemented in the SDK is listed in the table below:

| Feature                                     | Status             |
|-------------------------------------------- | ------------------ |
| Parse workflow JSON and YAML definitions    | :heavy_check_mark: | 
| Programmatically build workflow definitions | :no_entry_sign:    |
| Validate workflow definitions (Schema)      | :heavy_check_mark: |
| Validate workflow definitions (Integrity)   | :heavy_check_mark: |
| Generate workflow diagram (SVG)             | :no_entry_sign:    |


## Releases
|                              Latest Releases                               | Conformance to spec version |
|:--------------------------------------------------------------------------:| :---: |
| [v1.0.0](https://github.com/serverlessworkflow/sdk-go/releases/tag/v1.0.0) | [v0.5](https://github.com/serverlessworkflow/specification/tree/0.5.x) |
| [v2.0.1](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.0.1) | [v0.6](https://github.com/serverlessworkflow/specification/tree/0.6.x) |
| [v2.1.2](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.1.2) | [v0.7](https://github.com/serverlessworkflow/specification/tree/0.7.x) |
| [v2.2.4](https://github.com/serverlessworkflow/sdk-go/releases/tag/v2.2.4) | [v0.8](https://github.com/serverlessworkflow/specification/tree/0.8.x) |

## How to use

Run the following command in the root of your Go's project:

```shell script
$ go get github.com/serverlessworkflow/sdk-go/v2
```

Your `go.mod` file should be updated to add a dependency from the Serverless Workflow specification.

To use the generated types, import the package in your go file like this:

```go
import "github.com/serverlessworkflow/sdk-go/v2/model"
```

Then just reference the package in your Go file like `myfunction := model.Function{}`.

### Parsing Serverless Workflow files

Serverless Workflow Specification supports YAML and JSON files for Workflow definitions.
To transform such files into a Go data structure, use:

```go
package sw

import (
    "github.com/serverlessworkflow/sdk-go/v2/model"
    "github.com/serverlessworkflow/sdk-go/v2/parser"
)

func ParseWorkflow(filePath string) (*model.Workflow, error) {
    workflow, err := parser.FromFile(filePath)
    if err != nil {
        return nil, err
    } 
    return workflow, nil
} 
```

The `Workflow` structure then can be used in your application. 

## Slack Channel

Join us at [CNCF Slack](https://communityinviter.com/apps/cloud-native/cncf), channel `#serverless-workflow-sdk` and say hello 🙋.

## Contributors Guide

This guide aims to guide newcomers to getting started with the project standards.


### Code Style

For this project we use basically the default configuration for most used IDEs.
For the configurations below, make sure to properly configure your IDE:

- **imports**: goimports

This should be enough to get you started.

If you are unsure that your IDE is not correctly configured, you can run the lint checks:

```bash
make lint
```

If something goes wrong, the error will be printed, e.g.:
```bash
$ make lint
make addheaders
make fmt
./hack/go-lint.sh
util/floatstr/floatstr_test.go:19: File is not `goimports`-ed (goimports)
        "k8s.io/apimachinery/pkg/util/yaml"
make: *** [lint] Error 1
```

Lint issues can be fixed with the `--fix` flag, this command can be used:
```bash
make lint params=--fix
```


### EditorConfig
For IntelliJ you can find an example `editorconfig` file [here](contrib/intellij.editorconfig). To use it please visit
the Jetbrains [documentation](https://www.jetbrains.com/help/idea/editorconfig.html).


### Known Issues

On MacOSX/darwin you might get this issue:
```
 goimports: can't extract issues from gofmt diff output
```
To solve install the `diffutils` package:

```bash
 brew install diffutils
```

