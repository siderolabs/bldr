package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/otiai10/copy"
	"github.com/spf13/cobra"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha1"
	"github.com/talos-systems/bldr/internal/upgrade"
	"gopkg.in/yaml.v2"
)

var convertSource string

const pkgFileShebang = "# syntax = docker.io/%s/bldr:%s-frontend\n"

// convertCmd represents the convert command
var convertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Convert v1alpha1 pkg.yaml in branches to subdirs in v1alpha2",
	Long: `This command scans a branched repository in v1alpha1 format and
creates new v1alpha2 directory structure. Output goes to --root (v1alpha2),
source is --source (v1alpha1).
`,
	Run: func(cmd *cobra.Command, args []string) {
		sourceDir, err := filepath.Abs(convertSource)
		if err != nil {
			log.Fatal(err)
		}

		destDir, err := filepath.Abs(pkgRoot)
		if err != nil {
			log.Fatal(err)
		}

		topTarget := filepath.Base(sourceDir)

		branches, err := listBranches(sourceDir, "origin")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Discovered branches: %q\n", branches)
		fmt.Printf("Top target: %q\n", topTarget)

		for _, branch := range branches {
			fmt.Printf("Processing %q...\n", branch)
			stageName := branch
			if branch == "master" {
				stageName = topTarget
			}

			if err = checkoutBranch(sourceDir, "origin", branch); err != nil {
				log.Fatal(err)
			}

			oldPkg, err := v1alpha1.NewPkg(filepath.Join(sourceDir, constants.PkgYaml), nil)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Printf("Skipping %q, no pkg.yaml found\n", stageName)
					continue
				}
				log.Fatal(err)
			}

			stageDir := filepath.Join(destDir, stageName)
			if err = os.MkdirAll(stageDir, os.ModePerm); err != nil {
				log.Fatal(err)
			}

			newPkg := upgrade.FromV1Alpha1(oldPkg, branches)

			destPkg, err := os.Create(filepath.Join(stageDir, constants.PkgYaml))
			if err != nil {
				log.Fatal(err)
			}

			enc := yaml.NewEncoder(destPkg)

			if err = enc.Encode(newPkg); err != nil {
				log.Fatal(err)
			}
			if err = enc.Close(); err != nil {
				log.Fatal(err)
			}
			if err = destPkg.Close(); err != nil {
				log.Fatal(err)
			}

			sourceDirF, err := os.Open(sourceDir)
			if err != nil {
				log.Fatal(err)
			}
			files, err := sourceDirF.Readdir(-1)
			if err != nil {
				log.Fatal(err)
			}
			sourceDirF.Close()

			for _, srcFile := range files {
				if !srcFile.IsDir() {
					continue
				}

				if strings.HasPrefix(srcFile.Name(), ".") {
					continue
				}

				fmt.Printf("Copy %s -> %s\n", filepath.Join(sourceDir, srcFile.Name()), filepath.Join(stageDir, srcFile.Name()))
				if err = copy.Copy(filepath.Join(sourceDir, srcFile.Name()), filepath.Join(stageDir, srcFile.Name())); err != nil {
					log.Fatal(err)
				}
			}
		}

		pkgFile, err := os.Create(filepath.Join(destDir, constants.Pkgfile))
		if err != nil {
			log.Fatal(err)
		}
		if _, err = pkgFile.Write([]byte(fmt.Sprintf(pkgFileShebang, constants.DefaultOrganization, constants.Version))); err != nil {
			log.Fatal(err)
		}
		if err = pkgFile.Close(); err != nil {
			log.Fatal(err)
		}
	},
}

func listBranches(repoPath, remote string) ([]string, error) {
	cmd := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/remotes/"+remote)
	cmd.Dir = repoPath

	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(output), "\n")
	branches := []string{}
	for _, line := range lines {
		if line == "" || line == remote+"/HEAD" {
			continue
		}
		branches = append(branches, strings.ReplaceAll(line, remote+"/", ""))
	}

	return branches, nil
}

func checkoutBranch(repoPath, remote, branch string) error {
	cmd := exec.Command("git", "checkout", remote+"/"+branch)
	cmd.Dir = repoPath
	return cmd.Run()
}

func init() {
	convertCmd.Flags().StringVar(&convertSource, "source", "", "Path to git repository with v1alpha1 source")
	rootCmd.AddCommand(convertCmd)
}
