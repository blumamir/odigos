/*
Copyright 2022.

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

package v1alpha1

import (
	"context"
	json "encoding/json"
	"fmt"
	"time"

	v1alpha1 "github.com/odigos-io/odigos/api/actions/v1alpha1"
	actionsv1alpha1 "github.com/odigos-io/odigos/api/generated/actions/applyconfiguration/actions/v1alpha1"
	scheme "github.com/odigos-io/odigos/api/generated/actions/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// AddClusterInfosGetter has a method to return a AddClusterInfoInterface.
// A group's client should implement this interface.
type AddClusterInfosGetter interface {
	AddClusterInfos(namespace string) AddClusterInfoInterface
}

// AddClusterInfoInterface has methods to work with AddClusterInfo resources.
type AddClusterInfoInterface interface {
	Create(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.CreateOptions) (*v1alpha1.AddClusterInfo, error)
	Update(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.UpdateOptions) (*v1alpha1.AddClusterInfo, error)
	UpdateStatus(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.UpdateOptions) (*v1alpha1.AddClusterInfo, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v1alpha1.AddClusterInfo, error)
	List(ctx context.Context, opts v1.ListOptions) (*v1alpha1.AddClusterInfoList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddClusterInfo, err error)
	Apply(ctx context.Context, addClusterInfo *actionsv1alpha1.AddClusterInfoApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AddClusterInfo, err error)
	ApplyStatus(ctx context.Context, addClusterInfo *actionsv1alpha1.AddClusterInfoApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AddClusterInfo, err error)
	AddClusterInfoExpansion
}

// addClusterInfos implements AddClusterInfoInterface
type addClusterInfos struct {
	client rest.Interface
	ns     string
}

// newAddClusterInfos returns a AddClusterInfos
func newAddClusterInfos(c *ActionsV1alpha1Client, namespace string) *addClusterInfos {
	return &addClusterInfos{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the addClusterInfo, and returns the corresponding addClusterInfo object, and an error if there is any.
func (c *addClusterInfos) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.AddClusterInfo, err error) {
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of AddClusterInfos that match those selectors.
func (c *addClusterInfos) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.AddClusterInfoList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1alpha1.AddClusterInfoList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("addclusterinfos").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested addClusterInfos.
func (c *addClusterInfos) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("addclusterinfos").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a addClusterInfo and creates it.  Returns the server's representation of the addClusterInfo, and an error, if there is any.
func (c *addClusterInfos) Create(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.CreateOptions) (result *v1alpha1.AddClusterInfo, err error) {
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("addclusterinfos").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(addClusterInfo).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a addClusterInfo and updates it. Returns the server's representation of the addClusterInfo, and an error, if there is any.
func (c *addClusterInfos) Update(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.UpdateOptions) (result *v1alpha1.AddClusterInfo, err error) {
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(addClusterInfo.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(addClusterInfo).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *addClusterInfos) UpdateStatus(ctx context.Context, addClusterInfo *v1alpha1.AddClusterInfo, opts v1.UpdateOptions) (result *v1alpha1.AddClusterInfo, err error) {
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(addClusterInfo.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(addClusterInfo).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the addClusterInfo and deletes it. Returns an error if one occurs.
func (c *addClusterInfos) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *addClusterInfos) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("addclusterinfos").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched addClusterInfo.
func (c *addClusterInfos) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.AddClusterInfo, err error) {
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// Apply takes the given apply declarative configuration, applies it and returns the applied addClusterInfo.
func (c *addClusterInfos) Apply(ctx context.Context, addClusterInfo *actionsv1alpha1.AddClusterInfoApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AddClusterInfo, err error) {
	if addClusterInfo == nil {
		return nil, fmt.Errorf("addClusterInfo provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(addClusterInfo)
	if err != nil {
		return nil, err
	}
	name := addClusterInfo.Name
	if name == nil {
		return nil, fmt.Errorf("addClusterInfo.Name must be provided to Apply")
	}
	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(*name).
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}

// ApplyStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating ApplyStatus().
func (c *addClusterInfos) ApplyStatus(ctx context.Context, addClusterInfo *actionsv1alpha1.AddClusterInfoApplyConfiguration, opts v1.ApplyOptions) (result *v1alpha1.AddClusterInfo, err error) {
	if addClusterInfo == nil {
		return nil, fmt.Errorf("addClusterInfo provided to Apply must not be nil")
	}
	patchOpts := opts.ToPatchOptions()
	data, err := json.Marshal(addClusterInfo)
	if err != nil {
		return nil, err
	}

	name := addClusterInfo.Name
	if name == nil {
		return nil, fmt.Errorf("addClusterInfo.Name must be provided to Apply")
	}

	result = &v1alpha1.AddClusterInfo{}
	err = c.client.Patch(types.ApplyPatchType).
		Namespace(c.ns).
		Resource("addclusterinfos").
		Name(*name).
		SubResource("status").
		VersionedParams(&patchOpts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}