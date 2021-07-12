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
	Ctx     context.Context

	pkgFile *v1alpha2.Pkgfile
}

type packageProcess func(baseDir string, contents []byte) error

func (bkfl *BuildkitFrontendLoader) walk(path string, process packageProcess) error {
	entries, err := bkfl.Ref.ReadDir(bkfl.Ctx, client.ReadDirRequest{
		Path: path,
	})

	if err != nil {
		return fmt.Errorf("error readdir %q: %w", path, err)
	}

	for _, entry := range entries {
		if os.FileMode(entry.GetMode())&os.ModeDir > 0 {
			if err = bkfl.walk(filepath.Join(path, entry.GetPath()), process); err != nil {
				return err
			}
		} else if entry.GetPath() == constants.PkgYaml {
			contents, err := bkfl.Ref.ReadFile(bkfl.Ctx, client.ReadRequest{
				Filename: filepath.Join(path, entry.GetPath()),
			})
			if err != nil {
				return fmt.Errorf("error reading %q under %q: %w", entry.GetPath(), path, err)
			}

			err = process(path, contents)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Load implements PackageLoader.
func (bkfl *BuildkitFrontendLoader) Load() (*LoadResult, error) {
	if bkfl.Logger == nil {
		bkfl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

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

	process := func(baseDir string, contents []byte) error {
		pkg, err2 := v1alpha2.NewPkg(baseDir, "", contents, bkfl.Context)
		if err2 != nil {
			log.Printf("error loading %q: %s", baseDir, err2)
			multiErr = multierror.Append(multiErr, fmt.Errorf("error loading %q: %w", baseDir, err2))

			return nil
		}

		log.Printf("loaded pkg %q from %q", pkg.Name, baseDir)
		pkgs = append(pkgs, pkg)

		return nil
	}

	err = bkfl.walk("/", process)

	return &LoadResult{
		Pkgfile: bkfl.pkgFile,
		Pkgs:    pkgs,
	}, multierror.Append(multiErr, err).ErrorOrNil()
}
