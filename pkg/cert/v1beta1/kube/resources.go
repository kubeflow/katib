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
	"context"
	"sigs.k8s.io/controller-runtime/pkg/client"
)
func (c *Client) CreateResources(ctx context.Context, obj client.Object) error {
	kubeClient, err := c.getKubeClient()
	if err != nil {
		return err
	}
	if err = kubeClient.Create(ctx, obj); err != nil {
		return err
	}
	return nil
}

func (c *Client) GetResources(ctx context.Context, key client.ObjectKey, obj client.Object) error {
	kubeClient, err := c.getKubeClient()
	if err != nil {
		return err
	}
	if err = kubeClient.Get(ctx, key, obj); err != nil {
		return err
	}
	return nil
}

func (c *Client) DeleteResources(ctx context.Context, obj client.Object) error {
	kubeClient, err := c.getKubeClient()
	if err != nil {
		return err
	}
	if err = kubeClient.Delete(ctx, obj); err != nil {
		return err
	}
	return nil
}

func (c *Client) PatchResources(ctx context.Context, obj client.Object, patch client.Patch) error {
	kubeClient, err := c.getKubeClient()
	if err != nil {
		return err
	}
	if err = kubeClient.Patch(ctx, obj, patch); err != nil {
		return err
	}
	return nil
}
