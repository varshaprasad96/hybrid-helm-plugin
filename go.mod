module github.com/varshaprasad96/hybrid-helm-plugin

go 1.16

require (
	github.com/blang/semver/v4 v4.0.0
	github.com/sirupsen/logrus v1.7.0
	github.com/spf13/afero v1.2.2
	github.com/spf13/pflag v1.0.5
	helm.sh/helm/v3 v3.5.0
	k8s.io/api v0.20.4
	k8s.io/apimachinery v0.20.4
	rsc.io/letsencrypt v0.0.3 // indirect
	sigs.k8s.io/kubebuilder/v3 v3.1.0
	sigs.k8s.io/yaml v1.2.0
)
