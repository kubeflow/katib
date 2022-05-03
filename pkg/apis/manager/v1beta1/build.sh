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

docker run -i --rm -v "$PWD:$PWD" -w "$PWD" znly/protoc --python_out=plugins=grpc:./python --go_out=plugins=grpc:. -I. $proto
docker run -i --rm -v "$PWD:$PWD" -w "$PWD" znly/protoc --plugin=protoc-gen-grpc=/usr/bin/grpc_python_plugin --python_out=./python --grpc_out=./python -I. $proto

docker build -t protoc-gen-doc gen-doc/
docker run --rm -v "$PWD/gen-doc:/out" -v "$PWD:/apiprotos" protoc-gen-doc --doc_opt=markdown,api.md -I /protobuf -I /apiprotos $proto
docker run --rm -v "$PWD/gen-doc:/out" -v "$PWD:/apiprotos" protoc-gen-doc --doc_opt=html,index.html -I /protobuf -I /apiprotos $proto
