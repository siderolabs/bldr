// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package solver

import (
	"fmt"
	"strings"

	"github.com/emicklei/dot"

	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

// PackageDependency wraps v1alpha2.Depency with resolved internal dependencies.
type PackageDependency struct {
	// Pkg is set only for Internal dependencies.
	Node *PackageNode
	v1alpha2.Dependency
}

// ID returns unique string for dependency.
func (dep PackageDependency) ID() string {
	return fmt.Sprintf("%s-%s-%s", dep.Image, dep.Stage, dep.To)
}

// PackageNode is a Pkg with associated dependencies.
type PackageNode struct {
	Pkg          *v1alpha2.Pkg
	Name         string
	Dependencies []PackageDependency
}

// DumpDot dumps node and dependencies.
func (node *PackageNode) DumpDot(g *dot.Graph) dot.Node {
	n := g.Node(node.Name)

	for _, dep := range node.Dependencies {
		var depNode dot.Node

		if dep.IsInternal() {
			depNode = g.Node(dep.Stage)
		} else {
			imageRef := dep.Image
			// cut the digest
			imageRef, _, _ = strings.Cut(imageRef, "@")

			depNode = g.Node(imageRef)
			depNode.Box()
			depNode.Attr("fillcolor", "lemonchiffon")
			depNode.Attr("style", "filled")
		}

		edge := depNode.Edge(n)

		if dep.Runtime {
			edge.Attr("style", "bold")
			edge.Attr("color", "forestgreen")
		}
	}

	if node.Pkg.Variant == v1alpha2.Alpine {
		packageNode := g.Node("alpine")
		packageNode.Box()
		packageNode.Attr("fillcolor", "aquamarine")
		packageNode.Attr("style", "filled")

		packageNode.Edge(n)
	}

	for _, dep := range node.Pkg.Install {
		packageNode := g.Node("Alpine: " + dep)
		packageNode.Box()
		packageNode.Attr("fillcolor", "aquamarine")
		packageNode.Attr("style", "filled")

		packageNode.Edge(n)
	}

	return n
}

// RuntimeDependencies returns (recursively) all the runtime dependencies for the package.
func (node *PackageNode) RuntimeDependencies() (deps []PackageDependency) {
	for _, dep := range node.Dependencies {
		if !dep.Runtime {
			continue
		}

		deps = append(deps, dep)
		if dep.Node != nil {
			deps = append(deps, dep.Node.RuntimeDependencies()...)
		}
	}

	return
}

// PackageGraph capture root of the DAG.
type PackageGraph struct {
	Root *PackageNode
}

func (graph *PackageGraph) flatten(set PackageSet, node *PackageNode, skip map[*PackageNode]struct{}) PackageSet {
	if _, exists := skip[node]; exists {
		return set
	}

	set = append(set, node)
	skip[node] = struct{}{}

	for _, dep := range node.Dependencies {
		if dep.Node != nil {
			set = graph.flatten(set, dep.Node, skip)
		}
	}

	return set
}

// ToSet converts graph to set of nodes.
func (graph *PackageGraph) ToSet() PackageSet {
	return graph.flatten(nil, graph.Root, make(map[*PackageNode]struct{}))
}
