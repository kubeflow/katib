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

# Build Katib images
build: 
	bash scripts/v1alpha1/build.sh

# Deploy katib manifests into a k8s cluster
deploy: 
	bash scripts/v1alpha1/deploy.sh

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...
