/*
Copyright 2020 The Knative Authors

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
	"fmt"

	corev1 "k8s.io/api/core/v1"

	k8snetworkingv1 "k8s.io/api/networking/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/client-go/kubernetes"
	k8snetworkingv1lister "k8s.io/client-go/listers/networking/v1"

	"knative.dev/pkg/reconciler"

	"knative.dev/pkg/controller"

	"knative.dev/avi-controller/pkg/reconciler/kingress/resources"
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"

	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
)

type Reconciler struct {
	kubeclient       kubernetes.Interface
	k8sIngressLister k8snetworkingv1lister.IngressLister
}

// Check that our Reconciler implements Interface
var _ kingressreconciler.Interface = (*Reconciler)(nil)

// FinalizeKind implements Interface.FinalizeKind.
func (r *Reconciler) FinalizeKind(ctx context.Context, o *netv1alpha1.Ingress) reconciler.Event {
	if err := r.kubeclient.NetworkingV1().Ingresses("tanzu-system-ingress").Delete(ctx, o.Name, metav1.DeleteOptions{}); err != nil && !apierrs.IsNotFound(err) {
		return err
	}
	
	return nil
}

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *netv1alpha1.Ingress) reconciler.Event {
	if err := r.reconcileIngress(ctx, o); err != nil {
		return err
	}

	return nil
}

// Reconcile a k8s ingress for each external host in a knative ingress
func (r *Reconciler) reconcileIngress(ctx context.Context, ing *netv1alpha1.Ingress) error {
	recorder := controller.GetEventRecorder(ctx)

	k8sIngress, err := r.k8sIngressLister.Ingresses("tanzu-system-ingress").Get(ing.Name)
	// k8sIngress, err := r.k8sIngressLister.Ingresses(ing.Namespace).Get(ing.Name)
	if apierrs.IsNotFound(err) {
		desired := resources.MakeK8sIngress(ing)
		k8sIngress, err := r.kubeclient.NetworkingV1().Ingresses("tanzu-system-ingress").Create(ctx, desired, metav1.CreateOptions{})
		// k8sIngress, err := r.kubeclient.NetworkingV1().Ingresses(ing.Namespace).Create(ctx, desired, metav1.CreateOptions{})
		if err != nil {
			recorder.Eventf(ing, corev1.EventTypeWarning, "CreationFailed", "Failed to create Ingress: %v", err)
			return fmt.Errorf("failed to create Ingress: %w", err)
		}
		recorder.Eventf(ing, corev1.EventTypeNormal, "Created", "Created K8s Ingress %q", k8sIngress.GetName())
		return nil
	} else if err != nil {
		return err
	}

	return r.reconcileIngressUpdate(ctx, ing, k8sIngress)
}

func (r *Reconciler) reconcileIngressUpdate(ctx context.Context, ing *netv1alpha1.Ingress, k8sIngress *k8snetworkingv1.Ingress) error {
	original := k8sIngress.DeepCopy()
	desired := resources.MakeK8sIngress(ing)

	recorder := controller.GetEventRecorder(ctx)

	if !equality.Semantic.DeepEqual(original.Spec, desired.Spec) ||
		!equality.Semantic.DeepEqual(original.Annotations, desired.Annotations) ||
		!equality.Semantic.DeepEqual(original.Labels, desired.Labels) {

		// Don't modify the informers copy.
		original.Spec = desired.Spec
		original.Annotations = desired.Annotations
		original.Labels = desired.Labels
		_, err := r.kubeclient.NetworkingV1().Ingresses(original.Namespace).Update(ctx, original, metav1.UpdateOptions{})

		if err != nil {
			recorder.Eventf(ing, corev1.EventTypeWarning, "UpdateFailed", "Failed to update k8s Ingress: %v", err)
			return fmt.Errorf("failed to update k8s Ingress: %w", err)
		}
	}
	return nil
}
