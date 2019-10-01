package types

// Variables presents generic variables for templating/environment
type Variables map[string]string

// Merge two Variables
func (v Variables) Merge(other Variables) Variables {
	for key := range other {
		v[key] = other[key]
	}

	return v
}
