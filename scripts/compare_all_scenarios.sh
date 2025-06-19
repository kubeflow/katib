#!/usr/bin/env bash
# Compare Helm vs Kustomize manifests for all Katib scenarios

set -euo pipefail

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
ROOT_DIR="$(dirname "$SCRIPT_DIR")"

# All scenarios to test
SCENARIOS=(
    "standalone"
    "cert-manager"
    "external-db"
    "leader-election"
    "openshift"
    "standalone-postgres"
    "with-kubeflow"
)

echo "Total scenarios to test: ${#SCENARIOS[@]}"

declare -a PASSED_SCENARIOS=()
declare -a FAILED_SCENARIOS=()

for scenario in "${SCENARIOS[@]}"; do
    echo "Testing scenario: $scenario"
    
    if "$SCRIPT_DIR/compare_helm_kustomize.sh" "$scenario"; then
        echo "PASSED: $scenario"
        PASSED_SCENARIOS+=("$scenario")
    else
        echo "FAILED: $scenario"
        FAILED_SCENARIOS+=("$scenario")
    fi
done

echo "Passed scenarios (${#PASSED_SCENARIOS[@]}/${#SCENARIOS[@]}):"
for scenario in "${PASSED_SCENARIOS[@]}"; do
    echo "  $scenario"
done

if [ ${#FAILED_SCENARIOS[@]} -gt 0 ]; then
    echo "Failed scenarios (${#FAILED_SCENARIOS[@]}/${#SCENARIOS[@]}):"
    for scenario in "${FAILED_SCENARIOS[@]}"; do
        echo "  $scenario"
    done
    
    echo "OVERALL RESULT: FAILED"
    echo "Some scenarios have differences between Helm and Kustomize manifests."
    exit 1
else
    echo "OVERALL RESULT: SUCCESS"
    echo "All scenarios passed! Helm and Kustomize manifests are equivalent."
    exit 0
fi 