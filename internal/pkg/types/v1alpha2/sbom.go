// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

// SBOMStep is a step with data to generate an SBOM (Software Bill of Materials) for the package.
type SBOMStep struct {
	OutputPath string   `yaml:"outputPath,omitempty"`
	Name       string   `yaml:"name,omitempty"`
	Version    string   `yaml:"version,omitempty"`
	CPEs       []string `yaml:"cpes,omitempty"`
	PURL       string   `yaml:"purl,omitempty"`
	Licenses   []string `yaml:"licenses,omitempty"`
}
