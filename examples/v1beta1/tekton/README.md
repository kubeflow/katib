# Katib Examples with Tekton Pipelines Integration

Here you can find examples of using Katib with [Tekton](https://github.com/tektoncd/pipeline).

## Installation

### Tekton Pipelines

To deploy Tekton Pipelines `v0.26.0`, run the following command:

```bash
kubectl apply -f https://storage.googleapis.com/tekton-releases/pipeline/previous/v0.26.0/release.yaml
```

Check that Tekton Pipelines components are running:

```bash
$ kubectl get pods -n tekton-pipelines

NAME                                           READY   STATUS    RESTARTS   AGE
tekton-pipelines-controller-799cdc78fc-sm4vl   1/1     Running   0          50s
tekton-pipelines-webhook-79d8f4f9bc-qmk97      1/1     Running   0          50s
```

**Note:** You must modify Tekton [`nop`](https://github.com/tektoncd/pipeline/tree/master/cmd/nop)
image to run Tekton Pipelines. `Nop` image is used to stop sidecar containers after main container
is completed. Since Katib is using Metrics Collector sidecar container
and Tekton Pipelines Controller should not kill sidecar containers, you have to
set this `nop` image to Metrics Collector image.

For example, if you are using
[StdOut](https://www.kubeflow.org/docs/components/katib/experiment/#metrics-collector) Metrics Collector,
`nop` image must be equal to `docker.io/kubeflowkatib/file-metrics-collector`.

Run the following command to modify the `nop` image:

```bash
kubectl patch deploy tekton-pipelines-controller -n tekton-pipelines --type='json' \
  -p='[{"op": "replace", "path": "/spec/template/spec/containers/0/args/9", "value": "docker.io/kubeflowkatib/file-metrics-collector"}]'
```

Check that Tekton Pipelines Controller's pod was restarted:

```bash
$ kubectl get pods -n tekton-pipelines

NAME                                           READY   STATUS    RESTARTS   AGE
tekton-pipelines-controller-7fcb6c6cd4-p8zf2   1/1     Running   0          2m2s
tekton-pipelines-webhook-7f9888f9b-7d6mr       1/1     Running   0          3m
```

Verify that `nop` image was modified:

```bash
$ kubectl get $(kubectl get pods -o name -n tekton-pipelines | grep tekton-pipelines-controller) -n tekton-pipelines -o yaml | grep katib

   - docker.io/kubeflowkatib/file-metrics-collector
```

### Katib Controller

To run Tekton Pipelines within Katib Trials you have to update Katib
[ClusterRole's rules](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/components/controller/rbac.yaml#L5)
with the appropriate permission:

```yaml
- apiGroups:
    - tekton.dev
  resources:
    - pipelineruns
    - taskruns
  verbs:
    - "get"
    - "list"
    - "watch"
    - "create"
    - "delete"
```

Run the following command to update Katib ClusterRole:

```bash
kubectl patch ClusterRole katib-controller -n kubeflow --type=json \
  -p='[{"op": "add", "path": "/rules/-", "value": {"apiGroups":["tekton.dev"],"resources":["pipelineruns", "taskruns"],"verbs":["get", "list", "watch", "create", "delete"]}}]'
```

Modify Katib Config [controller parameters](https://github.com/kubeflow/katib/blob/fc858d15dd41ff69166a2020efa200199063f9ba/manifests/v1beta1/installs/katib-standalone/katib-config.yaml#L9-L15) with the new entity:

```
trialResources:
 - <object-kind>.<object-API-version>.<object-API-group>
```

For example, to support Tekton Pipelines:

```
trialResources:
  - PipelineRun.v1beta1.tekton.dev
```

After these changes, deploy Katib as described in the [install guide](https://www.kubeflow.org/docs/components/katib/installation/) and wait until the katib-controller Pod is created. You can check logs from the Katib controller to verify your resource integration:

```bash
$ kubectl logs $(kubectl get pods -n kubeflow -o name | grep katib-controller) -n kubeflow | grep '"CRD Kind":"PipelineRun"'

{"level":"info","ts":1628032648.6285546,"logger":"trial-controller","msg":"Job watch added successfully","CRD Group":"tekton.dev","CRD Version":"v1beta1","CRD Kind":"PipelineRun"}
```

If you ran the above steps successfully, you should be able to run Tekton Pipelines examples.

Learn more about using custom Kubernetes resource as a Trial template in the
[official Kubeflow guides](https://www.kubeflow.org/docs/components/katib/user-guides/trial-template/#use-crds-with-trial-template).
