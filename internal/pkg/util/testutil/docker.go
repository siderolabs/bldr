// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testutil

import (
	"os/exec"
	"sync"
	"testing"
)

// DockerRunner runs bldr via docker buildx.
type DockerRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner DockerRunner) Run(t *testing.T) {
	if err := IsDockerAvailable(); err != nil {
		t.Skipf("docker buildx is not available: %q", err)
	}

	args := []string{
		"buildx",
		"build",
		"-f", "./Pkgfile",
		"--target", runner.Target,
		"--build-arg", "TAG=testtag",
	}

	if runner.Platform != "" {
		args = append(args, "--platform", runner.Platform)
	}

	cmd := exec.Command("docker", append(args, ".")...)

	runner.run(t, cmd, "docker buildx")
}

var (
	dockerCheckOnce sync.Once
	//nolint:errname
	dockerCheckError error
)

// IsDockerAvailable returns nil if docker buildx is ready to use.
func IsDockerAvailable() error {
	dockerCheckOnce.Do(func() {
		dockerCheckError = exec.Command("docker", "buildx", "ls").Run()
	})

	return dockerCheckError
}
