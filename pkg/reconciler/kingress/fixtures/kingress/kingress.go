package kingress_test

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"
	"knative.dev/pkg/kmeta"
)

type IngressOption func(*v1alpha1.Ingress)

func Ingress(name string, namespace string, opts ...IngressOption) *v1alpha1.Ingress {
	i := &v1alpha1.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   namespace,
			Labels:      map[string]string{},
			Annotations: map[string]string{},
		},
		Spec: v1alpha1.IngressSpec{
			Rules: []v1alpha1.IngressRule{},
		},
	}

	for _, opt := range opts {
		opt(i)
	}

	return i
}

func WithGeneration(gen int64) IngressOption {
	return func(i *v1alpha1.Ingress) {
		i.Generation = gen
	}
}

func WithLabels(l map[string]string) IngressOption {
	return func(i *v1alpha1.Ingress) {
		i.Labels = kmeta.UnionMaps(i.Labels, l)
	}
}

func WithAnnotations(ann map[string]string) IngressOption {
	return func(i *v1alpha1.Ingress) {
		i.Annotations = kmeta.UnionMaps(i.Annotations, ann)
	}
}

func WithRules(rules ...v1alpha1.IngressRule) IngressOption {
	return func(i *v1alpha1.Ingress) {
		i.Spec.Rules = append(i.Spec.Rules, rules...)
	}
}

func WithInitialConditions(i *v1alpha1.Ingress) {
	i.Status.InitializeConditions()
}

type RuleOption func(*v1alpha1.IngressRule)

func BasicRule(host string, visibility v1alpha1.IngressVisibility, opts ...RuleOption) v1alpha1.IngressRule {
	rule := v1alpha1.IngressRule{
		Hosts: []string{
			host,
		},
		Visibility: visibility,
		HTTP: &v1alpha1.HTTPIngressRuleValue{
			Paths: []v1alpha1.HTTPIngressPath{},
		},
	}

	for _, opt := range opts {
		opt(&rule)
	}

	return rule
}
