/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// Package cmd contains definitions of CLI commands.
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/environment"
)

var (
	pkgRoot string
	options = &environment.Options{
		BuildPlatform:  environment.LinuxAmd64,
		TargetPlatform: environment.LinuxAmd64,
	}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "bldr",
	Short: "A tool to build and manage software via Pkgfile and pkg.yaml",
	Long: `bldr usually works in buildkit frontend mode when it's not directly
exposed as a CLI tool. In that mode of operation bldr loads root Pkgfile and
a set of pkg.yamls, processes them, builds dependency graph and outputs it
as LLB graph to buildkit backend.

bldr can be also used to produce graph of dependencies between build steps and
output LLB directly which is useful for development or debugging.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	rootCmd.PersistentFlags().StringVarP(&pkgRoot, "root", "", ".", "The path to a pkg root")
}
