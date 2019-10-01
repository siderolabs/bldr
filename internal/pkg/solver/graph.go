package solver

import (
	"io"

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
	for _, dep := range node.DependsOn {
		n.Edge(dep.DumpDot(g))
	}

	return n
}

// PackageGraph capture root of the DAG
type PackageGraph struct {
	Root *PackageNode
}

// DumpDot dumps whole graph in dot format
func (graph *PackageGraph) DumpDot(w io.Writer) {
	g := dot.NewGraph(dot.Directed)
	graph.Root.DumpDot(g)

	g.Write(w)
}
