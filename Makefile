HAS_LINT := $(shell command -v golangci-lint;)
HAS_YAMLLINT := $(shell command -v yamllint;)
HAS_SHELLCHECK := $(shell command -v shellcheck;)
HAS_SETUP_ENVTEST := $(shell command -v setup-envtest;)
HAS_MOCKGEN := $(shell command -v mockgen;)

COMMIT := v1beta1-$(shell git rev-parse --short=7 HEAD)
KATIB_REGISTRY := docker.io/kubeflowkatib
CPU_ARCH ?= linux/amd64,linux/arm64
ENVTEST_K8S_VERSION ?= 1.31
MOCKGEN_VERSION ?= $(shell grep 'go.uber.org/mock' go.mod | cut -d ' ' -f 2)
GO_VERSION=$(shell grep '^go' go.mod | cut -d ' ' -f 2)
GOPATH ?= $(shell go env GOPATH)

TEST_TENSORFLOW_EVENT_FILE_PATH ?= $(CURDIR)/test/unit/v1beta1/metricscollector/testdata/tfevent-metricscollector/logs

# Run tests
.PHONY: test
test: envtest
	KUBEBUILDER_ASSETS="$(shell setup-envtest use $(ENVTEST_K8S_VERSION) -p path)" go test ./pkg/... ./cmd/... -coverprofile coverage.out

envtest:
ifndef HAS_SETUP_ENVTEST
	go install sigs.k8s.io/controller-runtime/tools/setup-envtest@release-0.19
	$(info "setup-envtest has been installed")
endif
	$(info "setup-envtest has already installed")

check: generated-codes go-mod fmt vet lint

fmt:
	hack/verify-gofmt.sh

lint:
ifndef HAS_LINT
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.57.2
	$(info "golangci-lint has been installed")
endif
	hack/verify-golangci-lint.sh

yamllint:
ifndef HAS_YAMLLINT
	pip install --prefer-binary yamllint
	$(info "yamllint has been installed")
endif
	hack/verify-yamllint.sh

vet:
	go vet ./pkg/... ./cmd/...

shellcheck:
ifndef HAS_SHELLCHECK
	bash hack/install-shellcheck.sh
	$(info "shellcheck has been installed")
endif
	hack/verify-shellcheck.sh

update:
	hack/update-gofmt.sh

# Deploy Katib v1beta1 manifests using Kustomize into a k8s cluster.
deploy:
	bash scripts/v1beta1/deploy.sh $(WITH_DATABASE_TYPE)

# Undeploy Katib v1beta1 manifests using Kustomize from a k8s cluster
undeploy:
	bash scripts/v1beta1/undeploy.sh

generated-codes: generate
ifneq ($(shell bash hack/verify-generated-codes.sh '.'; echo $$?),0)
	$(error 'Please run "make generate" to generate codes')
endif

go-mod: sync-go-mod
ifneq ($(shell bash hack/verify-generated-codes.sh 'go.*'; echo $$?),0)
	$(error 'Please run "go mod tidy -go $(GO_VERSION)" to sync Go modules')
endif

sync-go-mod:
	go mod tidy -go $(GO_VERSION)

.PHONY: go-mod-download
go-mod-download:
	go mod download

CONTROLLER_GEN = $(shell pwd)/bin/controller-gen
.PHONY: controller-gen
controller-gen:
	@GOBIN=$(shell pwd)/bin GO111MODULE=on go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.16.5

# Run this if you update any existing controller APIs.
# 1. Generate deepcopy, clientset, listers, informers for the APIs (hack/update-codegen.sh)
# 2. Generate open-api for the APIs (hack/update-openapigen)
# 3. Generate Python SDK for Katib (hack/gen-python-sdk/gen-sdk.sh)
# 4. Generate gRPC manager APIs (pkg/apis/manager/v1beta1/build.sh and pkg/apis/manager/health/build.sh)
# 5. Generate Go mock codes
generate: go-mod-download controller-gen
ifndef HAS_MOCKGEN
	go install go.uber.org/mock/mockgen@$(MOCKGEN_VERSION)
	$(info "mockgen has been installed")
endif
	go generate ./pkg/... ./cmd/...
	hack/gen-python-sdk/gen-sdk.sh
	hack/update-proto.sh
	hack/update-mockgen.sh

