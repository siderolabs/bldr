// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package convert

import (
	"context"
	"sort"

	"github.com/moby/buildkit/client/llb"

	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

// GraphLLB wraps PackageGraph to provide LLB conversion.
//
// GraphLLB caches common images used in the build.
type GraphLLB struct {
	*solver.PackageGraph
	solverFn SolverFunc

	Options *environment.Options

	BaseImages   map[v1alpha2.Variant]llb.State
	Checksummer  llb.State
	LocalContext llb.State

	baseImageProcessor llbProcessor
	cache              map[*solver.PackageNode]llb.State

	commonRunOptions []llb.RunOption
}

type llbProcessor func(llb.State) llb.State

// NewGraphLLB creates new GraphLLB and initializes shared images.
func NewGraphLLB(graph *solver.PackageGraph, solverFn SolverFunc, options *environment.Options) *GraphLLB {
	result := &GraphLLB{
		PackageGraph: graph,
		Options:      options,
		solverFn:     solverFn,
		cache:        make(map[*solver.PackageNode]llb.State),
	}

	if options.ProxyEnv != nil {
		result.commonRunOptions = append(result.commonRunOptions, llb.WithProxy(*options.ProxyEnv))
	}

	result.buildBaseImages()
	result.buildChecksummer()
	result.buildLocalContext()

	return result
}

func (graph *GraphLLB) buildBaseImages() {
	graph.BaseImages = make(map[v1alpha2.Variant]llb.State)

	addPkg := func(root llb.State) llb.State {
		return root.File(
			llb.Mkdir(pkgDir, constants.DefaultDirMode),
			llb.WithCustomNamef("%smkdir %s", graph.Options.CommonPrefix, pkgDir),
		).Dir(pkgDir)
	}

	addEnv := func(root llb.State) llb.State {
		vars := graph.Options.GetVariables()
		keys := make([]string, 0, len(vars))

		for key := range vars {
			keys = append(keys, key)
		}

		sort.Strings(keys)

		for _, key := range keys {
			root = root.AddEnv(key, vars[key])
		}

		return root
	}

	graph.baseImageProcessor = func(root llb.State) llb.State {
		return addEnv(addPkg(root))
	}

	graph.BaseImages[v1alpha2.Alpine] = graph.baseImageProcessor(llb.Image(
		constants.DefaultBaseImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"base"),
	).Run(
		append(graph.commonRunOptions,
			llb.Shlex("apk --no-cache --update add bash"),
			llb.WithCustomName(graph.Options.CommonPrefix+"base-apkinstall"),
		)...,
	).Run(
		append(graph.commonRunOptions,
			llb.Args([]string{"ln", "-svf", "/bin/bash", "/bin/sh"}),
			llb.WithCustomName(graph.Options.CommonPrefix+"base-symlink"),
		)...,
	).Root())

	graph.BaseImages[v1alpha2.Scratch] = graph.baseImageProcessor(llb.Scratch())
}

func (graph *GraphLLB) buildChecksummer() {
	graph.Checksummer = llb.Image(
		constants.StageXBusyboxImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"cksum"),
	)
}

func (graph *GraphLLB) buildLocalContext() {
	graph.LocalContext = llb.Local(
		"context",
		llb.ExcludePatterns(
			[]string{
				"**/.*",
				"**/" + constants.PkgYaml,
				"**/" + constants.VarsYaml,
			},
		),
		llb.WithCustomName(graph.Options.CommonPrefix+"context"),
	)
}

// Build converts package graph to LLB.
func (graph *GraphLLB) Build(ctx context.Context) (llb.State, error) {
	return NewNodeLLB(graph.Root, graph).Build(ctx)
}

// Marshal returns marshaled LLB.
func (graph *GraphLLB) Marshal(ctx context.Context) (*llb.Definition, error) {
	out, err := graph.Build(ctx)
	if err != nil {
		return nil, err
	}

	out = out.SetMarshalDefaults(graph.Options.BuildPlatform.LLBPlatform)

	return out.Marshal(ctx)
}
