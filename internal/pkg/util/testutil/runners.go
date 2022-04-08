/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package testutil

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"testing"
)

// Run the integration test.
type Run interface {
	Run(t *testing.T)
}

// CommandRunner is an abstract runner mix-in which processes command result.
type CommandRunner struct {
	Expect       string
	ExpectStdout *string
}

func (runner CommandRunner) run(t *testing.T, cmd *exec.Cmd, title string) {
	var stdout bytes.Buffer

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout

	if runner.ExpectStdout != nil {
		cmd.Stdout = &stdout
	}

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

	if runner.ExpectStdout != nil {
		if *runner.ExpectStdout != stdout.String() {
			t.Fatalf("%s stdout mismatch: %q != %q", title, *runner.ExpectStdout, stdout.String())
		}
	}
}

func getRunner(manifest RunManifest) (Run, error) {
	switch manifest.Runner {
	case "docker":
		return DockerRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target:   manifest.Target,
			Platform: manifest.Platform,
		}, nil
	case "buildkit":
		return BuildkitRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target:   manifest.Target,
			Platform: manifest.Platform,
		}, nil
	case "eval":
		return EvalRunner{
			CommandRunner: CommandRunner{
				Expect:       manifest.Expect,
				ExpectStdout: manifest.ExpectStdout,
			},
			Target:   manifest.Target,
			Template: manifest.Template,
		}, nil
	case "llb":
		return LLBRunner{
			CommandRunner: CommandRunner{
				Expect: manifest.Expect,
			},
			Target:   manifest.Target,
			Platform: manifest.Platform,
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
