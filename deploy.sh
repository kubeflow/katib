#/bin/bash
set -x
set -e
kubectl apply -f manifests/0-namespace.yaml
kubectl apply -f manifests/modeldb/db
kubectl apply -f manifests/modeldb/backend
kubectl apply -f manifests/modeldb/frontend
kubectl apply -f manifests/dlk
kubectl apply -f manifests/vizier/db
kubectl apply -f manifests/vizier/core
kubectl apply -f manifests/vizier/suggestion/random
