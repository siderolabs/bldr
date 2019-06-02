package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha1"
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack a build as a container",
	Long: `This command provides a wrapper around docker buildx. It is opinionated
about a number of things.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg, err := v1alpha1.NewPkg(pkgFile)
		if err != nil {
			log.Fatal(err)
		}

		if err := pkg.Pack(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(packCmd)
}
