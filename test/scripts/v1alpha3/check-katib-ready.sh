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

#This is required. But I don't know why.
# VERSION=${VERSION/%?/}

echo "Install Katib "
echo "REGISTRY ${REGISTRY}"
echo "REPO_NAME ${REPO_NAME}"
echo "VERSION ${VERSION}"

# Katib controller
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/katib-controller@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha3\/katib-controller:${VERSION}@" manifests/v1alpha3/katib-controller/katib-controller.yaml

# Metrics collector
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/file-metrics-collector@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/file-metrics-collector:${VERSION}@" manifests/v1alpha3/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/tfevent-metrics-collector@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/tfevent-metrics-collector:${VERSION}@" manifests/v1alpha3/katib-controller/katib-config.yaml

# Katib manager
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/katib-manager@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha3\/katib-manager:${VERSION}@" manifests/v1alpha3/manager/deployment.yaml

# UI
sed -i -e "s@image: gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/katib-ui@image: ${REGISTRY}\/${REPO_NAME}\/v1alpha3\/katib-ui:${VERSION}@" manifests/v1alpha3/ui/deployment.yaml

# Suggestion algorithms
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/suggestion-nasrl@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/suggestion-nasrl@" manifests/v1alpha3/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/suggestion-hyperband@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/suggestion-hyperband@" manifests/v1alpha3/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/suggestion-chocolate@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/suggestion-chocolate@" manifests/v1alpha3/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/suggestion-hyperopt@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/suggestion-hyperopt@" manifests/v1alpha3/katib-controller/katib-config.yaml
sed -i -e "s@gcr.io\/kubeflow-images-public\/katib\/v1alpha3\/suggestion-skopt@${REGISTRY}\/${REPO_NAME}\/v1alpha3\/suggestion-skopt@" manifests/v1alpha3/katib-controller/katib-config.yaml

cat manifests/v1alpha3/katib-controller/katib-config.yaml


mkdir -p ${GO_DIR}
cp -r . ${GO_DIR}/
cp -r pkg/apis/manager/v1alpha3/python/* ${GO_DIR}/test/e2e/v1alpha3


echo "Deploying tf-operator crds from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/tf-training/tf-job-crds/base"
kustomize build . | kubectl apply -f -

echo "Deploying pytorch-operator crds from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-job-crds/base"
kustomize build . | kubectl apply -f -

cd ${GO_DIR}
./scripts/v1alpha3/deploy.sh

echo "Deploying tf-operator from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/tf-training/tf-job-operator/base"
kustomize build . | kubectl apply -n kubeflow -f - --validate=false

echo "Deploying pytorch-operator from kubeflow/manifests master"
cd "${MANIFESTS_DIR}/pytorch-job/pytorch-operator/base/"
kustomize build . | kubectl apply -n kubeflow -f - --validate=false


TIMEOUT=120
PODNUM=$(kubectl get deploy -n kubeflow | grep -v NAME | wc -l)
until kubectl get pods -n kubeflow | grep Running | [[ $(wc -l) -eq $PODNUM ]]; do
    echo Pod Status $(kubectl get pods -n kubeflow | grep "1/1" | wc -l)/$PODNUM
    
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
kubectl cluster-info
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

cd ${GO_DIR}/test/e2e/v1alpha3

echo "Building run-e2e-experiment for e2e test cases"
go build -o run-e2e-experiment github.com/kubeflow/katib/test/e2e/v1alpha3/run-e2e-experiment.go
go build -o resume-e2e-experiment github.com/kubeflow/katib/test/e2e/v1alpha3/resume-e2e-experiment.go

kubectl apply -f valid-experiment.yaml
kubectl delete -f valid-experiment.yaml
set +o errexit
kubectl apply -f invalid-experiment.yaml
if [ $? -ne 1 ]; then
  echo "Failed to create invalid-experiment: return code $?"
  exit 1
fi
set -o errexit
kubectl -n kubeflow port-forward $(kubectl -n kubeflow get pod -o=name | grep katib-manager |sed -e "s@pods\/@@") 6789:6789 &
echo "kubectl port-forward start"
sleep 5
TIMEOUT=120
until curl localhost:6789 || [ $TIMEOUT -eq 0 ]; do
    sleep 5
    TIMEOUT=$(( TIMEOUT - 1 ))
done

exit 0
