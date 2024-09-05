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

# This shell script is used to setup Katib deployment.
set -o errexit
set -o pipefail
set -o nounset
cd "$(dirname "$0")"

DEPLOY_KATIB_UI=${1:-false}
DEPLOY_TRAINING_OPERATOR=${2:-false}
WITH_DATABASE_TYPE=${3:-mysql}

E2E_TEST_IMAGE_TAG="e2e-test"
TRAINING_OPERATOR_VERSION="v1.6.0-rc.0"

echo "Start to install Katib"

# Update Katib images with `e2e-test`.
cd ../../../../../ && make update-images OLD_PREFIX="docker.io/kubeflowkatib/" NEW_PREFIX="docker.io/kubeflowkatib/" TAG="$E2E_TEST_IMAGE_TAG" && cd -

# first declare the which kustomization file to use, by default use mysql.
KUSTOMIZATION_FILE="../../../../../manifests/v1beta1/installs/katib-standalone/kustomization.yaml"
PVC_FILE="../../../../../manifests/v1beta1/components/mysql/pvc.yaml"

# If the database type is postgres, then use postgres.
if [ "$WITH_DATABASE_TYPE" == "postgres" ]; then
  KUSTOMIZATION_FILE="../../../../../manifests/v1beta1/installs/katib-standalone-postgres/kustomization.yaml"
  PVC_FILE="../../../../../manifests/v1beta1/components/postgres/pvc.yaml"
fi

# If the user wants to deploy Katib UI, then use the kustomization file for Katib UI.
if ! "$DEPLOY_KATIB_UI"; then
  index="$(yq eval '.resources.[] | select(. == "../../components/ui/") | path | .[-1]' $KUSTOMIZATION_FILE)"
  index="$index" yq eval -i 'del(.resources.[env(index)])' $KUSTOMIZATION_FILE
fi

# Since e2e test doesn't need to large storage, we use a small PVC for Katib.
yq eval -i '.spec.resources.requests.storage|="2Gi"' $PVC_FILE

echo -e "\n The Katib will be deployed with the following configs"
cat $KUSTOMIZATION_FILE
cat ../../../../../manifests/v1beta1/installs/katib-standalone/katib-config.yaml

# If the user wants to deploy training operator, then use the kustomization file for training operator.
if "$DEPLOY_TRAINING_OPERATOR"; then
  echo "Deploying Training Operator $TRAINING_OPERATOR_VERSION"
  kustomize build "github.com/kubeflow/training-operator/manifests/overlays/standalone?ref=$TRAINING_OPERATOR_VERSION" | kubectl apply -f -
fi

echo "Deploying Katib"
cd ../../../../../ && WITH_DATABASE_TYPE=$WITH_DATABASE_TYPE make deploy && cd -

# Wait until all Katib pods is running.
TIMEOUT=120s

kubectl wait --for=condition=ContainersReady=True --timeout=${TIMEOUT} -l "katib.kubeflow.org/component in ($WITH_DATABASE_TYPE,controller,db-manager,ui)" -n kubeflow pod ||
  (kubectl get pods -n kubeflow && kubectl describe pods -n kubeflow && exit 1)

echo "All Katib components are running."
echo "Katib deployments"
kubectl -n kubeflow get deploy
echo "Katib services"
kubectl -n kubeflow get svc
echo "Katib pods"
kubectl -n kubeflow get pod

# Check that Katib is working with 2 Experiments.
kubectl apply -f ../../testdata/valid-experiment.yaml
kubectl delete -f ../../testdata/valid-experiment.yaml

# Check the ValidatingWebhookConfiguration works well.
set +o errexit
kubectl apply -f ../../testdata/invalid-experiment.yaml
if [ $? -ne 1 ]; then
  echo "Failed to create invalid-experiment: return code $?"
  exit 1
fi
set -o errexit

exit 0
