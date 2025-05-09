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

test: deepcopy buildergen
	@echo "🧪 Running tests..."
	@go test ./...
	@echo "✅  Tests completed!"

deepcopy: $(DEEPCOPY_GEN) ## Download deepcopy-gen locally if necessary.
	@echo "📦 Running deepcopy-gen..."
	@./hack/deepcopy-gen.sh deepcopy > /dev/null
	@make lint
	@echo "✅  Deepcopy generation and linting completed!"

buildergen: $(BUILDER_GEN) ## Download builder-gen locally if necessary.
	@echo "📦 Running builder-gen..."
	@./hack/builder-gen.sh buildergen > /dev/null
	@make lint
	@echo "✅  Builder generation and linting completed!"

.PHONY: kube-integration
kube-integration: controller-gen
	@echo "📦 Generating Kubernetes objects..."
	@$(CONTROLLER_GEN) object:headerFile="./hack/boilerplate.txt" paths="./kubernetes/api/..."
	@echo "📦 Generating Kubernetes CRDs..."
	@$(CONTROLLER_GEN) rbac:roleName=manager-role crd:allowDangerousTypes=true webhook paths="./kubernetes/..." output:crd:artifacts:config=config/crd/bases
	@make lint
	@echo "✅  Kubernetes integration completed!"


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

