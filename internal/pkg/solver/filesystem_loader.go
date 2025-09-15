// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package solver

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/go-multierror"

	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/types"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

// FilesystemPackageLoader loads packages by walking file system tree.
type FilesystemPackageLoader struct {
	*log.Logger
	Context types.Variables

	HookOnLoad      func(path string, contents []byte)
	HookOnVariables func(path string, vars types.Variables)

	pathContexts      map[string]types.Variables
	multiErr          *multierror.Error
	pkgFile           *v1alpha2.Pkgfile
	Root              string
	absRootPath       string
	pkgFilePaths      []string
	varFilePaths      []string
	templateFilePaths []string
	pkgs              []*v1alpha2.Pkg
}

func (fspl *FilesystemPackageLoader) walkFunc() filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fspl.Printf("error walking %q: %s", path, err)

			return nil
		}

		if info.Name() != "." && strings.HasPrefix(info.Name(), ".") && info.IsDir() {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		switch {
		case info.Name() == constants.PkgYaml:
			fspl.pkgFilePaths = append(fspl.pkgFilePaths, path)
		case info.Name() == constants.VarsYaml:
			fspl.varFilePaths = append(fspl.varFilePaths, path)
		case strings.HasSuffix(info.Name(), constants.TemplateExt):
			fspl.templateFilePaths = append(fspl.templateFilePaths, path)
		}

		return nil
	}
}

// Load implements PackageLoader.
func (fspl *FilesystemPackageLoader) Load() (*LoadResult, error) {
	if fspl.Logger == nil {
		fspl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

	if fspl.Root == "" {
		fspl.Root = "."
	}

	fspl.pathContexts = make(map[string]types.Variables)

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
	if err == nil {
		sort.Slice(fspl.varFilePaths, func(i, j int) bool {
			return filepath.Dir(fspl.varFilePaths[i]) < filepath.Dir(fspl.varFilePaths[j])
		})

		for _, path := range fspl.varFilePaths {
			if err = fspl.loadVariables(path); err != nil {
				fspl.Printf("error loading variables %q: %s", path, err)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error loading variables %q: %w", path, err))
			}

			fspl.Printf("loaded variables from %q", path)
		}

		for _, path := range fspl.pkgFilePaths {
			var pkg *v1alpha2.Pkg

			pkg, err = fspl.loadPkg(path)
			if err != nil {
				fspl.Printf("error loading %q: %s", path, err)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error loading %q: %w", path, err))

				continue
			}

			fspl.Printf("loaded pkg %q from %q", pkg.Name, path)
			fspl.pkgs = append(fspl.pkgs, pkg)
		}

		for _, path := range fspl.templateFilePaths {
			var pkg *v1alpha2.Pkg

			if pkg, err = fspl.attachTemplate(path); err != nil {
				fspl.Printf("error attaching template %q: %s", path, err)
				fspl.multiErr = multierror.Append(fspl.multiErr, fmt.Errorf("error attaching template %q: %w", path, err))
			} else {
				fspl.Printf("attached template %q to %q", path, pkg.Name)
			}
		}
	}

	return &LoadResult{
		Pkgfile: fspl.pkgFile,
		Pkgs:    fspl.pkgs,
	}, multierror.Append(fspl.multiErr, err).ErrorOrNil()
}

func (fspl *FilesystemPackageLoader) resolveContext(basePath string) types.Variables {
	context := fspl.Context.Copy()

	dirs := strings.Split(basePath, string(filepath.Separator))

	for i := 0; i <= len(dirs); i++ {
		var subPath string

		if i == 0 {
			subPath = "."
		} else {
			subPath = strings.Join(dirs[:i], string(filepath.Separator))
		}

		if subcontext, ok := fspl.pathContexts[subPath]; ok {
			context.Merge(subcontext)
		}
	}

	return context
}

func (fspl *FilesystemPackageLoader) loadVariables(path string) error {
	absFile, err := filepath.Abs(path)
	if err != nil {
		return err
	}

	basePath, err := filepath.Rel(fspl.absRootPath, absFile)
	if err != nil {
		return err
	}

	baseContext := fspl.resolveContext(filepath.Dir(basePath))

	if fspl.HookOnLoad != nil {
		var contents []byte

		contents, err = os.ReadFile(path)
		if err != nil {
			return err
		}

		fspl.HookOnLoad(basePath, contents)
	}

	var vars types.Variables

	if err = vars.Load(path, baseContext); err != nil {
		return err
	}

	if fspl.HookOnVariables != nil {
		fspl.HookOnVariables(basePath, vars)
	}

	fspl.pathContexts[filepath.Dir(basePath)] = vars

	return nil
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

	context := fspl.resolveContext(filepath.Dir(basePath))

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	defer f.Close() //nolint:errcheck

	contents, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return v1alpha2.NewPkg(filepath.Dir(basePath), path, contents, context)
}

func (fspl *FilesystemPackageLoader) attachTemplate(path string) (*v1alpha2.Pkg, error) {
	// find the closest pkgs in relative path
	var (
		closestPkg  *v1alpha2.Pkg
		shortestRel string
	)

	absFile, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	relPath, err := filepath.Rel(fspl.absRootPath, absFile)
	if err != nil {
		return nil, err
	}

	for _, pkg := range fspl.pkgs {
		var rel string

		rel, err = filepath.Rel(pkg.BaseDir, filepath.Dir(relPath))
		if err != nil || strings.HasPrefix(rel, "..") {
			continue
		}

		if shortestRel == "" || len(rel) < len(shortestRel) {
			closestPkg = pkg
			shortestRel = rel
		}
	}

	if closestPkg == nil {
		return nil, fmt.Errorf("no suitable package found for template %q", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	// attach the template to the closest package
	return closestPkg, closestPkg.AttachTemplatedFile(filepath.Join(shortestRel, filepath.Base(path)), content)
}

func (fspl *FilesystemPackageLoader) loadPkgfile() error {
	f, err := os.Open(filepath.Join(fspl.Root, constants.Pkgfile))
	if err != nil {
		if os.IsNotExist(err) {
			fspl.Printf("skipping %q: %s", constants.Pkgfile, err)

			return nil
		}

		return err
	}

	defer f.Close() //nolint:errcheck

	contents, err := io.ReadAll(f)
	if err != nil {
		return err
	}

	if fspl.HookOnLoad != nil {
		fspl.HookOnLoad(constants.Pkgfile, contents)
	}

	fspl.pkgFile, err = v1alpha2.NewPkgfile(contents)
	if err != nil {
		return fmt.Errorf("error parsing %q: %w", constants.Pkgfile, err)
	}

	if fspl.HookOnVariables != nil {
		fspl.HookOnVariables(constants.Pkgfile, fspl.pkgFile.Vars)
	}

	fspl.Context.Merge(fspl.pkgFile.Vars)
	fspl.Printf("loaded %q", constants.Pkgfile)

	return nil
}
