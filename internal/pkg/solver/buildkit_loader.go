/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/hashicorp/go-multierror"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// BuildkitFrontendLoader loads packages from buildkit client.Reference.
type BuildkitFrontendLoader struct {
	*log.Logger
	Context types.Variables
	Ref     client.Reference
	//nolint:containedctx
	Ctx context.Context

	pathContexts map[string]types.Variables
	pkgFile      *v1alpha2.Pkgfile
}

type processor func(baseDir string, contents []byte) error

//nolint:gocyclo,cyclop
func (bkfl *BuildkitFrontendLoader) walk(path string, processVars, processPkgs processor) error {
	entries, err := bkfl.Ref.ReadDir(bkfl.Ctx, client.ReadDirRequest{
		Path: path,
	})
	if err != nil {
		return fmt.Errorf("error readdir %q: %w", path, err)
	}

	// 1. find and load variables
	for _, entry := range entries {
		if entry.GetPath() == constants.VarsYaml {
			var contents []byte

			contents, err = bkfl.Ref.ReadFile(bkfl.Ctx, client.ReadRequest{
				Filename: filepath.Join(path, entry.GetPath()),
			})
			if err != nil {
				return fmt.Errorf("error reading %q under %q: %w", entry.GetPath(), path, err)
			}

			err = processVars(path, contents)
			if err != nil {
				return err
			}
		}
	}

	// 2. find and load package files
	for _, entry := range entries {
		if entry.GetPath() == constants.PkgYaml {
			var contents []byte

			contents, err = bkfl.Ref.ReadFile(bkfl.Ctx, client.ReadRequest{
				Filename: filepath.Join(path, entry.GetPath()),
			})
			if err != nil {
				return fmt.Errorf("error reading %q under %q: %w", entry.GetPath(), path, err)
			}

			err = processPkgs(path, contents)
			if err != nil {
				return err
			}
		}
	}

	// 3. descend into subdirectories
	for _, entry := range entries {
		if os.FileMode(entry.GetMode())&os.ModeDir > 0 {
			if err = bkfl.walk(filepath.Join(path, entry.GetPath()), processVars, processPkgs); err != nil {
				return err
			}
		}
	}

	return nil
}

func (bkfl *BuildkitFrontendLoader) resolveContext(basePath string) types.Variables {
	context := bkfl.Context.Copy()

	dirs := strings.Split(basePath, string(filepath.Separator))

	for i := 0; i <= len(dirs); i++ {
		var subPath string

		if i == 0 {
			subPath = "/"
		} else {
			subPath = strings.Join(dirs[:i], string(filepath.Separator))
		}

		if subcontext, ok := bkfl.pathContexts[subPath]; ok {
			context.Merge(subcontext)
		}
	}

	return context
}

func (bkfl *BuildkitFrontendLoader) loadVariables(baseDir string, contents []byte) error {
	baseContext := bkfl.resolveContext(baseDir)

	var vars types.Variables

	if err := vars.LoadContents(contents, baseContext); err != nil {
		return fmt.Errorf("error loading variables at %q: %w", baseDir, err)
	}

	log.Printf("loaded variables from %q", baseDir)

	bkfl.pathContexts[baseDir] = vars

	return nil
}

// Load implements PackageLoader.
func (bkfl *BuildkitFrontendLoader) Load() (*LoadResult, error) {
	if bkfl.Logger == nil {
		bkfl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

	bkfl.pathContexts = make(map[string]types.Variables)

	contents, err := bkfl.Ref.ReadFile(bkfl.Ctx, client.ReadRequest{
		Filename: constants.Pkgfile,
	})
	if err != nil {
		return nil, fmt.Errorf("error loading %q: %w", constants.Pkgfile, err)
	}

	bkfl.pkgFile, err = v1alpha2.NewPkgfile(contents)
	if err != nil {
		return nil, fmt.Errorf("error parsing %q: %w", constants.Pkgfile, err)
	}

	log.Printf("loaded %q", constants.Pkgfile)

	bkfl.Context.Merge(bkfl.pkgFile.Vars)

	var (
		pkgs     []*v1alpha2.Pkg
		multiErr *multierror.Error
	)

	processPackage := func(baseDir string, contents []byte) error {
		pkg, err2 := v1alpha2.NewPkg(baseDir, "", contents, bkfl.resolveContext(baseDir))
		if err2 != nil {
			log.Printf("error loading %q: %s", baseDir, err2)
			multiErr = multierror.Append(multiErr, fmt.Errorf("error loading %q: %w", baseDir, err2))

			return nil
		}

		log.Printf("loaded pkg %q from %q", pkg.Name, baseDir)
		pkgs = append(pkgs, pkg)

		return nil
	}

	err = bkfl.walk("/", bkfl.loadVariables, processPackage)

	return &LoadResult{
		Pkgfile: bkfl.pkgFile,
		Pkgs:    pkgs,
	}, multierror.Append(multiErr, err).ErrorOrNil()
}
