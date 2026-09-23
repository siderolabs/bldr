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

	// CustomLicenses are additional declared licenses, combined with Licenses using AND.
	CustomLicenses []CustomLicense `yaml:"customLicenses,omitempty"`
}

// CustomLicense provides the text of a license not on the SPDX license list.
type CustomLicense struct {
	// ID is an SPDX LicenseRef- identifier, unique within the SBOM step.
	ID string `yaml:"id"`
	// Text is the full, non-empty license text, preserved verbatim in the SBOM.
	Text string `yaml:"text"`
}
