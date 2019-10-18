/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"os/exec"
	"testing"
)

// ValidateRunner runs bldr validate
type ValidateRunner struct {
	CommandRunner
}

// Run implements Run interface
func (runner ValidateRunner) Run(t *testing.T) {
	cmd := exec.Command("bldr", "validate")

	runner.run(t, cmd, "bldr validate")
}
