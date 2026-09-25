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
TRAINING_OPERATOR_VERSION="v1.9.0"
KUBECTL_REQUEST_TIMEOUT="30s"
WEBHOOK_TIMEOUT_SECONDS=120
WEBHOOK_TIMEOUT="${WEBHOOK_TIMEOUT_SECONDS}s"

dump_katib_webhook_diagnostics() {
  echo "Katib controller readiness diagnostics"
  echo "Katib controller deployment"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get deploy/katib-controller -o wide || true
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow describe deploy/katib-controller || true

  echo "Katib controller service"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get service/katib-controller -o wide || true
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow describe service/katib-controller || true

  echo "Katib controller endpoints"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get endpoints/katib-controller -o wide || true
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow describe endpoints/katib-controller || true

  echo "Katib controller EndpointSlices"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get endpointslice -l kubernetes.io/service-name=katib-controller -o wide --show-labels || true
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get endpointslice -l kubernetes.io/service-name=katib-controller -o yaml || true

  echo "Katib webhook configurations"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" get mutatingwebhookconfiguration,validatingwebhookconfiguration -o yaml || true

  echo "Katib controller logs"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow logs deploy/katib-controller --all-containers --tail=500 || true

  echo "Kubeflow namespace events"
  kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get events --sort-by=.lastTimestamp || true
}

wait_for_katib_webhook_endpoint() {
  echo "Waiting for Katib controller rollout and webhook endpoint readiness."

  if ! kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow rollout status deployment/katib-controller --timeout="${WEBHOOK_TIMEOUT}"; then
    echo "Katib controller deployment did not finish rollout."
    dump_katib_webhook_diagnostics
    exit 1
  fi

  if ! kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get service/katib-controller -o wide; then
    echo "Katib controller service is not available."
    dump_katib_webhook_diagnostics
    exit 1
  fi

  local deadline
  deadline=$((SECONDS + WEBHOOK_TIMEOUT_SECONDS))
  local endpoint_addresses
  local endpointslice_addresses

  while true; do
    endpoint_addresses="$(kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get endpoints katib-controller -o jsonpath='{.subsets[*].addresses[*].ip}' || true)"
    endpointslice_addresses="$(kubectl --request-timeout="${KUBECTL_REQUEST_TIMEOUT}" -n kubeflow get endpointslice -l kubernetes.io/service-name=katib-controller -o jsonpath='{.items[*].endpoints[*].addresses[*]}' || true)"

    if [ -n "${endpoint_addresses}" ] && [ -n "${endpointslice_addresses}" ]; then
      break
    fi

    if [ "${SECONDS}" -ge "${deadline}" ]; then
      echo "Katib controller service endpoints were not populated before timeout."
      dump_katib_webhook_diagnostics
      exit 1
    fi

    echo "Waiting for katib-controller Endpoints and EndpointSlices to be populated."
    sleep 5
  done

  echo "Katib controller webhook endpoint is ready."
}

echo "Start to install Katib"

# Update Katib images with `e2e-test`.
cd ../../../../../ && make update-images OLD_PREFIX="ghcr.io/kubeflow/katib/" NEW_PREFIX="ghcr.io/kubeflow/katib/" TAG="$E2E_TEST_IMAGE_TAG" && cd -

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
  kubectl apply --server-side -k "github.com/kubeflow/training-operator/manifests/overlays/standalone?ref=$TRAINING_OPERATOR_VERSION"
fi

echo "Deploying Katib"
cd ../../../../../ && WITH_DATABASE_TYPE=$WITH_DATABASE_TYPE make deploy && cd -

# Wait until all Katib pods is running.
TIMEOUT=120s

kubectl wait --for=condition=ContainersReady=True --timeout=${TIMEOUT} -l "katib.kubeflow.org/component in ($WITH_DATABASE_TYPE,controller,db-manager,ui)" -n kubeflow pod ||
  (kubectl get pods -n kubeflow && kubectl describe pods -n kubeflow && exit 1)

wait_for_katib_webhook_endpoint

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
