# Simple GKE Demo
You can deploy katib components and try simple mnist demo on the cloud!

## build your training docker image
The sample code is docker-image directory.
It is based on [KubeFlow's Github issue summaraize example](https://github.com/kubeflow/examples/tree/master/github_issue_summarization/notebooks) image.

## deploy Cluster
This is grid parameter search demo for [KubeFlow's Github issue summaraize example](https://github.com/kubeflow/examples/tree/master/github_issue_summarization)

Let's deploy GKE cluster by gcloud commnad.
```
gcloud container clusters create --machine-type n1-highmem-4 --num-nodes 2 katib  
```


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

### Push data to GCS
If you want to pull input data and push logs to google cloud storage, You need to confidurate GCP service account.

See also https://github.com/kubeflow/examples/blob/master/github_issue_summarization/training_the_model_tfjob.md

## Create Study
Run study with `go run`. you should set endpoint of vizier-core.

```
go run git-issue-summarize-demo.go -s {{ip-addr of node vizier-core deployed}}:30678
```
Katib will make 2 girds of learling-rate parameter from 0.005 to 0.5.

The `git-issue-summarize-demo.go` is a controller of this study.
It sets hyper parameter config, suggestion config, k8s job config and then, create study and run trial with katib-API.

The training logic is `docker-image/train.py`

## UI
You can check your Model with Web UI.

Acsess to `http://{{node address}}:30080/`

The Results will be saved automatically.
