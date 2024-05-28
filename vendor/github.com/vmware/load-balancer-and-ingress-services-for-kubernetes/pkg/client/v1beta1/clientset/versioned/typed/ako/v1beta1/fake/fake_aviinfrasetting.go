/*
Copyright The Kubernetes Authors.

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

package fake

import (
	"context"

	v1beta1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAviInfraSettings implements AviInfraSettingInterface
type FakeAviInfraSettings struct {
	Fake *FakeAkoV1beta1
}

var aviinfrasettingsResource = v1beta1.SchemeGroupVersion.WithResource("aviinfrasettings")

var aviinfrasettingsKind = v1beta1.SchemeGroupVersion.WithKind("AviInfraSetting")

// Get takes name of the aviInfraSetting, and returns the corresponding aviInfraSetting object, and an error if there is any.
func (c *FakeAviInfraSettings) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1beta1.AviInfraSetting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(aviinfrasettingsResource, name), &v1beta1.AviInfraSetting{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.AviInfraSetting), err
}

// List takes label and field selectors, and returns the list of AviInfraSettings that match those selectors.
func (c *FakeAviInfraSettings) List(ctx context.Context, opts v1.ListOptions) (result *v1beta1.AviInfraSettingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(aviinfrasettingsResource, aviinfrasettingsKind, opts), &v1beta1.AviInfraSettingList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.AviInfraSettingList{ListMeta: obj.(*v1beta1.AviInfraSettingList).ListMeta}
	for _, item := range obj.(*v1beta1.AviInfraSettingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested aviInfraSettings.
func (c *FakeAviInfraSettings) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(aviinfrasettingsResource, opts))
}

// Create takes the representation of a aviInfraSetting and creates it.  Returns the server's representation of the aviInfraSetting, and an error, if there is any.
func (c *FakeAviInfraSettings) Create(ctx context.Context, aviInfraSetting *v1beta1.AviInfraSetting, opts v1.CreateOptions) (result *v1beta1.AviInfraSetting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(aviinfrasettingsResource, aviInfraSetting), &v1beta1.AviInfraSetting{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.AviInfraSetting), err
}

// Update takes the representation of a aviInfraSetting and updates it. Returns the server's representation of the aviInfraSetting, and an error, if there is any.
func (c *FakeAviInfraSettings) Update(ctx context.Context, aviInfraSetting *v1beta1.AviInfraSetting, opts v1.UpdateOptions) (result *v1beta1.AviInfraSetting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(aviinfrasettingsResource, aviInfraSetting), &v1beta1.AviInfraSetting{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.AviInfraSetting), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeAviInfraSettings) UpdateStatus(ctx context.Context, aviInfraSetting *v1beta1.AviInfraSetting, opts v1.UpdateOptions) (*v1beta1.AviInfraSetting, error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceAction(aviinfrasettingsResource, "status", aviInfraSetting), &v1beta1.AviInfraSetting{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.AviInfraSetting), err
}

// Delete takes name of the aviInfraSetting and deletes it. Returns an error if one occurs.
func (c *FakeAviInfraSettings) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(aviinfrasettingsResource, name, opts), &v1beta1.AviInfraSetting{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAviInfraSettings) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(aviinfrasettingsResource, listOpts)

	_, err := c.Fake.Invokes(action, &v1beta1.AviInfraSettingList{})
	return err
}

// Patch applies the patch and returns the patched aviInfraSetting.
func (c *FakeAviInfraSettings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1beta1.AviInfraSetting, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(aviinfrasettingsResource, name, pt, data, subresources...), &v1beta1.AviInfraSetting{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.AviInfraSetting), err
}
