package solver

import (
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// PackageLoader implements some way to fetch collection of Pkgs
type PackageLoader interface {
	Load() ([]*v1alpha2.Pkg, error)
}
