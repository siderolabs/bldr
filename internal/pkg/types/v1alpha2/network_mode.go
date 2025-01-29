// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package v1alpha2

import "fmt"

// NetworkMode is a specification of the network mode.
//
//nolint:recvcheck
type NetworkMode int

const (
	// NetworkModeNone disables networking for the step.
	NetworkModeNone NetworkMode = iota
	// NetworkModeDefault leaves network in default (container mode).
	NetworkModeDefault
	// NetworkModeHost uses host network for the step.
	NetworkModeHost
)

func (m NetworkMode) String() string {
	return []string{"none", "default", "host"}[m]
}

// UnmarshalYAML implements yaml.Unmarshaller interface.
func (m *NetworkMode) UnmarshalYAML(unmarshal func(any) error) error {
	var aux string

	if err := unmarshal(&aux); err != nil {
		return err
	}

	var val NetworkMode

	switch aux {
	case NetworkModeNone.String():
		val = NetworkModeNone
	case NetworkModeDefault.String():
		val = NetworkModeDefault
	case NetworkModeHost.String():
		val = NetworkModeHost
	default:
		return fmt.Errorf("unknown networkmode %q", aux)
	}

	*m = val

	return nil
}

// MarshalYAML implements yaml.Marshaller interface.
func (m NetworkMode) MarshalYAML() (any, error) {
	return m.String(), nil
}
