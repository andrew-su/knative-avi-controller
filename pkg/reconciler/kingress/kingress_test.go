package kingress

import (
	"context"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"

	fakeingressclient "knative.dev/networking/pkg/client/injection/client/fake"
	kingressreconciler "knative.dev/networking/pkg/client/injection/reconciler/networking/v1alpha1/ingress"
	fakekubeclient "knative.dev/pkg/client/injection/kube/client/fake"

	k8snetworkingv1 "k8s.io/api/networking/v1"
	"knative.dev/networking/pkg/apis/networking"
	"knative.dev/networking/pkg/apis/networking/v1alpha1"

	clientgotesting "k8s.io/client-go/testing"
	. "knative.dev/avi-controller/pkg/reconciler/testing"
	. "knative.dev/pkg/reconciler/testing"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	fixtures_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures"
	ingress_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/ingress"
	kingress_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/kingress"
	"knative.dev/avi-controller/pkg/reconciler/kingress/resources"
)

func TestReconciler(t *testing.T) {
	table := TableTest{
		{
			Name: "bad workqueue key",
			Key:  "too/many/parts",
		}, {
			Name: "key not found",
			Key:  "foo/not-found",
		}, {
			Name: "skip ingress not matching class key",
			Key:  "ns/name",
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: "fake-controller",
					}),
					kingress_test.WithRules(kingress_test.BasicRule("lazy.my-ns.example.com", v1alpha1.IngressVisibilityExternalIP)),
				),
			},
		}, {
			Name:                    "single external rule",
			Key:                     "meadow/butterfly",
			SkipNamespaceValidation: true,
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					kingress_test.WithGeneration(4),
					kingress_test.WithRules(kingress_test.BasicRule("fuzzy.my-ns.example.com", v1alpha1.IngressVisibilityExternalIP)),
					kingress_test.WithInitialConditions,
				),
			},
			WantCreates: []runtime.Object{
				ingress_test.Ingress("butterfly", fixtures_test.DefaultConfig,
					ingress_test.WithLabels(map[string]string{
						resources.ParentNameKey:      "butterfly",
						resources.ParentNamespaceKey: "meadow",
						resources.GenerationKey:      "4",
					}),
					ingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					ingress_test.WithRule("fuzzy.my-ns.example.com", fixtures_test.DefaultConfig),
				),
			},
			WantPatches: []clientgotesting.PatchActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: "meadow",
				},
				Name:  "butterfly",
				Patch: []byte(`{"metadata":{"finalizers":["avi.ingresses.networking.internal.knative.dev"],"resourceVersion":""}}`),
			}},
			WantEvents: []string{
				Eventf(corev1.EventTypeNormal, "FinalizerUpdate", `Updated "butterfly" finalizers`),
			},
		}, {
			Name:                    "multiple external rule",
			Key:                     "meadow/butterfly",
			SkipNamespaceValidation: true,
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					kingress_test.WithGeneration(4),
					kingress_test.WithRules(
						kingress_test.BasicRule("fuzzy.my-ns.example.com", v1alpha1.IngressVisibilityExternalIP),
						kingress_test.BasicRule("muddy.my-ns.example.com", v1alpha1.IngressVisibilityExternalIP),
					),
					kingress_test.WithInitialConditions,
				),
			},
			WantCreates: []runtime.Object{
				ingress_test.Ingress("butterfly", fixtures_test.DefaultConfig,
					ingress_test.WithLabels(map[string]string{
						resources.ParentNameKey:      "butterfly",
						resources.ParentNamespaceKey: "meadow",
						resources.GenerationKey:      "4",
					}),
					ingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					ingress_test.WithRule("muddy.my-ns.example.com", fixtures_test.DefaultConfig),
				),
			},
			WantPatches: []clientgotesting.PatchActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: "meadow",
				},
				Name:  "butterfly",
				Patch: []byte(`{"metadata":{"finalizers":["avi.ingresses.networking.internal.knative.dev"],"resourceVersion":""}}`),
			}},
			WantEvents: []string{
				Eventf(corev1.EventTypeNormal, "FinalizerUpdate", `Updated "butterfly" finalizers`),
			},
		}, {
			Name:                    "cluster local only rule",
			Key:                     "meadow/butterfly",
			SkipNamespaceValidation: true,
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					kingress_test.WithGeneration(4),
					kingress_test.WithRules(kingress_test.BasicRule("fuzzy.svc.cluster.local", v1alpha1.IngressVisibilityClusterLocal)),
					kingress_test.WithInitialConditions,
				),
			},
			WantPatches: []clientgotesting.PatchActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: "meadow",
				},
				Name:  "butterfly",
				Patch: []byte(`{"metadata":{"finalizers":["avi.ingresses.networking.internal.knative.dev"],"resourceVersion":""}}`),
			}},
			WantEvents: []string{
				Eventf(corev1.EventTypeNormal, "FinalizerUpdate", `Updated "butterfly" finalizers`),
			},
		}, {
			Name:                    "change to clusterlocal",
			Key:                     "meadow/butterfly",
			SkipNamespaceValidation: true,
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					kingress_test.WithGeneration(4),
					kingress_test.WithRules(kingress_test.BasicRule("fuzzy.svc.cluster.local", v1alpha1.IngressVisibilityClusterLocal)),
					kingress_test.WithInitialConditions,
				),
				ingress_test.Ingress("butterfly", fixtures_test.DefaultConfig,
					ingress_test.WithLabels(map[string]string{
						resources.ParentNameKey:      "butterfly",
						resources.ParentNamespaceKey: "meadow",
						resources.GenerationKey:      "4",
					}),
					ingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					ingress_test.WithRule("fuzzy.my-ns.example.com", fixtures_test.DefaultConfig),
				),
			},
			WantDeletes: []clientgotesting.DeleteActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: fixtures_test.DefaultConfig.Avi.EnvoyService.Namespace,
					Resource:  k8snetworkingv1.SchemeGroupVersion.WithResource("ingresses"),
				},
				Name: "butterfly",
			}},
			WantPatches: []clientgotesting.PatchActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: "meadow",
				},
				Name:  "butterfly",
				Patch: []byte(`{"metadata":{"finalizers":["avi.ingresses.networking.internal.knative.dev"],"resourceVersion":""}}`),
			}},
			WantEvents: []string{
				Eventf(corev1.EventTypeNormal, "FinalizerUpdate", `Updated "butterfly" finalizers`),
			},
		}, {
			Name:                    "update external kingress",
			Key:                     "meadow/butterfly",
			SkipNamespaceValidation: true,
			Objects: []runtime.Object{
				kingress_test.Ingress("butterfly", "meadow",
					kingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					kingress_test.WithGeneration(4),
					kingress_test.WithRules(kingress_test.BasicRule("fuzzy.my-ns.fancy.com", v1alpha1.IngressVisibilityExternalIP)),
					kingress_test.WithInitialConditions,
				),
				ingress_test.Ingress("butterfly", fixtures_test.DefaultConfig,
					ingress_test.WithLabels(map[string]string{
						resources.ParentNameKey:      "butterfly",
						resources.ParentNamespaceKey: "meadow",
						resources.GenerationKey:      "4",
					}),
					ingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					ingress_test.WithRule("fuzzy.my-ns.example.com", fixtures_test.DefaultConfig),
				),
			},
			WantUpdates: []clientgotesting.UpdateActionImpl{{
				Object: ingress_test.Ingress("butterfly", fixtures_test.DefaultConfig,
					ingress_test.WithLabels(map[string]string{
						resources.ParentNameKey:      "butterfly",
						resources.ParentNamespaceKey: "meadow",
						resources.GenerationKey:      "4",
					}),
					ingress_test.WithAnnotations(map[string]string{
						networking.IngressClassAnnotationKey: IngressClassName,
					}),
					ingress_test.WithRule("fuzzy.my-ns.fancy.com", fixtures_test.DefaultConfig),
				),
			}},
			WantPatches: []clientgotesting.PatchActionImpl{{
				ActionImpl: clientgotesting.ActionImpl{
					Namespace: "meadow",
				},
				Name:  "butterfly",
				Patch: []byte(`{"metadata":{"finalizers":["avi.ingresses.networking.internal.knative.dev"],"resourceVersion":""}}`),
			}},
			WantEvents: []string{
				Eventf(corev1.EventTypeNormal, "FinalizerUpdate", `Updated "butterfly" finalizers`),
			},
		},
	}

	table.Test(t, MakeFactory(func(ctx context.Context, listers *Listers, cmw configmap.Watcher) controller.Reconciler {
		r := &Reconciler{
			kubeclient:       fakekubeclient.Get(ctx),
			k8sIngressLister: listers.GetIngressLister(),
		}

		ingr := kingressreconciler.NewReconciler(ctx, logging.FromContext(ctx), fakeingressclient.Get(ctx),
			listers.GetKIngressLister(), controller.GetEventRecorder(ctx), r, IngressClassName, controller.Options{
				SkipStatusUpdates: true,
				FinalizerName:     FinalizerName,
				ConfigStore: &testConfigStore{
					config: &fixtures_test.DefaultConfig,
				}},
		)

		return ingr
	}))
}

type testConfigStore struct {
	config *config.Config
}

func (t *testConfigStore) ToContext(ctx context.Context) context.Context {
	return config.ToContext(ctx, t.config)
}
