
#!/bin/bash
#
# A simple script to build the Docker images.
# This is intended to be invoked as a step in Argo to build the docker image.
#
# build_image.sh ${DOCKERFILE} ${IMAGE} ${TAG} (optional | ${TEST_REGISTRY})
set -ex

DOCKERFILE=$1
# Context directory should be root directory so we pull in all sources.
CONTEXT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/../../ >/dev/null && pwd )"
IMAGE=$2
TAG=$3

# Wait for the Docker daemon to be available.
until docker ps
do sleep 3
done

docker build -t ${TAG} -f ${DOCKERFILE} ${CONTEXT_DIR}

gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}
gcloud docker -- push "${IMAGE}:${TAG}"
docker tag "${IMAGE}:${TAG}" "${IMAGE}:latest"
gcloud docker -- push "${IMAGE}:latest"