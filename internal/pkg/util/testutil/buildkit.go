/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"os"
	"os/exec"
	"sync"
	"testing"
)

// BuildkitRunner runs bldr via buildctl/buildkit.
type BuildkitRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner BuildkitRunner) Run(t *testing.T) {
	if err := IsBuildkitAvailable(); err != nil {
		t.Skipf("buildkit is not available: %q", err)
	}

	args := append(getBuildkitGlobalFlags(),
		"build",
		"--frontend", "dockerfile.v0",
		"--local", "context=.",
		"--local", "dockerfile=.",
		"--opt", "filename=Pkgfile",
		"--opt", "target="+runner.Target,
	)

	if runner.Platform != "" {
		args = append(args, "--opt", "platform="+runner.Platform)
	}

	cmd := exec.Command("buildctl", args...)

	runner.run(t, cmd, "buildkit")
}

func getBuildkitGlobalFlags() []string {
	var globalOpts []string

	if buildkitHost, ok := os.LookupEnv("BUILDKIT_HOST"); ok {
		globalOpts = append(globalOpts, "--addr", buildkitHost)
	}

	return globalOpts
}

var (
	buildkitCheckOnce  sync.Once
	buildkitCheckError error
)

// IsBuildkitAvailable returns nil if buildkit is ready to use.
func IsBuildkitAvailable() error {
	buildkitCheckOnce.Do(func() {
		buildkitCheckError = exec.Command("buildctl", append(getBuildkitGlobalFlags(), "debug", "workers")...).Run()
	})

	return buildkitCheckError
}
