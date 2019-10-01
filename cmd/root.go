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
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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
