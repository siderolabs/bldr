// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package sbom_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.yaml.in/yaml/v4"

	"github.com/siderolabs/bldr/internal/pkg/sbom"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

func TestCustomLicenses(t *testing.T) {
	const recipe = `name: libtpms
version: 0.10.2
licenses:
  - BSD-3-Clause
customLicenses:
  - id: LicenseRef-TCGL
    text: |
      Custom license text.
      Preserve this second line.
  - id: LicenseRef-Another.1
    text: Another license.
`

	var metadata v1alpha2.SBOMStep

	require.NoError(t, yaml.Unmarshal([]byte(recipe), &metadata))

	created := time.Unix(0, 0).UTC()

	var previous string

	for range 3 {
		doc, err := sbom.CreatePackageSBOM(&v1alpha2.Pkg{Name: "libtpms"}, metadata)
		require.NoError(t, err)

		encoded, err := sbom.ToSpdxJSON(*doc, created)
		require.NoError(t, err)

		var decoded struct {
			Packages []struct {
				Name            string `json:"name"`
				LicenseDeclared string `json:"licenseDeclared"`
			} `json:"packages"`
			Licenses []struct {
				ID   string `json:"licenseId"`
				Text string `json:"extractedText"`
			} `json:"hasExtractedLicensingInfos"`
		}

		require.NoError(t, json.Unmarshal([]byte(encoded), &decoded))
		require.Len(t, decoded.Licenses, 2)

		texts := make(map[string]string)
		for _, entry := range decoded.Licenses {
			texts[entry.ID] = entry.Text
		}

		assert.Equal(t, map[string]string{
			"LicenseRef-TCGL":      "Custom license text.\nPreserve this second line.\n",
			"LicenseRef-Another.1": "Another license.",
		}, texts)

		found := false

		for _, entry := range decoded.Packages {
			if entry.Name == "libtpms" {
				found = true

				assert.Equal(t, "BSD-3-Clause AND LicenseRef-Another.1 AND LicenseRef-TCGL", entry.LicenseDeclared)
			}
		}

		assert.True(t, found)

		if previous != "" {
			assert.Equal(t, previous, encoded)
		}

		previous = encoded
	}
}

func TestLegacyLicenses(t *testing.T) {
	doc, err := sbom.CreatePackageSBOM(&v1alpha2.Pkg{Name: "example"}, v1alpha2.SBOMStep{
		Licenses: []string{"MIT OR Apache-2.0", "BSD-3-Clause"},
	})
	require.NoError(t, err)

	encoded, err := sbom.ToSpdxJSON(*doc, time.Unix(0, 0).UTC())
	require.NoError(t, err)
	assert.Contains(t, encoded, `"licenseDeclared": "BSD-3-Clause AND (MIT OR Apache-2.0)"`)
	assert.NotContains(t, encoded, "hasExtractedLicensingInfos")
}

func TestDistinctCustomLicenseReferences(t *testing.T) {
	for _, expression := range []string{
		"MIT OR LicenseRef-Custom-Other",
		"MIT OR DocumentRef-other:LicenseRef-Custom",
		"legacy non-SPDX license",
	} {
		t.Run(expression, func(t *testing.T) {
			_, err := sbom.CreatePackageSBOM(&v1alpha2.Pkg{Name: "example"}, v1alpha2.SBOMStep{
				Licenses:       []string{expression},
				CustomLicenses: []v1alpha2.CustomLicense{{ID: "LicenseRef-Custom", Text: "license text"}},
			})
			require.NoError(t, err)
		})
	}
}

func TestInvalidCustomLicenses(t *testing.T) {
	for _, test := range []struct {
		name   string
		recipe string
	}{
		{"missing ID", "customLicenses: [{text: license text}]"},
		{"missing prefix", "customLicenses: [{id: TCGL, text: license text}]"},
		{"empty suffix", "customLicenses: [{id: LicenseRef-, text: license text}]"},
		{"invalid ID", "customLicenses: [{id: LicenseRef-invalid_id, text: license text}]"},
		{"expression", "customLicenses: [{id: 'LicenseRef-A OR MIT', text: license text}]"},
		{"missing text", "customLicenses: [{id: LicenseRef-TCGL}]"},
		{"blank text", "customLicenses: [{id: LicenseRef-TCGL, text: '  '}]"},
		{"legacy overlap", "licenses: [LicenseRef-TCGL]\ncustomLicenses: [{id: LicenseRef-TCGL, text: license text}]"},
		{"expression overlap", "licenses: ['MIT OR LicenseRef-TCGL']\ncustomLicenses: [{id: LicenseRef-TCGL, text: license text}]"},
		{"nested expression overlap", "licenses: ['MIT AND (LicenseRef-TCGL OR Apache-2.0)']\ncustomLicenses: [{id: LicenseRef-TCGL, text: license text}]"},
		{"duplicate ID", "customLicenses: [{id: LicenseRef-TCGL, text: first}, {id: LicenseRef-TCGL, text: second}]"},
	} {
		t.Run(test.name, func(t *testing.T) {
			var metadata v1alpha2.SBOMStep

			require.NoError(t, yaml.Unmarshal([]byte(test.recipe), &metadata))

			doc, err := sbom.CreatePackageSBOM(&v1alpha2.Pkg{Name: "example"}, metadata)
			require.Error(t, err)
			assert.Nil(t, doc)
		})
	}
}
