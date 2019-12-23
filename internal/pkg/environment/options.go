/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package environment

import (
	"github.com/containerd/containerd/platforms"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/pkg/errors"
)

// Options for bldr
type Options struct {
	Platform          specs.Platform
	ToolchainPlatform ToolchainPlatform
	Target            string
}

// NewOptions initializes and returns default options for the runtime
// environment.
func NewOptions() (*Options, error) {
	defaultSpec := platforms.DefaultSpec()

	options := &Options{}
	if err := options.Set(defaultSpec.OS + "/" + defaultSpec.Architecture); err != nil {
		return nil, err
	}

	return options, nil
}

// Set sets the options based on the specified platform.
func (options *Options) Set(v string) error {
	platform, err := parsePlatform(v)
	if err != nil {
		return err
	}

	options.Platform = *platform

	if err := options.ToolchainPlatform.Set(v); err != nil {
		return err
	}

	return nil
}

func parsePlatform(v string) (*specs.Platform, error) {
	p, err := platforms.Parse(v)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse target platform %s", v)
	}

	p = platforms.Normalize(p)

	return &p, nil
}
