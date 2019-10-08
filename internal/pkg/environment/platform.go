package environment

import (
	"fmt"

	"github.com/moby/buildkit/client/llb"

	"github.com/talos-systems/bldr/internal/pkg/types"
)

// Platform describes build & target platforms
type Platform struct {
	ID          string
	Arch        string
	Target      string
	Build       string
	Host        string
	LLBPlatform llb.ConstraintsOpt
}

// BuildVariables returns build env variables
func (p Platform) BuildVariables() types.Variables {
	return types.Variables{
		"BUILD": p.Build,
		"HOST":  p.Host,
	}
}

// TargetVariables returns target env variables
func (p Platform) TargetVariables() types.Variables {
	return types.Variables{
		"ARCH":   p.Arch,
		"TARGET": p.Target,
	}
}

func (p Platform) String() string {
	return p.ID
}

// Set implements pflag.Value interface
func (p *Platform) Set(id string) error {
	if _, exists := Platforms[id]; !exists {
		return fmt.Errorf("platform %q is not defined", id)
	}

	*p = Platforms[id]
	return nil
}

// Type implements pflag.Value interface
func (p *Platform) Type() string {
	return "platform"
}

// Platform definitions
var (
	LinuxAmd64 = Platform{
		ID:          "linux/amd64",
		Arch:        "x86_64",
		Target:      "x86_64-talos-linux-musl",
		Build:       "x86_64-linux-musl",
		Host:        "x86_64-linux-musl",
		LLBPlatform: llb.LinuxAmd64,
	}

	LinuxArm64 = Platform{
		ID:          "linux/arm64",
		Arch:        "aarch64",
		Target:      "aarch64-talos-linux-musl",
		Build:       "aarch64-linux-musl",
		Host:        "aarch64-linux-musl",
		LLBPlatform: llb.LinuxArm64,
	}

	LinuxArmv7 = Platform{
		ID:          "linux/armv7",
		Arch:        "armv7",
		Target:      "armv7-talos-linux-musl",
		Build:       "armv7-linux-musl",
		Host:        "armv7-linux-musl",
		LLBPlatform: llb.LinuxArmhf,
	}
)

// Platforms is mapping of platform ID to Platform
var Platforms = map[string]Platform{}

func init() {
	for _, platform := range []Platform{
		LinuxAmd64,
		LinuxArm64,
		LinuxArmv7,
	} {
		Platforms[platform.ID] = platform
	}
}
