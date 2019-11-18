module github.com/talos-systems/bldr

go 1.13

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible
	github.com/alessio/shellescape v0.0.0-20190409004728-b115ca0f9053
	github.com/emicklei/dot v0.10.1
	github.com/google/uuid v1.1.1 // indirect
	github.com/hashicorp/go-multierror v1.0.0
	github.com/huandu/xstrings v1.2.0 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	github.com/moby/buildkit v0.6.2-0.20190921015714-10cef0c6e178
	github.com/opencontainers/go-digest v1.0.0-rc1
	github.com/otiai10/copy v1.0.2
	github.com/spf13/cobra v0.0.4
	golang.org/x/text v0.3.2 // indirect
	golang.org/x/xerrors v0.0.0-20190717185122-a985d3407aa7
	gopkg.in/yaml.v2 v2.2.2
)

replace github.com/docker/docker v1.14.0-0.20190319215453-e7b5f7dbe98c => github.com/docker/docker v1.4.2-0.20190319215453-e7b5f7dbe98c

replace golang.org/x/crypto v0.0.0-20190129210102-0709b304e793 => golang.org/x/crypto v0.0.0-20180904163835-0709b304e793
