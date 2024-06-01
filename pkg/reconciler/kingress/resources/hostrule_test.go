package resources_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"

	"knative.dev/avi-controller/pkg/reconciler/kingress/config"
	fixtures_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures"
	hostrule_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/hostrule"
	kingress_test "knative.dev/avi-controller/pkg/reconciler/kingress/fixtures/kingress"
	. "knative.dev/avi-controller/pkg/reconciler/kingress/resources"
	netv1alpha1 "knative.dev/networking/pkg/apis/networking/v1alpha1"

	aviv1beta1 "github.com/vmware/load-balancer-and-ingress-services-for-kubernetes/pkg/apis/ako/v1beta1"
)

func TestMakeHostRule_Labels(t *testing.T) {
	ctx := config.ToContext(context.Background(), &fixtures_test.DefaultConfig)

	testcases := []struct {
		name     string
		input    *netv1alpha1.Ingress
		expected *aviv1beta1.HostRule
	}{
		{
			name: "default labels",
			input: kingress_test.Ingress("river", "water",
				kingress_test.WithGeneration(5),
				kingress_test.WithRules(kingress_test.BasicRule("foo.bar.example.com", netv1alpha1.IngressVisibilityExternalIP)),
			),
			expected: hostrule_test.HostRule("river", fixtures_test.DefaultConfig,
				hostrule_test.WithLabels(map[string]string{
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
			expected: hostrule_test.HostRule("river", fixtures_test.DefaultConfig,
				hostrule_test.WithLabels(map[string]string{
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

			output := MakeHostRule(ctx, tc.input)
			if !cmp.Equal(tc.expected.Labels, output.Labels) {
				t.Error("MakeHostRule labels (-want, +got) =", cmp.Diff(tc.expected.Labels, output.Labels))
			}
		})
	}
}

func TestMakeHostRule_HasSingleExternal(t *testing.T) {
	ctx := config.ToContext(context.Background(), &fixtures_test.DefaultConfig)
	ingress := kingress_test.Ingress("river", "water",
		kingress_test.WithGeneration(4),
		kingress_test.WithRules(kingress_test.BasicRule("foo.default.example.com", netv1alpha1.IngressVisibilityExternalIP)),
	)

	expected := hostrule_test.HostRule("river", fixtures_test.DefaultConfig,
		hostrule_test.WithBasicSpec("foo.default.example.com"),
		hostrule_test.WithLabels(map[string]string{
			ParentNameKey:      "river",
			ParentNamespaceKey: "water",
			GenerationKey:      "4",
		}),
		hostrule_test.WithDefaultTLS,
	)

	hostrule := MakeHostRule(ctx, ingress)

	if !cmp.Equal(expected, hostrule) {
		t.Error("MakeHostRule (-want, +got) =", cmp.Diff(expected, hostrule))
	}
}
