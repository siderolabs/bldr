// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package convert

import (
	"context"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"

	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/solver"
)

// SolverFunc can be called to solve the package into the llb state via buildkit.
type SolverFunc func(ctx context.Context, platform environment.Platform, target string) (*client.Result, error)

// MarshalLLB translates package graph into LLB DAG and marshals it.
func MarshalLLB(ctx context.Context, graph *solver.PackageGraph, solver SolverFunc, options *environment.Options) (*llb.Definition, error) {
	return NewGraphLLB(graph, solver, options).Marshal(ctx)
}
