addheaders:
	@command -v addlicense > /dev/null || go install -modfile=tools.mod -v github.com/google/addlicense
	@addlicense -c "The Serverless Workflow Specification Authors" -l apache .

fmt:
	@go vet ./...
	@go fmt ./...

goimports:
	@command -v goimports > /dev/null || go install golang.org/x/tools/cmd/goimports@latest
	@goimports -w .


lint:
	@command -v golangci-lint > /dev/null || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "${GOPATH}/bin"
	make addheaders
	make goimports
	make fmt
	./hack/go-lint.sh ${params}

.PHONY: test
coverage="false"

test: deepcopy buildergen
	make lint
	@go test ./...

.PHONY: deepcopy buildergen
deepcopy: $(DEEPCOPY_GEN) ## Download deepcopy-gen locally if necessary.
	./hack/deepcopy-gen.sh deepcopy

buildergen: $(BUILDER_GEN) ## Download builder-gen locally if necessary.
	./hack/builder-gen.sh buildergen

.PHONY: kube-integration
kube-integration: controller-gen
	$(CONTROLLER_GEN) rbac:roleName=manager-role crd:allowDangerousTypes=true webhook paths="./..." output:crd:artifacts:config=config/crd/bases


####################################
# install controller-gen tool
## Location to install dependencies to
LOCALBIN ?= $(shell pwd)/bin
$(LOCALBIN):
	mkdir -p $(LOCALBIN)

CONTROLLER_TOOLS_VERSION ?= v0.16.3
CONTROLLER_GEN ?= $(LOCALBIN)/controller-gen
.PHONY: controller-gen
controller-gen: $(CONTROLLER_GEN) ## Download controller-gen locally if necessary.
$(CONTROLLER_GEN): $(LOCALBIN)
	test -s $(LOCALBIN)/controller-gen || GOBIN=$(LOCALBIN) go install sigs.k8s.io/controller-tools/cmd/controller-gen@$(CONTROLLER_TOOLS_VERSION)

