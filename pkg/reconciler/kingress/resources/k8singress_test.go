package resources_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	v1 "k8s.io/api/networking/v1"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	. "knative.dev/avi-controller/pkg/reconciler/kingress/resources"
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"

	fixtures_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures"
	ingress_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/ingress"
	kingress_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/kingress"
)

func TestMakeK8sIngress_Labels(t *testing.T) {
	ctx := config.ToContext(context.Background(), &fixtures_test.DefaultConfig)

	testcases := []struct {
		name     string
		input    *netv1alpha1.Ingress
		expected *v1.Ingress
	}{
		{
			name: "no new labels",
			input: kingress_test.Ingress("river", "water",
				kingress_test.WithGeneration(5),
				kingress_test.WithRules(kingress_test.BasicRule("foo.bar.example.com", netv1alpha1.IngressVisibilityExternalIP)),
			),
			expected: ingress_test.Ingress("river", fixtures_test.DefaultConfig,
				ingress_test.WithLabels(map[string]string{
					ParentNameKey:      "river",
					ParentNamespaceKey: "water",
					GenerationKey:      "5",
				}),
			),
		},
		{
			name: "with additional labels",
			input: kingress_test.Ingress("river", "water",
				kingress_test.WithGeneration(5),
				kingress_test.WithLabels(map[string]string{
					"foo": "bar",
				}),
				kingress_test.WithRules(kingress_test.BasicRule("foo.bar.example.com", netv1alpha1.IngressVisibilityExternalIP)),
			),
			expected: ingress_test.Ingress("river", fixtures_test.DefaultConfig,
				ingress_test.WithLabels(map[string]string{
					ParentNameKey:      "river",
					ParentNamespaceKey: "water",
					GenerationKey:      "5",
					"foo":              "bar",
				}),
			),
		},
	}

	for _, tc := range testcases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			output := MakeK8sIngress(ctx, tc.input)
			if !cmp.Equal(tc.expected.Labels, output.Labels) {
				t.Error("MakeK8sIngress labels (-want, +got) =", cmp.Diff(tc.expected.Labels, output.Labels))
			}
		})
	}
}

func TestMakeK8sIngress_HasSingleExternal(t *testing.T) {
	ctx := config.ToContext(context.Background(), &fixtures_test.DefaultConfig)
	ingress := kingress_test.Ingress("river", "water",
		kingress_test.WithGeneration(4),
		kingress_test.WithRules(kingress_test.BasicRule("foo.default.example.com", netv1alpha1.IngressVisibilityExternalIP)),
	)

	expected := ingress_test.Ingress("river", fixtures_test.DefaultConfig,
		ingress_test.WithLabels(map[string]string{
			ParentNameKey:      "river",
			ParentNamespaceKey: "water",
			GenerationKey:      "4",
		}),
		ingress_test.WithRule("foo.default.example.com", fixtures_test.DefaultConfig),
	)

	k8singress := MakeK8sIngress(ctx, ingress)

	if !cmp.Equal(expected, k8singress) {
		t.Error("MakeK8sIngress (-want, +got) =", cmp.Diff(expected, k8singress))
	}
}
