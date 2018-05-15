#/bin/bash
set -x
gcloud container clusters create --machine-type n1-highmem-4 --num-nodes 2 katib  
SERVICE_ACCOUNT=github-issue-summarization
set -e
kubectl config set-credentials temp-admin --username=admin --password=$(gcloud container clusters describe katib --format="value(masterAuth.password)")
kubectl config set-context temp-context --cluster=$(kubectl config get-clusters | grep katib) --user=temp-admin
kubectl config use-context temp-context
kubectl apply -f manifests/0-namespace.yaml
kubectl apply -f manifests/modeldb/db
kubectl apply -f manifests/modeldb/backend
kubectl apply -f manifests/modeldb/frontend
kubectl apply -f manifests/vizier/db
kubectl apply -f manifests/vizier/core
kubectl apply -f manifests/vizier/suggestion/random
kubectl apply -f manifests/vizier/suggestion/grid
kubectl apply -f manifests/vizier/earlystopping/medianstopping
gcloud compute firewall-rules create katibservice --allow tcp:30080,tcp:30678
gcloud compute instances list
