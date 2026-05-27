// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package environment

import (
	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/types"
)

// Default returns default values for environment variables.
//
// CFLAGS/CXXFLAGS are platform-specific and contributed by [Platform.TargetVariables].
func Default() types.Variables {
	return types.Variables{
		"LDFLAGS": "-s",
		"VENDOR":  "talos",
		"SYSROOT": "/talos",
		"PATH":    constants.DefaultPath,
	}
}
