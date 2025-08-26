// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert/yaml"

	"github.com/siderolabs/bldr/internal/pkg/solver"
)

var dumpCmdFlags struct {
	templatePath string
	buildArgs    []string
}

// dumpCmd represents the graph command.
var dumpCmd = &cobra.Command{
	Use:   "dump",
	Short: "Dump information about all packages.",
	Long: `This command outputs the set of all packages as JSON
output, or runs a Go template against a set of all packages
and prints the template output.
`,
	Args: cobra.NoArgs,
	Run: func(_ *cobra.Command, _ []string) {
		context := options.GetVariables().Copy()

		for _, buildArg := range dumpCmdFlags.buildArgs {
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

		var packageSet solver.PackageSet

		if options.Target != "" {
			var graph *solver.PackageGraph

			graph, err = packages.Resolve(options.Target)
			if err != nil {
				log.Fatal(err)
			}

			packageSet = graph.ToSet()
		} else {
			packageSet = packages.ToSet()
		}

		packageSet = packageSet.Sorted()

		if dumpCmdFlags.templatePath == "" {
			if err = packageSet.DumpJSON(os.Stdout); err != nil {
				log.Fatal(err)
			}

			return
		}

		funcs := sprig.HermeticTxtFuncMap()
		funcs["mustFromYAML"] = mustFromYAML

		tmpl, err := template.New(filepath.Base(dumpCmdFlags.templatePath)).
			Funcs(funcs).
			ParseFiles(dumpCmdFlags.templatePath)
		if err != nil {
			log.Fatal(err)
		}

		if err := packageSet.Template(os.Stdout, tmpl); err != nil {
			log.Fatal(err)
		}
	},
}

func mustFromYAML(v string) (any, error) {
	var output any

	err := yaml.Unmarshal([]byte(v), &output)

	return output, err
}

func init() {
	dumpCmd.Flags().StringVarP(&options.Target, "target", "t", "", "Target image to dump, if not set - dump all stages")
	dumpCmd.Flags().StringSliceVar(&dumpCmdFlags.buildArgs, "build-arg", nil, "Build arguments to pass similar to docker buildx")
	dumpCmd.Flags().StringVarP(&dumpCmdFlags.templatePath, "template", "", "", "Path to Go template file to use for output formatting")
	rootCmd.AddCommand(dumpCmd)
}
