package resources

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	k8snetworkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"knative.dev/pkg/kmeta"

	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func MakeK8sIngress(ing *netv1alpha1.Ingress) *k8snetworkingv1.Ingress {
	ingressRules := []k8snetworkingv1.IngressRule{}

	for _, rule := range ing.Spec.Rules {
		if rule.Visibility != netv1alpha1.IngressVisibilityExternalIP {
			continue
		}

		k8sRuleValue := k8snetworkingv1.IngressRuleValue{
			HTTP: &k8snetworkingv1.HTTPIngressRuleValue{
				Paths: []k8snetworkingv1.HTTPIngressPath{
					{
						Path:     "/",
						PathType: ptr.To(k8snetworkingv1.PathTypeExact),
						Backend: k8snetworkingv1.IngressBackend{
							Service: &k8snetworkingv1.IngressServiceBackend{
								Name: "envoy-clusterip",
								Port: k8snetworkingv1.ServiceBackendPort{
									Number: 80,
								},
							},
						},
					},
				},
			},
		}

		for _, host := range rule.Hosts {
			k8sRule := k8snetworkingv1.IngressRule{
				Host:             host,
				IngressRuleValue: k8sRuleValue,
			}
			ingressRules = append(ingressRules, k8sRule)
		}
	}

	k8sIngress := &k8snetworkingv1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name: ing.Name,
			Labels: kmeta.UnionMaps(ing.Labels, map[string]string{
				ParentNameKey:      ing.Name,
				ParentNamespaceKey: ing.Namespace,
				GenerationKey:      fmt.Sprintf("%d", ing.Generation),
				"app":              "gslb",
			}),
			Annotations: kmeta.FilterMap(ing.GetAnnotations(), func(key string) bool {
				return key == corev1.LastAppliedConfigAnnotation
			}),
		},
		Spec: k8snetworkingv1.IngressSpec{
			IngressClassName: ptr.To("avi-lb"),
			Rules:            ingressRules,
		},
	}
	return k8sIngress
}
