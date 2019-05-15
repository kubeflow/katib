#!/bin/bash
#
# A simple script to build the Docker images.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG} {ROOT_DIR}
set -ex

DOCKERFILE=$1
CONTEXT_DIR=$(dirname "$DOCKERFILE")
IMAGE=$2
TAG=$3
ROOT_DIR=$4
export GOPATH=${ROOT_DIR}
GO_DIR=${GOPATH}/src/github.com/kubeflow/pytorch-operator
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

export PATH=${GOPATH}/bin:/usr/local/go/bin:${PATH}
echo "Create symlink to GOPATH"
mkdir -p ${GOPATH}/src/github.com/kubeflow
ln -s ${CONTEXT_DIR} ${GO_DIR}
cd ${GO_DIR}
echo "Build pytorch operator v1beta1 binary"
go build github.com/kubeflow/pytorch-operator/cmd/pytorch-operator.v1beta1
echo "Build pytorch operator v1beta2 binary"
go build github.com/kubeflow/pytorch-operator/cmd/pytorch-operator.v1beta2

echo "Building container in gcloud"
gcloud builds submit . --tag=${IMAGE}:${TAG}
echo "Image built successfully"
