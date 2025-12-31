# syntax = docker/dockerfile-upstream:1.20.0-labs

# THIS FILE WAS AUTOMATICALLY GENERATED, PLEASE DO NOT EDIT.
#
# Generated on 2025-12-31T05:18:26Z by kres 26be706.

ARG TOOLCHAIN=scratch

FROM --platform=linux/amd64 ghcr.io/siderolabs/ca-certificates:v1.12.0 AS image-ca-certificates

FROM --platform=linux/amd64 ghcr.io/siderolabs/fhs:v1.12.0 AS image-fhs

# runs markdownlint
FROM docker.io/oven/bun:1.3.4-alpine AS lint-markdown
WORKDIR /src
RUN bun i markdownlint-cli@0.47.0 sentences-per-line@0.3.0
COPY .markdownlint.json .
COPY ./CHANGELOG.md ./CHANGELOG.md
COPY ./README.md ./README.md
RUN bunx markdownlint --ignore "CHANGELOG.md" --ignore "**/node_modules/**" --ignore '**/hack/chglog/**' --rules sentences-per-line .

# base toolchain image
FROM --platform=${BUILDPLATFORM} ${TOOLCHAIN} AS toolchain
RUN apk --update --no-cache add bash build-base curl jq protoc protobuf-dev

# build tools
FROM --platform=${BUILDPLATFORM} toolchain AS tools
ENV GO111MODULE=on
ARG CGO_ENABLED
ENV CGO_ENABLED=${CGO_ENABLED}
ARG GOTOOLCHAIN
ENV GOTOOLCHAIN=${GOTOOLCHAIN}
ARG GOEXPERIMENT
ENV GOEXPERIMENT=${GOEXPERIMENT}
ENV GOPATH=/go
ARG DEEPCOPY_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go install github.com/siderolabs/deep-copy@${DEEPCOPY_VERSION} \
	&& mv /go/bin/deep-copy /bin/deep-copy
ARG GOLANGCILINT_VERSION
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCILINT_VERSION} \
	&& mv /go/bin/golangci-lint /bin/golangci-lint
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& mv /go/bin/govulncheck /bin/govulncheck
ARG GOFUMPT_VERSION
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
	&& mv /go/bin/gofumpt /bin/gofumpt

# tools and sources
FROM tools AS base
WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
RUN cd .
RUN --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go mod verify
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go list -mod=readonly all >/dev/null

FROM tools AS embed-generate
ARG SHA
ARG TAG
WORKDIR /src
RUN mkdir -p internal/version/data && \
    echo -n ${SHA} > internal/version/data/sha && \
    echo -n ${TAG} > internal/version/data/tag

# builds the integration test binary
FROM base AS integration-build
ARG REGISTRY
ARG USERNAME
ARG TAG
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/pkg/constants"
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg go test -c -tags integration -ldflags "-s -w -X ${VERSION_PKG}.Version=${TAG} -X ${VERSION_PKG}.DefaultOrganization=${USERNAME} -X ${VERSION_PKG}.DefaultRegistry=${REGISTRY}" ./internal/pkg/integration

# runs gofumpt
FROM base AS lint-gofumpt
RUN FILES="$(gofumpt -l .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumpt -w .':\n${FILES}"; exit 1)

# runs golangci-lint
FROM base AS lint-golangci-lint
WORKDIR /src
COPY .golangci.yml .
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint,id=bldr/root/.cache/golangci-lint,sharing=locked --mount=type=cache,target=/go/pkg,id=bldr/go/pkg golangci-lint run --config .golangci.yml

# runs golangci-lint fmt
FROM base AS lint-golangci-lint-fmt-run
WORKDIR /src
COPY .golangci.yml .
ENV GOGC=50
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint,id=bldr/root/.cache/golangci-lint,sharing=locked --mount=type=cache,target=/go/pkg,id=bldr/go/pkg golangci-lint fmt --config .golangci.yml
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint,id=bldr/root/.cache/golangci-lint,sharing=locked --mount=type=cache,target=/go/pkg,id=bldr/go/pkg golangci-lint run --fix --issues-exit-code 0 --config .golangci.yml

