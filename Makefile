.PHONY: gen-types
gen-types:
	./hack/generate-spec-types.sh

addheaders:
	@which addlicense > /dev/null || go get -u github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .

test:
	make addheaders
	go vet ./...
	go fmt ./...
	go test ./...