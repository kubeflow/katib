# Katib config overview

Here you can find information about [Katib config](https://github.com/kubeflow/katib/blob/master/manifests/v1alpha3/katib-controller/katib-config.yaml).

Katib config is the Kubernetes [Config Map](https://kubernetes.io/docs/tasks/configure-pod-container/configure-pod-configmap/) that contains information about current [Metrics Collectors](https://github.com/kubeflow/katib/blob/master/docs/proposals/metrics-collector.md) (`key = metrics-collector-sidecar`) and current [Suggestions](https://github.com/kubeflow/katib/blob/master/docs/proposals/suggestion.md) (`key = suggestion`).

Katib Config Map must be deployed in [KATIB_CORE_NAMESPACE](https://github.com/kubeflow/katib/blob/master/docs/components-env-variables.md#katib-controller) with `katib-config` name. Katib controller parses Katib config when you submit Experiment.

You can edit this Config Map even after deploying Katib.

For example, if you deploy Katib in Kubeflow namespace, to edit Katib config run this:

`kubectl edit configMap katib-config -n kubeflow`

## Metrics Collector Sidecar settings

These settings related to Katib Metrics Collectors, where `key = metrics-collector-sidecar` and value = corresponding settings JSON for each Metrics Collector.

Example for `File` Metrics Collector with all settings:

```json
{
  "File": {
    "image": "gcr.io/kubeflow-images-public/katib/v1alpha3/file-metrics-collector",
    "imagePullPolicy": "Always",
    "resources": {
      "requests": {
        "memory": "200Mi",
        "cpu": "250m",
        "ephemeral-storage": "200Mi"
      },
      "limits": {
        "memory": "1Gi",
        "cpu": "500m",
        "ephemeral-storage": "2Gi"
      }
    }
  }
}
```

All of these settings except `image` can be omitted. If you don't specify any other settings, default value will be set.

1. `image` - Docker image name for the `File` Metrics Collector.

   **Must be specified**.

2. `imagePullPolicy` - `File` Metrics Collector container [image pull policy](https://kubernetes.io/docs/concepts/configuration/overview/#container-images).

   Default value is `IfNotPresent`.

3. `resources` - `File` Metrics Collector [container resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container). In the above example you can see how to specify `limits` and `requests`. Currently, you can specify only `memory`, `cpu` and `ephemeral-storage`.

   Default values for `requests` are:

   - `memory = 10Mi`.
   - `cpu = 50m`.
   - `ephemeral-storage = 500Mi`.

   Default values for `limits` are:

   - `memory = 100Mi`.
   - `cpu = 500m`.
   - `ephemeral-storage = 5Gi`.

## Suggestion settings

These settings related to Katib suggestions, where `key = suggestion` and value = corresponding settings JSON for each algorithm. If you want to use new algorithm, you must update Katib config with the new suggestion.

Example for `random` suggestion with all settings:

```json
{
  "random": {
    "image": "gcr.io/kubeflow-images-public/katib/v1alpha3/suggestion-hyperopt",
    "imagePullPolicy": "Always",
    "resources": {
      "requests": {
        "memory": "100Mi",
        "cpu": "100m",
        "ephemeral-storage": "100Mi"
      },
      "limits": {
        "memory": "500Mi",
        "cpu": "500m",
        "ephemeral-storage": "3Gi"
      }
    },
    "serviceAccountName": "suggestion-serviceaccount"
  }
}
```

All of these settings except `image` can be omitted. If you don't specify any other settings, default value will be set.

1. `image` - Docker image name for the `random` suggestion.

   **Must be specified**.

2. `imagePullPolicy` - `Random` suggestion container [image pull policy](https://kubernetes.io/docs/concepts/configuration/overview/#container-images).

   Default value is `IfNotPresent`.

3. `resources` - `Random` suggestion [container resources](https://kubernetes.io/docs/concepts/configuration/manage-compute-resources-container/#resource-requests-and-limits-of-pod-and-container). In the above example you can see how to specify `limits` and `requests`. Currently, you can specify only `memory`, `cpu` and `ephemeral-storage`.

   Default values for `requests` are:

   - `memory = 10Mi`.
   - `cpu = 50m`.
   - `ephemeral-storage = 500Mi`.

   Default values for `limits` are:

   - `memory = 100Mi`.
   - `cpu = 500m`.
   - `ephemeral-storage = 5Gi`.

4. `serviceAccountName` - `Random` suggestion container [service account](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/). In the above example, `suggestion-serviceaccount` service account will be used for each Experiment with `random` algorithm, until you change or delete this service account from Katib config.

   By default suggestion pod will not have specific service account. Suggestion pod will use [default](https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/#use-the-default-service-account-to-access-the-api-server) service account.
