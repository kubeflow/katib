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

package v1beta1

import (
	"github.com/kubeflow/katib/pkg/cert-generator/v1beta1/consts"
	"github.com/kubeflow/katib/pkg/cert-generator/v1beta1/generate"
	"github.com/spf13/cobra"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// NewKatibCertGeneratorCmd sets up `katib-cert-generator` command.
func NewKatibCertGeneratorCmd(kubeClient client.Client) (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   consts.JobName,
		Short: consts.JobName,
		Long:  consts.JobName,
	}
	cmd.AddCommand(generate.NewGenerateCmd(kubeClient))
	return cmd, nil
}
