/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/environment"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

// graphCmd represents the graph command
var graphCmd = &cobra.Command{
	Use:   "graph",
	Short: "Graph dependencies between pkgs",
	Long: `This command outputs 'dot' formatted DAG of dependencies
starting from target to all the dependencies.

Typical usage:

  bldr graph | dot -Tpng > graph.png
`,
	Run: func(cmd *cobra.Command, args []string) {
		options, err := environment.NewOptions()
		if err != nil {
			log.Fatal(err)
		}

		if target != "" {
			options.Target = target
		}

		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.ToolchainPlatform.GetVariables(),
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
	graphCmd.Flags().StringVarP(&target, "target", "t", "", "Target image to graph, if not set - graph all stages")
	rootCmd.AddCommand(graphCmd)
}
