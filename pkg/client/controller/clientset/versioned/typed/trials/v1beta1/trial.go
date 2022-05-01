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

// Code generated by client-gen. DO NOT EDIT.

package v1beta1

import (
	"context"
	"time"

	v1beta1 "github.com/kubeflow/katib/pkg/apis/controller/trials/v1beta1"
	scheme "github.com/kubeflow/katib/pkg/client/controller/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// TrialsGetter has a method to return a TrialInterface.
// A group's client should implement this interface.
type TrialsGetter interface {
	Trials(namespace string) TrialInterface
}

// TrialInterface has methods to work with Trial resources.
type TrialInterface interface {
	Create(ctx context.Context, trial *v1beta1.Trial, opts v1.CreateOptions) (*v1beta1.Trial, error)
	Update(ctx context.Context, trial *v1beta1.Trial, opts v1.UpdateOptions) (*v1beta1.Trial, error)
	UpdateStatus(ctx context.Context, trial *v1beta1.Trial, opts v1.UpdateOptions) (*v1beta1.Trial, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1beta1.Trial, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1beta1.TrialList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Trial, err error)
	TrialExpansion
}

// trials implements TrialInterface
type trials struct {
	client rest.Interface
	ns     string
}

// newTrials returns a Trials
func newTrials(c *TrialV1beta1Client, namespace string) *trials {
	return &trials{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the trial, and returns the corresponding trial object, and an error if there is any.
func (c *trials) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.Trial, err error) {
	result = &v1beta1.Trial{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("trials").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Trials that match those selectors.
func (c *trials) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.TrialList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1beta1.TrialList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("trials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested trials.
func (c *trials) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("trials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a trial and creates it.  Returns the server's representation of the trial, and an error, if there is any.
func (c *trials) Create(ctx context.Context, trial *v1beta1.Trial, opts v1.CreateOptions) (result *v1beta1.Trial, err error) {
	result = &v1beta1.Trial{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("trials").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(trial).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a trial and updates it. Returns the server's representation of the trial, and an error, if there is any.
func (c *trials) Update(ctx context.Context, trial *v1beta1.Trial, opts v1.UpdateOptions) (result *v1beta1.Trial, err error) {
	result = &v1beta1.Trial{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("trials").
		Name(trial.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(trial).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *trials) UpdateStatus(ctx context.Context, trial *v1beta1.Trial, opts v1.UpdateOptions) (result *v1beta1.Trial, err error) {
	result = &v1beta1.Trial{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("trials").
		Name(trial.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(trial).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the trial and deletes it. Returns an error if one occurs.
func (c *trials) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("trials").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *trials) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("trials").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched trial.
func (c *trials) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.Trial, err error) {
	result = &v1beta1.Trial{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("trials").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
