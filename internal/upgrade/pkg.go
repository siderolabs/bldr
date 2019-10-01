package upgrade

import (
	"strings"

	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha1"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

// FromV1Alpha1 upgrades v1alpha1 format -> v1alpha2
func FromV1Alpha1(oldPkg *v1alpha1.Pkg, stageNames []string) *v1alpha2.Pkg {
	convertDeps := func(old []*v1alpha1.Dependency) v1alpha2.Dependencies {
		new := v1alpha2.Dependencies{}
		for _, dep := range old {
			src := strings.SplitN(dep.Image, ":", 2)[0]
			parts := strings.Split(src, "/")
			src = parts[len(parts)-1]

			isStage := false
			for _, stageName := range stageNames {
				if src == stageName {
					isStage = true
					break
				}
			}

			if isStage {
				new = append(new, v1alpha2.Dependency{
					Stage: src,
					To:    dep.To,
				})
			} else {
				new = append(new, v1alpha2.Dependency{
					Image: dep.Image,
					To:    dep.To,
				})
			}
		}

		return new
	}

	convertSteps := func(old []*v1alpha1.Step) []v1alpha2.Step {
		new := []v1alpha2.Step{}
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
			new = append(new, newStep)
		}

		return new
	}

	convertFinalize := func(old []*v1alpha1.Finalize) []v1alpha2.Finalize {
		new := []v1alpha2.Finalize{}
		for _, f := range old {
			new = append(new, v1alpha2.Finalize(*f))
		}

		return new
	}

	if oldPkg.Shell == "/bin/sh" {
		oldPkg.Shell = ""
	}

	return &v1alpha2.Pkg{
		Name:         oldPkg.Name,
		Install:      v1alpha2.Install(oldPkg.Install),
		Dependencies: convertDeps(oldPkg.Dependencies),
		Steps:        convertSteps(oldPkg.Steps),
		Finalize:     convertFinalize(oldPkg.Finalize),
		Variant:      v1alpha2.Variant(oldPkg.Variant),
		Shell:        v1alpha2.Shell(oldPkg.Shell),
	}
}
