#/bin/bash
set -x
set -e
gcloud container clusters create --machine-type n1-highmem-2 katib  
SERVICE_ACCOUNT=github-issue-summarization
PROJECT=katib-202401 # The GCP Project name
NAMESPACE=katib
gcloud iam service-accounts --project=${PROJECT} create ${SERVICE_ACCOUNT} \
      --display-name "GCP Service Account for use with kubeflow examples"

gcloud projects add-iam-policy-binding ${PROJECT} --member \
      serviceAccount:${SERVICE_ACCOUNT}@${PROJECT}.iam.gserviceaccount.com --role=roles/storage.admin

KEY_FILE=/home/user/secrets/${SERVICE_ACCOUNT}@${PROJECT}.iam.gserviceaccount.com.json
gcloud iam service-accounts keys create ${KEY_FILE} \
      --iam-account ${SERVICE_ACCOUNT}@${PROJECT}.iam.gserviceaccount.com
set -e
kubectl config set-credentials temp-admin --username=admin --password=$(gcloud container clusters describe katib --format="value(masterAuth.password)")
kubectl config set-context temp-context --cluster=$(kubectl config get-clusters | grep katib) --user=temp-admin
kubectl config use-context temp-context
kubectl apply -f manifests/0-namespace.yaml
kubectl --namespace=${NAMESPACE} create secret generic gcp-credentials --from-file=key.json="${KEY_FILE}"
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
