/*
Copyright 2021 The Knative Authors

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

package config

import (
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"knative.dev/pkg/configmap"
)

const (
	// GatewayConfigName is the config map name for the gateway configuration.
	AviConfigName = "config-avi"

	aviIngressClassNameKey = "avi-ingress-class"
	envoyServiceKey        = "envoy-service"
	gslbDomainLabelsKey    = "gslb-domain-labels"
	gslbSelectorKey        = "gslb-selector"
)

// Avi contains configuration defined in the config map.
type Avi struct {
	AviIngressClassName string
	EnvoyService        *types.NamespacedName
}

// FromConfigMap creates a Avi config from the supplied ConfigMap
func FromConfigMap(cm *corev1.ConfigMap) (*Avi, error) {
	var (
		err    error
		config = &Avi{
			AviIngressClassName: "avi-lb",
			EnvoyService: &types.NamespacedName{
				Namespace: "tanzu-system-ingress",
				Name:      "envoy-clusterip",
			},
		}
	)

	err = configmap.Parse(cm.Data,
		configmap.AsString(config.AviIngressClassName, &config.AviIngressClassName),
		configmap.AsNamespacedName(envoyServiceKey, config.EnvoyService),
	)

	return config, err
}
