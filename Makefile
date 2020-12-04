.PHONY: gen-types
gen-types:
	./hack/generate-spec-types.sh

addheaders:
	@command -v addlicense > /dev/null || go install -v github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .
	@go mod tidy

fmt:
	@go vet ./...
	@go fmt ./...

lint:
	@command -v golint > /dev/null || go install -v golang.org/x/lint/golint
	@command -v golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin"
	@go mod tidy
	make addheaders
	make fmt
	./hack/go-lint.sh

.PHONY: test
coverage="false"
test:
	make lint
	@go test ./...
