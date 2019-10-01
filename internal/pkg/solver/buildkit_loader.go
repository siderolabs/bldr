package solver

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// BuildkitFrontendLoader loads packages from buildkit client.Reference
type BuildkitFrontendLoader struct {
	*log.Logger
	Context types.Variables
	Ref     client.Reference
	Ctx     context.Context
}

type packageProcess func(baseDir string, contents []byte)

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

			process(path, contents)
		}
	}

	return nil
}

// Load implements PackageLoader
func (bkfl *BuildkitFrontendLoader) Load() ([]*v1alpha2.Pkg, error) {
	if bkfl.Logger == nil {
		bkfl.Logger = log.New(log.Writer(), "[loader] ", log.Flags())
	}

	var pkgs []*v1alpha2.Pkg

	process := func(baseDir string, contents []byte) {
		pkg, err := v1alpha2.NewPkg(baseDir, contents, bkfl.Context)
		if err != nil {
			log.Printf("skipping %q: %s", baseDir, err)
		}
		log.Printf("loaded pkg %q from %q", pkg.Name, baseDir)
		pkgs = append(pkgs, pkg)
	}

	err := bkfl.walk("/", process)

	return pkgs, err
}
