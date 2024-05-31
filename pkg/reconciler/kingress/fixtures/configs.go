package fixtures_test

import (
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
)

var (
	DefaultConfig = config.Config{
		Avi: &config.Avi{
			AviIngressClassName: "avi-lb",
			EnvoyService: &types.NamespacedName{
				Namespace: "where-envoy-lives",
				Name:      "envoy-service-name",
			},
		},
	}
)
