HAS_LINT := $(shell command -v golint;)

# Run tests
.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile coverage.out

check: generate fmt vet lint

fmt:
	hack/verify-gofmt.sh

lint:
ifndef HAS_LINT
	go get -u golang.org/x/lint/golint
	echo "installing golint"
endif
	hack/verify-golint.sh

vet:
	go vet ./pkg/... ./cmd/...

update:
	hack/update-gofmt.sh

# Deploy Katib v1beta1 manifests into a k8s cluster
deploy:
	bash scripts/v1beta1/deploy.sh

# Undeploy Katib v1beta1 manifests from a k8s cluster
undeploy:
	bash scripts/v1beta1/undeploy.sh

# Generate deepcopy, clientset, listers, informers, open-api and python SDK for APIs.
# Run this if you update any existing controller APIs.
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...
	hack/gen-python-sdk/gen-sdk.sh

# Build images for Katib v1beta1 components
build: generate
ifeq ($(and $(REGISTRY),$(TAG)),)
	$(error REGISTRY and TAG must be set. Usage make build REGISTRY=<registry> TAG=<TAG>)
endif
	bash scripts/v1beta1/build.sh -r $(REGISTRY) -t $(TAG)

# Prettier UI format check for Katib v1beta1
prettier-check:
	npm run format:check --prefix pkg/ui/v1beta1/frontend
