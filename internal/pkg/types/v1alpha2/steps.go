// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import "github.com/hashicorp/go-multierror"

// Environment is a set of environment variables to be set in the step.
type Environment map[string]string

// Steps is a collection of Step.
type Steps []Step

// Validate steps.
func (steps Steps) Validate() error {
	var multiErr *multierror.Error

	for _, step := range steps {
		multiErr = multierror.Append(multiErr, step.Validate())
	}

	return multiErr.ErrorOrNil()
}

// Step describes a single build step.
//
// Steps are executed sequentially, each step runs in its own
// empty temporary directory.
type Step struct {
	Env        Environment  `yaml:"env,omitempty"`
	CachePaths []string     `yaml:"cachePaths,omitempty"`
	TmpDir     string       `yaml:"-"`
	Sources    Sources      `yaml:"sources,omitempty"`
	Prepare    Instructions `yaml:"prepare,omitempty"`
	Build      Instructions `yaml:"build,omitempty"`
	Install    Instructions `yaml:"install,omitempty"`
	Test       Instructions `yaml:"test,omitempty"`
}

// Validate the step.
func (step *Step) Validate() error {
	return step.Sources.Validate()
}
