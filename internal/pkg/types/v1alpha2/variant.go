package v1alpha2

import "fmt"

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
		return fmt.Errorf("unknown variant %q", aux)
	}
	*v = val

	return nil
}

func (v Variant) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}
