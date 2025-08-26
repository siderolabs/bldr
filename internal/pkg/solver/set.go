// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package solver

import (
	"cmp"
	"encoding/json"
	"io"
	"slices"
	"text/template"

	"github.com/emicklei/dot"
)

// PackageSet is a list of PackageNodes.
type PackageSet []*PackageNode

// DumpDot dumps nodes and deps in dot format.
func (set PackageSet) DumpDot(w io.Writer) {
	g := dot.NewGraph(dot.Directed)

	for _, node := range set {
		node.DumpDot(g)
	}

	g.Write(w)
}

// Sorted returns a new set which is sorted by name package set.
func (set PackageSet) Sorted() PackageSet {
	return slices.SortedFunc(slices.Values(set), func(a, b *PackageNode) int {
		return cmp.Compare(a.Name, b.Name)
	})
}

// DumpJSON dumps the package set as JSON.
func (set PackageSet) DumpJSON(w io.Writer) error {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")

	return encoder.Encode(set)
}

// Template dumps the package set as a Go template.
func (set PackageSet) Template(w io.Writer, tmpl *template.Template) error {
	return tmpl.Execute(w, set)
}
