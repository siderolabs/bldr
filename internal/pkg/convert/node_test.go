// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package convert_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/siderolabs/bldr/internal/pkg/convert"
	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

func TestBuildSBOMErrors(t *testing.T) {
	for _, test := range []struct {
		name      string
		wantError string
		metadata  v1alpha2.SBOMStep
	}{
		{"invalid custom license", "invalid custom license ID", v1alpha2.SBOMStep{
			CustomLicenses: []v1alpha2.CustomLicense{{ID: "invalid", Text: "license text"}},
		}},
		{"empty license text", "requires non-empty text", v1alpha2.SBOMStep{
			CustomLicenses: []v1alpha2.CustomLicense{{ID: "LicenseRef-Custom"}},
		}},
		{"overlapping licenses", "duplicate custom license ID", v1alpha2.SBOMStep{
			Licenses:       []string{"MIT OR LicenseRef-Custom"},
			CustomLicenses: []v1alpha2.CustomLicense{{ID: "LicenseRef-Custom", Text: "license text"}},
		}},
		{"invalid CPE", "failed to parse CPE", v1alpha2.SBOMStep{CPEs: []string{"invalid"}}},
		{"valid", "", v1alpha2.SBOMStep{Licenses: []string{"MIT"}}},
	} {
		t.Run(test.name, func(t *testing.T) {
			test.metadata.OutputPath = "/sbom/example.spdx.json"
			graph := convert.NewGraphLLB(&solver.PackageGraph{Root: &solver.PackageNode{
				Name: "example",
				Pkg: &v1alpha2.Pkg{
					Name: "example", Variant: v1alpha2.Scratch,
					Steps:    []v1alpha2.Step{{SBOM: test.metadata}},
					Finalize: []v1alpha2.Finalize{{From: "/sbom", To: "/sbom"}},
				},
			}}, nil, &environment.Options{})

			// Repeat to ensure a failed build is not cached as a success.
			for range 2 {
				_, err := graph.Build(t.Context())
				if test.wantError == "" {
					require.NoError(t, err)
				} else {
					require.ErrorContains(t, err, test.wantError)
				}
			}
		})
	}
}
