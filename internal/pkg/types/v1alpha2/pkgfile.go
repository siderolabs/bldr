/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

import (
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/talos-systems/bldr/internal/pkg/types"
)

// Pkgfile describes structure of 'Pkgfile'
type Pkgfile struct {
	Format string          `yaml:"format"`
	Vars   types.Variables `yaml:"vars,omitempty"`
}

// NewPkgfile loads Pkgfile from `[]byte` contents
func NewPkgfile(contents []byte) (*Pkgfile, error) {
	var pkgfile Pkgfile

	if err := yaml.Unmarshal(contents, &pkgfile); err != nil {
		return nil, err
	}

	// TODO: this might be used in the future to pick correct format
	//       based on Pkgfile, leave it simple for now
	if pkgfile.Format != "v1alpha2" {
		return nil, fmt.Errorf("unsupported format: %q, supported formats: %q", pkgfile.Format, []string{"v1alpha2"})
	}

	return &pkgfile, nil
}
