package v1alpha3

import (
	"github.com/kubeflow/katib/pkg/util/v1alpha3/env"
	"strings"
)

const (
	// AvailableNameSpaceEnvName is the env name of available namespaces
	AvailableNameSpaceEnvName = "KATIB_UI_AVAILABLE_NS"
	// ClusterRoleKey is a key represents having access to cluster role
	ClusterRoleKey = "KATIB_UI_CLUSTER_ROLE"
)

var (
	availableNameSpaces []string
	hasClusterRole      bool
)

func init() {
	ns := env.GetEnvOrDefault(AvailableNameSpaceEnvName, ClusterRoleKey)
	if ns == ClusterRoleKey {
		// no namespace restriction when working with kubernetes client
		availableNameSpaces = []string{""}
		hasClusterRole = true
	} else {
		availableNameSpaces = strings.Split(ns, " ")
		hasClusterRole = false
	}
}
