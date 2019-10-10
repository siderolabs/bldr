/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package cmd

import (
	"context"
	"log"

	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/frontend/gateway/grpcclient"
	"github.com/moby/buildkit/util/appcontext"
	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/pkgfile"
)

// frontendCmd represents the frontend command
var frontendCmd = &cobra.Command{
	Use:   "frontend",
	Short: "Buildkit frontend for Pkgfile",
	Long: `This command implements buildkit frontend.
`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := grpcclient.RunFromEnvironment(
			appcontext.Context(),
			func(ctx context.Context, c client.Client) (*client.Result, error) {
				return pkgfile.Build(ctx, c, options)
			},
		); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(frontendCmd)
}
