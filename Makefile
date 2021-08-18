HAS_LINT := $(shell command -v golint;)
COMMIT := v1beta1-$(shell git rev-parse --short=7 HEAD)
KATIB_REGISTRY := docker.io/kubeflowkatib

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

# Deploy Katib v1beta1 manifests using Kustomize into a k8s cluster.
deploy:
	bash scripts/v1beta1/deploy.sh

# Undeploy Katib v1beta1 manifests using Kustomize from a k8s cluster
undeploy:
	bash scripts/v1beta1/undeploy.sh

# Run this if you update any existing controller APIs.
# 1. Genereate deepcopy, clientset, listers, informers for the APIs (hack/update-codegen.sh)
# 2. Generate open-api for the APIs (hack/update-openapigen)
# 3. Generate Python SDK for Katib (hack/gen-python-sdk/gen-sdk.sh)
# 4. Generate gRPC manager APIs (pkg/apis/manager/v1beta1/build.sh and pkg/apis/manager/health/build.sh)
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...
	hack/gen-python-sdk/gen-sdk.sh
	cd ./pkg/apis/manager/v1beta1 && ./build.sh
	cd ./pkg/apis/manager/health && ./build.sh

# Build images for the Katib v1beta1 components.
build: generate
ifeq ($(and $(REGISTRY),$(TAG)),)
	$(error REGISTRY and TAG must be set. Usage: make build REGISTRY=<registry> TAG=<tag>)
endif
	bash scripts/v1beta1/build.sh $(REGISTRY) $(TAG)

# Build and push Katib images from the latest master commit.
push-latest: generate
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) latest
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) $(COMMIT)
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) latest
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) $(COMMIT)

# Build and push Katib images for the given tag.
push-tag: generate
ifeq ($(TAG),)
	$(error TAG must be set. Usage: make push-tag TAG=<release-tag>)
endif
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) $(TAG)
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) $(COMMIT)
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) $(TAG)
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) $(COMMIT)

# Release a new version of Katib.
release:
ifeq ($(and $(BRANCH),$(TAG)),)
	$(error BRANCH and TAG must be set. Usage: make release BRANCH=<branch> TAG=<tag>)
endif
	bash scripts/v1beta1/release.sh $(BRANCH) $(TAG)

# Prettier UI format check for Katib v1beta1.
prettier-check:
	npm run format:check --prefix pkg/new-ui/v1beta1/frontend

# Update boilerplate for the source code.
update-boilerplate:
	./hack/boilerplate/update-boilerplate.sh
