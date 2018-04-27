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
You don't worry if the `vizier-core` get an error. It will be recovered after DB will be prepared.
Wait until all components are Running status.

## Create Study
### Random Suggestion Demo
You can run rundom suggesiton demo.
```
go run random/radom-suggest-demo.go
```
In this demo, 2 random learning rate parameters generated randomly between Min 0.03 and Max 0.07.
Logs
```
2018/04/26 17:43:26 Study ID n9debe3de9ef67c8
2018/04/26 17:43:26 Study ID n9debe3de9ef67c8 StudyConfname:"random-demo" owner:"katib" optimization_type:MAXIMIZE optimization_goal:0.99 parameter_configs:<configs:<name:"--lr" parameter_type:DOUBLE feasible:<max:"0,03" min:"0.07" > > > default_suggestion_algorithm:"random" default_early_stopping_algorithm:"medianstopping" objective_value_name:"Validation-accuracy" metrics:"accuracy" metrics:"Validation-accuracy"
2018/04/26 17:43:26 Get Random Suggestions [trial_id:"i988add515f1ca4c" study_id:"n9debe3de9ef67c8" parameter_set:<name:"--lr" parameter_type:DOUBLE value:"0.0611" >  trial_id:"g7afad58be7da888" study_id:"n9debe3de9ef67c8" parameter_set:<name:"--lr" parameter_type:DOUBLE value:"0.0444" > ]
2018/04/26 17:43:26 WorkerID p4482bfb5cdc17ee start
2018/04/26 17:43:26 WorkerID c19ca08ca4e6aab1 start
```

### Grid Demo
Same as random.
```
go run grid/grid-suggest-demo.go
```
In this demo, make 4 grids Min 0.03 and Max 0.07.

## UI
You can check your Model with Web UI.
Acsess to `http://192.168.99.100:30080/`
The Results will be saved automatically.
