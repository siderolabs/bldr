/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package constants

import "os"

// DefaultBaseImage for non-scratch builds
const DefaultBaseImage = "docker.io/alpine:3.10"

// DefaultDirMode is UNIX file mode for mkdir
const DefaultDirMode os.FileMode = 0755

// DefaultPath is default value for PATH environment variable
const DefaultPath = "/bin:/usr/bin:/sbin:/usr/sbin"

// PkgYaml is the filename of 'pkg.yaml'
const PkgYaml = "pkg.yaml"

// Pkgfile is the filename of 'Pkgfile'
const Pkgfile = "Pkgfile"
