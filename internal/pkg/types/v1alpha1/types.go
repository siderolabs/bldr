// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package v1alpha1 is the v1alpha1 version of the API.
package v1alpha1

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Options are the options for the build.
type Options struct {
	CacheTo      string
	CacheFrom    string
	Organization string
	Platform     string
	Progress     string
	Push         string
	Registry     string
}

// Install configures the install packages.
type Install []string

// Dependency defined the dependency for a package.
type Dependency struct {
	Image string `yaml:"image,omitempty"`
	To    string `yaml:"to,omitempty"`
}

// Instruction defines the instruction for a package.
type Instruction string

// Source define a package source options.
type Source struct {
	URL         string `yaml:"url,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	SHA256      string `yaml:"sha256,omitempty"`
	SHA512      string `yaml:"sha512,omitempty"`
}

// Step defines a step in the build process.
type Step struct {
	Prepare *Instruction `yaml:"prepare,omitempty"`
	Build   *Instruction `yaml:"build,omitempty"`
	Install *Instruction `yaml:"install,omitempty"`
	Test    *Instruction `yaml:"test,omitempty"`
	Sources []*Source    `yaml:"sources,omitempty"`
}

// Finalize the package.
type Finalize struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}

// Variant defines the variant of the package.
type Variant int

const (
	// Alpine is the Alpine variant.
	Alpine Variant = iota
	// Scratch is the Scratch variant.
	Scratch
)

func (v Variant) String() string {
	return []string{"alpine", "scratch"}[v]
}

// UnmarshalYAML implements the yaml.Unmarshaler interface.
func (v *Variant) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string

	if err := unmarshal(&aux); err != nil {
		return err
	}

	var val Variant

	switch aux {
	case Alpine.String():
		val = Alpine
	case Scratch.String():
		val = Scratch
	default:
		return fmt.Errorf("unknown variant %q", aux)
	}

	*v = val

	return nil
}

// Shell is the shell to use for the package.
type Shell string

// Pkg defines the package to build.
type Pkg struct {
	Options      *Options
	Name         string `yaml:"name,omitempty"`
	Bldr         string `yaml:"bldr,omitempty"`
	Shell        Shell  `yaml:"shell,omitempty"`
	Version      string
	Install      Install       `yaml:"install,omitempty"`
	Dependencies []*Dependency `yaml:"dependencies,omitempty"`
	Steps        []*Step       `yaml:"steps,omitempty"`
	Finalize     []*Finalize   `yaml:"finalize,omitempty"`
	Variant      Variant       `yaml:"variant,omitempty"`
}

// NewPkg initializes a new Pkg.
func NewPkg(file string, options *Options) (*Pkg, error) {
	p := &Pkg{}

	b, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(b, p); err != nil {
		return nil, err
	}

	if p.Shell == "" {
		p.Shell = "/bin/sh"
	}

	p.Options = options

	return p, nil
}
