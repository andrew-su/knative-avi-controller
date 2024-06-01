package hostrule_test

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1beta1"
	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	"knative.dev/pkg/kmeta"
)

type HostRuleOption func(*v1beta1.HostRule)

func HostRule(name string, cfg config.Config, opts ...HostRuleOption) *v1beta1.HostRule {
	hr := &v1beta1.HostRule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   cfg.Avi.EnvoyService.Namespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
	}

	for _, opt := range opts {
		opt(hr)
	}

	return hr
}

func WithLabels(l map[string]string) HostRuleOption {
	return func(hr *v1beta1.HostRule) {
		hr.Labels = kmeta.UnionMaps(hr.Labels, l)
	}
}

func WithAnnotations(ann map[string]string) HostRuleOption {
	return func(hr *v1beta1.HostRule) {
		hr.Annotations = kmeta.UnionMaps(hr.Annotations, ann)
	}
}

func WithBasicSpec(fqdn string) HostRuleOption {
	return func(hr *v1beta1.HostRule) {
		spec := v1beta1.HostRuleSpec{
			VirtualHost: v1beta1.HostRuleVirtualHost{
				EnableVirtualHost: ptr.To(true),
				Fqdn:              fqdn,
				FqdnType:          v1beta1.Exact,
			},
		}
		hr.Spec = spec
	}
}

func WithGSLBConfig(gslbDN string) HostRuleOption {
	return func(hr *v1beta1.HostRule) {
		hr.Spec.VirtualHost.Gslb = v1beta1.HostRuleGSLB{
			Fqdn: gslbDN,
		}
	}
}

func WithDefaultTLS(hr *v1beta1.HostRule) {
	hr.Spec.VirtualHost.TLS = v1beta1.HostRuleTLS{
		SSLKeyCertificate: v1beta1.HostRuleSSLKeyCertificate{
			Name: "System-Default-Cert",
			Type: v1beta1.HostRuleSecretTypeAviReference,
			AlternateCertificate: v1beta1.HostRuleSecret{
				Name: "System-Default-Cert-EC",
				Type: v1beta1.HostRuleSecretTypeAviReference,
			},
		},
		SSLProfile:  "System-Standard",
		Termination: "edge",
	}
}
