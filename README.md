# Go SDK for Serverless Workflow

Here you will find all the [specification types](https://github.com/serverlessworkflow/specification/blob/master/schema/workflow.json) defined by our Json Schemas, in Go.

Some types defined by the specification can be generic objects (such as [`Extensions`](https://github.com/serverlessworkflow/specification/blob/master/extending/README.md)) 
or share a minimum interface, like [`States`](https://github.com/serverlessworkflow/specification/tree/master/README.md#State-Definition). 

In cases like these, we've decided to represent their types as the Kubernetes type [`RawExtension`](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/#rawextension),
that can be found on `kuberbetes/serverlessworkflow` package. 

This way the Serverless Workflow types can be Kubernetes friendly, to make it easy for one to
use the package when developing Kubernetes applications.

## How to use

Run the following command in the root of your Go's project:

```shell script
$ go get -u github.com/serverlessworkflow/sdk-go
```

Your `go.mod` file should be updated to add a dependency from the Serverless Workflow specification.

To use the generated types, import the package in your go file like this:

```go
package mypackage

import "github.com/serverlessworkflow/sdk-go/pkg/apis/sw"
```

Then just reference the package within your file like `myfunction := sw.Function{}`.

If you wish to use the Kubernetes friendly types, import the `kubernetes` package instead:

```go
package mypackage

import "github.com/serverlessworkflow/sdk-go/pkg/apis/kubernetes/sw"
```
