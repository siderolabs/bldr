/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha1

import (
	"io/ioutil"

	"golang.org/x/xerrors"

	"gopkg.in/yaml.v2"
)

type Options struct {
	CacheTo      string
	CacheFrom    string
	Organization string
	Platform     string
	Progress     string
	Push         string
	Registry     string
}

type Install []string

type Dependency struct {
	Image string `yaml:"image,omitempty"`
	To    string `yaml:"to,omitempty"`
}

type Instruction string

type Source struct {
	URL         string `yaml:"url,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	SHA256      string `yaml:"sha256,omitempty"`
	SHA512      string `yaml:"sha512,omitempty"`
}

type Step struct {
	Prepare *Instruction `yaml:"prepare,omitempty"`
	Build   *Instruction `yaml:"build,omitempty"`
	Install *Instruction `yaml:"install,omitempty"`
	Test    *Instruction `yaml:"test,omitempty"`
	Sources []*Source    `yaml:"sources,omitempty"`
}

type Finalize struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}

type Variant int

const (
	Alpine Variant = iota
	Scratch
)

func (v Variant) String() string {
	return []string{"alpine", "scratch"}[v]
}

func (v *Variant) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string
	err := unmarshal(&aux)
	if err != nil {
		return err
	}

	var val Variant
	switch aux {
	case Alpine.String():
		val = Alpine
	case Scratch.String():
		val = Scratch
	default:
		return xerrors.Errorf("unknown variant %q")
	}
	*v = val

	return nil
}

type Shell string

type Pkg struct {
	Name         string        `yaml:"name,omitempty"`
	Bldr         string        `yaml:"bldr,omitempty"`
	Install      Install       `yaml:"install,omitempty"`
	Dependencies []*Dependency `yaml:"dependencies,omitempty"`
	Steps        []*Step       `yaml:"steps,omitempty"`
	Finalize     []*Finalize   `yaml:"finalize,omitempty"`
	Variant      Variant       `yaml:"variant,omitempty"`
	Shell        Shell         `yaml:"shell,omitempty"`

	Version string
	Options *Options
}

func NewPkg(file string, options *Options) (*Pkg, error) {
	p := &Pkg{}
	b, err := ioutil.ReadFile(file)
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
