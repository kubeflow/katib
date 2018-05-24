# Getting Start

## Requirements

- Docker
- kubernetes cluster ( kubectlable from your PC and if you want to use GPU, set up k8s [GPU scheduling]( https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/ ))
- Ingress controller (e.g. Nginx)

## Install the system and CLI

First, Copy CLI tool.

For Linux
```bash
$ curl -Lo katib-cli https://github.com/kubeflow/katib/releases/download/v0.1.1-alpha/katib-cli-linux-amd64 && chmod +x katib-cli && sudo mv katib-cli /usr/local/bin/
```

For Mac
```bash
$ curl -Lo katib-cli https://github.com/kubeflow/katib/releases/download/v0.1.1-alpha/katib-cli-darwin-amd64 && chmod +x katib-cli && sudo mv katib-cli /usr/local/bin/
```

The cli tool will be put `/usr/local/bin/` directory.

Let's deploy Katib on your cluster.
Kubernetes manifests are in `manifests` directory.
Set the environment of your cluster(Ingress, Persistent Volumes).

```bash
$ ./scripts/deploy.sh
```

## Use CLI

Check which node the vizier-core was deployed.
Then access vizier API.

```bash
$ kubectl get -n katib pod -o wide
NAME                                        READY     STATUS    RESTARTS   AGE       IP          NODE
dlk-manager-6d8886f988-m485v                1/1       Running   0          11m       10.44.0.4   node2
modeldb-backend-57667f44f6-5cl8k            1/1       Running   0          11m       10.35.0.4   gpu-node2
modeldb-db-6fc46458f6-fv2mn                 1/1       Running   0          11m       10.47.0.4   gpu-node3
modeldb-frontend-5f6cf5c496-m7gxc           1/1       Running   0          11m       10.39.0.4   gpu-node1
vizier-core-864dd6fdd4-r55qv                1/1       Running   0          11m       10.35.0.5   gpu-node2
vizier-db-7b6f8c59bc-mjhh4                  1/1       Running   0          11m       10.36.0.4   node1
vizier-suggestion-random-5895dc79b4-pbqkc   1/1       Running   0          11m       10.47.0.5   gpu-node3

$ katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:14:49 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
```

## Create Example Study

Try Createstudy. Study will be created and start hyperparameter search.

```bash
$ katib-cli -s gpu-node2:30678 -f ../examples/random.yml Createstudy
2018/04/03 05:16:37 connecting gpu-node2:30678
2018/04/03 05:16:37 study conf{cifer10 root MAXIMIZE 0 configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0.07" min:"0.03" > > configs:<name:"--lr-factor" parameter_type:DOUBLE feasible:<max:"0.2" min:"0.05" > > configs:<name:"--max-random-h" parameter_type:INT feasible:<max:"46" min:"26" > > configs:<name:"--max-random-l" parameter_type:INT feasible:<max:"75" min:"25" > > configs:<name:"--num-epochs" parameter_type:INT feasible:<max:"3" min:"3" > >  [] random median  [name:"SuggestionNum" value:"2"  name:"MaxParallel" value:"2" ] [] Validation-accuracy [accuracy] mxnet/python:gpu [python /mxnet/example/image-classification/train_cifar10.py --batch-size=512 --gpus=0,1] 2 default-scheduler <nil> }
2018/04/03 05:16:37 req Createstudy
2018/04/03 05:16:37 CreateStudy: study_id:"fef3711aa343fae6"
```

You can check the job is running with `kubectl` command.

```bash
$ katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:19:49 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
fef3711aa343fae6        cifer10 root    2       0

$ kubectl get -n katib job
NAME                        DESIRED   SUCCESSFUL   AGE
b325ec8d96ce16df-worker-0   1         0            1m
wbe8aabd6ad4f50e-worker-0   1         0            1m
```

Check the status of jobs with `katib-cli` command.

```bash
$ katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:26:20 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
fef3711aa343fae6        cifer10 root    1       1
```

When some trials are completed, you can check the result of completed trials.
See endpoint of Katib UI ingress.
In this example, the endpoint is `k-cluster.example.net/katib`

