/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"
	"os"

	"github.com/moby/buildkit/client/llb"
	"github.com/spf13/cobra"

	"github.com/talos-systems/bldr/internal/pkg/convert"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

// llbCmd represents the llb command
var llbCmd = &cobra.Command{
	Use:   "llb",
	Short: "Dump buildkit LLB for the build",
	Long: `This command parses build instructions from pkg.yaml files,
and outputs buildkit LLB to stdout. This can be used as 'bldr pack ... | buildctl ...'.
`,
	Run: func(cmd *cobra.Command, args []string) {
		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		packages, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}

		graph, err := packages.Resolve(options.Target)
		if err != nil {
			log.Fatal(err)
		}

		out, err := convert.BuildLLB(graph, options)
		if err != nil {
			log.Fatal(err)
		}

		dt, err := out.Marshal(options.BuildPlatform.LLBPlatform)
		if err != nil {
			log.Fatal(err)
		}

		err = llb.WriteTo(dt, os.Stdout)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	llbCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to build")
	llbCmd.MarkFlagRequired("target") //nolint: errcheck
	llbCmd.Flags().Var(&options.BuildPlatform, "build-platform", "Build platform")
	llbCmd.Flags().Var(&options.TargetPlatform, "target-platform", "Target platform")
	rootCmd.AddCommand(llbCmd)
}
