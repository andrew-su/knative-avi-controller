# Knative Avi Controller

Knative `avi-controller` reconciles knative's `Ingress` resource to generate the
appropriate kubernetes' `Ingress` and Avi's `HostRule` resources.

Usage requirements:
* Knative deployed with an Ingress which runs on the same cluster.
* Service pointing to the ingress should be running in ClusterIP mode.
* Avi's AKO and AMKO is installed and configured correctly.

TODO:
* Add configuration for GSLB configuration.

If you are interested in contributing, see [CONTRIBUTING.md](./CONTRIBUTING.md)
and [DEVELOPMENT.md](./DEVELOPMENT.md).