# Build images for the Katib v1beta1 components.
build: generate
ifeq ($(and $(REGISTRY),$(TAG),$(CPU_ARCH)),)
	$(error REGISTRY and TAG must be set. Usage: make build REGISTRY=<registry> TAG=<tag> CPU_ARCH=<cpu-architecture>)
endif
	bash scripts/v1beta1/build.sh $(REGISTRY) $(TAG) $(CPU_ARCH)

# Build and push Katib images from the latest master commit.
push-latest: generate
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) latest $(CPU_ARCH)
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) $(COMMIT) $(CPU_ARCH)
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) latest
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) $(COMMIT)

# Build and push Katib images for the given tag.
push-tag:
ifeq ($(TAG),)
	$(error TAG must be set. Usage: make push-tag TAG=<release-tag>)
endif
	bash scripts/v1beta1/build.sh $(KATIB_REGISTRY) $(TAG) $(CPU_ARCH)
	bash scripts/v1beta1/push.sh $(KATIB_REGISTRY) $(TAG)

# Release a new version of Katib.
release:
ifeq ($(and $(BRANCH),$(TAG)),)
	$(error BRANCH and TAG must be set. Usage: make release BRANCH=<branch> TAG=<tag>)
endif
	bash scripts/v1beta1/release.sh $(BRANCH) $(TAG)

# Update all Katib images.
update-images:
ifeq ($(and $(OLD_PREFIX),$(NEW_PREFIX),$(TAG)),)
	$(error OLD_PREFIX, NEW_PREFIX, and TAG must be set. \
	Usage: make update-images OLD_PREFIX=<old-prefix> NEW_PREFIX=<new-prefix> TAG=<tag> \
	For more information, check this file: scripts/v1beta1/update-images.sh)
endif
	bash scripts/v1beta1/update-images.sh $(OLD_PREFIX) $(NEW_PREFIX) $(TAG)

# Prettier UI format check for Katib v1beta1.
prettier-check:
	npm run format:check --prefix pkg/ui/v1beta1/frontend

# Update boilerplate for the source code.
update-boilerplate:
	./hack/boilerplate/update-boilerplate.sh

prepare-pytest:
	pip install --prefer-binary -r test/unit/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/hyperopt/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/optuna/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/hyperband/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/nas/enas/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/nas/darts/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/suggestion/pbt/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/earlystopping/medianstop/v1beta1/requirements.txt
	pip install --prefer-binary -r cmd/metricscollector/v1beta1/tfevent-metricscollector/requirements.txt
	# The sqlalchemy on which optuna depends requires typing-extensions>=4.6.0.
	# REF: https://github.com/kubeflow/katib/pull/2251
	# TODO (tenzen-y): Once we upgrade libraries depended on typing-extensions==4.5.0, we can remove this line.
	pip install typing-extensions==4.6.3

prepare-pytest-testdata:
ifeq ("$(wildcard $(TEST_TENSORFLOW_EVENT_FILE_PATH))", "")
	python examples/v1beta1/trial-images/tf-mnist-with-summaries/mnist.py --epochs 5 --batch-size 200 --log-path $(TEST_TENSORFLOW_EVENT_FILE_PATH)
endif

# TODO(Electronic-Waste): Remove the import rewrite when protobuf supports `python_package` option.
# REF: https://github.com/protocolbuffers/protobuf/issues/7061
pytest: prepare-pytest prepare-pytest-testdata
	pytest ./test/unit/v1beta1/suggestion
	pytest ./test/unit/v1beta1/earlystopping
	pytest ./test/unit/v1beta1/metricscollector
	cp ./pkg/apis/manager/v1beta1/python/api_pb2.py ./sdk/python/v1beta1/kubeflow/katib/katib_api_pb2.py
	cp ./pkg/apis/manager/v1beta1/python/api_pb2_grpc.py ./sdk/python/v1beta1/kubeflow/katib/katib_api_pb2_grpc.py
	sed -i "s/api_pb2/kubeflow\.katib\.katib_api_pb2/g" ./sdk/python/v1beta1/kubeflow/katib/katib_api_pb2_grpc.py
	pytest ./sdk/python/v1beta1/kubeflow/katib
	rm ./sdk/python/v1beta1/kubeflow/katib/katib_api_pb2.py ./sdk/python/v1beta1/kubeflow/katib/katib_api_pb2_grpc.py

