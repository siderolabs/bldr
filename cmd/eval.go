/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"
	"os"
	"strings"
	"text/template"

	"github.com/spf13/cobra"

	"github.com/talos-systems/bldr/internal/pkg/solver"
)

var evalCmdFlags struct {
	buildArgs []string
}

// evalCmd represents the eval command.
var evalCmd = &cobra.Command{
	Use:   "eval",
	Short: "Evaluate a Go template using the variables defined in the vars.yaml and Pkgfile.",
	Long: `This command prints the result of evaluating a Go template give as the argument.
 Variables are looked up for the target specified as the '--target' flag.'.
 `,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		context := options.GetVariables().Copy()

		for _, buildArg := range evalCmdFlags.buildArgs {
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

		graph, err := packages.Resolve(options.Target)
		if err != nil {
			log.Fatal(err)
		}

		tmpl, err := template.New("arg").Parse(args[0])
		if err != nil {
			log.Fatal(err)
		}

		if err = tmpl.Execute(os.Stdout, graph.Root.Pkg.Context); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	evalCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to build")
	evalCmd.Flags().StringSliceVar(&evalCmdFlags.buildArgs, "build-arg", nil, "Build arguments to pass similar to docker buildx")
	evalCmd.MarkFlagRequired("target") //nolint:errcheck
	evalCmd.Flags().Var(&options.BuildPlatform, "build-platform", "Build platform")
	evalCmd.Flags().Var(&options.TargetPlatform, "target-platform", "Target platform")
	rootCmd.AddCommand(evalCmd)
}
