/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

import (
	"bytes"
	"errors"
	"text/template"

	"gopkg.in/yaml.v2"

	"github.com/Masterminds/sprig"
	"github.com/hashicorp/go-multierror"

	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types"
)

// Pkg represents build instructions for a single package
type Pkg struct {
	Name         string       `yaml:"name,omitempty"`
	Variant      Variant      `yaml:"variant,omitempty"`
	Shell        Shell        `yaml:"shell,omitempty"`
	Install      Install      `yaml:"install,omitempty"`
	Dependencies Dependencies `yaml:"dependencies,omitempty"`
	Steps        Steps        `yaml:"steps,omitempty"`
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

	tmpl, err := template.New(constants.PkgYaml).
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(string(contents))
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

	if err := p.Validate(); err != nil {
		return nil, err
	}

	return p, nil
}

// Validate the Pkg
func (p *Pkg) Validate() error {
	var multiErr *multierror.Error

	if p.Name == "" {
		multiErr = multierror.Append(multiErr, errors.New("package name can't be empty"))
	}

	if len(p.Steps) > 0 && len(p.Finalize) == 0 {
		multiErr = multierror.Append(multiErr, errors.New("finalize steps are missing, this is going to lead to empty build"))
	}

	multiErr = multierror.Append(multiErr, p.Steps.Validate(), p.Dependencies.Validate())

	return multiErr.ErrorOrNil()
}
