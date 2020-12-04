.PHONY: gen-types
gen-types:
	./hack/generate-spec-types.sh

addheaders:
<<<<<<< HEAD
	@command -v addlicense > /dev/null || go install -modfile=tools.mod -v github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .
=======
>>>>>>> changed the path of addlicense binary and also clean the go package installation

fmt:
	@go vet ./...
	@go fmt ./...

lint:
	@command -v golint > /dev/null || go install -modfile=tools.mod -v golang.org/x/lint/golint
	@command -v golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin"
	make addheaders
	make fmt
	./hack/go-lint.sh

.PHONY: test
coverage="false"
test:
	make lint
	@go test ./...
