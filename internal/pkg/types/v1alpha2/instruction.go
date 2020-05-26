/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

// Instructions is a list of shell commands.
type Instructions []Instruction

// Instruction is a single shell command.
type Instruction string

// Script formats Instruction for /bin/sh -c execution.
func (ins Instruction) Script() string {
	return "set -eou pipefail\n" + string(ins)
}
