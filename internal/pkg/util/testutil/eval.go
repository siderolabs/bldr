// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package testutil

import (
	"os/exec"
	"testing"
)

// EvalRunner runs bldr eval.
type EvalRunner struct {
	CommandRunner

	Target   string
	Template string
}

// Run implements Run interface.
func (runner EvalRunner) Run(t *testing.T) {
	cmd := exec.Command("bldr", "eval", "--target", runner.Target, runner.Template)

	runner.run(t, cmd, "bldr eval")
}
