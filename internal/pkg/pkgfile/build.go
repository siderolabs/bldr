/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package pkgfile

import (
	"context"
	"fmt"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/convert"
	"github.com/talos-systems/bldr/internal/pkg/environment"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

const (
	keyTarget         = "target"
	keyBuildPlatform  = "build-platform"
	keyTargetPlatform = "target-platform"

	localNameDockerfile = "dockerfile"
	sharedKeyHint       = constants.PkgYaml
)

// Build is an entrypoint for buildkit frontend
func Build(ctx context.Context, c client.Client, options *environment.Options) (*client.Result, error) {
	opts := c.BuildOpts().Opts
	options.Target = opts[keyTarget]

	if opts[keyBuildPlatform] != "" {
		options.BuildPlatform.Set(opts[keyBuildPlatform]) //nolint: errcheck
	}

	if opts[keyTargetPlatform] != "" {
		options.TargetPlatform.Set(opts[keyTargetPlatform]) //nolint: errcheck
	}

	pkgRef, err := fetchPkgs(ctx, c)
	if err != nil {
		return nil, err
	}

	loader := solver.BuildkitFrontendLoader{
		Context: options.GetVariables(),
		Ref:     pkgRef,
		Ctx:     ctx,
	}

	packages, err := solver.NewPackages(&loader)
	if err != nil {
		return nil, err
	}

	graph, err := packages.Resolve(options.Target)
	if err != nil {
		return nil, err
	}

	out, err := convert.BuildLLB(graph, options)
	if err != nil {
		return nil, err
	}

	def, err := out.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal local source: %q", err)
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve dockerfile: %q", err)
	}

	return res, nil
}

func fetchPkgs(ctx context.Context, c client.Client) (client.Reference, error) {
	name := fmt.Sprintf("load %s and %ss", constants.Pkgfile, constants.Pkgfile)

	src := llb.Local(localNameDockerfile,
		llb.IncludePatterns([]string{
			constants.Pkgfile,
			"**/" + constants.PkgYaml,
			"*/",
		}),
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint(sharedKeyHint),
		llb.WithCustomName(name),
	)

	def, err := src.Marshal()
	if err != nil {
		return nil, fmt.Errorf("failed to marshal local source: %q", err)
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve pkgfile: %q", err)
	}

	return res.SingleRef()
}
