/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package v1alpha2_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

//nolint:lll
func TestSourceValidateChecksums(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping in -short mode")
	}

	const (
		expectedSHA256 = "2aa5f088cbb332e73fc3def546800616b38d3bfe6b8713b8a6404060f22503e8"
		expectedSHA512 = "ce64105ff71615f9d235cc7c8656b6409fc40cc90d15a28d355fadd9072d2eab842af379dd8bba0f1181715753143e4a07491e0f9e5f8df806327d7c95a34fae"
	)

	source := v1alpha2.Source{
		URL:         "https://dl.google.com/go/go1.12.5.src.tar.gz",
		Destination: "go1.12.5.src.tar.gz",
		SHA256:      expectedSHA256,
		SHA512:      expectedSHA512,
	}

	actualSHA256, actualSHA512, err := source.ValidateChecksums(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedSHA256, actualSHA256)
	assert.Equal(t, expectedSHA512, actualSHA512)

	source.SHA256 = strings.Repeat("0", 64)
	source.SHA512 = strings.Repeat("1", 64)

	actualSHA256, actualSHA512, err = source.ValidateChecksums(context.Background())
	assert.EqualError(t, err, `2 errors occurred:
	* go1.12.5.src.tar.gz sha256 does not match: expected 0000000000000000000000000000000000000000000000000000000000000000, got 2aa5f088cbb332e73fc3def546800616b38d3bfe6b8713b8a6404060f22503e8
	* go1.12.5.src.tar.gz sha512 does not match: expected 1111111111111111111111111111111111111111111111111111111111111111, got ce64105ff71615f9d235cc7c8656b6409fc40cc90d15a28d355fadd9072d2eab842af379dd8bba0f1181715753143e4a07491e0f9e5f8df806327d7c95a34fae

`)
	assert.Equal(t, expectedSHA256, actualSHA256)
	assert.Equal(t, expectedSHA512, actualSHA512)
}
