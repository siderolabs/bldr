// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package update //nolint:testpackage // to test unexported function

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExtractVersion(t *testing.T) {
	t.Parallel()

	for s, expected := range map[string]string{
		"https://dl.google.com/go/go1.15.7.src.tar.gz":    "1.15.7",
		"https://ftp.pcre.org/pub/pcre/pcre-8.43.tar.bz2": "8.43.0",

		"https://ftp.gnu.org/gnu/automake/automake-1.16.tar.gz":   "1.16.0",
		"https://ftp.gnu.org/gnu/automake/automake-1.16.1.tar.xz": "1.16.1",

		"https://github.com/pullmoll/void-linux/archive/refs/tags/v1.2.7.tar.gz":                           "1.2.7",
		"https://github.com/protocolbuffers/protobuf/releases/download/v3.17.3/protobuf-cpp-3.17.3.tar.gz": "3.17.3",
		"https://github.com/opencontainers/runc/releases/download/v1.0.1/runc.tar.xz":                      "1.0.1",
	} {
		t.Run(s, func(t *testing.T) {
			t.Parallel()

			actualV, err := extractVersion(s)
			require.NoError(t, err)
			require.NotNil(t, actualV)
			assert.Equal(t, expected, actualV.String())
		})
	}
}
