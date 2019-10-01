package v1alpha2

type Install []string

type Environment map[string]string

type Step struct {
	Sources []Source     `yaml:"sources,omitempty"`
	Env     Environment  `yaml:"env,omitempty"`
	Prepare Instructions `yaml:"prepare,omitempty"`
	Build   Instructions `yaml:"build,omitempty"`
	Install Instructions `yaml:"install,omitempty"`
	Test    Instructions `yaml:"test,omitempty"`

	TmpDir string `yaml:"-"`
}

type Finalize struct {
	From string `yaml:"from,omitempty"`
	To   string `yaml:"to,omitempty"`
}
