package solver

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// FilesystemPackageLoader loads packages by walking file system tree
type FilesystemPackageLoader struct {
	*log.Logger
	Root    string
	Context types.Variables

	absRootPath string
}

// Load implements PackageLoader
func (fspl *FilesystemPackageLoader) Load() ([]*v1alpha2.Pkg, error) {
	if fspl.Logger == nil {
		fspl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

	if fspl.Root == "" {
		fspl.Root = "."
	}

	var err error
	fspl.absRootPath, err = filepath.Abs(fspl.Root)
	if err != nil {
		return nil, err
	}

	var pkgs []*v1alpha2.Pkg

	err = filepath.Walk(fspl.Root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fspl.Logger.Printf("error walking %q: %s", path, err)
			return nil
		}

		if info.Name() != "." && strings.HasPrefix(info.Name(), ".") && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		if info.Name() == constants.PkgYaml {
			pkg, e := fspl.loadPkg(path)
			if e != nil {
				fspl.Logger.Printf("skipping %q: %s", path, e)
				return nil
			}
			fspl.Logger.Printf("loaded pkg %q from %q", pkg.Name, path)
			pkgs = append(pkgs, pkg)
		}

		return nil
	})

	return pkgs, err
}

func (fspl *FilesystemPackageLoader) loadPkg(path string) (*v1alpha2.Pkg, error) {
	absFile, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	basePath, err := filepath.Rel(fspl.absRootPath, absFile)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return v1alpha2.NewPkg(filepath.Dir(basePath), contents, fspl.Context)
}
