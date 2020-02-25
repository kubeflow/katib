# How Katib v1alpha3 tunes hyperparameter automatically in a Kubernetes native way

See the following guides in the Kubeflow documentation:

* [Concepts](https://www.kubeflow.org/docs/components/hyperparameter-tuning/overview/) 
  in Katib, hyperparameter tuning, and neural architecture search.
* [Getting started with Katib](https://kubeflow.org/docs/components/hyperparameter-tuning/hyperparameter/).
* Detailed guide to [configuring and running a Katib 
  experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/).

## Example and Illustration

After install Katib v1alpha3, you can run `kubectl apply -f katib/examples/v1alpha3/random-example.yaml` to try the first example of Katib.
Then you can get the new `Experiment` as below. Katib concepts will be introduced based on this example.
```
# kubectl get experiment random-example -n kubeflow -o yaml
apiVersion: kubeflow.org/v1alpha3
kind: Experiment
metadata:
  ...
  name: random-example
  namespace: kubeflow
spec:
  algorithm:
    algorithmName: random
  maxFailedTrialCount: 3
  maxTrialCount: 12
  metricsCollectorSpec:
    collector:
      kind: StdOut
  objective:
    additionalMetricNames:
    - accuracy
    goal: 0.99
    objectiveMetricName: Validation-accuracy
    type: maximize
  parallelTrialCount: 3
  parameters:
  - feasibleSpace:
      max: "0.03"
      min: "0.01"
    name: --lr
    parameterType: double
  - feasibleSpace:
      max: "5"
      min: "2"
    name: --num-layers
    parameterType: int
  trialTemplate:
    goTemplate:
      rawTemplate: |-
        apiVersion: batch/v1
        kind: Job
        metadata:
          name: {{.Trial}}
          namespace: {{.NameSpace}}
        spec:
          template:
            spec:
              containers:
              - name: {{.Trial}}
                image: docker.io/kubeflowkatib/mxnet-mnist
                command:
                - "python3"
                - "/opt/mxnet-mnist/mnist.py"
                - "--batch-size=64"
                {{- with .HyperParameters}}
                {{- range .}}
                - "{{.Name}}={{.Value}}"
                {{- end}}
                {{- end}}
              restartPolicy: Never
status:
  ...
```
#### Experiment
When you want to tune hyperparameters for your machine learning model before 
training it further, you just need to create an `Experiment` CR like above. To
learn what fields are included in the `Experiment.spec`, see
the detailed guide to [configuring and running a Katib 
experiment](https://kubeflow.org/docs/components/hyperparameter-tuning/experiment/).

#### Trial
For each set of hyperparameters, Katib will internally generate a `Trial` CR with the hyperparameters key-value pairs, job manifest string with parameters instantiated and some other fields like below. `Trial` CR is used for internal logic control, and end user can even ignore it.
```
# kubectl get trial -n kubeflow
NAME                      STATUS      AGE
random-example-fm2g6jpj   Succeeded   4h
random-example-hhzm57bn   Succeeded   4h
random-example-n8whlq8g   Succeeded   4h

# kubectl get trial random-example-fm2g6jpj -o yaml -n kubeflow
apiVersion: kubeflow.org/v1alpha3
kind: Trial
metadata:
  ...
  name: random-example-fm2g6jpj
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1alpha3
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-example
    uid: c7bbb111-de6b-11e9-a6cc-00163e01b303
spec:
  metricsCollector:
    collector:
      kind: StdOut
  objective:
    additionalMetricNames:
    - accuracy
    goal: 0.99
    objectiveMetricName: Validation-accuracy
    type: maximize
  parameterAssignments:
  - name: --lr
    value: "0.027435456064371484"
  - name: --num-layers
    value: "4"
  - name: --optimizer
    value: sgd
  runSpec: |-
    apiVersion: batch/v1
    kind: Job
    metadata:
      name: random-example-fm2g6jpj
      namespace: kubeflow
    spec:
      template:
        spec:
          containers:
          - name: random-example-fm2g6jpj
            image: docker.io/kubeflowkatib/mxnet-mnist
            command:
            - "python3"
            - "/opt/mxnet-mnist/mnist.py"
            - "--batch-size=64"
            - "--lr=0.027435456064371484"
            - "--num-layers=4"
            - "--optimizer=sgd"
          restartPolicy: Never
status:
  completionTime: 2019-09-24T01:38:39Z
  conditions:
  - lastTransitionTime: 2019-09-24T01:37:26Z
    lastUpdateTime: 2019-09-24T01:37:26Z
    message: Trial is created
    reason: TrialCreated
    status: "True"
    type: Created
  - lastTransitionTime: 2019-09-24T01:38:39Z
    lastUpdateTime: 2019-09-24T01:38:39Z
    message: Trial is running
    reason: TrialRunning
    status: "False"
    type: Running
  - lastTransitionTime: 2019-09-24T01:38:39Z
    lastUpdateTime: 2019-09-24T01:38:39Z
    message: Trial has succeeded
    reason: TrialSucceeded
    status: "True"
    type: Succeeded
  observation:
    metrics:
    - name: Validation-accuracy
      value: 0.981489
  startTime: 2019-09-24T01:37:26Z
```
#### Suggestion
Katib will internally create a `Suggestion` CR for each `Experiment` CR. `Suggestion` CR includes the hyperparameter algorithm name by `algorithmName` field and how many sets of hyperparameter Katib asks to be generated by `requests` field. The CR also traces all already generated sets of hyperparameter in `status.suggestions`. Same as `Trial`, `Suggestion` CR is used for internal logic control and end user can even ignore it.
```
# kubectl get suggestion random-example -n kubeflow -o yaml
apiVersion: kubeflow.org/v1alpha3
kind: Suggestion
metadata:
  ...
  name: random-example
  namespace: kubeflow
  ownerReferences:
  - apiVersion: kubeflow.org/v1alpha3
    blockOwnerDeletion: true
    controller: true
    kind: Experiment
    name: random-example
    uid: c7bbb111-de6b-11e9-a6cc-00163e01b303
spec:
  algorithmName: random
  requests: 3
status:
  ...
  suggestions:
  - name: random-example-fm2g6jpj
    parameterAssignments:
    - name: --lr
      value: "0.027435456064371484"
    - name: --num-layers
      value: "4"
    - name: --optimizer
      value: sgd
  - name: random-example-n8whlq8g
    parameterAssignments:
    - name: --lr
      value: "0.013743390382347042"
    - name: --num-layers
      value: "3"
    - name: --optimizer
      value: sgd
  - name: random-example-hhzm57bn
    parameterAssignments:
    - name: --lr
      value: "0.012495283371215943"
    - name: --num-layers
      value: "2"
    - name: --optimizer
      value: sgd
```
## What happens after an `Experiment` CR created
When a user created an `Experiment` CR, Katib controllers including experiment controller, trial controller and suggestion controller will work together to achieve hyperparameters tuning for user Machine learning model.
<center>
<img width="100%" alt="image" src="images/katib-workflow.png">
</center>

1. A `Experiment` CR is submitted to Kubernetes API server, Katib experiment mutating and validating webhook will be called to set default value for the `Experiment` CR and validate the CR separately.
2. Experiment controller create a `Suggestion` CR.
3. Suggestion controller create the algorithm deployment and service based on the new `Suggestion` CR.
4. When Suggestion controller verifies that the algorithm service is ready, it calls the service to generate `spec.request - len(status.suggestions)` sets of hyperparamters and append them into `status.suggestions`
5. Experiment controller finds that `Suggestion` CR had been updated, then generate each `Trial` for each new hyperparamters set. 
6. Trial controller generates job based on `runSpec` manifest with the new hyperparamters set.
7. Related job controller (Kubernetes batch Job, kubeflow PytorchJob or kubeflow TFJob) generated Pods.
8. Katib Pod mutating webhook is called to inject metrics collector sidecar container to the candidate Pod.
9. During the ML model container runs, metrics collector container in the same Pod tries to collect metrics from it and persists them into Katib DB backend.
10. When the ML model Job ends, Trial controller will update status of the corresponding `Trial` CR.
11. When a `Trial` CR goes to end, Experiment controller will increase `request` field of corresponding 
`Suggestion` CR if in need, then everything goes to `step 4` again. Of course, if `Trial` CRs meet one of `end` condition (exceeds `maxTrialCount`, `maxFailedTrialCount` or `goal`), Experiment controller will take everything done.