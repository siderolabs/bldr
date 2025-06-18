// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/cobra"

	"github.com/siderolabs/bldr/internal/pkg/sbom"
	"github.com/siderolabs/bldr/internal/pkg/solver"
)

// sbomCmd represents the graph command.
var sbomCmd = &cobra.Command{
	Use:   "sbom",
	Short: "Generate an SBOM for a package",
	Long: `This command outputs evaluates the package build instructions
and outputs a Software Bill of Materials (SBOM) for it in SPDX format.
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

		graph, err := packages.Resolve(options.Target)
		if err != nil {
			log.Fatal(err)
		}

		pkg := graph.Root.Pkg

		sbomDoc, err := sbom.CreatePackageSBOM(pkg, options.TargetPlatform.Arch)
		if err != nil {
			log.Fatalf("failed to create SBOM for package %q: %v", pkg.Name, err)
		}

		s, err := sbom.ToSpdxJSON(*sbomDoc, time.Unix(1, 0))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(s)
	},
}

func init() {
	sbomCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to describe")
	sbomCmd.MarkFlagRequired("target") //nolint:errcheck
	sbomCmd.Flags().Var(&options.TargetPlatform, "target-platform", "Target platform")
	rootCmd.AddCommand(sbomCmd)
}
