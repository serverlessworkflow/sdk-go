addheaders:
	@command -v addlicense > /dev/null || go install -modfile=tools.mod -v github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .

fmt:
	@go vet ./...
	@go fmt ./...

lint:
	@command -v golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin"
	make addheaders
	make fmt
	./hack/go-lint.sh ${params}

.PHONY: test
coverage="false"
test: deepcopy
	make lint
	@go test ./...

.PHONY: deepcopy
deepcopy: $(DEEPCOPY_GEN) ## Download deeepcopy-gen locally if necessary.
	./hack/deepcopy-gen.sh deepcopy