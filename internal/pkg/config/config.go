package config

import (
	"io/ioutil"
	"path/filepath"

	"github.com/mitchellh/go-homedir"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Repository string `yaml:"repository,omitempty"`
	Username   string `yaml:"username,omitempty"`
}

func Open() (*Config, error) {
	c := &Config{}
	home, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	expanded, err := homedir.Expand(home)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadFile(filepath.Join(expanded, ".bldr.yaml"))
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(b, c); err != nil {
		return nil, err
	}

	return c, nil
}
