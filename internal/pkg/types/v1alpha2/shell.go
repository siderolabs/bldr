// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

// Shell is a path to the shell used to execute Instructions.
type Shell string

// Get returns current shell.
func (sh Shell) Get() string {
	if sh == "" {
		return "/bin/sh"
	}

	return string(sh)
}
