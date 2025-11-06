// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package cmd

import (
	"bufio"
	"bytes"
	"context"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"slices"
	"strings"

	"github.com/cheggaaa/pb/v3"
	"github.com/spf13/cobra"

	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

// updateCmd represents the `update` command.
var updateCmd = &cobra.Command{
	Use:   "update <file> <variable>",
	Short: "Update checksums for a given version variable change.",
	Long: `This command updates checksums of all sources that reference given variable
in their URL.

Example:

  bldr update Pkgfile linux_version

Another way to call this command is to pass the output of "git diff -U0" to the standard input,
then this command will try to bump every changed variable in the diff:

  git diff -U0 | bldr update
`,
	Run: func(_ *cobra.Command, args []string) {
		ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
		defer cancel()

		switch len(args) {
		case 0:
			diffUpdater(ctx)
		case 2:
			varsPath, varName := args[0], args[1]

			singleVariableUpdater(ctx, varsPath, varName)
		default:
			log.Fatalf("expected 0 or 2 arguments, got %d", len(args))
		}
	},
}

func diffUpdater(ctx context.Context) {
	fmt.Fprintf(os.Stderr, "reading git diff from stdin\n")

	scanner := bufio.NewScanner(os.Stdin)

	var currentPath string

	for scanner.Scan() {
		line := scanner.Text()

		switch {
		case strings.HasPrefix(line, "+++ b/"):
			currentPath = strings.TrimPrefix(line, "+++ b/")
		case strings.HasPrefix(line, "+++ w/"):
			currentPath = strings.TrimPrefix(line, "+++ w/")
		case strings.HasPrefix(line, "+"):
			variableName, _, ok := strings.Cut(strings.TrimPrefix(line, "+"), ": ")
			if !ok {
				continue
			}

			variableName = strings.TrimSpace(variableName)

			if variableName == "" || strings.Contains(strings.ToLower(variableName), "sha256") || strings.Contains(strings.ToLower(variableName), "sha512") {
				continue
			}

			fmt.Fprintf(os.Stderr, "updating %s in %s\n", variableName, currentPath)

			singleVariableUpdater(ctx, currentPath, variableName)
		}
	}
}

//nolint:gocognit
func singleVariableUpdater(ctx context.Context, varsPath, varName string) {
	const phase1Repl = "XXXXXXXXXXXXX^^^^YYYYYYYYYYYYYY"

	var (
		affectedSources []v1alpha2.Source
		sourceVars      = map[string][]byte{}
	)

	for _, targetPlatform := range []environment.Platform{
		environment.LinuxAmd64,
		environment.LinuxArm64,
	} {
		options.TargetPlatform = targetPlatform

		loaderHooked := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
			HookOnVariables: func(path string, vars types.Variables) {
				if path == varsPath {
					if _, ok := vars[varName]; ok {
						vars[varName] = phase1Repl
					}
				}
			},
			HookOnLoad: func(path string, contents []byte) {
				sourceVars[path] = contents
			},
		}

		loaderClean := solver.FilesystemPackageLoader{
			Root:    pkgRoot,
			Context: options.GetVariables(),
		}

		packagesHooked, err := solver.NewPackages(&loaderHooked)
		if err != nil {
			log.Fatal(err)
		}

		packagesClean, err := solver.NewPackages(&loaderClean)
		if err != nil {
			log.Fatal(err)
		}

		cleanSorted := packagesClean.ToSet().Sorted()

		for pkgIdx, pkg := range packagesHooked.ToSet().Sorted() {
			for stepIdx, step := range pkg.Pkg.Steps {
				for sourceIdx, src := range step.Sources {
					if strings.Contains(src.URL, phase1Repl) {
						origSrc := cleanSorted[pkgIdx].Pkg.Steps[stepIdx].Sources[sourceIdx]

						if !slices.Contains(affectedSources, origSrc) {
							affectedSources = append(affectedSources, origSrc)
						}
					}
				}
			}
		}
	}

	for _, src := range affectedSources {
		fmt.Fprintf(os.Stderr, "affected source: %s (%s/%s)\n", src.URL, src.SHA256, src.SHA512)
	}

	newChecksums := make(map[v1alpha2.Source]v1alpha2.Source, len(affectedSources))

	for _, src := range affectedSources {
		var err error

		newChecksums[src], err = downloadAndChecksum(ctx, src.URL)
		if err != nil {
			log.Fatalf("error processing %q: %s", src.URL, err)
		}
	}

	newContents := map[string][]byte{}

	for oldSrc, newSrc := range newChecksums {
		fmt.Printf("updating %s, sha256 %s -> %s\n", oldSrc.URL, oldSrc.SHA256, newSrc.SHA256)

		for path, contents := range sourceVars {
			origContent, ok := newContents[path]
			if !ok {
				origContent = contents
			}

			replacedContents := bytes.ReplaceAll(origContent, []byte(oldSrc.SHA256), []byte(newSrc.SHA256))
			replacedContents = bytes.ReplaceAll(replacedContents, []byte(oldSrc.SHA512), []byte(newSrc.SHA512))

			if !bytes.Equal(origContent, replacedContents) {
				fmt.Fprintf(os.Stderr, "updating for %s in %s\n", oldSrc.URL, path)

				newContents[path] = replacedContents
			}
		}
	}

	for path, contents := range newContents {
		err := os.WriteFile(path, contents, 0o644)
		if err != nil {
			log.Fatalf("error writing %q: %s", path, err)
		}
	}
}

func downloadAndChecksum(ctx context.Context, url string) (v1alpha2.Source, error) {
	fmt.Fprintf(os.Stderr, "downloading %s\n", url)

	var result v1alpha2.Source

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return result, fmt.Errorf("error creating request for %q: %w", url, err)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return result, fmt.Errorf("error downloading %q: %w", url, err)
	}
	defer resp.Body.Close() //nolint:errcheck

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return result, fmt.Errorf("error downloading %q: status code %d", url, resp.StatusCode)
	}

	bar := pb.Full.Start64(resp.ContentLength)
	defer bar.Finish()

	s256 := sha256.New()
	s512 := sha512.New()

	_, err = io.Copy(io.MultiWriter(s256, s512), bar.NewProxyReader(resp.Body))
	if err != nil {
		return result, fmt.Errorf("error reading %q: %w", url, err)
	}

	result.URL = url
	result.SHA256 = hex.EncodeToString(s256.Sum(nil))
	result.SHA512 = hex.EncodeToString(s512.Sum(nil))

	return result, nil
}

func init() {
	rootCmd.AddCommand(updateCmd)
}
