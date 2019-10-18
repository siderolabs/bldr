/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

import (
	"errors"
	"fmt"
	"net/url"

	"github.com/hashicorp/go-multierror"
)

// Sources is a collection of Source
type Sources []Source

// Validate sources
func (sources Sources) Validate() error {
	var multiErr *multierror.Error

	for _, source := range sources {
		multiErr = multierror.Append(multiErr, source.Validate())
	}

	return multiErr.ErrorOrNil()
}

// Source describe build source to be downloaded
type Source struct {
	URL         string `yaml:"url,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	SHA256      string `yaml:"sha256,omitempty"`
	SHA512      string `yaml:"sha512,omitempty"`
}

// ToSHA512Sum returns in format of line expected by 'sha512sum'
func (source *Source) ToSHA512Sum() []byte {
	return []byte(source.SHA512 + " *" + source.Destination + "\n")
}

// Validate source
func (source *Source) Validate() error {
	var multiErr *multierror.Error

	if source.URL == "" {
		multiErr = multierror.Append(multiErr, errors.New("source.url can't be empty"))
	} else if _, err := url.Parse(source.URL); err != nil {
		multiErr = multierror.Append(multiErr, fmt.Errorf("error parsing source.url %q: %w", source.URL, err))
	}

	if source.Destination == "" {
		multiErr = multierror.Append(multiErr, errors.New("source.destination can't be empty"))
	}

	if source.SHA256 == "" {
		multiErr = multierror.Append(multiErr, errors.New("source.sha256 can't be empty"))
	}

	if source.SHA512 == "" {
		multiErr = multierror.Append(multiErr, errors.New("source.sha512 can't be empty"))
	}

	return multiErr.ErrorOrNil()
}
