/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

import (
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/hashicorp/go-multierror"
)

// Sources is a collection of Source.
type Sources []Source

// Validate sources.
func (sources Sources) Validate() error {
	var multiErr *multierror.Error

	for _, source := range sources {
		multiErr = multierror.Append(multiErr, source.Validate())
	}

	return multiErr.ErrorOrNil()
}

// Source describe build source to be downloaded.
type Source struct {
	URL         string `yaml:"url,omitempty"`
	Destination string `yaml:"destination,omitempty"`
	SHA256      string `yaml:"sha256,omitempty"`
	SHA512      string `yaml:"sha512,omitempty"`
}

// ToSHA512Sum returns in format of line expected by 'sha512sum'.
func (source *Source) ToSHA512Sum() []byte {
	return []byte(source.SHA512 + " *" + source.Destination + "\n")
}

// Validate source.
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

	switch len(source.SHA256) {
	case 0:
		multiErr = multierror.Append(multiErr, errors.New("source.sha256 can't be empty"))
	case 64: //nolint:gomnd
		// nothing
	default:
		multiErr = multierror.Append(multiErr, errors.New("source.sha256 should be 64 chars long"))
	}

	switch len(source.SHA512) {
	case 0:
		multiErr = multierror.Append(multiErr, errors.New("source.sha512 can't be empty"))
	case 128: //nolint:gomnd
		// nothing
	default:
		multiErr = multierror.Append(multiErr, errors.New("source.sha512 should be 128 chars long"))
	}

	return multiErr.ErrorOrNil()
}

// ValidateChecksums downloads the source, validates checksums,
// and returns actual checksums and validation error, if any.
func (source *Source) ValidateChecksums(ctx context.Context) (string, string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", source.URL, nil)
	if err != nil {
		return "", "", err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	defer resp.Body.Close() //nolint:errcheck

	s256 := sha256.New()
	s512 := sha512.New()

	if _, err = io.Copy(io.MultiWriter(s256, s512), resp.Body); err != nil {
		return "", "", err
	}

	var (
		actualSHA256, actualSHA512 string
		multiErr                   *multierror.Error
	)

	if actualSHA256 = hex.EncodeToString(s256.Sum(nil)); source.SHA256 != actualSHA256 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("source.sha256 does not match: expected %s, got %s", source.SHA256, actualSHA256))
	}

	if actualSHA512 = hex.EncodeToString(s512.Sum(nil)); source.SHA512 != actualSHA512 {
		multiErr = multierror.Append(multiErr, fmt.Errorf("source.sha512 does not match: expected %s, got %s", source.SHA512, actualSHA512))
	}

	return actualSHA256, actualSHA512, multiErr.ErrorOrNil()
}
