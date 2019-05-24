HAS_DEP := $(shell command -v dep;)
HAS_LINT := $(shell command -v golint;)

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

# Build Katib images for v1alpha2
build: 
	bash scripts/v1alpha2/build.sh

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
