package v1alpha2

type Shell string

func (sh Shell) Get() string {
	if sh == "" {
		return "/bin/sh"
	}

	return string(sh)
}
