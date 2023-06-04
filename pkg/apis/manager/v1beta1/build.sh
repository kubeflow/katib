#!/usr/bin/env bash

# Copyright 2022 The Kubeflow Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
set -e

cd "$(dirname "$0")"

proto="api.proto"

# Get host paths for kubernetes api modules
GO_MOD_K8S_API=$(go list -m -f '{{.Dir}}' k8s.io/api)
GO_MOD_K8S_APIMACHINERY=$(go list -m -f '{{.Dir}}' k8s.io/apimachinery)

docker run -i --rm \
	-v "$PWD:$PWD" \
	-v "$GO_MOD_K8S_API:$GOPATH/pkg/mod/k8s.io/api" \
	-v "$GO_MOD_K8S_APIMACHINERY:$GOPATH/pkg/mod/k8s.io/apimachinery" \
	-w "$PWD" \
	znly/protoc --python_out=plugins=grpc:./python --go_out=plugins=grpc:. -I="$GOPATH/pkg/mod/" -I. $proto
docker run -i --rm \
	-v "$PWD:$PWD" \
	-v "$GO_MOD_K8S_API:$GOPATH/pkg/mod/k8s.io/api" \
	-v "$GO_MOD_K8S_APIMACHINERY:$GOPATH/pkg/mod/k8s.io/apimachinery" \
	-w "$PWD" \
	znly/protoc --plugin=protoc-gen-grpc=/usr/bin/grpc_python_plugin --python_out=./python --grpc_out=./python -I "$GOPATH/pkg/mod/" -I. $proto

docker build -t protoc-gen-doc gen-doc/
docker run --rm \
	-v "$PWD/gen-doc:/out" \
	-v "$PWD:/apiprotos" \
	-v "$GO_MOD_K8S_API:$GOPATH/pkg/mod/k8s.io/api" \
	-v "$GO_MOD_K8S_APIMACHINERY:$GOPATH/pkg/mod/k8s.io/apimachinery" \
	protoc-gen-doc --doc_opt=markdown,api.md -I "$GOPATH/pkg/mod/" -I /protobuf -I /apiprotos $proto
docker run --rm \
	-v "$PWD/gen-doc:/out" \
	-v "$PWD:/apiprotos" \
	-v "$GO_MOD_K8S_API:$GOPATH/pkg/mod/k8s.io/api" \
	-v "$GO_MOD_K8S_APIMACHINERY:$GOPATH/pkg/mod/k8s.io/apimachinery" \
	protoc-gen-doc --doc_opt=html,index.html -I "$GOPATH/pkg/mod/" -I /protobuf -I /apiprotos $proto
