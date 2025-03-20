/*
Copyright 2022 The Kubeflow Authors.
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
package util

import (
	"github.com/kubeflow/katib/pkg/controller.v1beta1/consts"
)

// AppendIstioSidecarLabel adds the Istio sidecar injection label to a labels map
func AppendIstioSidecarLabel(labels map[string]string) map[string]string {
	res := make(map[string]string)
	for k, v := range labels {
		res[k] = v
	}
	res[consts.LabelIstioSidecarInjectName] = consts.LabelIstioSidecarInjectValue
	return res
}
