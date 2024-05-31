/*
Copyright 2021 The Knative Authors

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

package ingress

import (
	"k8s.io/apimachinery/pkg/runtime"
	fakekubeclientset "k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/tools/cache"

	fakeservingclientset "knative.dev/networking/pkg/client/clientset/versioned/fake"

	networking "knative.dev/networking/pkg/apis/networking/v1alpha1"
	networkinglisters "knative.dev/networking/pkg/client/listers/networking/v1alpha1"

	k8snetworking "k8s.io/api/networking/v1"
	k8snetworkinglister "k8s.io/client-go/listers/networking/v1"

	"knative.dev/pkg/reconciler/testing"
)

var clientSetSchemes = []func(*runtime.Scheme) error{
	fakeservingclientset.AddToScheme,
	fakekubeclientset.AddToScheme,
}

type Listers struct {
	sorter testing.ObjectSorter
}

func NewListers(objs []runtime.Object) Listers {
	scheme := NewScheme()

	ls := Listers{
		sorter: testing.NewObjectSorter(scheme),
	}

	ls.sorter.AddObjects(objs...)

	return ls
}

func NewScheme() *runtime.Scheme {
	scheme := runtime.NewScheme()

	for _, addTo := range clientSetSchemes {
		addTo(scheme)
	}
	return scheme
}

func (*Listers) NewScheme() *runtime.Scheme {
	return NewScheme()
}

// IndexerFor returns the indexer for the given object.
func (l *Listers) IndexerFor(obj runtime.Object) cache.Indexer {
	return l.sorter.IndexerForObjectType(obj)
}

func (l *Listers) GetNetworkingObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakeservingclientset.AddToScheme)
}

func (l *Listers) GetKubeObjects() []runtime.Object {
	return l.sorter.ObjectsForSchemeFunc(fakekubeclientset.AddToScheme)
}

// GetIngressLister get lister for Ingress resource.
func (l *Listers) GetIngressLister() k8snetworkinglister.IngressLister {
	return k8snetworkinglister.NewIngressLister(l.IndexerFor(&k8snetworking.Ingress{}))
}

// GetKIngressLister get lister for Knative Ingress resource.
func (l *Listers) GetKIngressLister() networkinglisters.IngressLister {
	return networkinglisters.NewIngressLister(l.IndexerFor(&networking.Ingress{}))
}
