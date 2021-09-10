/*
Copyright 2021 The Kubeflow Authors.

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

package v1beta1

import (
	"github.com/kubeflow/katib/pkg/cert/v1beta1/generate"
	"github.com/kubeflow/katib/pkg/cert/v1beta1/kube"
	"github.com/spf13/cobra"
)

func NewKatibCertGeneratorCmd(kubeClient *kube.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "katib-cert-generator",
		Short: "katib-cert-generator",
		Long:  "katib-cert-generator",
	}
	cmd.AddCommand(generate.NewGenerateCmd(kubeClient))
	return cmd, nil
}