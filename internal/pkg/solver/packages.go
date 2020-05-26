/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package solver

import (
	"fmt"

	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// Packages is a collect of Pkg objects with dependencies tracked.
type Packages struct {
	packages map[string]*v1alpha2.Pkg
}

// NewPackages builds Packages using PackageLoader.
func NewPackages(loader PackageLoader) (*Packages, error) {
	pkgs, err := loader.Load()
	if err != nil {
		return nil, err
	}

	result := &Packages{
		packages: make(map[string]*v1alpha2.Pkg),
	}

	for _, pkg := range pkgs {
		name := pkg.Name

		if dup, exists := result.packages[name]; exists {
			return nil, fmt.Errorf("package %q already exists, duplicate in dirs %q and %q", name, pkg.BaseDir, dup.BaseDir)
		}

		result.packages[name] = pkg
	}

	return result, nil
}

func (pkgs *Packages) resolve(name string, path []string, cache map[string]*PackageNode) (*PackageNode, error) {
	if node := cache[name]; node != nil {
		return node, nil
	}

	pkg := pkgs.packages[name]
	if pkg == nil {
		return nil, fmt.Errorf("package %q not defined", name)
	}

	for _, pathName := range path {
		if pathName == name {
			return nil, fmt.Errorf("circular dependency detected %v -> %q", path, name)
		}
	}

	path = append(path, name)

	node := &PackageNode{
		Pkg:  pkg,
		Name: name,
	}

	for _, dep := range pkg.Dependencies {
		nodeDep := PackageDependency{
			Dependency: dep,
		}

		if dep.IsInternal() {
			depPkg, err := pkgs.resolve(dep.Stage, path, cache)
			if err != nil {
				return nil, fmt.Errorf("error resolving dependency %q of %q: %w", dep.Stage, name, err)
			}

			nodeDep.Node = depPkg
		}

		node.Dependencies = append(node.Dependencies, nodeDep)
	}

	cache[name] = node

	return node, nil
}

// Resolve trims down the package tree to have only deps of the target.
func (pkgs *Packages) Resolve(target string) (*PackageGraph, error) {
	root, err := pkgs.resolve(target, nil, make(map[string]*PackageNode))
	if err != nil {
		return nil, err
	}

	return &PackageGraph{root}, nil
}

// ToSet converts to set of package nodes.
func (pkgs *Packages) ToSet() (set PackageSet) {
	for name, pkg := range pkgs.packages {
		dependencies := make([]PackageDependency, len(pkg.Dependencies))
		for i := range pkg.Dependencies {
			dependencies[i].Dependency = pkg.Dependencies[i]
		}

		set = append(set, &PackageNode{
			Name:         name,
			Pkg:          pkg,
			Dependencies: dependencies,
		})
	}

	return
}
