/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// LoadResult is a result of PackageLoader.Load function.
type LoadResult struct {
	Pkgfile *v1alpha2.Pkgfile
	Pkgs    []*v1alpha2.Pkg
}

// PackageLoader implements some way to fetch collection of Pkgs.
type PackageLoader interface {
	Load() (*LoadResult, error)
}
