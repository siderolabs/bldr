// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testutil

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"

	"github.com/alessio/shellescape"
)

// LLBRunner runs bldr via bldr llb | buildctl.
type LLBRunner struct {
	CommandRunner
	Target   string
	Platform string
}

// Run implements Run interface.
func (runner LLBRunner) Run(t *testing.T) {
	if err := IsBuildkitAvailable(); err != nil {
		t.Skipf("buildkit is not available: %q", err)
	}

	args := getBuildkitGlobalFlags()
	for i := range args {
		args[i] = shellescape.Quote(args[i])
	}

	platformArgs := ""
	if runner.Platform != "" {
		platformArgs = fmt.Sprintf("--build-platform=%s --target-platform=%s", shellescape.Quote(runner.Platform), shellescape.Quote(runner.Platform))
	}

	cmd := exec.Command("/bin/sh", "-c",
		fmt.Sprintf("bldr llb --target=%s %s | buildctl %s build --local context=.", shellescape.Quote(runner.Target), platformArgs, strings.Join(args, " ")),
	)

	runner.run(t, cmd, "bldr llb")
}
