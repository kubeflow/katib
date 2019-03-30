# Run tests
test:
	go test ./pkg/... ./cmd/...

build: 
	sh scripts/build.sh

# deploy katib manifests into a k8s cluster
deploy: 
	sh scripts/deploy.sh

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
