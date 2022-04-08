/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"runtime"
	"sort"
	"sync"
	"text/tabwriter"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/solver"
	"github.com/talos-systems/bldr/internal/pkg/update"
)

type pkgInfo struct {
	file   string
	source string
}

type updateInfo struct {
	file string
	*update.LatestInfo
}

//nolint:gocyclo,cyclop
func checkUpdates(ctx context.Context, set solver.PackageSet, l *log.Logger) error {
	var (
		wg          sync.WaitGroup
		concurrency = runtime.GOMAXPROCS(-1)
		sources     = make(chan *pkgInfo)
		updates     = make(chan *updateInfo)
	)

	// start updaters
	for i := 0; i < concurrency; i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for src := range sources {
				res, e := update.Latest(ctx, src.source)
				if e != nil {
					l.Print(e)

					continue
				}

				updates <- &updateInfo{
					file:       src.file,
					LatestInfo: res,
				}
			}
		}()
	}

	var (
		res  []updateInfo
		done = make(chan struct{})
	)

	// start results reader
	go func() {
		for update := range updates {
			res = append(res, *update)
		}

		close(done)
	}()

	// send work to updaters
	for _, node := range set {
		for _, step := range node.Pkg.Steps {
			for _, src := range step.Sources {
				sources <- &pkgInfo{
					file:   node.Pkg.FileName,
					source: src.URL,
				}
			}
		}
	}

	close(sources)
	wg.Wait()
	close(updates)
	<-done

	sort.Slice(res, func(i, j int) bool { return res[i].file < res[j].file })

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "%s\t%s\t%s\n", "File", "Update", "URL")

	for _, info := range res {
		if updateCmdFlag.all || info.HasUpdate {
			url := info.LatestURL
			if url == "" {
				url = info.BaseURL
			}

			fmt.Fprintf(w, "%s\t%t\t%s\n", info.file, info.HasUpdate, url)
		}
	}

	return w.Flush()
}

var updateCmdFlag struct {
	all bool
	dry bool
}

// updateCmd represents the `update` command.
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update pkgs",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if !updateCmdFlag.dry {
			log.Fatal("Real update is not implemented yet; pass `--dry` flag.")
		}

		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		packages, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}

		l := log.New(log.Writer(), "[update] ", log.Flags())
		if !debug {
			l.SetOutput(ioutil.Discard)
		}

		if err = checkUpdates(context.TODO(), packages.ToSet(), l); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	updateCmd.Flags().BoolVarP(&updateCmdFlag.all, "all", "a", false, "List all packages, not only updated")
	updateCmd.Flags().BoolVar(&updateCmdFlag.dry, "dry", false, "Dry run: check for updates, but not actually update pkgs")
	rootCmd.AddCommand(updateCmd)
}