# runs govulncheck
FROM base AS lint-govulncheck
WORKDIR /src
COPY --chmod=0755 hack/govulncheck.sh ./hack/govulncheck.sh
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg ./hack/govulncheck.sh -exclude 'GO-2025-4020' ./...

# runs unit-tests with race detector
FROM base AS unit-tests-race
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg --mount=type=cache,target=/tmp,id=bldr/tmp CGO_ENABLED=1 go test -race ${TESTPKGS}

# runs unit-tests
FROM base AS unit-tests-run
WORKDIR /src
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg --mount=type=cache,target=/tmp,id=bldr/tmp go test -covermode=atomic -coverprofile=coverage.txt -coverpkg=${TESTPKGS} ${TESTPKGS}

FROM embed-generate AS embed-abbrev-generate
WORKDIR /src
ARG ABBREV_TAG
RUN echo -n 'undefined' > internal/version/data/sha && \
    echo -n ${ABBREV_TAG} > internal/version/data/tag

# copies out the integration test binary
FROM scratch AS integration.test
COPY --from=integration-build /src/integration.test /integration.test

# clean golangci-lint fmt output
FROM scratch AS lint-golangci-lint-fmt
COPY --from=lint-golangci-lint-fmt-run /src .

FROM scratch AS unit-tests
COPY --from=unit-tests-run /src/coverage.txt /coverage-unit-tests.txt

# cleaned up specs and compiled versions
FROM scratch AS generate
COPY --from=embed-abbrev-generate /src/internal/version internal/version

# builds bldr-darwin-amd64
FROM base AS bldr-darwin-amd64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg GOARCH=amd64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-darwin-amd64

# builds bldr-darwin-arm64
FROM base AS bldr-darwin-arm64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg GOARCH=arm64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-darwin-arm64

# builds bldr-linux-amd64
FROM base AS bldr-linux-amd64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-linux-amd64

# builds bldr-linux-arm64
FROM base AS bldr-linux-arm64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-linux-arm64

# builds bldr-linux-riscv64
FROM base AS bldr-linux-riscv64-build
COPY --from=generate / /
COPY --from=embed-generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build,id=bldr/root/.cache/go-build --mount=type=cache,target=/go/pkg,id=bldr/go/pkg GOARCH=riscv64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-linux-riscv64

FROM scratch AS bldr-darwin-amd64
COPY --from=bldr-darwin-amd64-build /bldr-darwin-amd64 /bldr-darwin-amd64

FROM scratch AS bldr-darwin-arm64
COPY --from=bldr-darwin-arm64-build /bldr-darwin-arm64 /bldr-darwin-arm64

FROM scratch AS bldr-linux-amd64
COPY --from=bldr-linux-amd64-build /bldr-linux-amd64 /bldr-linux-amd64

FROM scratch AS bldr-linux-arm64
COPY --from=bldr-linux-arm64-build /bldr-linux-arm64 /bldr-linux-arm64

FROM scratch AS bldr-linux-riscv64
COPY --from=bldr-linux-riscv64-build /bldr-linux-riscv64 /bldr-linux-riscv64

FROM bldr-linux-${TARGETARCH} AS bldr

FROM scratch AS bldr-all
COPY --from=bldr-darwin-amd64 / /
COPY --from=bldr-darwin-arm64 / /
COPY --from=bldr-linux-amd64 / /
COPY --from=bldr-linux-arm64 / /
COPY --from=bldr-linux-riscv64 / /

FROM scratch AS image-bldr
ARG TARGETARCH
COPY --from=bldr bldr-linux-${TARGETARCH} /bldr
COPY --from=image-fhs / /
COPY --from=image-ca-certificates / /
LABEL org.opencontainers.image.source=https://github.com/siderolabs/bldr
ENTRYPOINT ["/bldr","frontend"]

