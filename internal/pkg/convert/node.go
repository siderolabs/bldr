/* This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at http://mozilla.org/MPL/2.0/. */

package convert

import (
	"fmt"
	"path/filepath"
	"sort"

	"github.com/moby/buildkit/client/llb"
	"github.com/opencontainers/go-digest"
	"github.com/talos-systems/bldr/internal/pkg/constants"
	"github.com/talos-systems/bldr/internal/pkg/solver"
	"github.com/talos-systems/bldr/internal/pkg/types/v1alpha2"
)

const (
	tmpDirTemplate = "/tmp/build/%d"
	pkgDir         = "/pkg"
)

var defaultCopyOptions = &llb.CopyInfo{
	CopyDirContentsOnly: true,
	CreateDestPath:      true,
	FollowSymlinks:      true,
}

// NodeLLB wraps PackageNode to provide LLB conversion
type NodeLLB struct {
	*solver.PackageNode

	Graph  *GraphLLB
	Prefix string
}

// NewNodeLLB wraps PackageNode for LLB conversion
func NewNodeLLB(node *solver.PackageNode, graph *GraphLLB) *NodeLLB {
	return &NodeLLB{
		PackageNode: node,

		Graph:  graph,
		Prefix: node.Name + ":",
	}
}

func (node *NodeLLB) base() llb.State {
	return node.Graph.BaseImages[node.Pkg.Variant]
}

func (node *NodeLLB) install(root llb.State) llb.State {
	if len(node.Pkg.Install) > 0 {
		root = root.Run(
			llb.Args(
				append([]string{"/sbin/apk", "add", "--no-cache"}, node.Pkg.Install...)),
			llb.WithCustomName(node.Prefix+"apk-install"),
		).Root()
	}

	return root
}

func (node *NodeLLB) context(root llb.State) llb.State {
	relPath := node.Pkg.BaseDir

	return root.File(
		llb.Copy(node.Graph.LocalContext, filepath.Join("/", relPath), pkgDir, defaultCopyOptions),
		llb.WithCustomNamef(node.Prefix+"context %s -> %s", relPath, pkgDir),
	)
}

func (node *NodeLLB) dependencies(root llb.State) (llb.State, error) {
	for _, dep := range node.Pkg.Dependencies {
		var (
			depState llb.State
			srcName  string
		)

		if dep.IsInternal() {
			for _, depNode := range node.DependsOn {
				if depNode.Name != dep.Stage {
					continue
				}

				var err error

				depState, err = NewNodeLLB(depNode, node.Graph).Build()
				if err != nil {
					return llb.Scratch(), err
				}

				srcName = depNode.Name

				break
			}
		} else {
			depState = llb.Image(dep.Image)
			srcName = dep.Image
		}

		root = root.File(
			llb.Copy(depState, dep.Src(), dep.Dest(), defaultCopyOptions),
			llb.WithCustomNamef(node.Prefix+"copy --from %s %s -> %s", srcName, dep.Src(), dep.Dest()))
	}

	return root, nil
}

func (node *NodeLLB) stepTmpDir(root llb.State, i int, step *v1alpha2.Step) llb.State {
	if step.TmpDir == "" {
		step.TmpDir = fmt.Sprintf(tmpDirTemplate, i)
	}

	return root.File(
		llb.Mkdir(step.TmpDir, constants.DefaultDirMode, llb.WithParents(true)),
		llb.WithCustomName(node.Prefix+"mkdir "+step.TmpDir),
	).Dir(step.TmpDir)
}

func (node *NodeLLB) stepDownload(root llb.State, step v1alpha2.Step) llb.State {
	for _, source := range step.Sources {
		download := llb.HTTP(
			source.URL,
			llb.Filename(filepath.Join("/", source.Destination)),
			llb.Checksum(digest.NewDigestFromEncoded(digest.SHA256, source.SHA256)),
			llb.WithCustomNamef(node.Prefix+"download %s -> %s", source.URL, source.Destination),
		)

		checksummer := node.Graph.Checksummer.File(
			llb.Mkfile("/checksums", 0644, source.ToSHA512Sum()).
				Copy(download, "/", "/", defaultCopyOptions).
				Mkdir("/empty", constants.DefaultDirMode),
			llb.WithCustomName(node.Prefix+"cksum-prepare"),
		).Run(
			llb.Shlex("sha512sum -c --strict /checksums"),
			llb.WithCustomName(node.Prefix+"cksum-verify"),
		).Root()

		root = root.File(
			llb.Copy(download, "/", step.TmpDir, defaultCopyOptions).
				Copy(checksummer, "/empty", "/", defaultCopyOptions), // TODO: this is "fake" dependency on checksummer
			llb.WithCustomName(node.Prefix+"download finalize"),
		)
	}

	return root
}

func (node *NodeLLB) stepEnvironment(root llb.State, step v1alpha2.Step) llb.State {
	vars := step.Env
	keys := make([]string, 0, len(vars))

	for key := range vars {
		keys = append(keys, key)
	}

	sort.Strings(keys)

	for _, key := range keys {
		root = root.AddEnv(key, vars[key])
	}

	return root
}

func (node *NodeLLB) stepScripts(root llb.State, i int, step v1alpha2.Step) llb.State {
	for _, script := range []struct {
		Desc         string
		Instructions v1alpha2.Instructions
	}{
		{"prepare", step.Prepare},
		{"build", step.Build},
		{"install", step.Install},
	} {
		for _, instruction := range script.Instructions {
			root = root.Run(
				llb.Args([]string{
					node.Pkg.Shell.Get(),
					"-c",
					instruction.Script(),
				}),
				llb.WithCustomName(fmt.Sprintf("%s%s-%d", node.Prefix, script.Desc, i)),
			).Root()
		}
	}

	return root
}

func (node *NodeLLB) step(root llb.State, i int, step v1alpha2.Step) llb.State {
	root = node.stepTmpDir(root, i, &step)
	root = node.stepDownload(root, step)
	root = node.stepEnvironment(root, step)
	root = node.stepScripts(root, i, step)

	return root
}

func (node *NodeLLB) finalize(root llb.State) llb.State {
	newroot := llb.Scratch()

	for _, fin := range node.Pkg.Finalize {
		newroot = newroot.File(
			llb.Copy(root, fin.From, fin.To, defaultCopyOptions),
			llb.WithCustomNamef(node.Prefix+"finalize %s -> %s", fin.From, fin.To),
		)
	}

	return newroot
}

// Build converts PackageNode to buildkit LLB
func (node *NodeLLB) Build() (llb.State, error) {
	var err error

	root := node.base()
	root = node.install(root)
	root = node.context(root)

	root, err = node.dependencies(root)
	if err != nil {
		return llb.Scratch(), err
	}

	for i, step := range node.Pkg.Steps {
		root = node.step(root, i, step)
	}

	root = node.finalize(root)

	return root, nil
}
