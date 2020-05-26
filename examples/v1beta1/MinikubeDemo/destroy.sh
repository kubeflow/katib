#/bin/bash
set -x
set -e
minikube delete
pkill kubectl
