# Simple Minikube Demo
You can deploy katib components and try a simple mnist demo on your laptop!

## Requirement
* VirtualBox
* Minikube
* kubectl

## deploy
Only type command `./deploy`.

A Minikube cluster and Katib components will be deployed!

You can check them with `kubectl get -n katib get pods`.
Don't worry if the `vizier-core` get an error. 
It will be recovered after DB will be prepared.
Wait until all components will be Running status.

You need a katib client.
It will call API of Katib to create study, get suggestions, run trial, and get metrics.
The details of the system flow for the client and katib components is [here](https://docs.google.com/presentation/d/1Dk4XxKfVncb2v2CUDAd3OhM7XLyG3B9jEuuVmMiCIMg/edit#slide=id.g3b424a8f63_2_146).

The [client example](./client-example.go) will read three config files.
* study-config.yml: Define study property and feasible space of parameters.
* suggesiton-config.yml: Define suggesiton parameter for each study and suggestion service. In this file, the config is for grid suggestion service.
* worker-config.yml: Define config for evaluation worker.

## Create Study
### Random Suggestion Demo
You can run rundom suggesiton demo.
```
go run client-example.go -s 192.168.99.100:30678 -a random
```
In this demo, 2 random parameters in
* Learning Rate (--lr) - type: double
* Number of NN Layer (--num-layers) - type: int
* optimizer (--optimizer) - type: categorical

Logs
```
2018/04/26 17:43:26 Study ID n9debe3de9ef67c8
2018/04/26 17:43:26 Study ID n9debe3de9ef67c8 StudyConfname:"random-demo" owner:"katib" optimization_type:MAXIMIZE optimization_goal:0.99 parameter_configs:<configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0,03" min:"0.07" > > > default_suggestion_algorithm:"random" default_early_stopping_algorithm:"medianstopping" objective_value_name:"Validation-accuracy" metrics:"accuracy" metrics:"Validation-accuracy"
2018/04/26 17:43:26 Get Random Suggestions [trial_id:"i988add515f1ca4c" study_id:"n9debe3de9ef67c8" parameter_set:<name:"--lr" parameter_type:DOUBLE value:"0.0611" >  trial_id:"g7afad58be7da888" study_id:"n9debe3de9ef67c8" parameter_set:<name:"--lr" parameter_type:DOUBLE value:"0.0444" > ]
2018/04/26 17:43:26 WorkerID p4482bfb5cdc17ee start
2018/04/26 17:43:26 WorkerID c19ca08ca4e6aab1 start
```

### Grid Demo
Almost same as random suggestion.
```
go run client-example.go -s 192.168.99.100:30678 -a grid
```
In this demo, make 4 grids for learning rate (--lr) Min 0.03 and Max 0.07.

## UI
You can check your Model with Web UI.
Acsess to `http://192.168.99.100:30080/`
The Results will be saved automatically.

## ModelManagement
You can export model data to yaml file with CLI.
```
katib-cli -s {{server-cli}} pull study {{study ID or name}}  -o {{filename}}
```

And you can push your existing models to Katib with CLI.
`mnist-models.yaml` is traind 22 models using random suggestion with this Parameter Config.

```
configs:
    - name: --lr
      parametertype: 1
      feasible:
        max: "0.07"
        min: "0.03"
        list: []
    - name: --lr-factor
      parametertype: 1
      feasible:
        max: "0.05"
        min: "0.005"
        list: []
    - name: --lr-step
      parametertype: 2
      feasible:
        max: "20"
        min: "5"
        list: []
    - name: --optimizer
      parametertype: 4
      feasible:
        max: ""
        min: ""
        list:
        - sgd
        - adam
        - ftrl
```
You can easy to explore the model on ModelDB.

```
katib-cli -s 192.168.99.100:30678 push md -f mnist-models.yaml
```
