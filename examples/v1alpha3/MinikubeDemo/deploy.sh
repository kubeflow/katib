#/bin/bash
set -x
set -e
minikube start --disk-size 50g --memory 4096 --cpus 4
bash ../../../scripts/v1alpha3/deploy.sh
