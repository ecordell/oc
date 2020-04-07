package catalog

import (
	"fmt"
	"github.com/openshift/oc/pkg/cli/image/imagesource"
	"strings"

	"github.com/alicebob/sqlittle"
	"github.com/docker/distribution/reference"
	"k8s.io/apimachinery/pkg/util/errors"
)

type Mirrorer interface {
	Mirror() (map[string]string, error)
}

// DatabaseExtractor knows how to pull an index image and extract its database
type DatabaseExtractor interface {
	Extract(from imagesource.TypedImageReference) (string, error)
}

type DatabaseExtractorFunc func(from imagesource.TypedImageReference) (string, error)

func (f DatabaseExtractorFunc) Extract(from imagesource.TypedImageReference) (string, error) {
	return f(from)
}

// ImageMirrorer knows how to mirror an image from one registry to another
type ImageMirrorer interface {
	Mirror(mapping map[string]string) error
}

type ImageMirrorerFunc func(mapping map[string]string) error

func (f ImageMirrorerFunc) Mirror(mapping map[string]string) error {
	return f(mapping)
}

type IndexImageMirrorer struct {
	ImageMirrorer     ImageMirrorer
	DatabaseExtractor DatabaseExtractor

	// options
	Source, Dest imagesource.TypedImageReference
}

var _ Mirrorer = &IndexImageMirrorer{}

func NewIndexImageMirror(options ...ImageIndexMirrorOption) (*IndexImageMirrorer, error) {
	config := DefaultImageIndexMirrorerOptions()
	config.Apply(options)
	if err := config.Complete(); err != nil {
		return nil, err
	}
	if err := config.Validate(); err != nil {
		return nil, err
	}
	return &IndexImageMirrorer{
		ImageMirrorer:     config.ImageMirrorer,
		DatabaseExtractor: config.DatabaseExtractor,
		Source:            config.Source,
		Dest:              config.Dest,
	}, nil
}

func (b *IndexImageMirrorer) Mirror() (map[string]string, error) {
	dbFile, err := b.DatabaseExtractor.Extract(b.Source)
	if err != nil {
		return nil, err
	}

	sqlDb, err := sqlittle.Open(dbFile)
	if err != nil {
		return nil, err
	}

	// TODO: what is the minimum required db migration?

	// get all images
	var images = make(map[string]struct{}, 0)
	var errs = make([]error, 0)
	columns := []string{"image"}
	table := "related_image"
	reader := func(r sqlittle.Row) {
		var image string
		if err := r.Scan(&image); err != nil {
			errs = append(errs, err)
			return
		}
		images[image] = struct{}{}
	}
	if err := sqlDb.Select(table, reader, columns...); err != nil {
		errs = append(errs, err)
		return nil, errors.NewAggregate(errs)
	}

	// get all bundlepaths
	columns = []string{"bundlepath"}
	table = "operatorbundle"
	reader = func(r sqlittle.Row) {
		var bundlePath string
		if err := r.Scan(&bundlePath); err != nil {
			errs = append(errs, err)
			return
		}
		images[bundlePath] = struct{}{}
	}
	if err := sqlDb.Select(table, reader, columns...); err != nil {
		errs = append(errs, err)
		return nil, errors.NewAggregate(errs)
	}

	// TODO: build mapping options for quay
	mapping := map[string]string{}
	for img := range images {
		if img == "" {
			continue
		}
		ref, err := reference.ParseNormalizedNamed(img)
		if err != nil {
			errs = append(errs, fmt.Errorf("couldn't mustParse image for mirroring (%s), skipping mirror: %s", img, err.Error()))
			continue
		}
		domain := reference.Domain(ref)
		mapping[ref.String()] = b.Dest.String() + strings.TrimPrefix(ref.String(), domain)
	}

	if err := b.ImageMirrorer.Mirror(mapping); err != nil {
		errs = append(errs, fmt.Errorf("mirroring failed: %s", err.Error()))
	}

	return mapping, errors.NewAggregate(errs)
}
