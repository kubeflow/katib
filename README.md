# Katib
HyperParamete Tuning on Kubernetes.
This project is [Google vizier](https://static.googleusercontent.com/media/research.google.com/ja//pubs/archive/bcb15507f4b52991a0783013df4222240e942381.pdf) inspired.
Katib is a scalable and flexible hyperparameter tuning framework and  tightly integrate with kubernetes.
And it does not depend on a specific DL framework.
There are examples of three frameworks ( tensorflow, mxnet, and  pytorch).
## Name
Katib stands for `secretary` in Arabic. 
As Vizier stands for high official or prime minister in Arabic, I named this project Katib in honor of Vizier.

## Vizier compatible
Katib has Study, Trial and Suggestion that are defined in Goodle vizier.

### Study
Represents a single optimization run over a feasible space. 
Each Study contains a configuration describing the feasible space, as well as a set of Trials. 
It is assumed that objective function f(x) does not change in the course of a Study.

### Trial
A  list of parameter values, x, that will lead to a single evaluation of f(x). 
A trial can be “Completed”, which means that it has been evaluated and the objective value f(x) has been assigned to it, otherwise it is “Pending”.
One trial corresponding to one k8s Job.

### Suggestion 
An algorithm to make parameter set. 
Currently parameter expolalation algorithms Katib supported are

* random
* grid 
* [hyperband](https://arxiv.org/pdf/1603.06560.pdf)

## Components
Katib consists of several components as below.
Each component is running on k8s as a deployment.
And each component communicates with grpc, the API is defined at `API/api.proto`.

- vizier: main components.
    - vizier-core : API server of vizier.
    - vizier-db
- dlk-manager : a interface of kubernetes.
- suggesiont : implimentations of each expolalation algorithm.
    - vizier-suggestion-random
    - vizier-suggestion-grid
    - vizier-suggestion-hyperband
- modeldb : WebUI
    - modeldb-frontend
    - modeldb-backend
    - modeldb-db

## StudyConfig
In Study config file, you define the feasible space of parameters and configuration of kubernetes job.
Examples of Study config are in `conf` directory.
The configuration items are as follows.

- name: Study name
- owner: Owner
- objectivevaluename: Name of the objective value. Your evaluated software should be print log `{objectivevaluename}={objective value}` in std-io.
- optimizationtype: Optimization direction of the objective value. 1=maximize 2=minimize
- suggestalgorithm: [random, grid, hyperband] now
- suggestionparameters: Parameter of the algorithm. Set name-value style.
    - In random suggestion
        - SuggestionNum: How many suggestions will Katib create.
        - MaxParallel: Max number of run on kubernetes
    - In grid suggestion
        - MaxParallel: Max number of run on kubernetes
        - GridDefault: default number of grid
        - name: [parameter name] grid number of specified parameter.
- metrics: The value you want to save to modeldb besides objectivevaluename.
- image: docker image name
- mount
    - pvc: pvc
    - path: MountPath in container
- pullsecret: Name of Image pull secret
- gpu: number of GPU (If you want to run cpu task, set 0 or delete this parameter)
- command: commands
- parameterconfigs: define feasible space
    - configs
        - name : parameter space
        - parametertype: 1=float, 2=int, 4=categorical
        - feasible 
            - min
            - max
            - list (for categorical)

## Web UI
Katib provide Web UI based on ModelDB( https://github.com/mitdbg/modeldb ).
The ingress setting is defined in manifests/modeldb/frontend/ingress.yaml

## Getting Start

### Requirements
- docker
- kubernetes cluster ( kubectlable from your PC and if you want to use GPU, set up k8s [GPU scheduling]( https://kubernetes.io/docs/tasks/manage-gpus/scheduling-gpus/ ))
- Ingress controller (e.g. nginx)

### Install
First, Copy CLI tool.

```
docker pull katib/katib-cli
docker run --name katib-cli -itd katib/katib-cli sh
docker cp katib-cli:/go/src/github.com/mlkube/katib/cli/katib-cli bin/katib-cli
docker rm -f katib-cli
```
The cli tool will be put `bin` directory.

Let's deploy Katib on your cluster.
Kubernetes manifests are in `manifests` directory.
Set the environment of your cluster(Ingress, Persistent Volumes).

```
$ ./deploy
```
### Use CLI
Check which node the vizier-core was deployed.
Then access vizier API.
```
$ kubectl get -n katib pod -o wide
NAME                                        READY     STATUS    RESTARTS   AGE       IP          NODE
dlk-manager-6d8886f988-m485v                1/1       Running   0          11m       10.44.0.4   node2
modeldb-backend-57667f44f6-5cl8k            1/1       Running   0          11m       10.35.0.4   gpu-node2
modeldb-db-6fc46458f6-fv2mn                 1/1       Running   0          11m       10.47.0.4   gpu-node3
modeldb-frontend-5f6cf5c496-m7gxc           1/1       Running   0          11m       10.39.0.4   gpu-node1
vizier-core-864dd6fdd4-r55qv                1/1       Running   0          11m       10.35.0.5   gpu-node2
vizier-db-7b6f8c59bc-mjhh4                  1/1       Running   0          11m       10.36.0.4   node1
vizier-suggestion-random-5895dc79b4-pbqkc   1/1       Running   0          11m       10.47.0.5   gpu-node3

$ ./katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:14:49 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
```
### Create Example Study
Try Createstudy. Study will be created and start hyperparameter search.

```
$ ./katib-cli -s gpu-node2:30678 -f ../conf/random.yml Createstudy
2018/04/03 05:16:37 connecting gpu-node2:30678
2018/04/03 05:16:37 study conf{cifer10 root MAXIMIZE 0 configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0.07" min:"0.03" > > configs:<name:"--lr-factor" parameter_type:DOUBLE feasible:<max:"0.2" min:"0.05" > > configs:<name:"--max-random-h" parameter_type:INT feasible:<max:"46" min:"26" > > configs:<name:"--max-random-l" parameter_type:INT feasible:<max:"75" min:"25" > > configs:<name:"--num-epochs" parameter_type:INT feasible:<max:"3" min:"3" > >  [] random median  [name:"SuggestionNum" value:"2"  name:"MaxParallel" value:"2" ] [] Validation-accuracy [accuracy] mxnet/python:gpu [python /mxnet/example/image-classification/train_cifar10.py --batch-size=512 --gpus=0,1] 2 default-scheduler <nil> }
2018/04/03 05:16:37 req Createstudy
2018/04/03 05:16:37 CreateStudy: study_id:"fef3711aa343fae6"
```

You can check the job is running with `kubectl` command.

```
$ ./katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:19:49 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
fef3711aa343fae6        cifer10 root    2       0

$ kubectl get -n katib job
NAME                        DESIRED   SUCCESSFUL   AGE
b325ec8d96ce16df-worker-0   1         0            1m
wbe8aabd6ad4f50e-worker-0   1         0            1m
```

Check the status of jobs with `katib-cli` command.

```
$ ./katib-cli -s gpu-node2:30678 Getstudies
2018/04/03 05:26:20 connecting gpu-node2:30678
StudyID                 Name    Owner   RunningTrial    CompletedTrial
fef3711aa343fae6        cifer10 root    1       1
```

When some trials are completed, you can check the result of completed trials.
See endpoint of Katib UI ingress.
In this example, the endpoint is `k-cluster.example.net/katib`

```
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

### Use Persistent Volume
Create PV and PVC in katib namespace.
For example,

`pv_nfs.yml`

```
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
```
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

```
$ kubectl apply -f pv_nfs.yml
persistentvolume "nfs" created

$ kubectl apply -f pvc_nfs.yml
persistentvolumeclaim "nfs" created
```

Then set up mount config in StudyConfig like below.

```
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

```
$ ./katib-cli -s gpu-node2:30678 -f ../conf/random-pv.yml Createstudy
2018/04/03 05:49:47 connecting gpu-node2:30678
2018/04/03 05:49:47 study conf{cifer10-pv-test root MAXIMIZE 0 configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0.07" min:"0.03" > > configs:<name:"--lr-factor" parameter_type:DOUBLE feasible:<max:"0.2" min:"0.05" > > configs:<name:"--max-random-h" parameter_type:INT feasible:<max:"46" min:"26" > > configs:<name:"--max-random-l" parameter_type:INT feasible:<max:"75" min:"25" > > configs:<name:"--num-epochs" parameter_type:INT feasible:<max:"3" min:"3" > >  [] random median  [name:"SuggestionNum" value:"2"  name:"MaxParallel" value:"2" ] [] Validation-accuracy [accuracy] mxnet/python:gpu [python /mxnet/example/image-classification/train_cifar10.py --batch-size=512 --gpus=0,1] 2 default-scheduler pvc:"nfs" path:"/nfs-mnt"  }
2018/04/03 05:49:47 req Createstudy
2018/04/03 05:49:47 CreateStudy: study_id:"p6ee7933f2b62f30"
```
Now the jobs will use the input files in the nfs.


## TensorBoard Integration
Not only TensorFlow but also several DL flameworks (e.g. PyTorch, MxNet) support TnsorBoard format logging.
Katib can integrate TensorBoard easily.
To use TensorBoard from Katib, you should define persistent volume clame and set mount config for the Study.
Katib search each trial log in `{pvc mount path}/logs/{Study ID}/{Trial ID}`.
`{{STUDY_ID}}` and  `{{TRIAL_ID}}` in the Studyconfig file are replaced the corresponding value when creating each job.
See example `conf/tf-nmt.yml` that is a config for parameter tuning of [tensorflow/nmt](https://github.com/tensorflow/nmt).

```
./katib-cli -s gpu-node2:30678 -f ../conf/tf-nmt.yml Createstudy
2018/04/03 05:52:11 connecting gpu-node2:30678
2018/04/03 05:52:11 study conf{tf-nmt root MINIMIZE 0 configs:<name:"--num_train_steps" parameter_type:INT feasible:<max:"1000" min:"1000" > > configs:<name:"--dropout" parameter_type:DOUBLE feasible:<max:"0.3" min:"0.1" > > configs:<name:"--beam_width" parameter_type:INT feasible:<max:"15" min:"5" > > configs:<name:"--num_units" parameter_type:INT feasible:<max:"1026" min:"256" > > configs:<name:"--attention" parameter_type:CATEGORICAL feasible:<list:"luong" list:"scaled_luong" list:"bahdanau" list:"normed_bahdanau" > > configs:<name:"--decay_scheme" parameter_type:CATEGORICAL feasible:<list:"luong234" list:"luong5" list:"luong10" > > configs:<name:"--encoder_type" parameter_type:CATEGORICAL feasible:<list:"bi" list:"uni" > >  [] random median  [name:"SuggestionNum" value:"10"  name:"MaxParallel" value:"6" ] [] test_ppl [ppl bleu_dev bleu_test] yujioshima/tf-nmt:latest-gpu [python -m nmt.nmt --src=vi --tgt=en --out_dir=/nfs-mnt/logs/{{STUDY_ID}}_{{TRIAL_ID}} --vocab_prefix=/nfs-mnt/learndatas/wmt15_en_vi/vocab --train_prefix=/nfs-mnt/learndatas/wmt15_en_vi/train --dev_prefix=/nfs-mnt/learndatas/wmt15_en_vi/tst2012 --test_prefix=/nfs-mnt/learndatas/wmt15_en_vi/tst2013 --attention_architecture=standard --attention=normed_bahdanau --batch_size=128 --colocate_gradients_with_ops=true --eos=</s> --forget_bias=1.0 --init_weight=0.1 --learning_rate=1.0 --max_gradient_norm=5.0 --metrics=bleu --share_vocab=false --num_buckets=5 --optimizer=sgd --sos=<s> --steps_per_stats=100 --time_major=true --unit_type=lstm --src_max_len=50 --tgt_max_len=50 --infer_batch_size=32] 1 default-scheduler pvc:"nfs" path:"/nfs-mnt"  }
2018/04/03 05:52:11 req Createstudy
2018/04/03 05:52:11 CreateStudy: study_id:"n5c80f4af709a70d"
```
Make TensorBord deployments, services, and ingress automatically and you can access from Web UI.

![katib-demo](https://user-images.githubusercontent.com/10014831/38241910-64fb0646-376e-11e8-8b98-c26e577f3935.gif)


## CLI
### katib
##### options
- s
Set address of vizier-core. {IP Addr}:{Port}. default localhost:6789
Katib API is grpc.
Unfortunately, nginx ingress controller does not support grpc now ( next version it will support! )
So vizier-core expose port as NodePort(30678}.

#### Getstudys
Get list of studys and their status.

#### Createstudy
Send create new study request to katib api server.

##### options
- f
Specify the config file of your study.

### Stopstudy [Study_ID]
Delete specified study from API server.
But the results of trials in modelDB won't be deleted.

## Implement new suggestion algorithm
Suggestion API is defined as grpc service at `API/api.proto`.
You can attach new algorithm easily.

- implement suggestion API
- make k8s service named vizier-suggestion-{ algorithm-name } and expose port 6789

And to add new suggestion service, you don't need to stop components ( vizier-core, modeldb, and anything) that are already running.

## Build from source
You can build all images from source.
```
./build
```

## Uninstall
Delete `katib` namespace from your kubernetes cluster.
All components will be deleted
```
kubectl delete ns katib
```

## TODOs
* Integrate KubeFlow (tf/pytorch/caffe2/-operator) 
* Support Early Stopping
* Enrich the GUI
