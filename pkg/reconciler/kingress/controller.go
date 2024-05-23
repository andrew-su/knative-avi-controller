/*
Copyright 2019 The Knative Authors

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

package kingress

import (
	"context"

	"k8s.io/client-go/tools/cache"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"

	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"

	kingressinformer "knative.dev/networking/pkg/client/injection/informers/networking/v1alpha1/ingress"
	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
	kubeclient "knative.dev/pkg/client/injection/kube/client"
	ingressinformer "knative.dev/pkg/client/injection/kube/informers/networking/v1/ingress"
)

const (
	ContourIngressClassName = "contour.ingress.networking.knative.dev"

	finalizerName = "avi.ingresses.networking.internal.knative.dev"
)

// NewController creates a Reconciler and returns the result of NewImpl.
func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {
	// Obtain an informer to both the main and child resources. These will be started by
	// the injection framework automatically. They'll keep a cached representation of the
	// cluster's state of the respective resource at all times.
	kingressInformer := kingressinformer.Get(ctx)
	ingressInformer := ingressinformer.Get(ctx)

	r := &Reconciler{
		// The client will be needed to create/delete Pods via the API.
		kubeclient:       kubeclient.Get(ctx),
		k8sIngressLister: ingressInformer.Lister(),
	}
	impl := kingressreconciler.NewImpl(ctx, r, ContourIngressClassName, func(impl *controller.Impl) controller.Options {
		return controller.Options{
			FinalizerName: finalizerName,
		}
	})

	// Listen for events on the main resource and enqueue themselves.
	kingressInformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	// Listen for events on the child resources and enqueue the owner of them.
	ingressInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterController(&netv1alpha1.Ingress{}),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
