# 1. Metrics Collector Proposal

- [1. Metrics Collector Proposal](#1-metrics-collector-proposal)
  - [1.1. Motivation](#11-motivation)
  - [1.2. Goal](#12-goal)
  - [1.3. API](#13-api)
    - [1.3.1. Metric Collector](#131-metric-collector)
  - [1.4. Implementation](#14-implementation)
    - [1.4.1. Mutating Webhook](#141-mutating-webhook)
    - [1.4.2. Metric Collector](#142-metric-collector)
    - [1.4.3. Collection of Final Metrics](#143-collection-of-final-metrics)

## 1.1. Motivation

[Katib](https://github.com/kubeflow/katib) is a hyperparameter tuning (HPT) and neural architecture search (NAS) system based on Kubernetes.
During the auto-training, the metrics collection is an essential step.
In the current design, the metrics collector is pulled-based.
Katib runs a metrics collector cron job for each Trial.
The cron job pulls the targeted pod logs periodically and then persist the logs into MySQL.
However, the pulled-based design has [some problems](https://github.com/kubeflow/tf-operator/issues/722#issuecomment-405669269), such as, at what frequency should we scrape the metrics and so on.

To enhance the extensibility and support EarlyStopping, we propose a new design of the metrics collector.
In the new design, katib use mutating webhook to inject metrics collector container as a sidecar into Job/Tfjob/PytorchJob pod.
The sidecar collects metrics of the master and then store them on the persistent layer (e.x. katib-manager and metadata server).

(need a figure here)

Fig. 1 Architecture of the new design

## 1.2. Goal

1. **A mutating webhook**: inject metrics collector as a sidecar into master pod.
2. **A metric collector**: collect metrics and store them on the persistent layer (katib-manager).
3. **The final metrics** of worker pods should be collected by trail controller and then be stored into trial status.

## 1.3. API

### 1.3.1. Metric Collector

For more detail, see [here](https://github.com/kubeflow/katib/pull/697#issuecomment-516264282).

    type MetricsCollectorSpec struct {
        Retain     bool       `json:"retain,omitempty"`
        // Deprecated Retain
        Retain bool `json:"retain,omitempty"`
        // Deprecated GoTemplate
        GoTemplate GoTemplate `json:"goTemplate,omitempty"`

        Source    *SourceSpec    `json:"source,omitempty"`
        Collector *CollectorSpec `json:"collector,omitempty"`
    }

    type SourceSpec struct {
        // Model-train source code can expose metrics by http, such as HTTP endpoint in
        // prometheus metric format
        HttpGet *v1.HTTPGetAction `json:"httpGet,omitempty"`
        // During training model, metrics may be persisted into local file in source
        // code, such as tfEvent use case
        FileSystemPath *FileSystemPath `json:"fileSystemPath,omitempty"`
        // Default metric output format is {"metric": "<metric_name>",
        // "value": <int_or_float>, "epoch": <int>, "step": <int>}, but if the output doesn't
        // follow default format, please extend it here
        Filter *FilterSpec `json:"filter,omitempty"`
        }

    type FilterSpec struct {
        // When the metrics output follows format as this field specified, metricsCollector
        // collects it and reports to metrics server, it can be "<metric_name>: <float>" or else
        MetricsFormat []string `json:"metricsFormat,omitempty"`
    }

    type FileSystemKind string

    const (
        DirectoryKind FileSystemKind = "diretory"
        FileKind      FileSystemKind = "file"
    )

    type FileSystemPath struct {
        Path string         `json:"path,omitempty"`
        Kind FileSystemKind `json:"kind,omitempty"`
    }

    type CollectorKind string

    const (
        StdOutCollector           CollectorKind = "stdOutCollector"
        FileCollector             CollectorKind = "fileCollector"
        TfEventCollector          CollectorKind = "tfEventCollector"
        PrometheusMetricCollector CollectorKind = "prometheusMetricCollector"
        CustomCollector           CollectorKind = "customCollector"
        // When model training source code persists metrics into persistent layer
        // directly, metricsCollector isn't in need, and its kind is "noneCollector"
        NoneCollector CollectorKind = "noneCollector"
    )

    type CollectorSpec struct {
        Kind CollectorKind `json:"kind"`
        // When kind is "customCollector", this field will be used
        CustomCollector *v1.Container `json:"customCollector,omitempty"`
    }

## 1.4. Implementation

### 1.4.1. Mutating Webhook

To avoid collecting duplicated metrics, as we discuss in [kubeflow/katib#685](https://github.com/kubeflow/katib/issues/685), only one metrics collector sidecar will be injected into the master pod during one Experiment.
In the new design, there are two modes for katib mutating webhook to inject the sidecar: **Pod Level Injecting** and **Job Level Injecting**.

The webhook decides which mode to be used based on the `katib-metrics-collector-injection=enabled` label tagged on the namespace.
In the namespace with `katib-metrics-collector-injection=enabled` label, the webhook inject the sidecar in the pod level. Otherwise, without this label, injecting in the job level.

In **Pod Level Injecting**,

1. Job operators (_e.x. TFjob/PyTorchjob_) tag a specific label on the master pod, for example, the label like `kubeflow.org/replica-role: master`.
2. The webhook inject the metric collector only if the webhook recognizes this label.
3. The webhook uses [ObjectSelector](https://github.com/kubernetes/kubernetes/pull/78505) to skip on irrelevant objects in order to optimize the performance.
4. ObjectSelector is only supported above _Kubernetes v1.15_. Without this new feature, there may be a [performance issue](https://github.com/kubeflow/katib/issues/685#issuecomment-516226070) in webhook. In this situation, the following **Job Level Injecting** mode may be a better option.

In **Job Level Injecting**,

1. The webhook use different strategies to inject sidecar according to different job operators. For now, the webhook support PytorchJob and TfJob.
2. For PytorchJob, the metrics collector sidecar is injected into master template.
3. For TfJob, the metrics collector sidecar is injected into master template if master exists. Otherwise, the sidecar is injected into worker template with 0 index.

After injecting, the sidecar collects metrics of the master and then store them on the persistent layer (e.x. katib-manager and metadata server).

### 1.4.2. Metric Collector

_#WIP_

### 1.4.3. Collection of Final Metrics

The final metrics of worker pods should be collected by trail controller and then be stored into trial status.

_#WIP_
