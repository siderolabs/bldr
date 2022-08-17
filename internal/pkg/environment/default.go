// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package environment

import (
	"fmt"

	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/types"
)

// Default returns default values for environment variables.
func Default() types.Variables {
	return types.Variables{
		"CFLAGS":    "-g0 -Os",
		"CXXFLAGS":  "-g0 -Os",
		"LDFLAGS":   "-s",
		"VENDOR":    "talos",
		"SYSROOT":   "/talos",
		"TOOLCHAIN": "/toolchain",
		"PATH":      fmt.Sprintf("/toolchain/bin:%s", constants.DefaultPath),
	}
}
