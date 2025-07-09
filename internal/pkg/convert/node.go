// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package convert

import (
	"context"
	"fmt"
	"path"
	"path/filepath"
	"sort"

	"github.com/moby/buildkit/client/llb"
	"github.com/moby/buildkit/frontend/gateway/client"
	"github.com/moby/buildkit/solver/pb"
	"github.com/opencontainers/go-digest"
	"github.com/siderolabs/gen/xslices"

	"github.com/siderolabs/bldr/internal/pkg/constants"
	"github.com/siderolabs/bldr/internal/pkg/environment"
	"github.com/siderolabs/bldr/internal/pkg/sbom"
	"github.com/siderolabs/bldr/internal/pkg/solver"
	"github.com/siderolabs/bldr/internal/pkg/types/v1alpha2"
)

const (
	tmpDir = "/tmp/build"
	pkgDir = "/pkg"
)

func defaultCopyOptions(options *environment.Options, reproducible bool) *llb.CopyInfo {
	copyOptions := &llb.CopyInfo{
		CopyDirContentsOnly: true,
		CreateDestPath:      true,
		FollowSymlinks:      true,
	}

	if reproducible {
		copyOptions.ChownOpt = &llb.ChownOpt{
			User: &llb.UserOpt{
				UID: 0,
			},
			Group: &llb.UserOpt{
				UID: 0,
			},
		}

		if !options.SourceDateEpoch.IsZero() {
			copyOptions.CreatedTime = &options.SourceDateEpoch
		}
	}

	return copyOptions
}

// NodeLLB wraps PackageNode to provide LLB conversion.
type NodeLLB struct {
	*solver.PackageNode

	Graph  *GraphLLB
	Prefix string
}

// NewNodeLLB wraps PackageNode for LLB conversion.
func NewNodeLLB(node *solver.PackageNode, graph *GraphLLB) *NodeLLB {
	return &NodeLLB{
		PackageNode: node,

		Graph:  graph,
		Prefix: graph.Options.CommonPrefix + node.Name + ":",
	}
}

func (node *NodeLLB) base() llb.State {
	return node.Graph.BaseImages[node.Pkg.Variant]
}

func (node *NodeLLB) install(root llb.State) llb.State {
	if len(node.Pkg.Install) > 0 {
		root = root.Run(
			append(node.Graph.commonRunOptions,
				llb.Args(
					append([]string{"/sbin/apk", "add", "--no-cache"}, node.Pkg.Install...)),
				llb.WithCustomName(node.Prefix+"apk-install"),
			)...,
		).Root()
	}

	return root
}

func (node *NodeLLB) context(root llb.State) llb.State {
	relPath := node.Pkg.BaseDir

	return root.File(
		llb.Copy(node.Graph.LocalContext, filepath.Join("/", relPath), pkgDir, defaultCopyOptions(node.Graph.Options, false)),
		llb.WithCustomNamef(node.Prefix+"context %s -> %s", relPath, pkgDir),
	)
}

func (node *NodeLLB) convertDependency(ctx context.Context, dep solver.PackageDependency) (depState llb.State, srcName string, err error) {
	if dep.IsInternal() {
		if dep.Platform != "" && dep.Platform != node.Graph.Options.BuildPlatform.ID {
			var res *client.Result

			res, err = node.Graph.solverFn(ctx, environment.Platforms[dep.Platform], dep.Node.Name)
			if err != nil {
				return llb.Scratch(), "", err
			}

			var ref client.Reference

			ref, err = res.SingleRef()
			if err != nil {
				return llb.Scratch(), "", err
			}

			depState, err = ref.ToState()
			if err != nil {
				return llb.Scratch(), "", err
			}
		} else {
			depState, err = NewNodeLLB(dep.Node, node.Graph).Build(ctx)
			if err != nil {
				return llb.Scratch(), "", err
			}
		}

		srcName = dep.Node.Name
	} else {
		depState = llb.Image(dep.Image)
		srcName = dep.Image

		if dep.Platform != "" {
			platform, ok := environment.Platforms[dep.Platform]
			if !ok {
				return llb.Scratch(), "", fmt.Errorf("platform %q not supported", dep.Platform)
			}

			depState = llb.Image(dep.Image, platform.LLBPlatform)
		}
	}

	return depState, srcName, nil
}

