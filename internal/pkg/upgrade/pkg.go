// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package upgrade contains function to upgrade between pkg formats
package upgrade

import (
	"slices"
	"strings"

	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha1"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

func convertDeps(stageNames []string, old []*v1alpha1.Dependency) v1alpha2.Dependencies {
	newDep := v1alpha2.Dependencies{}

	for _, dep := range old {
		src := strings.SplitN(dep.Image, ":", 2)[0]
		parts := strings.Split(src, "/")
		src = parts[len(parts)-1]

		isStage := slices.Contains(stageNames, src)

		if isStage {
			newDep = append(newDep, v1alpha2.Dependency{
				Stage: src,
				To:    dep.To,
			})
		} else {
			newDep = append(newDep, v1alpha2.Dependency{
				Image: dep.Image,
				To:    dep.To,
			})
		}
	}

	return newDep
}

func convertSteps(old []*v1alpha1.Step) []v1alpha2.Step {
	newSteps := []v1alpha2.Step{}

	for _, step := range old {
		newStep := v1alpha2.Step{}

		if step.Prepare != nil {
			newStep.Prepare = v1alpha2.Instructions{v1alpha2.Instruction(*step.Prepare)}
		}

		if step.Build != nil {
			newStep.Build = v1alpha2.Instructions{v1alpha2.Instruction(*step.Build)}
		}

		if step.Install != nil {
			newStep.Install = v1alpha2.Instructions{v1alpha2.Instruction(*step.Install)}
		}

		if step.Test != nil {
			newStep.Test = v1alpha2.Instructions{v1alpha2.Instruction(*step.Test)}
		}

		for _, src := range step.Sources {
			newStep.Sources = append(newStep.Sources, v1alpha2.Source(*src))
		}

		newSteps = append(newSteps, newStep)
	}

	return newSteps
}

func convertFinalize(old []*v1alpha1.Finalize) []v1alpha2.Finalize {
	newFinalize := []v1alpha2.Finalize{}
	for _, f := range old {
		newFinalize = append(newFinalize, v1alpha2.Finalize(*f))
	}

	return newFinalize
}

// FromV1Alpha1 upgrades v1alpha1 format -> v1alpha2.
func FromV1Alpha1(oldPkg *v1alpha1.Pkg, stageNames []string) *v1alpha2.Pkg {
	if oldPkg.Shell == "/bin/sh" {
		oldPkg.Shell = ""
	}

	return &v1alpha2.Pkg{
		Name:         oldPkg.Name,
		Install:      v1alpha2.Install(oldPkg.Install),
		Dependencies: convertDeps(stageNames, oldPkg.Dependencies),
		Steps:        convertSteps(oldPkg.Steps),
		Finalize:     convertFinalize(oldPkg.Finalize),
		Variant:      v1alpha2.Variant(oldPkg.Variant),
		Shell:        v1alpha2.Shell(oldPkg.Shell),
	}
}
