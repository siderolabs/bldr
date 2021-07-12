/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package convert

import (
	"context"
	"sort"

	"github.com/moby/buildkit/client/llb"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/environment"
	"github.com/talos-systems/bldr/internal/pkg/solver"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// GraphLLB wraps PackageGraph to provide LLB conversion.
//
// GraphLLB caches common images used in the build.
type GraphLLB struct {
	*solver.PackageGraph

	Options *environment.Options

	BaseImages   map[v1alpha2.Variant]llb.State
	Checksummer  llb.State
	LocalContext llb.State

	cache map[*solver.PackageNode]llb.State
}

// NewGraphLLB creates new GraphLLB and initializes shared images.
func NewGraphLLB(graph *solver.PackageGraph, options *environment.Options) *GraphLLB {
	result := &GraphLLB{
		PackageGraph: graph,
		Options:      options,
		cache:        make(map[*solver.PackageNode]llb.State),
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

	graph.BaseImages[v1alpha2.Alpine] = addEnv(addPkg(llb.Image(
		constants.DefaultBaseImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"base"),
	).Run(
		llb.Shlex("apk --no-cache --update add bash"),
		llb.WithCustomName(graph.Options.CommonPrefix+"base-apkinstall"),
	).Run(
		llb.Args([]string{"ln", "-svf", "/bin/bash", "/bin/sh"}),
		llb.WithCustomName(graph.Options.CommonPrefix+"base-symlink"),
	).Root()))

	graph.BaseImages[v1alpha2.Scratch] = addEnv(addPkg(llb.Scratch()))
}

func (graph *GraphLLB) buildChecksummer() {
	graph.Checksummer = llb.Image(
		constants.DefaultBaseImage,
		llb.WithCustomName(graph.Options.CommonPrefix+"cksum"),
	).Run(
		llb.Shlex("apk --no-cache --update add coreutils"),
		llb.WithCustomName(graph.Options.CommonPrefix+"cksum-apkinstall"),
	).Root()
}

func (graph *GraphLLB) buildLocalContext() {
	graph.LocalContext = llb.Local(
		"context",
		llb.ExcludePatterns(
			[]string{
				"**/.*",
				"**/" + constants.PkgYaml,
			},
		),
		llb.WithCustomName(graph.Options.CommonPrefix+"context"),
	)
}

// Build converts package graph to LLB.
func (graph *GraphLLB) Build() (llb.State, error) {
	return NewNodeLLB(graph.Root, graph).Build()
}

// Marshal returns marshaled LLB.
func (graph *GraphLLB) Marshal() (*llb.Definition, error) {
	out, err := graph.Build()
	if err != nil {
		return nil, err
	}

	out = out.SetMarshalDefaults(graph.Options.BuildPlatform.LLBPlatform)

	return out.Marshal(context.TODO())
}
