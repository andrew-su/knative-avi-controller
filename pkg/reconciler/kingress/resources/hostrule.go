package resources

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/utils/ptr"

	aviv1beta1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/kmeta"
)

func MakeHostRule(ctx context.Context, ing *netv1alpha1.Ingress) *aviv1beta1.HostRule {
	cfg := config.FromContext(ctx)

	namespace := cfg.Avi.EnvoyService.Namespace

	var hostname string
	for _, rule := range ing.Spec.Rules {
		if rule.Visibility != netv1alpha1.IngressVisibilityExternalIP {
			continue
		}

		hostname = rule.Hosts[0]
	}

	if hostname == "" {
		return nil
	}

	hostRule := &aviv1beta1.HostRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:      ing.Name,
			Namespace: namespace,
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
		Spec: aviv1beta1.HostRuleSpec{
			VirtualHost: aviv1beta1.HostRuleVirtualHost{
				FqdnType:          aviv1beta1.Exact,
				Fqdn:              hostname,
				EnableVirtualHost: ptr.To(true),
				// Gslb: aviv1beta1.HostRuleGSLB{
				// 	Fqdn: "onion2.suan.tanzu.biz",
				// },
				TLS: aviv1beta1.HostRuleTLS{
					SSLKeyCertificate: aviv1beta1.HostRuleSSLKeyCertificate{
						Name: "System-Default-Cert",
						Type: aviv1beta1.HostRuleSecretTypeAviReference,
						AlternateCertificate: aviv1beta1.HostRuleSecret{
							Name: "System-Default-Cert-EC",
							Type: aviv1beta1.HostRuleSecretTypeAviReference,
						},
					},
					SSLProfile:  "System-Standard",
					Termination: "edge",
				},
			},
		},
	}
	return hostRule
}
