/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"github.com/emicklei/dot"

	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// PackageNode is a Pkg with associated dependencies
type PackageNode struct {
	Pkg       *v1alpha2.Pkg
	Name      string
	DependsOn []*PackageNode
}

// DumpDot dumps node and dependencies
func (node *PackageNode) DumpDot(g *dot.Graph) dot.Node {
	n := g.Node(node.Name)

	for _, dep := range node.Pkg.InternalDependencies() {
		depNode := g.Node(dep)
		depNode.Edge(n)
	}

	for _, dep := range node.Pkg.ExternalDependencies() {
		imageNode := g.Node(dep)
		imageNode.Box()
		imageNode.Attr("fillcolor", "lemonchiffon")
		imageNode.Attr("style", "filled")

		imageNode.Edge(n)
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

// PackageGraph capture root of the DAG
type PackageGraph struct {
	Root *PackageNode
}

func (graph *PackageGraph) flatten(set PackageSet, node *PackageNode, skip map[*PackageNode]struct{}) PackageSet {
	if _, exists := skip[node]; exists {
		return set
	}

	set = append(set, node)
	skip[node] = struct{}{}

	for _, dep := range node.DependsOn {
		set = graph.flatten(set, dep, skip)
	}

	return set
}

// ToSet converts graph to set of nodes
func (graph *PackageGraph) ToSet() PackageSet {
	return graph.flatten(nil, graph.Root, make(map[*PackageNode]struct{}))
}