func (node *NodeLLB) dependencies(ctx context.Context, root llb.State) (llb.State, error) {
	deps := make([]solver.PackageDependency, 0, len(node.Dependencies))

	// collect all the dependencies including transitive runtime dependencies
	// into a list, and then build LLB deduplicating dependencies on the fly

	// order is preserved in general with runtime dependencies following direct dependency,
	// but due to deduplication all the duplicates are removed (only first appearance
	// stays in the list)

	for _, dep := range node.Dependencies {
		deps = append(deps, dep)
		if dep.Node != nil {
			deps = append(deps, dep.Node.RuntimeDependencies()...)
		}
	}

	seen := map[string]struct{}{}

	stages := []llb.State{root}

	for _, dep := range deps {
		if _, alreadyProcessed := seen[dep.ID()]; alreadyProcessed {
			continue
		}

		seen[dep.ID()] = struct{}{}

		depState, srcName, err := node.convertDependency(ctx, dep)
		if err != nil {
			return llb.Scratch(), err
		}

		if dep.Src() == "/" && dep.Dest() == "/" {
			// skip copying if the source and destination are "/"
			stages = append(stages, depState)
		} else {
			stages = append(stages,
				llb.Scratch().File(
					llb.Copy(depState, dep.Src(), dep.Dest(), defaultCopyOptions(node.Graph.Options, false)),
					llb.WithCustomNamef(node.Prefix+"copy --from %s %s -> %s", srcName, dep.Src(), dep.Dest()),
				),
			)
		}
	}

	if len(stages) == 1 {
		return root, nil
	}

	return root.WithOutput(llb.Merge(stages, llb.WithCustomName(node.Prefix+"copy")).Output()), nil
}

func (node *NodeLLB) stepTmpDir(root llb.State, step *v1alpha2.Step) llb.State {
	if step.TmpDir == "" {
		step.TmpDir = tmpDir
	}

	return root.File(
		llb.Mkdir(step.TmpDir, constants.DefaultDirMode, llb.WithParents(true)),
		llb.WithCustomName(node.Prefix+"mkdir "+step.TmpDir),
	).Dir(step.TmpDir)
}

func (node *NodeLLB) stepDownload(root llb.State, step v1alpha2.Step) llb.State {
	if len(step.Sources) == 0 {
		return root
	}

	stages := []llb.State{root}

	for _, source := range step.Sources {
		download := llb.HTTP(
			source.URL,
			llb.Filename(filepath.Join("/", source.Destination)),
			llb.Checksum(digest.NewDigestFromEncoded(digest.SHA256, source.SHA256)),
			llb.WithCustomNamef(node.Prefix+"download %s -> %s", source.URL, source.Destination),
		)

		checksummer := node.Graph.Checksummer.File(
			llb.Mkfile("/checksums", 0o644, source.ToSHA512Sum()).
				Copy(download, "/", "/", defaultCopyOptions(node.Graph.Options, false)).
				Mkdir("/empty", constants.DefaultDirMode),
			llb.WithCustomName(node.Prefix+"cksum-prepare"),
		).Run(
			append(node.Graph.commonRunOptions,
				llb.Shlex("sha512sum -c -w /checksums"),
				llb.WithCustomName(node.Prefix+"cksum-verify"),
				llb.Network(pb.NetMode_NONE),
			)...,
		).Root()

		stages = append(stages,
			llb.Scratch().File(
				llb.Copy(download, "/", step.TmpDir, defaultCopyOptions(node.Graph.Options, false)).
					Copy(checksummer, "/empty", "/", defaultCopyOptions(node.Graph.Options, false)), // TODO: this is "fake" dependency on checksummer
				llb.WithCustomName(node.Prefix+"download finalize"),
			),
		)
	}

	return root.WithOutput(llb.Merge(stages, llb.WithCustomName(node.Prefix+"download")).Output())
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
		// Detached script modifications to the files are  not propagated to the next steps.
		Detached bool
	}{
		{"prepare", step.Prepare, false},
		{"build", step.Build, false},
		{"install", step.Install, false},
		{"test", step.Test, true},
	} {
		if len(script.Instructions) == 0 {
			continue
		}

		scriptRoot := root

		for _, instruction := range script.Instructions {
			runOptions := append([]llb.RunOption(nil), node.Graph.commonRunOptions...)

			switch step.Network {
			case v1alpha2.NetworkModeNone:
				runOptions = append(runOptions, llb.Network(pb.NetMode_NONE))
			case v1alpha2.NetworkModeHost:
				runOptions = append(runOptions, llb.Network(pb.NetMode_HOST))
			case v1alpha2.NetworkModeDefault: // do nothing
			}

			runOptions = append(runOptions, xslices.Map(step.CachePaths, func(p string) llb.RunOption {
				return llb.AddMount(
					p,
					llb.Scratch(),
					llb.AsPersistentCacheDir(
						path.Clean(node.Graph.Options.CacheIDNamespace+"/"+p),
						llb.CacheMountShared,
					),
				)
			})...)

			runOptions = append(runOptions,
				llb.Args([]string{
					node.Pkg.Shell.Get(),
					"-c",
					instruction.Script(),
				}),
				llb.WithCustomName(fmt.Sprintf("%s%s-%d", node.Prefix, script.Desc, i)),
				llb.AddEnv("PKG_NAME", node.Pkg.Name),
			)

			if node.Graph.Options.NoCache {
				runOptions = append(runOptions, llb.IgnoreCache)
			}

			scriptRoot = scriptRoot.Run(runOptions...).Root()
		}

		if script.Detached {
			scriptRoot = scriptRoot.File(
				llb.Mkdir("/empty", constants.DefaultDirMode),
			)

			root = root.File(
				llb.Copy(scriptRoot, "/empty", "/", defaultCopyOptions(node.Graph.Options, false)),
			)
		} else {
			root = scriptRoot
		}
	}

	return root
}

