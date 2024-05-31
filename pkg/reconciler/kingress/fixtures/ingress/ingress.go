package ingress_test

import (
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	"knative.dev/pkg/kmeta"
)

type IngressOption func(*v1.Ingress)

func Ingress(name string, cfg config.Config, opts ...IngressOption) *v1.Ingress {
	i := &v1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   cfg.Avi.EnvoyService.Namespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: v1.IngressSpec{
			IngressClassName: &cfg.Avi.AviIngressClassName,
			Rules:            []v1.IngressRule{},
		},
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func WithGeneration(gen int64) IngressOption {
	return func(i *v1.Ingress) {
		i.Generation = gen
	}
}

func WithLabels(l map[string]string) IngressOption {
	return func(i *v1.Ingress) {
		i.Labels = kmeta.UnionMaps(i.Labels, l)
	}
}

func WithAnnotations(ann map[string]string) IngressOption {
	return func(i *v1.Ingress) {
		i.Annotations = kmeta.UnionMaps(i.Annotations, ann)
	}
}

func WithRule(host string, cfg config.Config) IngressOption {
	return func(i *v1.Ingress) {
		rule := v1.IngressRule{
			Host: host,
			IngressRuleValue: v1.IngressRuleValue{
				HTTP: &v1.HTTPIngressRuleValue{
					Paths: []v1.HTTPIngressPath{
						{
							Path:     "/",
							PathType: ptr.To(v1.PathTypeExact),
							Backend: v1.IngressBackend{
								Service: &v1.IngressServiceBackend{
									Name: cfg.Avi.EnvoyService.Name,
									Port: v1.ServiceBackendPort{
										Number: 80,
									},
								},
							},
						},
					},
				},
			},
		}

		i.Spec.Rules = append(i.Spec.Rules, rule)
	}
}
