# Run tests
test:
	go test ./pkg/... ./cmd/...

# Build Katib images
build: 
	sh scripts/v1alpha1/build.sh

# Deploy katib manifests into a k8s cluster
deploy: 
	sh scripts/v1alpha1/deploy.sh

# Run go fmt against code
fmt:
	go fmt ./pkg/... ./cmd/...

# Run go vet against code
vet:
	go vet ./pkg/... ./cmd/...

# Generate code
generate:
ifndef GOPATH
	$(error GOPATH not defined, please define GOPATH. Run "go help gopath" to learn more about GOPATH)
endif
	go generate ./pkg/... ./cmd/...