func (node *NodeLLB) stepSBOM(root llb.State, step v1alpha2.Step) llb.State {
	if step.SBOM.OutputPath == "" {
		return root
	}

	sbomDoc, err := sbom.CreatePackageSBOM(node.Pkg, step.SBOM)
	if err != nil {
		return root
	}

	sbomJSON, err := sbom.ToSpdxJSON(*sbomDoc, node.Graph.Options.SourceDateEpoch)
	if err != nil {
		return root
	}

	root = root.File(
		llb.Mkdir(filepath.Dir(step.SBOM.OutputPath), constants.DefaultDirMode, llb.WithParents(true)),
	)

	root = root.File(
		llb.Mkfile(step.SBOM.OutputPath, 0o644, []byte(sbomJSON)),
	)

	return root
}

func (node *NodeLLB) step(root llb.State, i int, step v1alpha2.Step) llb.State {
	root = node.stepTmpDir(root, &step)
	root = node.stepDownload(root, step)
	root = node.stepEnvironment(root, step)
	root = node.stepScripts(root, i, step)
	root = node.stepSBOM(root, step)

	return root
}

func (node *NodeLLB) finalize(root llb.State) llb.State {
	stages := make([]llb.State, 0, len(node.Pkg.Finalize))

	for _, fin := range node.Pkg.Finalize {
		stages = append(stages,
			llb.Scratch().File(
				llb.Copy(root, fin.From, fin.To, defaultCopyOptions(node.Graph.Options, true)),
				llb.WithCustomNamef(node.Prefix+"finalize %s -> %s", fin.From, fin.To),
			),
		)
	}

	return llb.Merge(stages, llb.WithCustomName(node.Prefix+"finalize"))
}

// Build converts PackageNode to buildkit LLB.
func (node *NodeLLB) Build(ctx context.Context) (llb.State, error) {
	if state, ok := node.Graph.cache[node.PackageNode]; ok {
		return state, nil
	}

	root := node.base()

	root, err := node.dependencies(ctx, root)
	if err != nil {
		return llb.Scratch(), err
	}

	root = node.install(root)
	root = node.context(root)

	for i, step := range node.Pkg.Steps {
		root = node.step(root, i, step)
	}

	root = node.finalize(root)

	node.Graph.cache[node.PackageNode] = root

	return root, nil
}
