package resources

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8snetworkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	"knative.dev/pkg/kmeta"

	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func MakeK8sIngress(ctx context.Context, ing *netv1alpha1.Ingress) *k8snetworkingv1.Ingress {
	cfg := config.FromContext(ctx)

	ingressClass := cfg.Avi.AviIngressClassName
	envoyServiceName := cfg.Avi.EnvoyService.Name
	envoyServiceNamespace := cfg.Avi.EnvoyService.Namespace // The resource must live in the same namespace as the service we target.

	hostname := ""

	for _, rule := range ing.Spec.Rules {
		// We only care about external services
		if rule.Visibility != netv1alpha1.IngressVisibilityExternalIP {
			continue
		}

		// We only need one "external" hostname
		hostname = rule.Hosts[0]
	}

	if hostname == "" {
		return nil
	}

	k8sIngress := &k8snetworkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ing.Name,
			Namespace: envoyServiceNamespace,
			Labels: kmeta.UnionMaps(ing.Labels, map[string]string{
				ParentNameKey:      ing.Name,
				ParentNamespaceKey: ing.Namespace,
				GenerationKey:      fmt.Sprintf("%d", ing.Generation),
				// "app":              "gslb",
			}),
			Annotations: kmeta.FilterMap(ing.GetAnnotations(), func(key string) bool {
				return key == corev1.LastAppliedConfigAnnotation
			}),
		},
		Spec: k8snetworkingv1.IngressSpec{
			IngressClassName: &ingressClass,
			Rules: []k8snetworkingv1.IngressRule{
				{
					Host: hostname,
					IngressRuleValue: k8snetworkingv1.IngressRuleValue{
						HTTP: &k8snetworkingv1.HTTPIngressRuleValue{
							Paths: []k8snetworkingv1.HTTPIngressPath{
								{
									Path:     "/",
									PathType: ptr.To(k8snetworkingv1.PathTypeExact),
									Backend: k8snetworkingv1.IngressBackend{
										Service: &k8snetworkingv1.IngressServiceBackend{
											Name: envoyServiceName,
											Port: k8snetworkingv1.ServiceBackendPort{
												Number: 80,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return k8sIngress
}
