/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

// Package types describes basic types which are not versioned.
package types

// Variables presents generic variables for templating/environment.
type Variables map[string]string

// Merge two Variables.
func (v Variables) Merge(other Variables) Variables {
	for key := range other {
		v[key] = other[key]
	}

	return v
}
