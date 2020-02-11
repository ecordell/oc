module github.com/openshift/oc

go 1.12

require (
	github.com/AaronO/go-git-http v0.0.0-20161214145340-1d9485b3a98f
	github.com/MakeNowJust/heredoc v0.0.0-20171113091838-e9091a26100e
	github.com/RangelReale/osincli v0.0.0-20160924135400-fababb0555f2
	github.com/alexbrainman/sspi v0.0.0-20180613141037-e580b900e9f5
	github.com/apcera/gssapi v0.0.0-00010101000000-000000000000
	github.com/aws/aws-sdk-go v1.17.7
	github.com/blang/semver v3.5.1+incompatible
	github.com/containerd/continuity v0.0.0-20191127005431-f65d91d395eb // indirect
	github.com/containers/image v0.0.0-00010101000000-000000000000
	github.com/containers/storage v0.0.0-20190726081758-912de200380a // indirect
	github.com/coreos/etcd v3.3.15+incompatible
	github.com/davecgh/go-spew v1.1.1
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.4.2-0.20181221150755-2cb26cfe9cbf
	github.com/docker/docker-credential-helpers v0.6.3 // indirect
	github.com/docker/go-units v0.4.0
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7
	github.com/elazarl/goproxy v0.0.0-20190911111923-ecfe977594f1
	github.com/fsnotify/fsnotify v1.4.7
	github.com/fsouza/go-dockerclient v0.0.0-20171004212419-da3951ba2e9e
	github.com/ghodss/yaml v1.0.0
	github.com/gonum/diff v0.0.0-20181124234638-500114f11e71 // indirect
	github.com/gonum/graph v0.0.0-20170401004347-50b27dea7ebb
	github.com/gonum/integrate v0.0.0-20181209220457-a422b5c0fdf2 // indirect
	github.com/gonum/mathext v0.0.0-20181121095525-8a4bf007ea55 // indirect
	github.com/gonum/stat v0.0.0-20181125101827-41a0da705a5b // indirect
	github.com/mitchellh/go-wordwrap v1.0.0
	github.com/moby/buildkit v0.0.0-20181107081847-c3a857e3fca0
	github.com/mtrmac/gpgme v0.0.0-20170102180018-b2432428689c // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/openshift/api v3.9.1-0.20190924102528-32369d4db2ad+incompatible
	github.com/openshift/build-machinery-go v0.0.0-20200210090402-3b072832771e
	github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240
	github.com/openshift/library-go v0.0.0-20200206134157-b4c763d94dcf
	github.com/operator-framework/operator-registry v1.5.4
	github.com/pkg/errors v0.8.1
	github.com/prometheus/client_golang v1.2.1
	github.com/russross/blackfriday v1.5.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	golang.org/x/crypto v0.0.0-20191028145041-f83a4685e152
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553
	golang.org/x/sys v0.0.0-20191028164358-195ce5e7f934
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0
	gopkg.in/ldap.v2 v2.5.1
	k8s.io/api v0.17.1
	k8s.io/apimachinery v0.17.1
	k8s.io/apiserver v0.17.1
	k8s.io/cli-runtime v0.17.0
	k8s.io/client-go v8.0.0+incompatible
	k8s.io/component-base v0.17.1
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.0.0
	k8s.io/kubernetes v1.16.0
	k8s.io/utils v0.0.0-20191114184206-e782cd3c129f
	sigs.k8s.io/yaml v1.1.0
)

replace (
	github.com/apcera/gssapi => github.com/openshift/gssapi v0.0.0-20161010215902-5fb4217df13b
	github.com/containers/image => github.com/openshift/containers-image v0.0.0-20190130162819-76de87591e9d
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309 // Required by Helm from Operator-regsitry
	github.com/openshift/api => github.com/openshift/api v0.0.0-20200206153619-c562e326a6ae // release-4.4
	github.com/openshift/client-go => github.com/openshift/client-go v0.0.0-20200116152001-92a2713fa240 // release-4.4
	github.com/operator-framework/operator-registry => github.com/operator-framework/operator-registry v1.5.9

	k8s.io/api => k8s.io/api v0.17.1
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.1
	k8s.io/apimachinery => github.com/openshift/kubernetes-apimachinery v0.0.0-20191211181342-5a804e65bdc1
	k8s.io/apiserver => k8s.io/apiserver v0.17.1
	k8s.io/cli-runtime => github.com/openshift/kubernetes-cli-runtime v0.0.0-20200114162348-c8810ef308ee
	k8s.io/client-go => github.com/openshift/kubernetes-client-go v0.0.0-20191211181558-5dcabadb2b45
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.17.1
	k8s.io/cluster-bootstrap => k8s.io/cluster-bootstrap v0.17.1
	k8s.io/code-generator => k8s.io/code-generator v0.17.1
	k8s.io/component-base => k8s.io/component-base v0.17.1
	k8s.io/cri-api => k8s.io/cri-api v0.17.1
	k8s.io/csi-translation-lib => k8s.io/csi-translation-lib v0.17.1
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.17.1
	k8s.io/kube-controller-manager => k8s.io/kube-controller-manager v0.17.1
	k8s.io/kube-proxy => k8s.io/kube-proxy v0.17.1
	k8s.io/kube-scheduler => k8s.io/kube-scheduler v0.17.1
	k8s.io/kubectl => github.com/openshift/kubernetes-kubectl v0.0.0-20200211153013-50adac736181
	k8s.io/kubelet => k8s.io/kubelet v0.17.1
	k8s.io/kubernetes => github.com/openshift/kubernetes v1.17.0-alpha.0.0.20191216151305-079984b0a154
	k8s.io/legacy-cloud-providers => k8s.io/legacy-cloud-providers v0.17.1
	k8s.io/metrics => k8s.io/metrics v0.17.1
	k8s.io/node-api => k8s.io/node-api v0.17.1
	k8s.io/sample-apiserver => k8s.io/sample-apiserver v0.17.1
	k8s.io/sample-cli-plugin => k8s.io/sample-cli-plugin v0.17.1
	k8s.io/sample-controller => k8s.io/sample-controller v0.17.1
)
