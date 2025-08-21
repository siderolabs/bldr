// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"

	"github.com/siderolabs/bldr/internal/pkg/constants"
)

// TemplatedFile is a file attached to the package context which went through the templating process.
type TemplatedFile struct {
	Path    string
	Content []byte
}

// AttachTemplatedFile attaches a templated file to the package context.
func (p *Pkg) AttachTemplatedFile(path string, content []byte) error {
	tmpl, err := template.New(path).
		Funcs(sprig.HermeticTxtFuncMap()).
		Parse(string(content))
	if err != nil {
		return fmt.Errorf("failed to parse templated file %s: %w", path, err)
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, p.Context); err != nil {
		return fmt.Errorf("failed to template file %s: %w", path, err)
	}

	p.templatedFiles = append(p.templatedFiles, TemplatedFile{
		Path:    strings.TrimSuffix(path, constants.TemplateExt),
		Content: buf.Bytes(),
	})

	return nil
}

// GetTemplatedFiles return a list of templated files in the package.
func (p *Pkg) GetTemplatedFiles() []TemplatedFile {
	return p.templatedFiles
}
