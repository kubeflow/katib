# Support custom CRD in Trial Job proposal

<!-- START doctoc generated TOC please keep comment here to allow auto update -->
<!-- DON'T EDIT THIS SECTION, INSTEAD RE-RUN doctoc TO UPDATE -->

## Table of Contents

- [Motivation](#motivation)
- [Goals](#goals)
- [Non-Goals](#non-goals)
- [Implementation](#implementation)
  - [API](#api)
  - [Trial controller watchers](#trial-controller-watchers)
  - [Primary pod label location](#primary-pod-label-location)
  - [Training container name](#training-container-name)
  - [Start metrics collector parser](#start-metrics-collector-parser)
  - [Succeeded status of running CRD](#succeeded-status-of-running-crd)
  - [Istio sidecar container](#istio-sidecar-container)

<!-- END doctoc generated TOC please keep comment here to allow auto update -->

Created by [doctoc](https://github.com/thlorenz/doctoc).

## Motivation

Running trial is one of the essential steps of executing Katib experiments.
We have implemented new trial template design in Katib v1beta1 ([katib/pull#1202](https://github.com/kubeflow/katib/pull/1202)
and [katib/pull#1215](https://github.com/kubeflow/katib/pull/1215)) to make
experiments valid YAML and make Katib more Kubernetes native.
After migrating to the new API, users still can run only [BatchJob](https://kubernetes.io/docs/concepts/workloads/controllers/job/),
[TFJob](https://github.com/kubeflow/tf-operator) or [PyTorchJob](https://github.com/kubeflow/pytorch-operator) in trial job.
If we want to support new CRD, we need to manually change Katib controller source code.

This approach makes impossible to use other CRDs in trial template, even if they satisfy trial job design.
The number of various Kubernetes CRDs grows significantly and many users would like to use them in Katib
(e.g, [katib/issue#1081, support Argo Workflow](https://github.com/kubeflow/katib/issues/1081)).
Another reason to design unified approach is that CRD controller can have `Go` package versions
that Katib controller doesn't support (e.g, [katib/issue#1081](https://github.com/kubeflow/katib/issues/1081#issuecomment-635338276)).

That is why we propose a new controller design to support custom CRD in trial job and make Katib usable for various Kubernetes resources.
To make this possible, we are changing API, trial controller, job provider, mutation webhook, metrics collector.

## Goals

1. Allow dynamic watchers for the custom CRD.
2. Inject Katib sidecar container on training pod.
3. Indicate training container for metrics collector execution.
4. Run metrics collector parser after all pod processes completion.
5. Get succeeded condition of running CRD.
6. Verify that `sidecar.istio.io/inject: false` label is added.

## Non-Goals

1. Inject Katib sidecar container on more than one pod simultaneously.
2. Specify list of succeeded conditions for the custom CRD.
3. Dynamically add new trial watcher for the custom CRD without Katib restart.

## Implementation

During implementation this feature, we should not brake current Katib controller logic.
Also, we need to make sure that CI is stable and it does not block other Katib work tasks.
After completion, we can clean-up redundant code.

### API

To achieve above goals, we introduce these `TrialTemplate` API changes.

```go

// TrialTemplate describes structure of Trial template
type TrialTemplate struct {
  // Retain indicates that Trial resources must be not cleanup
  Retain bool `json:"retain,omitempty"`

  // Source for Trial template (unstructured structure or config map)
  TrialSource `json:",inline"`

  // List of parameters that are used in Trial template
  TrialParameters []TrialParameterSpec `json:"trialParameters,omitempty"`

  // Label that determines if pod needs to be injected by Katib sidecar container
  PrimaryPodLabel map[string]string `json:"primaryPodLabel,omitempty"`

  // Name of training container where training is running
  PrimaryContainerName string `json:"primaryContainerName,omitempty"`

  // Name of condition when Trial custom resource is succeeded
  SucceededCondition string `json:"succeededCondition,omitempty"`

}
```

### Trial controller watchers

In the current design trial controller watches
[three supported resource](https://github.com/kubeflow/katib/blob/master/pkg/controller.v1beta1/trial/trial_controller.go#L94-L125).
To generate these parameters dynamically when Katib starts, we add additional flag (`-trial-resource`)
to Katib controller, which represents resources that can be used in trial template.
This flag contains custom CRD's `Group`, `Version`, `Kind` in `kind.version.group` format which needs to create controller watchers.
Trial controller iterates over these parameters and creates watchers.

For example, if trial can run TFJob, Argo Workflow and Kubernetes Batch Jobs, Katib controller flags must be:

```yaml
. . .
args:
  - "-webhook-port=8443"
  - "-trial-resource=TFJob.v1.kubeflow.org".
  - "-trial-resource=Workflow.v1alpha1.argoproj.io"
  - "-trial-resource=Job.v1.batch"
. . .
```

### Primary pod label location

Right now, we [inject](https://github.com/kubeflow/katib/blob/master/pkg/webhook/v1beta1/pod/utils.go#L58-L72)
metrics collector for TFJob and PyTorchJob only for _master_ pods with labels previously saved in controller constants.

We added a new `PrimaryPodLabel` parameter in `TrialTemplate` API to find primary pod that needs to be injected by Katib sidecar container.
User can define the key and value of the pod label which Katib must inject with sidecar container.

For example, for TFJob:

```yaml
. . .
PrimaryPodLabel:
  "training.kubeflow.org/job-role": "master"
. . .
```

### Training container name

In the current design we compare container name with
[default value](https://github.com/kubeflow/katib/blob/master/pkg/job/v1beta1/kubeflow.go#L63-L78) for TFJob and PyTorchJob
to find pod container where actual training is happening and metrics collector must parse metrics.

We introduce a new `PrimaryContainerName` field, where user can set container name with running training program to find proper training container.

For example, if training is running on container with `pytorch` name:

```yaml
. . .
PrimaryContainerName: "pytorch"
. . .
```

### Start metrics collector parser

As discussed in [katib/issue#1214](https://github.com/kubeflow/katib/issues/1214#issuecomment-642168716),
metrics collector starts parsing metrics only after all injected pod processes were finished.
That can avoid problems with other sidecar containers that various CRD can have.

We need to verify that [distributive training](https://docs.fast.ai/distributed.html#launch-your-training)
with more than one active process also works with this approach.

### Succeeded condition of running CRD

We have already [designed Kubeflow provider](https://github.com/kubeflow/katib/blob/master/pkg/job/v1alpha3/kubeflow.go#L27-L60)
to check succeeded conditions for the TFJob and PyTorchJob as `unstructured` objects by
[comparing](https://github.com/kubeflow/katib/blob/master/pkg/controller.v1beta1/trial/trial_controller_util.go#L161)
`.status.conditions[x].type` value with `Succeeded` value.

Different CRD can have unique status design (e.g, Kubernetes batch job succeeded condition is
[`Complete`](https://github.com/kubernetes/api/blob/master/batch/v1/types.go#L167-L173)).
We add a new parameter `SucceededCondition` to get CRD succeeded condition value and trigger trial controller.
Trial controller checks all running job conditions and verifies that running job has appropriate `type`
in `.status.conditions` with `status=True`.
We also should transform `reason` and `message` from custom CRD to the trial conditions, if it is available.

For example, for TFJob:

```yaml
. . .
SucceededCondition: Succeeded
. . .
```

### Istio sidecar container

Previously, we had problems with Istio sidecar containers,
check [kubeflow/issue#1081](https://github.com/kubeflow/kubeflow/issues/4742).
In some cases, it is unable to properly download datasets in training pod.
It was fixed by adding annotation `sidecar.istio.io/inject: false` to appropriate Trial job in Katib controller.

Various CRD can have unified design and it is hard to understand where annotation must be specified
to disable Istio injection for the running pods.
We need to update all Katib examples manually and add this annotation to every trial template.

This exception has to be documented and new Katib examples have to include this annotation in templates.
