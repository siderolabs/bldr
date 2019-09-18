package cmd

import (
	"log"

	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha1"
)

var (
	options = &v1alpha1.Options{}
)

// packCmd represents the pack command
var packCmd = &cobra.Command{
	Use:   "pack",
	Short: "Pack a build as a container",
	Long: `This command provides a wrapper around docker buildx. It is opinionated
about a number of things.
`,
	Run: func(cmd *cobra.Command, args []string) {
		pkg, err := v1alpha1.NewPkg(pkgFile, options)
		if err != nil {
			log.Fatal(err)
		}

		if err := pkg.Pack(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	packCmd.Flags().StringVarP(&options.Registry, "registry", "r", constants.DefaultRegistry, "Docker registry to tag the image with")
	packCmd.Flags().StringVarP(&options.Organization, "organization", "o", constants.DefaultOrganization, "Docker organization to tag the image with")
	packCmd.Flags().StringVarP(&options.Platform, "platform", "", "linux/amd64", "Passed through to docker build command")
	packCmd.Flags().StringVarP(&options.Progress, "progress", "", "auto", "Passed through to docker build command")
	packCmd.Flags().StringVarP(&options.Push, "push", "", "false", "Passed through to docker build command")
	packCmd.Flags().StringVarP(&options.CacheTo, "cache-to", "", "", "Passed through to docker build command")
	packCmd.Flags().StringVarP(&options.CacheFrom, "cache-from", "", "", "Passed through to docker build command")
	rootCmd.AddCommand(packCmd)
}
