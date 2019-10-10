/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package convert

import (
	"github.com/moby/buildkit/client/llb"

	"github.com/talos-systems/bldr/internal/pkg/environment"
	"github.com/talos-systems/bldr/internal/pkg/solver"
)

// BuildLLB translates package graph into LLB DAG
func BuildLLB(graph *solver.PackageGraph, options *environment.Options) (llb.State, error) {
	return NewGraphLLB(graph, options).Build()
}
