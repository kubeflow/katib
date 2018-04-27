# Simple GKE Demo
You can deploy katib components and try simple mnist demo on the cloud!

## deploy Cluster
This is grid parameter search demo for [KubeFlow's Github issue summaraize example](https://github.com/kubeflow/examples/tree/master/github_issue_summarization)

Let's deploy GKE cluster by gcloud commnad.

You need to confidurate GCP service account.
See also https://github.com/kubeflow/examples/blob/master/github_issue_summarization/training_the_model_tfjob.md

Then deploy Katib compenents.

```
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
```

In this demo, katib components export using NodePort.

So you should set firewall to allow the ports.

## Create Study
Please edit `git-issue-summarize-demo.go` the `manager` address to the node address that vizier-core deployed.
Then 
```
go run git-issue-summarize-demo.go
```
Katib will make 4 girds of learling-rate parameter from 0.005 to 0.5.

## UI
You can check your Model with Web UI.

Acsess to `http://{{node address}}:30080/`

The Results will be saved automatically.
