package namespace

var namespace string

// SetNamespace sets the global namespace
func SetNamespace(ns string) {
	namespace = ns
}

// GetNamespace returns the namespace variable
func GetNamespace() string {
	return namespace
}
