# Go SDK for Serverless Workflow

Here you will find all the [specification types](https://github.com/serverlessworkflow/specification/blob/master/schema/workflow.json) defined by our Json Schemas, in Go.

Current status of features implemented in the SDK is listed in the table below:

| Feature                                     | Status             |
|-------------------------------------------- | ------------------ |
| Parse workflow JSON and YAML definitions    | :heavy_check_mark: | 
| Programmatically build workflow definitions | :no_entry_sign:    |
| Validate workflow definitions (Schema)      | :heavy_check_mark: |
| Validate workflow definitions (Integrity)   | :no_entry_sign:    |
| Generate workflow diagram (SVG)             | :no_entry_sign:    |

## How to use

Run the following command in the root of your Go's project:

```shell script
$ go get github.com/serverlessworkflow/sdk-go
```

Your `go.mod` file should be updated to add a dependency from the Serverless Workflow specification.

To use the generated types, import the package in your go file like this:

```go
import "github.com/serverlessworkflow/sdk-go/model"
```

Then just reference the package in your Go file like `myfunction := model.Function{}`.

### Parsing Serverless Workflow files

Serverless Workflow Specification supports YAML and JSON files for Workflow definitions.
To transform such files into a Go data structure, use:

```go
package sw

import (

"github.com/serverlessworkflow/sdk-go/model"
"github.com/serverlessworkflow/sdk-go/parser"
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