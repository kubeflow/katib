#!/usr/bin/env bash
# Compare Helm vs Kustomize manifests for particular Katib components

set -euo pipefail

SCENARIO=${1:-"standalone"}
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"
CHART_DIR="$ROOT_DIR/charts/katib"
MANIFESTS_DIR="$ROOT_DIR/manifests/v1beta1"

declare -A KUSTOMIZE_PATHS=(
    ["standalone"]="$MANIFESTS_DIR/installs/katib-standalone"
    ["cert-manager"]="$MANIFESTS_DIR/installs/katib-cert-manager"
    ["external-db"]="$MANIFESTS_DIR/installs/katib-external-db"
    ["leader-election"]="$MANIFESTS_DIR/installs/katib-leader-election"
    ["openshift"]="$MANIFESTS_DIR/installs/katib-openshift"
    ["standalone-postgres"]="$MANIFESTS_DIR/installs/katib-standalone-postgres"
    ["with-kubeflow"]="$MANIFESTS_DIR/installs/katib-with-kubeflow"
)

declare -A HELM_VALUES=(
    ["standalone"]="$CHART_DIR/ci/values-standalone.yaml"
    ["cert-manager"]="$CHART_DIR/ci/values-cert-manager.yaml"
    ["external-db"]="$CHART_DIR/ci/values-external-db.yaml"
    ["leader-election"]="$CHART_DIR/ci/values-leader-election.yaml"
    ["openshift"]="$CHART_DIR/ci/values-openshift.yaml"
    ["standalone-postgres"]="$CHART_DIR/ci/values-postgres.yaml"
    ["with-kubeflow"]="$CHART_DIR/ci/values-kubeflow.yaml"
    ["enterprise"]="$CHART_DIR/ci/values-enterprise.yaml"
    ["production"]="$CHART_DIR/ci/values-production.yaml"
)

declare -A NAMESPACES=(
    ["standalone"]="kubeflow"
    ["cert-manager"]="kubeflow"
    ["external-db"]="kubeflow"
    ["leader-election"]="kubeflow"
    ["openshift"]="kubeflow"
    ["standalone-postgres"]="kubeflow"
    ["with-kubeflow"]="kubeflow"
    ["enterprise"]="kubeflow"
    ["production"]="kubeflow"
)

if [[ ! "${KUSTOMIZE_PATHS[$SCENARIO]:-}" ]]; then
    echo "ERROR: Unknown scenario: $SCENARIO"
    echo "Supported scenarios:"
    for scenario in "${!KUSTOMIZE_PATHS[@]}"; do
        echo "  - $scenario"
    done
    exit 1
fi

KUSTOMIZE_PATH="${KUSTOMIZE_PATHS[$SCENARIO]}"
HELM_VALUES_FILE="${HELM_VALUES[$SCENARIO]}"
NAMESPACE="${NAMESPACES[$SCENARIO]}"

echo "Comparing Katib manifests for scenario: $SCENARIO"

if [ ! -d "$KUSTOMIZE_PATH" ]; then
    echo "ERROR: Kustomize path does not exist: $KUSTOMIZE_PATH"
    exit 1
fi

if [ ! -f "$HELM_VALUES_FILE" ]; then
    echo "ERROR: Helm values file does not exist: $HELM_VALUES_FILE"
    exit 1
fi

if [ ! -d "$CHART_DIR" ]; then
    echo "ERROR: Helm chart directory does not exist: $CHART_DIR"
    exit 1
fi

KUSTOMIZE_OUTPUT="/tmp/kustomize-katib-${SCENARIO}.yaml"
HELM_OUTPUT="/tmp/helm-katib-${SCENARIO}.yaml"

cd "$ROOT_DIR"
kustomize build "$KUSTOMIZE_PATH" > "$KUSTOMIZE_OUTPUT"

cd "$CHART_DIR"
helm template katib . \
    --namespace "$NAMESPACE" \
    --include-crds \
    --values "$HELM_VALUES_FILE" > "$HELM_OUTPUT"

cd "$ROOT_DIR"
python3 "$SCRIPT_DIR/compare_manifests.py" \
    "$KUSTOMIZE_OUTPUT" \
    "$HELM_OUTPUT" \
    "$SCENARIO" \
    "$NAMESPACE" \
    ${VERBOSE:+--verbose}

COMPARISON_RESULT=$?

rm -f "$KUSTOMIZE_OUTPUT" "$HELM_OUTPUT"

if [ $COMPARISON_RESULT -eq 0 ]; then
    echo "SUCCESS: Manifests are equivalent for scenario '$SCENARIO'"
else
    echo "FAILED: Manifests are NOT equivalent for scenario '$SCENARIO'"
fi

exit $COMPARISON_RESULT 