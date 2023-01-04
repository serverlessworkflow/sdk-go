# Contributors Guide

This guide aims to guide newcomers to getting started with the project standards.


## Code Style

For this project we use basically the default configuration for most used IDEs.
For the configurations below, make sure to properly configure your IDE:

- **imports**: goimports

For IntelliJ you can find an example `editorconfig` file [here](contrib/intellij.editorconfig).

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

## Known Issues

On MacOSX/darwin you might get this issue:

```
 goimports: can't extract issues from gofmt diff output
```
To solve install the `diffutils` package:

```bash
 brew install diffutils
```