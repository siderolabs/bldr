/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
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
	multiErr    *multierror.Error
	pkgFile     *v1alpha2.Pkgfile
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
				fspl.Logger.Printf("error loading %q: %s", path, e)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error loading %q: %w", path, e))

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

	if err = fspl.loadPkgfile(); err != nil {
		return nil, err
	}

	fspl.pkgs = nil

	err = filepath.Walk(fspl.Root, fspl.walkFunc())

	return fspl.pkgs, multierror.Append(fspl.multiErr, err).ErrorOrNil()
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

func (fspl *FilesystemPackageLoader) loadPkgfile() error {
	f, err := os.Open(filepath.Join(fspl.Root, constants.Pkgfile))
	if err != nil {
		if os.IsNotExist(err) {
			fspl.Logger.Printf("skipping %q: %s", constants.Pkgfile, err)
			return nil
		}

		return err
	}

	defer f.Close() //nolint: errcheck

	contents, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fspl.pkgFile, err = v1alpha2.NewPkgfile(contents)
	if err != nil {
		return fmt.Errorf("error parsing %q: %w", constants.Pkgfile, err)
	}

	fspl.Context.Merge(fspl.pkgFile.Vars)
	fspl.Logger.Printf("loaded %q", constants.Pkgfile)

	return nil
}
