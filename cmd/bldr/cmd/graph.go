// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/siderolabs/bldr/internal/pkg/solver"
)

// graphCmd represents the graph command.
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Graph dependencies between pkgs",
	Long: `This command outputs 'dot' formatted DAG of dependencies
starting from target to all the dependencies.

Typical usage:

  bldr graph | dot -Tpng > graph.png
`,
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

		var packageSet solver.PackageSet

		if options.Target != "" {
			graph, err := packages.Resolve(options.Target)
			if err != nil {
				log.Fatal(err)
			}

			packageSet = graph.ToSet()
		} else {
			packageSet = packages.ToSet()
		}

		packageSet.DumpDot(os.Stdout)
	},
}

func init() {
	graphCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to graph, if not set - graph all stages")
	rootCmd.AddCommand(graphCmd)
}
