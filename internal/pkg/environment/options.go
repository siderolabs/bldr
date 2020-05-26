/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package environment

import "github.com/talos-systems/bldr/internal/pkg/types"

// Options for bldr.
type Options struct {
	BuildPlatform  Platform
	TargetPlatform Platform
	Target         string
	CommonPrefix   string
}

// GetVariables returns set of variables set for options.
func (options *Options) GetVariables() types.Variables {
	return Default().
		Merge(options.BuildPlatform.BuildVariables()).
		Merge(options.TargetPlatform.TargetVariables())
}
