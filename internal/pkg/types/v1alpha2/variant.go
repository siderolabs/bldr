// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import "fmt"

// Variant is a kind of base build image.
//
//nolint:recvcheck
type Variant int

const (
	// Alpine variant uses Alpine as base image for the build.
	Alpine Variant = iota
	// Scratch variant uses scratch image as base image for the build.
	Scratch
)

func (v Variant) String() string {
	return []string{"alpine", "scratch"}[v]
}

// UnmarshalYAML implements yaml.Unmarshaller interface.
func (v *Variant) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var aux string

	if err := unmarshal(&aux); err != nil {
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

// MarshalYAML implements yaml.Marshaller interface.
func (v Variant) MarshalYAML() (interface{}, error) {
	return v.String(), nil
}
