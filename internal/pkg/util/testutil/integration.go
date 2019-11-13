/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"testing"

	"github.com/otiai10/copy"
	"github.com/talos-systems/bldr/internal/pkg/constants"
)

// IntegrationTest describes single integration set (common testdata).
type IntegrationTest struct {
	Name     string
	Path     string
	Manifest TestManifest
}

// Run executes integration test.
func (test IntegrationTest) Run(t *testing.T) {
	// copy test data to temp directory
	tempDir, err := ioutil.TempDir("", "bldrtest")
	if err != nil {
		t.Fatalf("error creating temp directory: %v", err)
	}

	defer func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Fatalf("error cleaning up temp directory: %v", err)
		}
	}()

	if err = copy.Copy(test.Path, tempDir); err != nil {
		t.Fatalf("error copying to temp directory: %v", err)
	}

	var oldWd string

	oldWd, err = os.Getwd()
	if err != nil {
		t.Fatalf("error getting current directory: %v", err)
	}

	err = os.Chdir(tempDir)
	if err != nil {
		t.Fatalf("error changing working directory: %v", err)
	}

	defer func() {
		err = os.Chdir(oldWd)
		if err != nil {
			t.Fatalf("error restoring working directory: %v", err)
		}
	}()

	test.run(t)
}

func (test IntegrationTest) patch(t *testing.T) {
	pkgfile, err := os.OpenFile(constants.Pkgfile, os.O_RDWR, os.ModePerm)
	if err != nil {
		t.Fatalf("error opening %q: %v", constants.Pkgfile, err)
	}

	contents, err := ioutil.ReadAll(pkgfile)
	if err != nil {
		t.Fatalf("error reading %q: %v", constants.Pkgfile, err)
	}

	contents = bytes.ReplaceAll(contents, []byte("SHEBANG"), []byte(fmt.Sprintf("%s/%s/bldr:%s-frontend", constants.DefaultRegistry, constants.DefaultOrganization, constants.Version)))

	_, err = pkgfile.Seek(0, io.SeekStart)
	if err != nil {
		t.Fatalf("error seeking %q: %v", constants.Pkgfile, err)
	}

	err = pkgfile.Truncate(0)
	if err != nil {
		t.Fatalf("error truncating %q: %v", constants.Pkgfile, err)
	}

	_, err = pkgfile.Write(contents)
	if err != nil {
		t.Fatalf("error writing %q: %v", constants.Pkgfile, err)
	}

	if err = pkgfile.Close(); err != nil {
		t.Fatalf("error closing %q: %v", constants.Pkgfile, err)
	}
}

func (test IntegrationTest) run(t *testing.T) {
	test.patch(t)

	for _, runManifest := range test.Manifest.Runs {
		runner, err := getRunner(runManifest)
		if err != nil {
			t.Fatal(err)
		}

		t.Run(runManifest.Name, runner.Run)
	}
}
