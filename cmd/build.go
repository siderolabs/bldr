// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha1"
)

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build a package",
	Long: `This command assumes that the current directory contains a pkg.yaml file.
A pkg.yaml can contain any number of steps that will be executed in order.
A set of evnironment variables are available in each step:

- ARCH: The architecture of the current machine.
- BUILD: The target triple for the build machine.
- HOST: The target triple for the intended machine where the resulting binaries will be executed.
- TARGET: The target triple that the compiler will produce code for.
- VENDOR: An arbitrary string used to identify the vendor of the toolchain.
- SYSROOT:
- CFLAGS: A preset set of flags to optimize builds.
- CXXFLAGS: A preset set of flags to optimize builds.
- LDFLAGS: A preset set of flags to optimize builds.

The general format of a target triple is:
	<arch><sub>-<vendor>-<sys>-<abi>
See https://gcc.gnu.org/onlinedocs/gccint/Configure-Terms.html for a detailed description of target triples.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg, err := v1alpha1.NewPkg(pkgFile, nil)
		if err != nil {
			log.Fatal(err)
		}

		if err := pkg.Build(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)
}
