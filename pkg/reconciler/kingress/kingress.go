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

	"k8s.io/apimachinery/pkg/api/equality"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	k8snetworkingv1lister "k8s.io/client-go/listers/networking/v1"

	aviclient "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1beta1/clientset/versioned"
	avilister "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/client/v1beta1/listers/ako/v1beta1"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	"knative.dev/avi-controller/pkg/reconciler/kingress/resources"
	"knative.dev/pkg/reconciler"

	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
)

type Reconciler struct {
	kubeclient kubernetes.Interface
	akoclient  aviclient.Interface

	k8sIngressLister k8snetworkingv1lister.IngressLister
	hostRuleLister   avilister.HostRuleLister
}

// Check that our Reconciler implements Interface
var _ kingressreconciler.Interface = (*Reconciler)(nil)

// FinalizeKind implements Interface.FinalizeKind.
func (r *Reconciler) FinalizeKind(ctx context.Context, ing *netv1alpha1.Ingress) reconciler.Event {
	cfg := config.FromContext(ctx)

	if err := r.kubeclient.NetworkingV1().Ingresses(cfg.Avi.EnvoyService.Namespace).Delete(ctx, ing.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	if err := r.akoclient.AkoV1beta1().HostRules(cfg.Avi.EnvoyService.Namespace).Delete(ctx, ing.Name, metav1.DeleteOptions{}); err != nil {
		return err
	}

	return nil
}

// ReconcileKind implements Interface.ReconcileKind.
func (r *Reconciler) ReconcileKind(ctx context.Context, o *netv1alpha1.Ingress) reconciler.Event {
	if err := r.reconcileIngress(ctx, o); err != nil {
		return err
	}

	if err := r.reconcileHostRule(ctx, o); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) reconcileHostRule(ctx context.Context, ing *netv1alpha1.Ingress) error {
	cfg := config.FromContext(ctx)

	desired := resources.MakeHostRule(ctx, ing)

	if desired == nil {
		if _, err := r.akoclient.AkoV1beta1().HostRules(cfg.Avi.EnvoyService.Namespace).Get(ctx, ing.Name, metav1.GetOptions{}); apierrs.IsNotFound(err) {
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to get hostrule: %w", err)
		}

		if err := r.akoclient.AkoV1beta1().HostRules(cfg.Avi.EnvoyService.Namespace).Delete(ctx, ing.Name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to delete hostrule: %w", err)
		}
		return nil
	}

	existing, err := r.akoclient.AkoV1beta1().HostRules(cfg.Avi.EnvoyService.Namespace).Get(ctx, ing.Name, metav1.GetOptions{})
	if apierrs.IsNotFound(err) {
		_, err := r.akoclient.AkoV1beta1().HostRules(cfg.Avi.EnvoyService.Namespace).Create(ctx, desired, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ingress: %w", err)
		}
		return nil
	} else if err != nil {
		return err
	}

	original := existing.DeepCopy()
	if !equality.Semantic.DeepEqual(original.Spec, desired.Spec) ||
		!equality.Semantic.DeepEqual(original.Annotations, desired.Annotations) ||
		!equality.Semantic.DeepEqual(original.Labels, desired.Labels) {

		// Don't modify the informers copy.
		original.Spec = desired.Spec
		original.Annotations = desired.Annotations
		original.Labels = desired.Labels
		_, err := r.akoclient.AkoV1beta1().HostRules(original.Namespace).Update(ctx, original, metav1.UpdateOptions{})

		if err != nil {
			// recorder.Eventf(ing, corev1.EventTypeWarning, "UpdateFailed", "Failed to update K8s Ingress: %v", err)
			return fmt.Errorf("failed to update hostrule: %w", err)
		}
	}

	return nil
}

// Reconcile a k8s ingress for each external host in a knative ingress
func (r *Reconciler) reconcileIngress(ctx context.Context, ing *netv1alpha1.Ingress) error {
	cfg := config.FromContext(ctx)

	desired := resources.MakeK8sIngress(ctx, ing)

	if desired == nil {
		if _, err := r.kubeclient.NetworkingV1().Ingresses(cfg.Avi.EnvoyService.Namespace).Get(ctx, ing.Name, metav1.GetOptions{}); apierrs.IsNotFound(err) {
			return nil
		} else if err != nil {
			return fmt.Errorf("failed to get ingress: %w", err)
		}

		if err := r.kubeclient.NetworkingV1().Ingresses(cfg.Avi.EnvoyService.Namespace).Delete(ctx, ing.Name, metav1.DeleteOptions{}); err != nil {
			return fmt.Errorf("failed to delete ingress: %w", err)
		}
		return nil
	}

	existing, err := r.kubeclient.NetworkingV1().Ingresses(cfg.Avi.EnvoyService.Namespace).Get(ctx, ing.Name, metav1.GetOptions{})
	if apierrs.IsNotFound(err) {
		_, err := r.kubeclient.NetworkingV1().Ingresses(cfg.Avi.EnvoyService.Namespace).Create(ctx, desired, metav1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("failed to create ingress: %w", err)
		}
		return nil
	} else if err != nil {
		return err
	}

	original := existing.DeepCopy()
	if !equality.Semantic.DeepEqual(original.Spec, desired.Spec) ||
		!equality.Semantic.DeepEqual(original.Annotations, desired.Annotations) ||
		!equality.Semantic.DeepEqual(original.Labels, desired.Labels) {

		// Don't modify the informers copy.
		original.Spec = desired.Spec
		original.Annotations = desired.Annotations
		original.Labels = desired.Labels
		_, err := r.kubeclient.NetworkingV1().Ingresses(original.Namespace).Update(ctx, original, metav1.UpdateOptions{})

		if err != nil {
			// recorder.Eventf(ing, corev1.EventTypeWarning, "UpdateFailed", "Failed to update K8s Ingress: %v", err)
			return fmt.Errorf("failed to update ingress: %w", err)
		}
	}

	return nil
}
