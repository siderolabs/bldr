/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"io"

	"github.com/emicklei/dot"
)

// PackageSet is a list of PackageNodes
type PackageSet []*PackageNode

// DumpDot dumps nodes and deps in dot format
func (set PackageSet) DumpDot(w io.Writer) {
	g := dot.NewGraph(dot.Directed)

	for _, node := range set {
		node.DumpDot(g)
	}

	g.Write(w)
}