```bash
$ kubectl -n katib describe ingress katib-ui
Name:             katib-ui
Namespace:        katib
Address:
Default backend:  default-http-backend:80 (<none>)
Rules:
  Host                Path  Backends
  ----                ----  --------
  k-cluster.example.net
                      /katib   modeldb-frontend:3000 (<none>)
Annotations:
Events:
  Type    Reason  Age   From                      Message
  ----    ------  ----  ----                      -------
  Normal  CREATE  1m    nginx-ingress-controller  Ingress katib/katib-ui
  Normal  UPDATE  1m    nginx-ingress-controller  Ingress katib/katib-ui
```

## Use Persistent Volume

Create PV and PVC in katib namespace.
For example,

`pv_nfs.yml`

```yaml
#PV manifest
apiVersion: v1
kind: PersistentVolume
metadata:
  name: nfs
  namespace: katib
  labels:
    type: "nfs"
spec:
  capacity:
    storage: 300Gi
  accessModes:
    - ReadWriteMany
  nfs:
    server: 192.168.1.3
    path: "/nfs/"
```
`pvc_nfs.yml`
```yaml
#PVC manifest
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: nfs
  namespace: katib
spec:
  accessModes:
    - ReadWriteMany
  storageClassName: ""
  resources:
    requests:
      storage: 300Gi
  selector:
    matchLabels:
      type: "nfs"
```

```bash
$ kubectl apply -f pv_nfs.yml
persistentvolume "nfs" created

$ kubectl apply -f pvc_nfs.yml
persistentvolumeclaim "nfs" created
```

Then set up mount config in StudyConfig like below.

```yaml
name: cifer10
owner: root
optimizationtype: 2
suggestalgorithm: random
autostopalgorithm: median
objectivevaluename: Validation-accuracy
scheduler: default-scheduler
image: mxnet/python:gpu
mount:
    pvc: nfs
    path: /nfs-mnt
gpu: 1
suggestionparameters:
    -
      name: SuggestionNum
      value: 2
    -
      name: MaxParallel
      value: 2
command:
        - python
        - /mxnet/example/image-classification/train_cifar10.py
        - --batch-size=512
        - --gpus=0
        - --num-epochs=3
metrics:
    - accuracy
parameterconfigs:
    configs:
      -
        name: --lr
        parametertype: 1
        feasible:
            min: 0.03
            max: 0.07
      -
        name: --lr-factor
        parametertype: 1
        feasible:
            min: 0.05
            max: 0.2
      -
        name: --max-random-h
        parametertype: 2
        feasible:
            min: 26
            max: 46
      -
        name: --max-random-l
        parametertype: 2
        feasible:
            min: 25
            max: 75
```

```bash
$ katib-cli -s gpu-node2:30678 -f ../examples/random-pv.yml Createstudy
2018/04/03 05:49:47 connecting gpu-node2:30678
2018/04/03 05:49:47 study conf{cifer10-pv-test root MAXIMIZE 0 configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0.07" min:"0.03" > > configs:<name:"--lr-factor" parameter_type:DOUBLE feasible:<max:"0.2" min:"0.05" > > configs:<name:"--max-random-h" parameter_type:INT feasible:<max:"46" min:"26" > > configs:<name:"--max-random-l" parameter_type:INT feasible:<max:"75" min:"25" > > configs:<name:"--num-epochs" parameter_type:INT feasible:<max:"3" min:"3" > >  [] random median  [name:"SuggestionNum" value:"2"  name:"MaxParallel" value:"2" ] [] Validation-accuracy [accuracy] mxnet/python:gpu [python /mxnet/example/image-classification/train_cifar10.py --batch-size=512 --gpus=0,1] 2 default-scheduler pvc:"nfs" path:"/nfs-mnt"  }
2018/04/03 05:49:47 req Createstudy
2018/04/03 05:49:47 CreateStudy: study_id:"p6ee7933f2b62f30"
```
Now the jobs will use the input files in the nfs.

## Uninstall the system

Delete `katib` namespace from your kubernetes cluster.
All components will be deleted
```bash
kubectl delete ns katib
```
