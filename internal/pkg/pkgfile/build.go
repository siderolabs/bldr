/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package pkgfile

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	ctrplatforms "github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/frontend/dockerfile/dockerfile2llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"

	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/convert"
	"github.com/talos-systems/bldr/internal/pkg/environment"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

const (
	keyTarget         = "target"
	keyTargetPlatform = "platform"

	localNameDockerfile = "dockerfile"
	sharedKeyHint       = constants.PkgYaml
)

// Build is an entrypoint for buildkit frontend.
//
//nolint: gocyclo
func Build(ctx context.Context, c client.Client, options *environment.Options) (*client.Result, error) {
	opts := c.BuildOpts().Opts
	options.Target = opts[keyTarget]

	platforms := []environment.Platform{options.TargetPlatform}

	if opts[keyTargetPlatform] != "" {
		platforms = nil

		for _, p := range strings.Split(opts[keyTargetPlatform], ",") {
			var platform environment.Platform

			if err := platform.Set(p); err != nil {
				return nil, fmt.Errorf("unsupported platform %v: %w", p, err)
			}

			platforms = append(platforms, platform)
		}
	}

	exportMap := len(platforms) > 1

	expPlatforms := &exptypes.Platforms{
		Platforms: make([]exptypes.Platform, len(platforms)),
	}
	res := client.NewResult()

	var eg *errgroup.Group
	eg, ctx = errgroup.WithContext(ctx)

	for i, platform := range platforms {
		i := i
		platform := platform

		eg.Go(func() error {
			options := *options
			options.BuildPlatform = platform
			options.TargetPlatform = platform

			if len(platforms) > 1 {
				options.CommonPrefix = fmt.Sprintf("%s ", platform.ID)
			}

			pkgRef, err := fetchPkgs(ctx, c)
			if err != nil {
				return err
			}

			loader := solver.BuildkitFrontendLoader{
				Context: options.GetVariables(),
				Ref:     pkgRef,
				Ctx:     ctx,
			}

			packages, err := solver.NewPackages(&loader)
			if err != nil {
				return err
			}

			graph, err := packages.Resolve(options.Target)
			if err != nil {
				return err
			}

			def, err := convert.MarshalLLB(graph, &options)
			if err != nil {
				return err
			}

			r, err := c.Solve(ctx, client.SolveRequest{
				Definition: def.ToPB(),
			})
			if err != nil {
				return fmt.Errorf("failed to resolve dockerfile: %q", err)
			}

			ref, err := r.SingleRef()
			if err != nil {
				return err
			}

			img := dockerfile2llb.Image{
				Image: specs.Image{
					Architecture: platform.PlatformSpec.Architecture,
					OS:           platform.PlatformSpec.OS,
					RootFS: specs.RootFS{
						Type: "layers",
					},
				},
				Variant: platform.PlatformSpec.Variant,
			}

			config, err := json.Marshal(img)
			if err != nil {
				return fmt.Errorf("error marshaling image config: %w", err)
			}

			if !exportMap {
				res.AddMeta(exptypes.ExporterImageConfigKey, config)
				res.SetRef(ref)
			} else {
				k := ctrplatforms.Format(platform.PlatformSpec)
				res.AddMeta(fmt.Sprintf("%s/%s", exptypes.ExporterImageConfigKey, k), config)
				res.AddRef(k, ref)
				expPlatforms.Platforms[i] = exptypes.Platform{
					ID:       k,
					Platform: platform.PlatformSpec,
				}
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	if exportMap {
		dt, err := json.Marshal(expPlatforms)
		if err != nil {
			return nil, err
		}

		res.AddMeta(exptypes.ExporterPlatformsKey, dt)
	}

	return res, nil
}

func fetchPkgs(ctx context.Context, c client.Client) (client.Reference, error) {
	name := fmt.Sprintf("load %s and %ss", constants.Pkgfile, constants.PkgYaml)

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
