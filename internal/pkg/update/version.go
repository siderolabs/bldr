/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package update

import (
	"fmt"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/Masterminds/semver"
)

var (
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
	// extract file name
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}

	s = filepath.Base(u.Path)

	// remove common extensions
	found := true
	for found {
		ext := filepath.Ext(s)
		if _, found = commonExtensions[ext]; found {
			s = strings.TrimSuffix(s, ext)
		}
	}

	// remove package name, keep only version
	i := strings.IndexAny(s, "0123456789")
	if i < 0 {
		return nil, fmt.Errorf("failed to remove package name from %q", s)
	}

	s = s[i:]

	res, err := semver.NewVersion(s)
	if err != nil {
		return nil, fmt.Errorf("%q: %w", s, err)
	}

	return res, nil
}
