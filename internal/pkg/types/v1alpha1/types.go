package v1alpha1

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/hashicorp/go-multierror"
	"github.com/talos-systems/gitmeta/pkg/metadata"
	"golang.org/x/xerrors"

	"gopkg.in/yaml.v2"
)

type Options struct {
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

	Version  string
	Metadata *metadata.Metadata
	Options  *Options
}

func NewPkg(file string, options *Options) (*Pkg, error) {
	if err := os.Chdir(filepath.Dir(file)); err != nil {
		return nil, err
	}

	p := &Pkg{}
	b, err := ioutil.ReadFile(filepath.Base(file))
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

func (s *Step) download(dir string) error {
	if s.Sources == nil {
		return nil
	}

	var (
		result *multierror.Error
		wg     sync.WaitGroup
	)

	wg.Add(len(s.Sources))

	for _, src := range s.Sources {
		go func(src *Source) {
			log.Printf("Downloading %q", src.URL)
			defer wg.Done()
			if err := os.MkdirAll(filepath.Dir(src.Destination), 0700); err != nil {
				result = multierror.Append(result, err)
				return
			}

			dest := filepath.Join(dir, src.Destination)

			if err := download(src.URL, dest); err != nil {
				result = multierror.Append(result, err)
				return
			}

			log.Printf("Verifying %q", src.URL)
			if err := verify(dest, src.SHA256, src.SHA512); err != nil {
				result = multierror.Append(result, err)
				return
			}
		}(src)
	}

	wg.Wait()

	return result.ErrorOrNil()
}

func (s *Step) prepare(shell Shell) error {
	return s.Prepare.exec(shell)
}

func (s *Step) build(shell Shell) error {
	return s.Build.exec(shell)
}

func (s *Step) install(shell Shell) error {
	return s.Install.exec(shell)
}

func (s *Step) test(shell Shell) error {
	return s.Test.exec(shell)
}

func (i *Instruction) String(shell Shell) string {
	if i == nil {
		return ""
	}

	shebang := "#!" + string(shell)
	options := "set -eou pipefail"
	script := strings.Join([]string{shebang, options, string(*i)}, "\n")

	return script
}

func (i *Instruction) exec(shell Shell) error {
	if i == nil {
		return nil
	}

	return cmd(string(shell), "-c", i.String(shell))
}

func (f *Finalize) Normalize() {
	if f.From == "" {
		f.From = "/"
	}

	if f.To == "" {
		f.To = "/"
	}
}

func cmd(c string, args ...string) error {
	cmd := exec.Command(c, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return err
	}
	return cmd.Wait()
}
