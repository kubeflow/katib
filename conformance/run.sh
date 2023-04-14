#!/bin/sh

# Run conformance test and generate test report.
python test/e2e/v1beta1/scripts/gh-actions/run-e2e-experiment.py --experiment-path examples/v1beta1/hp-tuning/random.yaml --namespace kf-conformance \
--trial-pod-annotations '{"sidecar.istio.io/inject": "false"}' | tee /tmp/katib-conformance.log


# Create the done file.
touch /tmp/katib-conformance.done
echo "Done..."

# Keep the container running so the test logs can be downloaded.
while true; do sleep 10000; done