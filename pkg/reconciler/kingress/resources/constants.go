package resources

const (
	// GenerationKey holds the generation of the parent KIngress resource that the Ingress and
	// HostRule specs are derived from.
	GenerationKey = "avi-controller/generation"

	// ParentNameKey hold the name of the parent KIngress resource, since OwnerReferences cannot
	// be used for resources across namespace.
	ParentNameKey = "avi-controller/parent.name"

	// ParentNamespaceKey hold the namespace of the parent KIngress resource, since OwnerReferences
	// cannot be used for resources across namespace.
	ParentNamespaceKey = "avi-controller/parent.namespace"
)
