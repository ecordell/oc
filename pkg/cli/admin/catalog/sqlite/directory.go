package sqlite

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/openshift/oc/pkg/cli/admin/catalog/registry"
)

const ClusterServiceVersionKind = "ClusterServiceVersion"

type SQLPopulator interface {
	Populate() error
}

// DirectoryLoader loads a directory of resources into the database
type DirectoryLoader struct {
	store     registry.Load
	directory string
}

var _ SQLPopulator = &DirectoryLoader{}

func NewSQLLoaderForDirectory(store registry.Load, directory string) *DirectoryLoader {
	return &DirectoryLoader{
		store:     store,
		directory: directory,
	}
}

func (d *DirectoryLoader) Populate() error {
	log := logrus.WithField("dir", d.directory)

	log.Info("loading Bundles")
	if err := filepath.Walk(d.directory, d.LoadBundleWalkFunc); err != nil {
		return err
	}

	log.Info("loading Packages and Entries")
	if err := filepath.Walk(d.directory, d.LoadPackagesWalkFunc); err != nil {
		return err
	}

	return nil
}

// LoadBundleWalkFunc walks the directory. When it sees a `.clusterserviceversion.yaml` file, it
// attempts to load the surrounding files in the same directory as a bundle, and stores them in the
// db for querying
func (d *DirectoryLoader) LoadBundleWalkFunc(path string, f os.FileInfo, err error) error {
	if f == nil {
		return fmt.Errorf("Not a valid file")
	}

	log := logrus.WithFields(logrus.Fields{"dir": d.directory, "file": f.Name(), "load": "bundles"})

	if f.IsDir() {
		if strings.HasPrefix(f.Name(), ".") {
			log.Info("skipping hidden directory")
			return filepath.SkipDir
		}
		log.Info("directory")
		return nil
	}

	if strings.HasPrefix(f.Name(), ".") {
		log.Info("skipping hidden file")
		return nil
	}

	fileReader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to load file %s: %v", path, err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
	csv := unstructured.Unstructured{}

	if err = decoder.Decode(&csv); err != nil {
		return nil
	}

	if csv.GetKind() != ClusterServiceVersionKind {
		return nil
	}

	log.Info("found csv, loading bundle")

	bundle, err := d.LoadBundle(csv.GetName(), filepath.Dir(path))
	if err != nil {
		return fmt.Errorf("error loading objs in dir: %s", err.Error())
	}

	if bundle == nil || bundle.Size() == 0 {
		log.Warnf("no bundle objects found")
		return nil
	}

	if err := bundle.AllProvidedAPIsInBundle(); err != nil {
		return err
	}

	return d.store.AddOperatorBundle(bundle)
}

// LoadBundle takes the directory that a CSV is in and assumes the rest of the objects in that directory
// are part of the bundle.
func (d *DirectoryLoader) LoadBundle(csvName string, dir string) (*registry.Bundle, error) {
	bundle := &registry.Bundle{}
	log := logrus.WithFields(logrus.Fields{"dir": d.directory, "load": "bundle"})
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	// time="2019-07-19T07:50:30-04:00" level=info msg="found csv, loading bundle" dir=/var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830 file=akka-cluster-operator.v0.0.1.clusterserviceversion.yaml load=bundles
	// time="2019-07-19T07:50:30-04:00" level=info msg="loading bundle file" dir=/var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830 file=AkkaCluster-v1alpha1.crd.yaml load=bundle
	// time="2019-07-19T07:50:30-04:00" level=info msg="loading bundle file" dir=/var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830 file=akka-cluster-operator.package.yaml load=bundle
	// time="2019-07-19T07:50:30-04:00" level=info msg="could not decode contents of file /var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830/akka-cluster-operator/akka-cluster-operator.package.yaml into file: error unmarshaling JSON: Object 'Kind' is missing in '{\"channels\":[{\"currentCSV\":\"akka-cluster-operator.v0.0.1\",\"name\":\"alpha\"}],\"defaultChannel\":\"alpha\",\"packageName\":\"akka-cluster-operator\"}'" dir=/var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830 file=akka-cluster-operator.package.yaml load=bundle
	// time="2019-07-19T07:50:30-04:00" level=warning msg="no bundle objects found" dir=/var/folders/4m/pt431q9d2fsbk9zz2_vmg3tc0000gn/T/catalog-110954830 file=akka-cluster-operator.v0.0.1.clusterserviceversion.yaml load=bundles

	for _, f := range files {
		log = log.WithField("file", f.Name())
		if f.IsDir() {
			log.Info("skipping directory")
			continue
		}

		if strings.HasPrefix(f.Name(), ".") {
			log.Info("skipping hidden file")
			continue
		}

		log.Info("loading bundle file")
		path := filepath.Join(dir, f.Name())
		fileReader, err := os.Open(path)
		if err != nil {
			return nil, fmt.Errorf("unable to load file %s: %v", path, err)
		}
		decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
		obj := &unstructured.Unstructured{}

		if err = decoder.Decode(obj); err != nil {
			log.Infof("could not decode contents of file %s into file: %v", path, err)
			continue
		}

		// Don't include other CSVs in the bundle
		if obj.GetKind() == "ClusterServiceVersion" && obj.GetName() != csvName {
			continue
		}

		if obj.Object != nil {
			bundle.Add(obj)
		}

	}
	return bundle, nil
}

func (d *DirectoryLoader) LoadPackagesWalkFunc(path string, f os.FileInfo, err error) error {
	log := logrus.WithFields(logrus.Fields{"dir": d.directory, "file": f.Name(), "load": "package"})
	if f == nil {
		return fmt.Errorf("Not a valid file")
	}
	if f.IsDir() {
		if strings.HasPrefix(f.Name(), ".") {
			log.Info("skipping hidden directory")
			return filepath.SkipDir
		}
		log.Info("directory")
		return nil
	}

	if strings.HasPrefix(f.Name(), ".") {
		log.Info("skipping hidden file")
		return nil
	}

	fileReader, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("unable to load package from file %s: %v", path, err)
	}

	decoder := yaml.NewYAMLOrJSONDecoder(fileReader, 30)
	manifest := registry.PackageManifest{}
	if err = decoder.Decode(&manifest); err != nil {
		log.Infof("could not decode contents of file %s into package: %v", path, err)
		return nil
	}
	if manifest.PackageName == "" {
		return nil
	}

	if err := d.store.AddPackageChannels(manifest); err != nil {
		return fmt.Errorf("error loading package into db: %s", err.Error())
	}

	return nil
}
