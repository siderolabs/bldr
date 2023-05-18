// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package pkgfile

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	ctrplatforms "github.com/containerd/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/exporter/containerimage/exptypes"
	"github.com/moby/buildkit/exporter/containerimage/image"
	"github.com/moby/buildkit/frontend/gateway/client"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"golang.org/x/sync/errgroup"

	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/convert"
	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/solver"
)

const (
	keyTarget         = "target"
	keyTargetPlatform = "platform"
	keyMultiPlatform  = "multi-platform"
	keyNoCache        = "no-cache"

	buildArgPrefix          = "build-arg:"
	buildArgSourceDateEpoch = buildArgPrefix + "SOURCE_DATE_EPOCH"
	buildArgCacheNS         = buildArgPrefix + "BUILDKIT_CACHE_MOUNT_NS"

	localNameDockerfile = "dockerfile"
	sharedKeyHint       = constants.PkgYaml
)

// Build is an entrypoint for buildkit frontend.
//
//nolint:gocyclo,cyclop,gocognit
func Build(ctx context.Context, c client.Client, options *environment.Options) (*client.Result, error) {
	opts := c.BuildOpts().Opts

	options.Target = opts[keyTarget]
	options.ProxyEnv = proxyEnvFromBuildArgs(filter(opts, buildArgPrefix))
	_, options.NoCache = opts[keyNoCache]

	if sourceDateEpoch, ok := opts[buildArgSourceDateEpoch]; ok {
		timestamp, err := strconv.ParseInt(sourceDateEpoch, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q: %w", buildArgSourceDateEpoch, err)
		}

		options.SourceDateEpoch = time.Unix(timestamp, 0)
	}

	if cacheNS, ok := opts[buildArgCacheNS]; ok {
		options.CacheIDNamespace = cacheNS
	}

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

	if v := opts[keyMultiPlatform]; v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			return nil, fmt.Errorf("invalid boolean value %s", v)
		}

		if !b && exportMap {
			return nil, fmt.Errorf("returning multiple target platforms is not allowed")
		}

		exportMap = b
	}

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

			if exportMap {
				options.CommonPrefix = fmt.Sprintf("%s ", platform.ID)
			}

			pkgRef, err := fetchPkgs(ctx, c)
			if err != nil {
				return err
			}

			buildContext := options.GetVariables().Copy()
			// push build arguments as `BUILD_ARGS_` prefixed variables
			buildContext.Merge(prefix(filter(opts, buildArgPrefix), "BUILD_ARG_"))

			loader := solver.BuildkitFrontendLoader{
				Context: buildContext,
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

			def, err := convert.MarshalLLB(graph, &options) //nolint:contextcheck
			if err != nil {
				return err
			}

			r, err := c.Solve(ctx, client.SolveRequest{
				Definition: def.ToPB(),
			})
			if err != nil {
				return fmt.Errorf("failed to solve LLB: %w", err)
			}

			ref, err := r.SingleRef()
			if err != nil {
				return err
			}

			img := image.Image{
				Image: specs.Image{
					Platform: specs.Platform{
						Architecture: platform.PlatformSpec.Architecture,
						OS:           platform.PlatformSpec.OS,
						Variant:      platform.PlatformSpec.Variant,
					},
					RootFS: specs.RootFS{
						Type: "layers",
					},
				},
				Config: image.ImageConfig{
					ImageConfig: specs.ImageConfig{
						Labels: packages.ImageLabels(),
					},
				},
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
	name := fmt.Sprintf("load %s, %ss and %ss", constants.Pkgfile, constants.PkgYaml, constants.VarsYaml)

	src := llb.Local(localNameDockerfile,
		llb.IncludePatterns([]string{
			constants.Pkgfile,
			"**/" + constants.PkgYaml,
			"**/" + constants.VarsYaml,
			"*/",
		}),
		llb.SessionID(c.BuildOpts().SessionID),
		llb.SharedKeyHint(sharedKeyHint),
		llb.WithCustomName(name),
	)

	def, err := src.Marshal(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal local source: %w", err)
	}

	res, err := c.Solve(ctx, client.SolveRequest{
		Definition: def.ToPB(),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to resolve pkgfile: %w", err)
	}

	return res.SingleRef()
}

func proxyEnvFromBuildArgs(args map[string]string) *llb.ProxyEnv {
	pe := &llb.ProxyEnv{}
	isNil := true

	for k, v := range args {
		if strings.EqualFold(k, "http_proxy") {
			pe.HTTPProxy = v
			isNil = false
		}

		if strings.EqualFold(k, "https_proxy") {
			pe.HTTPSProxy = v
			isNil = false
		}

		if strings.EqualFold(k, "ftp_proxy") {
			pe.FTPProxy = v
			isNil = false
		}

		if strings.EqualFold(k, "no_proxy") {
			pe.NoProxy = v
			isNil = false
		}
	}

	if isNil {
		return nil
	}

	return pe
}

func filter(opt map[string]string, key string) map[string]string {
	m := map[string]string{}

	for k, v := range opt {
		if strings.HasPrefix(k, key) {
			m[strings.TrimPrefix(k, key)] = v
		}
	}

	return m
}

func prefix(opt map[string]string, prefix string) map[string]string {
	m := map[string]string{}

	for k, v := range opt {
		m[prefix+k] = v
	}

	return m
}
