/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"os"

	"gopkg.in/yaml.v2"
)

// TestManifest describes single integration test in test.ymal.
type TestManifest struct {
	Runs []RunManifest `yaml:"run"`
}

// RunManifest describes single run of integration test.
type RunManifest struct {
	Name   string `yaml:"name"`
	Runner string `yaml:"runner"`
	Target string `yaml:"target"`
	Expect string `yaml:"expect"`
}

// NewTestManifest loads TestManifest from test.yaml file.
func NewTestManifest(path string) (manifest TestManifest, err error) {
	var f *os.File

	f, err = os.Open(path)
	if err != nil {
		return
	}

	defer f.Close() //nolint: errcheck

	err = yaml.NewDecoder(f).Decode(&manifest)

	return
}
