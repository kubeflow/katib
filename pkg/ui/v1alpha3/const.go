package v1alpha3

import (
	"os"
)

const (
	// Name of environment variable indicates no ClusterRole permission.
	// In other words, katib only has access to the default namespace of katib deployment.
	NoClusterRole = "KATIB_NO_CLUSTER_ROLE"
)

var (
	hasClusterRole bool
)

func init() {
	_, ok := os.LookupEnv(NoClusterRole)
	hasClusterRole = !ok
}
