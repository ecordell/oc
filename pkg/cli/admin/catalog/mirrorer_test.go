package catalog

import (
	"testing"

	"github.com/openshift/oc/pkg/cli/image/imagesource"
)

func existingExtractor(dir string) DatabaseExtractorFunc {
	return func(from imagesource.TypedImageReference) (s string, e error) {
		return dir, nil
	}
}

func noopMirror(map[string]string) error {
	return nil
}

func mustParse(t *testing.T, img string) imagesource.TypedImageReference {
	imgRef, err := imagesource.ParseReference(img)
	if err != nil {
		t.Errorf("couldn't parse image ref %s: %v", img, err)
	}
	return imgRef
}

func TestMirror(t *testing.T) {
	type fields struct {
		ImageMirrorer     ImageMirrorerFunc
		DatabaseExtractor DatabaseExtractorFunc
		Source            imagesource.TypedImageReference
		Dest              imagesource.TypedImageReference
	}
	tests := []struct {
		name    string
		fields  fields
		want    map[string]string
		wantErr error
	}{
		{
			name: "maps related images and bundle images",
			fields: fields{
				ImageMirrorer:     noopMirror,
				DatabaseExtractor: existingExtractor("testdata/test.db"),
				Source:            mustParse(t, "quay.io/example/image:tag"),
				Dest:              mustParse(t, "localhost:5000"),
			},
			want: map[string]string{
				"quay.io/test/prometheus.0.14.0": "localhost:5000/test/prometheus.0.14.0",
				"quay.io/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8": "localhost:5000/coreos/etcd-operator@sha256:db563baa8194fcfe39d1df744ed70024b0f1f9e9b55b5923c2f3a413c44dc6b8",
				"quay.io/test/etcd.0.9.0": "localhost:5000/test/etcd.0.9.0",
				"quay.io/coreos/prometheus-operator@sha256:0e92dd9b5789c4b13d53e1319d0a6375bcca4caaf0d698af61198061222a576d": "localhost:5000/coreos/prometheus-operator@sha256:0e92dd9b5789c4b13d53e1319d0a6375bcca4caaf0d698af61198061222a576d",
				"quay.io/coreos/prometheus-operator@sha256:3daa69a8c6c2f1d35dcf1fe48a7cd8b230e55f5229a1ded438f687debade5bcf": "localhost:5000/coreos/prometheus-operator@sha256:3daa69a8c6c2f1d35dcf1fe48a7cd8b230e55f5229a1ded438f687debade5bcf",
				"quay.io/test/prometheus.0.22.2": "localhost:5000/test/prometheus.0.22.2",
				"quay.io/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2": "localhost:5000/coreos/etcd-operator@sha256:c0301e4686c3ed4206e370b42de5a3bd2229b9fb4906cf85f3f30650424abec2",
				"quay.io/coreos/prometheus-operator@sha256:5037b4e90dbb03ebdefaa547ddf6a1f748c8eeebeedf6b9d9f0913ad662b5731": "localhost:5000/coreos/prometheus-operator@sha256:5037b4e90dbb03ebdefaa547ddf6a1f748c8eeebeedf6b9d9f0913ad662b5731",
				"quay.io/test/etcd.0.9.2": "localhost:5000/test/etcd.0.9.2",
				"quay.io/test/prometheus.0.15.0": "localhost:5000/test/prometheus.0.15.0",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &IndexImageMirrorer{
				ImageMirrorer:     tt.fields.ImageMirrorer,
				DatabaseExtractor: tt.fields.DatabaseExtractor,
				Source:            tt.fields.Source,
				Dest:              tt.fields.Dest,
			}
			got, err := b.Mirror()
			if tt.wantErr != nil && tt.wantErr != err {
				t.Errorf("wanted err %v but got %v", tt.wantErr, err)
			}
			if tt.wantErr == nil && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			for k, v := range tt.want {
				w, ok := got[k]
				if !ok {
					t.Errorf("couldn't find wanted key %s", k)
					continue
				}
				if w != v {
					t.Errorf("incorrect mapping for %s. have %s, want %s", k, w, v)
				}
			}
			for k, v := range got {
				w, ok := tt.want[k]
				if !ok {
					t.Errorf("got unexpected key %s", k)
					continue
				}
				if w != v {
					t.Errorf("incorrect mapping for %s. have %s, want %s", k, v, w)
				}
			}
		})
	}
}
