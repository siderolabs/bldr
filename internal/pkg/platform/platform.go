// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package platform provides a function to convert a platform string to a v1.Platform.
package platform

import (
	"fmt"

	v1 "github.com/opencontainers/image-spec/specs-go/v1"
)

// ToV1Platform converts a platform string to a v1.Platform.
func ToV1Platform(platform, targetPlatform string) (v1.Platform, error) {
	if platform == "" {
		platform = targetPlatform
	}

	switch platform {
	case "linux/amd64":
		return v1.Platform{
			OS:           "linux",
			Architecture: "amd64",
		}, nil
	case "linux/arm64":
		return v1.Platform{
			OS:           "linux",
			Architecture: "arm64",
		}, nil
	default:
		return v1.Platform{}, fmt.Errorf("unknown platform %q", platform)
	}
}
