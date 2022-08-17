// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package update

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Masterminds/semver"
)

var (
	versionRE = regexp.MustCompile(semver.SemVerRegex) // the same as semver.versionRegex except "^" and "$"

	commonExtensions = map[string]struct{}{
		".bz2":  {},
		".diff": {},
		".gz":   {},
		".orig": {},
		".src":  {},
		".tar":  {},
		".xdp":  {},
		".xz":   {},
	}
)

// extractVersion extracts SemVer version from file name or URL.
func extractVersion(s string) (*semver.Version, error) {
	// remove common extensions like .bz2 that would confuse SemVer parser
	found := true
	for found {
		ext := filepath.Ext(s)
		if _, found = commonExtensions[ext]; found {
			s = strings.TrimSuffix(s, ext)
		}
	}

	matches := versionRE.FindAllString(s, -1)
	if len(matches) == 0 {
		return nil, fmt.Errorf("failed to find version in %q", s)
	}

	// use the last match to skip hostnames, folders, etc
	res, err := semver.NewVersion(matches[len(matches)-1])
	if err != nil {
		return nil, fmt.Errorf("%q: %w", s, err)
	}

	return res, nil
}
