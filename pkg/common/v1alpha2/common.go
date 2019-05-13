package v1alpha2

import (
	"os"

	"k8s.io/apimachinery/pkg/runtime/schema"

	experimentsv1alpha2 "github.com/kubeflow/katib/pkg/api/operators/apis/experiment/v1alpha2"
)

const (
	KatibManagerServiceIPEnvName        = "KATIB_MANAGER_PORT_6789_TCP_ADDR"
	KatibManagerServicePortEnvName      = "KATIB_MANAGER_PORT_6789_TCP_PORT"
	KatibManagerServiceNamespaceEnvName = "KATIB_MANAGER_NAMESPACE"
	KatibManagerService                 = "katib-manager"
	KatibManagerPort                    = "6789"
	ManagerAddr                         = KatibManagerService + ":" + KatibManagerPort
)

func GetManagerAddr() string {
	ns := os.Getenv(experimentsv1alpha2.DefaultKatibNamespaceEnvName)
	if len(ns) == 0 {
		addr := os.Getenv(KatibManagerServiceIPEnvName)
		port := os.Getenv(KatibManagerServicePortEnvName)
		if len(addr) > 0 && len(port) > 0 {
			return addr + ":" + port
		} else {
			return ManagerAddr
		}
	} else {
		return KatibManagerService + "." + ns + ":" + KatibManagerPort
	}
}

func GetSupportedJobList() []schema.GroupVersionKind {
	supportedJobList := []schema.GroupVersionKind{
		schema.GroupVersionKind{
			Group:   "batch",
			Version: "v1",
			Kind:    "Job",
		},
		schema.GroupVersionKind{
			Group:   "kubeflow.org",
			Version: "v1beta2",
			Kind:    "TFJob",
		},
		schema.GroupVersionKind{
			Group:   "kubeflow.org",
			Version: "v1beta2",
			Kind:    "PyTorchJob",
		},
	}
	return supportedJobList
}
