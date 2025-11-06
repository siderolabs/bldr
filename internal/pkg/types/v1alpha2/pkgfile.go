// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"fmt"

	"go.yaml.in/yaml/v4"

	"github.com/siderolabs/bldr/internal/pkg/types"
)

// Pkgfile describes structure of 'Pkgfile'.
type Pkgfile struct {
	Vars   types.Variables   `yaml:"vars,omitempty"`
	Labels map[string]string `yaml:"labels,omitempty"`
	Format string            `yaml:"format"`
}

// NewPkgfile loads Pkgfile from `[]byte` contents.
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
