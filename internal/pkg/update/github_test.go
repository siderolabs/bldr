// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package update_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/siderolabs/bldr/internal/pkg/update"
)

func TestLatestGithub(t *testing.T) {
	t.Parallel()

	if !testing.Short() {
		t.Skip("skipping in normal mode")
	}

	for source, expected := range map[string]*update.LatestInfo{
		// https://github.com/pullmoll/musl-fts/releases has only tags.
		"https://github.com/pullmoll/musl-fts/archive/refs/tags/v1.2.6.tar.gz": {
			HasUpdate: true,
			BaseURL:   "https://github.com/pullmoll/musl-fts/releases/",
			LatestURL: "https://github.com/pullmoll/musl-fts/archive/refs/tags/v1.2.7.tar.gz",
		},
		"https://github.com/pullmoll/musl-fts/archive/refs/tags/v1.2.7.tar.gz": {
			HasUpdate: false,
			BaseURL:   "https://github.com/pullmoll/musl-fts/releases/",
			LatestURL: "https://github.com/pullmoll/musl-fts/archive/refs/tags/v1.2.7.tar.gz",
		},

		// https://github.com/void-linux/musl-fts/releases has releases without extra assets.
		"https://github.com/void-linux/musl-fts/archive/refs/tags/v1.2.6.tar.gz": {
			HasUpdate: true,
			BaseURL:   "https://github.com/void-linux/musl-fts/releases/",
			LatestURL: "https://github.com/void-linux/musl-fts/archive/refs/tags/v1.2.7.tar.gz",
		},
		"https://github.com/void-linux/musl-fts/archive/refs/tags/v1.2.7.tar.gz": {
			HasUpdate: false,
			BaseURL:   "https://github.com/void-linux/musl-fts/releases/",
			LatestURL: "https://github.com/void-linux/musl-fts/archive/refs/tags/v1.2.7.tar.gz",
		},

		// https://github.com/protocolbuffers/protobuf/releases has releases with extra assets.
		"https://github.com/protocolbuffers/protobuf/releases/download/v3.15.6/protobuf-cpp-3.15.6.tar.gz": {
			HasUpdate: true,
			BaseURL:   "https://github.com/protocolbuffers/protobuf/releases/",
			LatestURL: "",
		},
		"https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protobuf-cpp-3.21.5.tar.gz": {
			HasUpdate: false,
			BaseURL:   "https://github.com/protocolbuffers/protobuf/releases/",
			LatestURL: "https://github.com/protocolbuffers/protobuf/releases/download/v21.5/protobuf-cpp-3.21.5.tar.gz",
		},

		// https://github.com/opencontainers/runc/releases has releases with extra assets and no version in the file name.
		"https://github.com/opencontainers/runc/releases/download/v1.0.0/runc.tar.xz": {
			HasUpdate: true,
			BaseURL:   "https://github.com/opencontainers/runc/releases/",
			LatestURL: "",
		},
		"https://github.com/opencontainers/runc/releases/download/v1.1.4/runc.tar.xz": {
			HasUpdate: false,
			BaseURL:   "https://github.com/opencontainers/runc/releases/",
			LatestURL: "https://github.com/opencontainers/runc/releases/download/v1.1.4/runc.tar.xz",
		},

		// https://github.com/containerd/containerd/releases has releases with extra assets that we don't use.
		"https://github.com/containerd/containerd/archive/refs/tags/v1.5.2.tar.gz": {
			HasUpdate: true,
			BaseURL:   "https://github.com/containerd/containerd/releases/",
			LatestURL: "",
		},
		"https://github.com/containerd/containerd/archive/refs/tags/v1.6.8.tar.gz": {
			HasUpdate: false,
			BaseURL:   "https://github.com/containerd/containerd/releases/",
			LatestURL: "https://github.com/containerd/containerd/archive/refs/tags/v1.6.8.tar.gz",
		},
	} {
		t.Run(source, func(t *testing.T) {
			t.Parallel()

			// check that source is actually working (with optional redirects)
			resp, err := http.Head(source) //nolint:noctx
			require.NoError(t, err)
			require.Equal(t, 200, resp.StatusCode)
			require.NoError(t, resp.Body.Close())

			actual, err := update.Latest(context.Background(), source)
			require.NoError(t, err)
			assert.Equal(t, expected, actual)
		})
	}
}
