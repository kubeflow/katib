# Katib examples with Tekton integration

Here you can find examples of using Katib with [Tekton](https://github.com/tektoncd/pipeline).
Check [here](https://github.com/tektoncd/pipeline/blob/master/docs/install.md#installing-tekton-pipelines-on-kubernetes) how to install Tekton on your cluster.

**Note** that you must modify Tekton [`nop`](https://github.com/tektoncd/pipeline/tree/master/cmd/nop) image to run Tekton pipelines. `Nop` images is used to stop sidecar containers after main container is completed. Metrics collector must be not stopped after training container is finished. To avoid this problem, `nop` image should be equal to metrics collector sidecar image.

For example, if you are using [StdOut](https://www.kubeflow.org/docs/components/hyperparameter-tuning/experiment/#metrics-collector) metrics collector, `nop` image must be equal to `gcr.io/kubeflow-images-public/katib/v1beta1/file-metrics-collector`.

After deploying Tekton on your cluster, run bellow command to modify `nop` image.

```bash
kubectl patch deploy tekton-pipelines-controller -n tekton-pipelines --type='json' \
-p='[{"op": "replace", "path": "/spec/template/spec/containers/0/args/9", "value": "gcr.io/kubeflow-images-public/katib/v1beta1/file-metrics-collector"}]'
```

Check that Tekton controller's pod was restarted:

```
kubectl get pods -n tekton-pipelines
```

Expected output:

```
NAME                                           READY   STATUS    RESTARTS   AGE
tekton-pipelines-controller-7fcb6c6cd4-p8zf2   1/1     Running   0          2m2s
tekton-pipelines-webhook-7f9888f9b-7d6mr       1/1     Running   0          12h
```

Check that `nop` image was modified:

```
kubectl get pod <tekton-controller-pod-name> -n tekton-pipelines- -o yaml | grep katib/v1beta1/file-metrics-collector
```

Expected output:

```
- gcr.io/kubeflow-images-public/katib/v1beta1/file-metrics-collector
```
