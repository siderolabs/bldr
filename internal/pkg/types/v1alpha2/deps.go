package v1alpha2

// Dependency on another image or stage
type Dependency struct {
	Image string `yaml:"image,omitempty"`
	Stage string `yaml:"stage,omitempty"`
	To    string `yaml:"to,omitempty"`
}

// IsInternal checks whether dependency is internal to some stage
func (d *Dependency) IsInternal() bool {
	return d.Stage != ""
}

// Src returns copy source (from dependency)
func (d *Dependency) Src() string {
	return "/"
}

// Dest returns copy destination (to base)
func (d *Dependency) Dest() string {
	if d.To != "" {
		return d.To
	}

	return "/"
}

// Dependencies is a list of Depency
type Dependencies []Dependency

// GetInternal returns list of all the internal dependencies
func (deps Dependencies) GetInternal() (internalDeps []string) {
	for _, dep := range deps {
		if dep.IsInternal() {
			internalDeps = append(internalDeps, dep.Stage)
		}
	}

	return
}
