/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

// Install is a list of Alpine package names to install.
type Install []string

// Environment is a set of environment variables to be set in the step.
type Environment map[string]string

// Step describes a single build step.
//
// Steps are executed sequentially, each step runs in its own
// empty temporary directory.
type Step struct {
	Sources []Source     `yaml:"sources,omitempty"`
	Env     Environment  `yaml:"env,omitempty"`
	Prepare Instructions `yaml:"prepare,omitempty"`
	Build   Instructions `yaml:"build,omitempty"`
	Install Instructions `yaml:"install,omitempty"`
	Test    Instructions `yaml:"test,omitempty"`

	TmpDir string `yaml:"-"`
}

// Finalize is a set of COPY instructions to finalize the build.
type Finalize struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}
