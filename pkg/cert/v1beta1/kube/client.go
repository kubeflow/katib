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

package kube

import (
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type Client struct {
	KubeClient client.Client
}

func (c *Client) getKubeClient() (client.Client, error){
	if c.KubeClient == nil {
		var err error
		if c.KubeClient, err = client.New(config.GetConfigOrDie(), client.Options{Scheme: scheme.Scheme}); err != nil {
			return nil, err
		}
	}
	return c.KubeClient, nil
}
