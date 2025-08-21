// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testutil

import (
	"os/exec"
	"sync"
	"testing"

	"github.com/siderolabs/bldr/internal/pkg/constants"
)

// DockerRunner runs bldr via docker buildx.
type DockerRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner DockerRunner) Run(t *testing.T) {
	if err := IsDockerAvailable(t); err != nil {
		t.Skipf("docker buildx is not available: %q", err)
	}

	args := []string{
		"buildx",
		"build",
		"--progress=plain",
		"-f", "./Pkgfile",
		"--target", runner.Target,
		"--build-arg", "TAG=testtag",
		"--build-arg", "BLDR_TAG=" + constants.Version,
	}
	if runner.Platform != "" {
		args = append(args, "--platform", runner.Platform)
	}

	cmd := exec.CommandContext(t.Context(), "docker", append(args, ".")...)

	runner.run(t, cmd, "docker buildx")
}

var (
	dockerCheckOnce sync.Once
	//nolint:errname
	dockerCheckError error
)

// IsDockerAvailable returns nil if docker buildx is ready to use.
func IsDockerAvailable(t *testing.T) error {
	dockerCheckOnce.Do(func() {
		dockerCheckError = exec.CommandContext(t.Context(), "docker", "buildx", "ls").Run()
	})

	return dockerCheckError
}
