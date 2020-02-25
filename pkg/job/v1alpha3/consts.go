package v1alpha3

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	// JobRoleMap is the map which is used to determin if the replica is master.
	// Katib will inject metrics collector into master replica.
	JobRoleMap = make(map[string][]string)
	// SupportedJobList returns the list of the supported jobs' GVK.
	SupportedJobList = make(map[string]schema.GroupVersionKind)
)
