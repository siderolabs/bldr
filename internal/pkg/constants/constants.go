// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package constants provides basic constants for the program
package constants

// Set of variables set during the build.
var (
	DefaultRegistry     string
	DefaultOrganization string
	Version             string

	// DefaultBaseImage for non-scratch builds.
	// renovate: datasource=docker versioning=docker depName=alpine
	DefaultBaseImage = "docker.io/alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c"

	// StageXBusyboxImage is the image name for busybox from stageX.
	// renovate: datasource=docker versioning=docker depName=siderolabs/stagex/core-busybox
	StageXBusyboxImage = "ghcr.io/siderolabs/stagex/core-busybox:1.36.1@sha256:c0b551b47d8f1ac2fd5f4712eafddb8717e6e563a47203e02f94f944f64c18b2"
)
