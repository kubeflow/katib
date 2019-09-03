/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package pod

import (
	"errors"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	commonv1alpha2 "github.com/kubeflow/katib/pkg/common/v1alpha2"
)

func getKabitJob(pod *v1.Pod) (string, string, error) {
	for _, gvk := range commonv1alpha2.GetSupportedJobList() {
		owners := pod.GetOwnerReferences()
		for _, owner := range owners {
			if isMatchGVK(owner, gvk) {
				return owner.Kind, owner.Name, nil
			}
		}
	}
	return "", "", errors.New("The Pod doesn't belong to Katib Job")
}

func isMatchGVK(owner metav1.OwnerReference, gvk schema.GroupVersionKind) bool {
	if owner.Kind != gvk.Kind {
		return false
	}
	gv := gvk.Group + "/" + gvk.Version
	if gv != owner.APIVersion {
		return false
	}
	return true
}

func isMasterRole(pod *v1.Pod, jobKind string) bool {
	if labels, ok := JobRoleMap[jobKind]; ok {
		if len(labels) == 0 {
			return true
		}
		for _, label := range labels {
			if v, err := getLabel(pod, label); err == nil {
				if v == MasterRole {
					return true
				}
			}
		}
	}
	return false
}

func getLabel(pod *v1.Pod, targetLabel string) (string, error) {
	labels := pod.Labels
	for k, v := range labels {
		if k == targetLabel {
			return v, nil
		}
	}
	return "", errors.New("Label " + targetLabel + " not found.")
}
