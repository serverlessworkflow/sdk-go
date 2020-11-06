.PHONY: gen-types
gen-types:
	./hack/generate-spec-types.sh

addheaders:
	@which addlicense > /dev/null || go get -u github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .

fmt:
	go vet ./...
	go fmt ./...

lint:
	go mod tidy
	make addheaders
	make fmt
	./hack/go-lint.sh

.PHONY: test
coverage="false"
test:
	make lint
	./hack/go-test.sh
