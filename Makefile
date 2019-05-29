HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

PREFIX ?= katib
CMD_PREFIX ?= cmd

# Run tests
.PHONY: test
test:
	go test ./pkg/... ./cmd/... -coverprofile coverage.out

depend:
ifndef HAS_DEP
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -v

check: fmt vet lint

fmt: depend generate
	hack/verify-gofmt.sh

lint: depend generate
ifndef HAS_LINT
	go get -u golang.org/x/lint/golint
	echo "installing golint"
endif
	hack/verify-golint.sh

vet: depend generate
	go vet ./pkg/... ./cmd/...

update:
	hack/update-gofmt.sh

# Deploy katib v1alpha2 manifests into a k8s cluster
deploy: 
	bash scripts/v1alpha2/deploy.sh

# Undeploy katib v1alpha2 manifests into a k8s cluster
undeploy:
	bash scripts/v1alpha2/undeploy.sh

# Build Katib images for v1alpha1
buildv1alpha1:
	bash scripts/v1alpha1/build.sh

# Deploy katib v1alpha1 manifests into a k8s cluster
deployv1alpha1:
	bash scripts/v1alpha1/deploy.sh

# Undeploy katib v1alpha1 manifests into a k8s cluster
undeployv1alpha1:
	bash scripts/v1alpha1/undeploy.sh

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...

############################################################
# Build docker image section for v1alpha2
############################################################
images: katib-controller katib-manager katib-manager-rest metrics-collector katib-ui tfevent-metrics-collector suggestion-random

katib-controller: depend generate
	docker build -t ${PREFIX}/v1alpha2/katib-controller -f ${CMD_PREFIX}/katib-controller/v1alpha2/Dockerfile .

katib-manager: depend generate
	docker build -t ${PREFIX}/v1alpha2/katib-manager -f ${CMD_PREFIX}/manager/v1alpha2/Dockerfile .

katib-manager-rest: depend generate
	docker build -t ${PREFIX}/v1alpha2/katib-manager-rest -f ${CMD_PREFIX}/manager-rest/v1alpha2/Dockerfile .

metrics-collector: depend generate
	docker build -t ${PREFIX}/v1alpha2/metrics-collector -f ${CMD_PREFIX}/metricscollector/v1alpha2/Dockerfile .

katib-ui: depend generate
	docker build -t ${PREFIX}/v1alpha2/katib-ui -f ${CMD_PREFIX}/ui/v1alpha2/Dockerfile .

tfevent-metrics-collector: depend generate
	docker build -t ${PREFIX}/v1alpha2/tfevent-metrics-collector -f ${CMD_PREFIX}/tfevent-metricscollector/v1alpha2/Dockerfile .

suggestion-random: depend generate
	docker build -t ${PREFIX}/v1alpha2/suggestion-random -f ${CMD_PREFIX}/suggestion/random/v1alpha2/Dockerfile .
