addheaders:
	@command -v addlicense > /dev/null || (echo "🚀 Installing addlicense..."; go install -modfile=tools.mod -v github.com/google/addlicense)
	@addlicense -c "The Open Workflow Specification Authors" -l apache .

fmt:
	@go vet ./...
	@go fmt ./...

goimports:
	@command -v goimports > /dev/null || (echo "🚀 Installing goimports..."; go install golang.org/x/tools/cmd/goimports@latest)
	@goimports -w .

lint:
	@echo "🚀 Installing/updating golangci-lint…"
	GO111MODULE=on go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

	@echo "🚀 Running lint…"
	@make addheaders
	@make goimports
	@make fmt
	@$(GOPATH)/bin/golangci-lint run ./... ${params}
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