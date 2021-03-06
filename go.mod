module github.com/redhat-integration/integration-operator

go 1.15

require (
	github.com/go-logr/logr v0.2.1
	github.com/go-logr/zapr v0.2.0 // indirect
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/operator-framework/api v0.3.17
	github.com/operator-framework/operator-lifecycle-manager v3.11.0+incompatible
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/client-go v0.19.2
	sigs.k8s.io/controller-runtime v0.6.3
)
