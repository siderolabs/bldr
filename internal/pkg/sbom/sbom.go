// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package sbom contains utilities for SBOM generation
package sbom

import (
	"fmt"
	"time"

	"github.com/anchore/syft/syft/artifact"
	"github.com/anchore/syft/syft/cpe"
	"github.com/anchore/syft/syft/file"
	"github.com/anchore/syft/syft/format"
	"github.com/anchore/syft/syft/format/spdxjson"
	"github.com/anchore/syft/syft/pkg"
	"github.com/anchore/syft/syft/sbom"
	"github.com/anchore/syft/syft/source"

	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
	"github.com/siderolabs/bldr/internal/version"
)

func parseCPEs(cpeStrings []string) ([]cpe.CPE, error) {
	cpes := []cpe.CPE{}

	for _, cpeStr := range cpeStrings {
		cpe, err := cpe.New(cpeStr, cpe.NVDDictionaryLookupSource)
		if err != nil {
			return nil, fmt.Errorf("failed to parse CPE %q: %w", cpeStr, err)
		}

		cpes = append(cpes, cpe)
	}

	return cpes, nil
}

func addPkgSources(sbomDoc *sbom.SBOM, bldrPkg *v1alpha2.Pkg, syftPkg pkg.Package) {
	// TODO: also add sources from the package itself, like the Pkgfile and patches
	for _, step := range bldrPkg.Steps {
		for _, source := range step.Sources {
			fileCoordinates := file.NewCoordinates(source.URL, "bldr sources")

			sbomDoc.Artifacts.FileDigests[fileCoordinates] = []file.Digest{
				{
					Algorithm: "sha256",
					Value:     source.SHA256,
				},
				{
					Algorithm: "sha512",
					Value:     source.SHA512,
				},
			}

			// Make sure the file is linked to the package
			sbomDoc.Relationships = append(sbomDoc.Relationships, artifact.Relationship{
				From: syftPkg,
				To:   fileCoordinates,
				Type: artifact.ContainsRelationship,
			})
		}
	}
}

// CreatePackageSBOM populates an SBOM document with data from the provided package.
func CreatePackageSBOM(bldrPkg *v1alpha2.Pkg) (*sbom.SBOM, error) {
	// Sample. Actually we would get metadata when instructed to generate SBOM
	sbomMetadata := bldrPkg.Steps[0].SBOM

	cpes, err := parseCPEs(sbomMetadata.CPEs)
	if err != nil {
		return nil, err
	}

	sbomDoc := &sbom.SBOM{
		Source: source.Description{
			ID:       "sidero-pkgs",
			Metadata: source.DirectoryMetadata{},
			Name:     "sidero-pkgs-" + bldrPkg.Name,
			Version:  sbomMetadata.Version,
		},
		Descriptor: sbom.Descriptor{
			Name:    "bldr",
			Version: version.Tag,
		},
		Artifacts: sbom.Artifacts{
			Packages:    pkg.NewCollection(),
			FileDigests: make(map[file.Coordinates][]file.Digest),
		},
		Relationships: []artifact.Relationship{},
	}

	syftPkg := pkg.Package{
		Name:    bldrPkg.Name,
		Version: sbomMetadata.Version,
		PURL:    sbomMetadata.PURL,
		Type:    pkg.Type("bldr-package"),
		FoundBy: "bldr",
		Locations: file.NewLocationSet(
			file.NewLocation("/" + bldrPkg.BaseDir + "/Pkgfile"),
		),
		CPEs:     cpes,
		Licenses: pkg.NewLicenseSet(pkg.NewLicensesFromValues(sbomMetadata.Licenses...)...),
	}

	sbomDoc.Artifacts.Packages.Add(syftPkg)

	addPkgSources(sbomDoc, bldrPkg, syftPkg)

	return sbomDoc, nil
}

// ToSpdxJSON formats the SBOM document into SPDX JSON format, using deterministic options.
func ToSpdxJSON(sbomDoc sbom.SBOM, createdTime time.Time) (string, error) {
	cfg := spdxjson.DefaultEncoderConfig()
	cfg.Pretty = true
	// Use UUIDv5 to make namespace deterministic based on the content
	cfg.DeterministicUUID = true
	cfg.CreatedTime = &createdTime

	enc, err := spdxjson.NewFormatEncoderWithConfig(cfg)
	if err != nil {
		return "", err
	}

	bytes, err := format.Encode(sbomDoc, enc)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}
