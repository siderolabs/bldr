// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"runtime"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/spf13/cobra"

	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

func validateChecksums(ctx context.Context, set solver.PackageSet, l *log.Logger) error {
	var (
		wg          sync.WaitGroup
		concurrency = runtime.GOMAXPROCS(-1)
		pkgs        = make(chan *v1alpha2.Pkg)
		errors      = make(chan error)
	)

	// start downloaders
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for pkg := range pkgs {
				for _, step := range pkg.Steps {
					for _, src := range step.Sources {
						l.Printf("downloading %s ...", src.URL)

						_, _, err := src.ValidateChecksums(ctx)
						if err != nil {
							errors <- fmt.Errorf("%s: %w", pkg.Name, err)
						}
					}
				}
			}
		}()
	}

	var (
		multiErr *multierror.Error
		done     = make(chan struct{})
	)

	// start results reader
	go func() {
		for err := range errors {
			multiErr = multierror.Append(multiErr, err)
		}

		close(done)
	}()

	// send work to downloaders
	for _, node := range set {
		pkgs <- node.Pkg
	}

	close(pkgs)
	wg.Wait()
	close(errors)
	<-done

	return multiErr.ErrorOrNil()
}

var validateCmdFlags struct {
	checksums bool
}

// validateCmd represents the validate command.
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate syntax of pkg.yaml files",
	Long: `This command scans directory tree for pkg.yaml files,
loads them and validates for errors. `,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		packages, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}

		if validateCmdFlags.checksums {
			l := log.New(log.Writer(), "[validate] ", log.Flags())
			if !debug {
				l.SetOutput(io.Discard)
			}

			if err = validateChecksums(context.TODO(), packages.ToSet(), l); err != nil {
				log.Fatal(err)
			}
		}
	},
}

func init() {
	validateCmd.Flags().BoolVar(&validateCmdFlags.checksums, "checksums", true, "validate checksums")
	rootCmd.AddCommand(validateCmd)
}
