// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package environment

import (
	"fmt"

	"github.com/containerd/platforms"
	"github.com/moby/buildkit/client/llb"
	specs "github.com/opencontainers/image-spec/specs-go/v1"

	"github.com/siderolabs/bldr/internal/pkg/types"
)

// Platform describes build & target platforms.
//
//nolint:recvcheck
type Platform struct {
	ID           string
	Arch         string
	Target       string
	Build        string
	Host         string
	LLBPlatform  llb.ConstraintsOpt
	PlatformSpec specs.Platform
}

// BuildVariables returns build env variables.
func (p Platform) BuildVariables() types.Variables {
	return types.Variables{
		"BUILD": p.Build,
		"HOST":  p.Host,
	}
}

// TargetVariables returns target env variables.
func (p Platform) TargetVariables() types.Variables {
	return types.Variables{
		"ARCH":   p.Arch,
		"TARGET": p.Target,
	}
}

func (p Platform) String() string {
	return p.ID
}

// Set implements pflag.Value interface.
func (p *Platform) Set(id string) error {
	if _, exists := Platforms[id]; !exists {
		return fmt.Errorf("platform %q is not defined", id)
	}

	*p = Platforms[id]

	return nil
}

// Type implements pflag.Value interface.
func (p *Platform) Type() string {
	return "platform"
}

// Platform definitions.
var (
	LinuxAmd64 = Platform{
		ID:           "linux/amd64",
		Arch:         "x86_64",
		Target:       "x86_64-talos-linux-musl",
		Build:        "x86_64-linux-musl",
		Host:         "x86_64-linux-musl",
		LLBPlatform:  llb.LinuxAmd64,
		PlatformSpec: platforms.MustParse("linux/amd64"),
	}

	LinuxArm64 = Platform{
		ID:           "linux/arm64",
		Arch:         "aarch64",
		Target:       "aarch64-talos-linux-musl",
		Build:        "aarch64-linux-musl",
		Host:         "aarch64-linux-musl",
		LLBPlatform:  llb.LinuxArm64,
		PlatformSpec: platforms.MustParse("linux/arm64"),
	}

	LinuxArmv7 = Platform{
		ID:           "linux/armv7",
		Arch:         "armv7",
		Target:       "armv7-talos-linux-musl",
		Build:        "armv7-linux-musl",
		Host:         "armv7-linux-musl",
		LLBPlatform:  llb.LinuxArmhf,
		PlatformSpec: platforms.MustParse("linux/arm7"),
	}

	LinuxRiscv64 = Platform{
		ID:           "linux/riscv64",
		Arch:         "riscv64",
		Target:       "riscv64-talos-linux-musl",
		Build:        "riscv64-linux-musl",
		Host:         "riscv64-linux-musl",
		LLBPlatform:  llb.Platform(specs.Platform{OS: "linux", Architecture: "riscv64"}),
		PlatformSpec: platforms.MustParse("linux/riscv64"),
	}
)

// Platforms is mapping of platform ID to Platform.
var Platforms = map[string]Platform{}

func init() {
	for _, platform := range []Platform{
		LinuxAmd64,
		LinuxArm64,
		LinuxArmv7,
		LinuxRiscv64,
	} {
		Platforms[platform.ID] = platform
	}
}
