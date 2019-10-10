/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

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
	pkgs        []*v1alpha2.Pkg
}

func (fspl *FilesystemPackageLoader) walkFunc() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
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
			fspl.pkgs = append(fspl.pkgs, pkg)
		}

		return nil
	}
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

	fspl.pkgs = nil

	err = filepath.Walk(fspl.Root, fspl.walkFunc())

	return fspl.pkgs, err
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

	defer f.Close() //nolint: errcheck

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return v1alpha2.NewPkg(filepath.Dir(basePath), contents, fspl.Context)
}
