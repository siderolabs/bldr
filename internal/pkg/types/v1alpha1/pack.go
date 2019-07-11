package v1alpha1

import (
	"bytes"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"text/template"

	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/gitmeta/pkg/git"
	"github.com/talos-systems/gitmeta/pkg/metadata"
	"golang.org/x/xerrors"
)

const tpl = `
{{- $metadata := .Metadata -}}
FROM {{ .Bldr }} AS build
SHELL [ "{{ .Shell }}", "-c" ]
{{ with .Install -}}
{{ .AsDockerRun }}
{{ end -}}{{ range $dep := .Dependencies -}}
{{ $dep.AsDockerCopy $metadata }}
{{ end -}}
COPY . .
RUN /bldr build

{{ with .Finalize -}}
FROM scratch
{{ range $_, $f := . -}}
COPY --from=build {{ $f.From }} {{ $f.To }}
{{ end -}}
{{ end -}}
`

func (p *Pkg) Pack() error {
	if p.Variant == Scratch && p.Install != nil {
		return xerrors.New("scratch variant does not support installs, you may want a dependency instead")
	}

	if p.Name == "" {
		return xerrors.New("missing package name")
	}

	if p.Finalize == nil {
		return xerrors.New("missing finalizer")
	}

	for _, f := range p.Finalize {
		f.Normalize()
	}

	git, err := git.NewGit()
	if err != nil {
		return err
	}
	m, err := metadata.NewMetadata(git)
	if err != nil {
		return err
	}

	p.Metadata = m
	if p.Bldr == "" {
		p.Bldr = fmt.Sprintf("%s/%s/%s:%s", p.Options.Registry, p.Options.Organization, "bldr", constants.Version)
	}
	p.Bldr = p.Bldr + "-" + p.Variant.String()

	dockerfile, err := render(p)
	if err != nil {
		return err
	}

	tmpfile, err := os.Create(".pkg.dockerfile")
	if err != nil {
		return err
	}

	defer os.Remove(tmpfile.Name())

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	go func() {
		<-interrupt
		os.Remove(tmpfile.Name())
	}()

	if _, err := tmpfile.Write([]byte(dockerfile)); err != nil {
		return err
	}
	if err := tmpfile.Close(); err != nil {
		return err
	}

	args := []string{
		"buildx",
		"build",
		"--platform=" + p.Options.Platform,
		"--progress=" + p.Options.Progress,
		"--push=" + p.Options.Push,
		"--cache-from=" + getenv("CACHE_FROM", strings.Join([]string{p.Options.Registry, p.Options.Organization, p.Name}, "/")+":cache"),
		"--cache-to=" + getenv("CACHE_TO", strings.Join([]string{p.Options.Registry, p.Options.Organization, p.Name}, "/")+":cache"),
		"--tag=" + strings.Join([]string{p.Options.Registry, p.Options.Organization, p.Name}, "/") + ":" + p.Metadata.Container.Image.Tag,
		"--file=" + tmpfile.Name(),
		".",
	}

	return cmd("docker", args...)
}

func (i Install) AsDockerRun() string {
	return fmt.Sprintf("RUN apk add --no-cache %s", strings.Join(i, " "))
}

func (d *Dependency) AsDockerCopy(m *metadata.Metadata) string {
	if d.To == "" {
		d.To = "/"
	}

	return fmt.Sprintf("COPY --from=%s / %s", d.Image, d.To)
}

func push(m *metadata.Metadata) string {
	value := "false"
	if m.Git.IsClean && m.Git.Branch == "master" {
		value = "true"
	} else if m.Git.IsClean && m.Git.IsTag {
		value = "true"
	}

	return value
}

func getenv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func render(p *Pkg) (string, error) {
	var wr bytes.Buffer
	tmpl, err := template.New("dockerfile").Parse(tpl)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&wr, &p)
	if err != nil {
		return "", err
	}

	return wr.String(), nil
}
