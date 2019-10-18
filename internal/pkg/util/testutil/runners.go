/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// Run the integration test
type Run interface {
	Run(t *testing.T)
}

// CommandRunner is an abstract runner mix-in which processes command result.
type CommandRunner struct {
	Expect string
}

func (runner CommandRunner) run(t *testing.T, cmd *exec.Cmd, title string) {
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	err := cmd.Run()

	switch runner.Expect {
	case "success":
		if err != nil {
			t.Fatalf("%s failed: %v", title, err)
		}
	case "fail":
		if err != nil {
			t.Fatalf("%s should have failed, but succeeded", title)
		}
	default:
		t.Fatalf("unsupported expect %q", runner.Expect)
	}
}

func getRunner(manifest RunManifest) (Run, error) {
	switch manifest.Runner {
	case "docker":
		return DockerRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target: manifest.Target,
		}, nil
	case "buildkit":
		return BuildkitRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target: manifest.Target,
		}, nil
	case "llb":
		return LLBRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target: manifest.Target,
		}, nil
	case "validate":
		return ValidateRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
		}, nil
	default:
		return nil, fmt.Errorf("unsupported runner: %q", manifest.Runner)
	}
}
