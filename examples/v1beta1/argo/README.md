# Katib Examples with Argo Workflows Integration

Here you can find examples of using Katib with [Argo Workflows](https://github.com/argoproj/argo-workflows).

**Note:** You have to install `Argo Workflows >= v3.1.3` to use it in Katib Experiments.

## Installation

### Argo Workflow

To deploy Argo Workflows `v3.1.3`, run the following commands:

```bash
kubectl create namespace argo
kubectl apply -n argo -f https://github.com/argoproj/argo-workflows/releases/download/v3.1.3/install.yaml
```

Check that Argo Workflow components are running:

```bash
$ kubectl get pods -n argo

NAME                                  READY   STATUS    RESTARTS   AGE
argo-server-5bbd69cc6b-6nvb6          1/1     Running   0          20s
workflow-controller-5f48fb7c8-vw9bp   1/1     Running   0          20s
```

After that, run below command to enable
[Katib Metrics Collector sidecar injection](https://www.kubeflow.org/docs/components/katib/experiment/#metrics-collector):

```bash
kubectl patch namespace argo -p '{"metadata":{"labels":{"katib.kubeflow.org/metrics-collector-injection":"enabled"}}}'
```

**Note:** Argo Workflows are using `docker` as a
[default container runtime executor](https://argoproj.github.io/argo-workflows/workflow-executors/#workflow-executors).
Since Katib is using Metrics Collector sidecar container and Argo Workflows controller
should not kill sidecar containers, you have to modify this
executor to [`emissary`](https://argoproj.github.io/argo-workflows/workflow-executors/#emissary-emissary).

Run the following command to change the `containerRuntimeExecutor` to `emissary` in the
Argo `workflow-controller-configmap`

```bash
kubectl patch ConfigMap -n argo workflow-controller-configmap --type='merge' -p='{"data":{"containerRuntimeExecutor":"emissary"}}'
```

Verify that `containerRuntimeExecutor` has been modified:

```bash
$ kubectl get ConfigMap -n argo workflow-controller-configmap -o yaml | grep containerRuntimeExecutor

  containerRuntimeExecutor: emissary
```

### Katib Controller

To run Argo Workflow within Katib Trials you have to update Katib
[ClusterRole's rules](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/components/controller/rbac.yaml#L5)
with the appropriate permission:

```yaml
- apiGroups:
    - argoproj.io
  resources:
    - workflows
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
  -p='[{"op": "add", "path": "/rules/-", "value": {"apiGroups":["argoproj.io"],"resources":["workflows"],"verbs":["get", "list", "watch", "create", "delete"]}}]'
```

In addition to that, you have to modify Katib
[Controller args](https://github.com/kubeflow/katib/blob/master/manifests/v1beta1/components/controller/controller.yaml#L27)
with the new flag `--trial-resources`.

Run the following command to update Katib Controller args:

```bash
kubectl patch Deployment katib-controller -n kubeflow --type=json \
  -p='[{"op": "add", "path": "/spec/template/spec/containers/0/args/-", "value": "--trial-resources=Workflow.v1alpha1.argoproj.io"}]'
```

Check that Katib Controller's pod was restarted:

```bash
$ kubectl get pods -n kubeflow

NAME                                         READY   STATUS      RESTARTS   AGE
katib-controller-784994d449-9bgj9            1/1     Running     0          28s
katib-db-manager-78697c7bd4-ck7l8            1/1     Running     0          6m13s
katib-mysql-854cdb87c4-krcm9                 1/1     Running     0          6m13s
katib-ui-57b9d7f6dd-cv6gn                    1/1     Running     0          6m13s
```

Check logs from Katib Controller to verify Argo Workflow integration:

```bash
$ kubectl logs $(kubectl get pods -n kubeflow -o name | grep katib-controller) -n kubeflow | grep '"CRD Kind":"Workflow"'

{"level":"info","ts":1628032648.6285546,"logger":"trial-controller","msg":"Job watch added successfully","CRD Group":"argoproj.io","CRD Version":"v1alpha1","CRD Kind":"Workflow"}
```

If you ran the above steps successfully, you should be able to run Argo Workflow examples.

Learn more about using custom Kubernetes resource as a Trial template in the
[official Kubeflow guides](https://www.kubeflow.org/docs/components/katib/trial-template/#use-custom-kubernetes-resource-as-a-trial-template).
