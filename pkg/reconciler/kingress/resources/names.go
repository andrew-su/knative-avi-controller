package resources

import (
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"
)

func GenerateName(ing *netv1alpha1.Ingress) string {
	return ing.Name
}
