module github.com/talos-systems/bldr

go 1.14

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/alessio/shellescape v1.2.2
	github.com/containerd/containerd v1.4.0-0
	github.com/emicklei/dot v0.11.0
	github.com/hashicorp/go-multierror v1.1.0
	github.com/huandu/xstrings v1.3.1 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/moby/buildkit v0.7.1
	github.com/opencontainers/go-digest v1.0.0
	github.com/opencontainers/image-spec v1.0.1
	github.com/otiai10/copy v1.2.0
	github.com/spf13/cobra v1.0.0
	golang.org/x/sync v0.0.0-20190911185100-cd5d95a43a6e
	gopkg.in/yaml.v2 v2.3.0
)

replace (
	github.com/containerd/containerd => github.com/containerd/containerd v1.3.1-0.20200512144102-f13ba8f2f2fd
	github.com/docker/docker => github.com/docker/docker v17.12.0-ce-rc1.0.20200310163718-4634ce647cf2+incompatible
)
