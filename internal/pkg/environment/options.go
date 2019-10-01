package environment

import "github.com/talos-systems/bldr/internal/pkg/types"

// Options for bldr
type Options struct {
	BuildPlatform  Platform
	TargetPlatform Platform
	Target         string
}

// GetVariables returns set of variables set for options
func (options *Options) GetVariables() types.Variables {
	return Default().
		Merge(options.BuildPlatform.BuildVariables()).
		Merge(options.TargetPlatform.TargetVariables())
}
