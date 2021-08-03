# Katib Examples with Argo Workflows Integration

Here you can find examples of using Katib with [Argo Workflows](https://github.com/argoproj/argo-workflows).

**Note:** You have to install `Argo >= v3.1` to use it in Katib Experiments.

## Installation

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
kubectl patch namespace argo -p '{"metadata":{"labels":{"katib-metricscollector-injection":"enabled"}}}'
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
