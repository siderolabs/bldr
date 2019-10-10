/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

import (
	"bytes"
	"text/template"

	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types"
	"gopkg.in/yaml.v2"
)

// Pkg represents build instructions for a single package
type Pkg struct {
	Name         string       `yaml:"name,omitempty"`
	Variant      Variant      `yaml:"variant,omitempty"`
	Shell        Shell        `yaml:"shell,omitempty"`
	Install      Install      `yaml:"install,omitempty"`
	Dependencies Dependencies `yaml:"dependencies,omitempty"`
	Steps        []Step       `yaml:"steps,omitempty"`
	Finalize     []Finalize   `yaml:"finalize,omitempty"`

	BaseDir string `yaml:"-"`
}

// NewPkg loads Pkg structure from file
func NewPkg(baseDir string, contents []byte, vars types.Variables) (*Pkg, error) {
	p := &Pkg{
		BaseDir: baseDir,
		Shell:   "/bin/sh",
		Variant: Alpine,
	}

	tmpl, err := template.New(constants.PkgYaml).Parse(string(contents))
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, vars); err != nil {
		return nil, err
	}

	if err := yaml.NewDecoder(&buf).Decode(p); err != nil {
		return nil, err
	}

	return p, nil
}

// InternalDependencies returns list of internal dependencies names
func (p *Pkg) InternalDependencies() []string {
	return p.Dependencies.GetInternal()
}

// ExternalDependencies returns list of external images
func (p *Pkg) ExternalDependencies() []string {
	return p.Dependencies.GetExternal()
}
