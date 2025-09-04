// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

var graphCmdFlags struct {
	buildArgs []string
}

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
		context := options.GetVariables().Copy()

		for _, buildArg := range graphCmdFlags.buildArgs {
			name, value, _ := strings.Cut(buildArg, "=")

			context["BUILD_ARG_"+name] = value
		}

		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: context,
		}

		packages, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}

		packages.FilterInPlace(func(pkg *v1alpha2.Pkg) bool {
			return pkg.Context["GRAPH_IGNORE"] != "true"
		})

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
	graphCmd.Flags().StringSliceVar(&graphCmdFlags.buildArgs, "build-arg", nil, "Build arguments to pass similar to docker buildx")
	rootCmd.AddCommand(graphCmd)
}
