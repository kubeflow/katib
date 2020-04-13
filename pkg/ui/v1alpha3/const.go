package v1alpha3

import (
	"github.com/kubeflow/katib/pkg/util/v1alpha3/env"
	"strings"
)

const (
	AvailableNameSpaceEnvName = "KATIB_UI_AVAILABLE_NS"
	ClusterRoleKey            = "KATIB_UI_CLUSTER_ROLE"
)

var (
	availableNameSpaces, hasClusterRole = func() ([]string, bool) {
		ns := env.GetEnvOrDefault(AvailableNameSpaceEnvName, ClusterRoleKey)
		if ns == ClusterRoleKey {
			return []string{""}, true
		}
		return strings.Split(ns, " "), false
	}()
)
