/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2

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
