#!/bin/bash

# Copyright 2018 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This shell script is used to build a cluster and create a namespace from our
# argo workflow


set -o errexit
set -o nounset
set -o pipefail

CLUSTER_NAME="${CLUSTER_NAME}"
ZONE="${GCP_ZONE}"
PROJECT="${GCP_PROJECT}"
NAMESPACE="${DEPLOY_NAMESPACE}"
REGISTRY="${GCP_REGISTRY}"
VERSION=$(git describe --tags --always --dirty)
GO_DIR=${GOPATH}/src/github.com/${REPO_OWNER}/${REPO_NAME}

echo "Activating service-account"
gcloud auth activate-service-account --key-file=${GOOGLE_APPLICATION_CREDENTIALS}

echo "Configuring kubectl"

echo "CLUSTER_NAME: ${CLUSTER_NAME}"
echo "ZONE: ${GCP_ZONE}"
echo "PROJECT: ${GCP_PROJECT}"

gcloud --project ${PROJECT} container clusters get-credentials ${CLUSTER_NAME} \
  --zone ${ZONE}
kubectl config set-context $(kubectl config current-context) --namespace=default
USER=`gcloud config get-value account`

kubectl apply -f - << EOF
kind: ClusterRoleBinding
apiVersion: rbac.authorization.k8s.io/v1
metadata:
  name: cluster-admins
subjects:
- kind: User
  name: $USER
roleRef:
  kind: ClusterRole
  name: cluster-admin
  apiGroup: ""
EOF

echo "Install Katib "
echo "REGISTRY ${REGISTRY}"
echo "REPO_NAME ${REPO_NAME}"
echo "VERSION ${VERSION}"

sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/vizier-core@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/vizier-core:${VERSION}@" manifests/v1alpha1/vizier/core/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/vizier-core-rest@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/vizier-core-rest:${VERSION}@" manifests/v1alpha1/vizier/core-rest/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/katib-ui@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/katib-ui:${VERSION}@" manifests/v1alpha1/vizier/ui/deployment.yaml
sed -i -e "s@type: NodePort@type: ClusterIP@" -e "/nodePort: 30678/d" manifests/v1alpha1/vizier/core/service.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/studyjob-controller@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/studyjob-controller:${VERSION}@" manifests/v1alpha1/studyjobcontroller/studyjobcontroller.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-random@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-random:${VERSION}@" manifests/v1alpha1/vizier/suggestion/random/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-grid@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-grid:${VERSION}@" manifests/v1alpha1/vizier/suggestion/grid/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-hyperband@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-hyperband:${VERSION}@" manifests/v1alpha1/vizier/suggestion/hyperband/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-bayesianoptimization@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-bayesianoptimization:${VERSION}@" manifests/v1alpha1/vizier/suggestion/bayesianoptimization/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-nasrl@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-nasrl:${VERSION}@" manifests/v1alpha1/vizier/suggestion/nasrl/deployment.yaml
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/suggestion-nasenvelopenet@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/suggestion-nasenvelopenet:${VERSION}@" manifests/v1alpha1/vizier/suggestion/nasenvelopenet/deployment.yaml

sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha1\/earlystopping-medianstopping@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha1\/earlystopping-medianstopping:${VERSION}@" manifests/v1alpha1/vizier/earlystopping/medianstopping/deployment.yaml
sed -i -e '/volumeMounts:/,$d' manifests/v1alpha1/vizier/db/deployment.yaml

cat manifests/v1alpha1/vizier/core/deployment.yaml
./scripts/v1alpha1/deploy.sh

TIMEOUT=120
PODNUM=$(kubectl get deploy -n kubeflow | grep -v NAME | wc -l)
until kubectl get pods -n kubeflow | grep Running | [[ $(wc -l) -eq $PODNUM ]]; do
    echo Pod Status $(kubectl get pods -n kubeflow | grep Running | wc -l)/$PODNUM
    sleep 10
    TIMEOUT=$(( TIMEOUT - 1 ))
    if [[ $TIMEOUT -eq 0 ]];then
        echo "NG"
        kubectl get pods -n kubeflow
        exit 1
    fi
done

echo "All Katib components are running."
kubectl version
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep vizier-core | grep -v vizier-core-rest | sed -e "s@pods\/@@") 6789:6789 &
echo "kubectl port-forward start"
sleep 5
TIMEOUT=120
until curl localhost:6789 || [ $TIMEOUT -eq 0 ]; do
    sleep 5
    TIMEOUT=$(( TIMEOUT - 1 ))
done

cp -r test ${GO_DIR}/test
cd ${GO_DIR}/test/e2e/v1alpha1
kubectl apply -f valid-studyjob.yaml
kubectl delete -f valid-studyjob.yaml
set +o errexit
kubectl apply -f invalid-studyjob.yaml
if [ $? -ne 1 ]; then
  exit 1
fi
set -o errexit
go run test-client.go -a random
go run test-client.go -a grid -c suggestion-config-grid.yml
#go run test-client.go -a hyperband -c suggestion-config-hyb.yml
