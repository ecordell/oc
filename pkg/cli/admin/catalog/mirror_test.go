package catalog

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestWriteToMapping(t *testing.T) {
	tests := []struct {
		name    string
		mapping map[string]Target
		wantErr bool
		want    string
	}{
		{
			name: "src is tagged",
			mapping: map[string]Target{
				"quay.io/halkyonio/operator:v0.1.8": {
					WithDigest: "",
					WithTag:    "quay.io/olmtest/halkyonio-operator:v0.1.8",
				},
			},
			want: "quay.io/halkyonio/operator:v0.1.8=quay.io/olmtest/halkyonio-operator:v0.1.8\n",
		},
		{
			name: "src has digest",
			mapping: map[string]Target{
				"docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4": {
					WithDigest: "quay.io/olmtest/strimzi-operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4",
					WithTag:    "quay.io/olmtest/strimzi-operator:2b13d275",
				},
			},
			want: "docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4=quay.io/olmtest/strimzi-operator:2b13d275\n",
		},
		{
			name: "multiple",
			mapping: map[string]Target{
				"docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4": {
					WithDigest: "quay.io/olmtest/strimzi-operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4",
					WithTag:    "quay.io/olmtest/strimzi-operator:2b13d275",
				},
				"quay.io/halkyonio/operator:v0.1.8": {
					WithDigest: "",
					WithTag:    "quay.io/olmtest/halkyonio-operator:v0.1.8",
				},
			},
			want: `docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4=quay.io/olmtest/strimzi-operator:2b13d275
quay.io/halkyonio/operator:v0.1.8=quay.io/olmtest/halkyonio-operator:v0.1.8
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := writeToMapping(&buf, tt.mapping); (err != nil) != tt.wantErr {
				t.Errorf("writeToMapping() error = %v, wantErr %v", err, tt.wantErr)
			}
			if buf.String() != tt.want {
				t.Errorf("incorrect mapping. have %s, want %s", buf.String(), tt.want)
			}
		})
	}
}

func TestGenerateICSP(t *testing.T) {
	type args struct {
		name    string
		mapping map[string]Target
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "src is tagged - skip mirror",
			args: args{
				name: "catalog",
				mapping: map[string]Target{
					"quay.io/halkyonio/operator:v0.1.8": {
						WithDigest: "",
						WithTag:    "quay.io/olmtest/halkyonio-operator:v0.1.8",
					},
				},
			},
			want: []byte(
				`apiVersion: operator.openshift.io/v1alpha1
kind: ImageContentSourcePolicy
metadata:
  name: catalog
spec:
  repositoryDigestMirrors: []
`,
			),
		},
		{
			name: "src has digest",
			args: args{
				name: "catalog",
				mapping: map[string]Target{
					"docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4": {
						WithDigest: "quay.io/olmtest/strimzi-operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4",
						WithTag:    "quay.io/olmtest/strimzi-operator:2b13d275",
					},
				},
			},
			want: []byte(
				`apiVersion: operator.openshift.io/v1alpha1
kind: ImageContentSourcePolicy
metadata:
  name: catalog
spec:
  repositoryDigestMirrors:
  - mirrors:
    - quay.io/olmtest/strimzi-operator
    source: docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4
`,
			),
		},
		{
			name: "multiple",
			args: args{
				name: "catalog",
				mapping: map[string]Target{
					"docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4": {
						WithDigest: "quay.io/olmtest/strimzi-operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4",
						WithTag:    "quay.io/olmtest/strimzi-operator:2b13d275",
					},
					"quay.io/halkyonio/operator:v0.1.8": {
						WithDigest: "",
						WithTag:    "quay.io/olmtest/halkyonio-operator:v0.1.8",
					},
				},
			},
			want: []byte(
				`apiVersion: operator.openshift.io/v1alpha1
kind: ImageContentSourcePolicy
metadata:
  name: catalog
spec:
  repositoryDigestMirrors:
  - mirrors:
    - quay.io/olmtest/strimzi-operator
    source: docker.io/strimzi/operator@sha256:d134a9865524c29fcf75bbc4469013bc38d8a15cb5f41acfddb6b9e492f556e4
`,
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateICSP(tt.args.name, tt.args.mapping)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateICSP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("generateICSP() got = %v, want %v, diff = %v", string(got), string(tt.want), cmp.Diff(got, tt.want))
			}
		})
	}
}
