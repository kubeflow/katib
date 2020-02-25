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
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/authn/k8schain"
	"github.com/google/go-containerregistry/pkg/name"
	crv1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"

	jobv1alpha3 "github.com/kubeflow/katib/pkg/job/v1alpha3"
)

func getKatibJob(pod *v1.Pod) (string, string, error) {
	for _, gvk := range jobv1alpha3.SupportedJobList {
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
	if labels, ok := jobv1alpha3.JobRoleMap[jobKind]; ok {
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

func getRemoteImage(pod *v1.Pod, namespace string, containerIndex int) (crv1.Image, error) {
	// verify the image name, then download the remote config file
	c := pod.Spec.Containers[containerIndex]
	ref, err := name.ParseReference(c.Image, name.WeakValidation)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse image %q: %v", c.Image, err)
	}
	imagePullSecrets := []string{}
	for _, s := range pod.Spec.ImagePullSecrets {
		imagePullSecrets = append(imagePullSecrets, s.Name)
	}
	kc, err := k8schain.NewInCluster(k8schain.Options{
		Namespace:          namespace,
		ServiceAccountName: pod.Spec.ServiceAccountName,
		ImagePullSecrets:   imagePullSecrets,
	})
	if err != nil {
		return nil, fmt.Errorf("Failed to create k8schain: %v", err)
	}

	mkc := authn.NewMultiKeychain(kc)
	img, err := remote.Image(ref, remote.WithAuthFromKeychain(mkc))
	if err != nil {
		return nil, fmt.Errorf("Failed to get container image %q info from registry: %v", c.Image, err)
	}

	return img, nil
}

func getContainerCommand(pod *v1.Pod, namespace string, containerIndex int) ([]string, error) {
	// https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#notes
	var err error
	var img crv1.Image
	var cfg *crv1.ConfigFile
	args := []string{}
	c := pod.Spec.Containers[containerIndex]
	if len(c.Command) != 0 {
		args = append(args, c.Command...)
	} else {
		img, err = getRemoteImage(pod, namespace, containerIndex)
		if err != nil {
			return nil, err
		}
		cfg, err = img.ConfigFile()
		if err != nil {
			return nil, fmt.Errorf("Failed to get config for image %q: %v", c.Image, err)
		}
		if len(cfg.Config.Entrypoint) != 0 {
			args = append(args, cfg.Config.Entrypoint...)
		}
	}
	if len(c.Args) != 0 {
		args = append(args, c.Args...)
	} else {
		if cfg != nil && len(cfg.Config.Cmd) != 0 {
			args = append(args, cfg.Config.Cmd...)
		}
	}
	return args, nil
}
