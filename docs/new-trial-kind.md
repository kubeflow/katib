# Document about how to support a new Kubernetes resource in Katib trial

## Update the supported list

First, `GetSupportedJobList` in [common.go](../pkg/common/v1alpha3/common.go) needs to be updated.

```go
func GetSupportedJobList() []schema.GroupVersionKind {
	supportedJobList := []schema.GroupVersionKind{
		{
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		},
		{
			Group:   "kubeflow.org",
			Version: "v1",
			Kind:    "TFJob",
		},
		{
			Group:   "kubeflow.org",
			Version: "v1",
			Kind:    "PyTorchJob",
		},
	}
	return supportedJobList
}
```

In this function, we define the Kubernetes `GroupVersionKind` that are supported in Katib. If you want to add a new kind, please append the `supportedJobList`.

## Update logic about status update

`GetDeployedJobStatus` in [trial_controller_util.go](../pkg/controller.v1alpha3/trial/trial_controller_util.go) needs to be updated.

It is used to determine if the trial is completed (Succeeded or Failed).

## Update logic about pod injection webhook

### Add logic to support mutating container

`isWorkerContainer` in [inject_webhook.go](../pkg/webhook/v1alpha3/pod/inject_webhook.go) needs to be updated.

```go
func isWorkerContainer(jobKind string, index int, c v1.Container) bool {
	switch jobKind {
	case BatchJob:
		if index == 0 {
			// for Job worker, the first container will be taken as worker container,
			// katib document should note it
			return true
		}
	case TFJob:
		if c.Name == TFJobWorkerContainerName {
			return true
		}
	case PyTorchJob:
		if c.Name == PyTorchJobWorkerContainerName {
			return true
		}
	default:
		log.Info("Invalid Katib worker kind", "JobKind", jobKind)
		return false
	}
	return false
}
```

The function is used to determine which container in the job is the actual main container.

### Add logic about how to determine the master pod

In Katib, we only inject metrics collector sidecar into the master pod (See [metrics-collector.md](./proposals/metrics-collector.md) for more details). Thus we need to update the `JobRoleMap` in [const.go](../pkg/webhook/v1alpha3/pod/const.go).

```go
var JobRoleMap = map[string][]string{
	"TFJob":      {JobRoleLabel, TFJobRoleLabel},
	"PyTorchJob": {JobRoleLabel, PyTorchJobRoleLabel},
	"Job":        {},
}
```