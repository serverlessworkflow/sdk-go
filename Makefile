addheaders:
	@command -v addlicense > /dev/null || (echo "🚀 Installing addlicense..."; go install -modfile=tools.mod -v github.com/google/addlicense)
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .

fmt:
	@go vet ./...
	@go fmt ./...

goimports:
	@command -v goimports > /dev/null || (echo "🚀 Installing goimports..."; go install golang.org/x/tools/cmd/goimports@latest)
	@goimports -w .

lint:
	@echo "🚀 Running lint..."
	@command -v golangci-lint > /dev/null || (echo "🚀 Installing golangci-lint..."; curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin")
	@make addheaders
	@make goimports
	@make fmt
	@./hack/go-lint.sh ${params}
	@echo "✅  Linting completed!"

.PHONY: test
coverage="false"

test:
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅  Tests completed!"

.PHONY: integration-test

integration-test:
	@echo "🔄 Running integration tests..."
	@./hack/integration-test.sh
	@echo "✅  Integration tests completed!"