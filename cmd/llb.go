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
	llbCmd.MarkFlagRequired("target")
	llbCmd.Flags().Var(&options.BuildPlatform, "build-platform", "Build platform")
	llbCmd.Flags().Var(&options.TargetPlatform, "target-platform", "Target platform")
	rootCmd.AddCommand(llbCmd)
}
