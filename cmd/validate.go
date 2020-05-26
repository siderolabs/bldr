/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

// validateCmd represents the validate command.
var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate syntax of pkg.yaml files",
	Long: `This command scans directory tree for pkg.yaml files,
loads them and validates for errors. `,
	Run: func(cmd *cobra.Command, args []string) {
		loader := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		_, err := solver.NewPackages(&loader)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(validateCmd)
}
