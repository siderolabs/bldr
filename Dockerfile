# syntax = docker/dockerfile-upstream:1.2.0-labs

# THIS FILE WAS AUTOMATICALLY GENERATED, PLEASE DO NOT EDIT.
#
# Generated on 2022-11-21T11:25:20Z by kres latest.

ARG TOOLCHAIN

# cleaned up specs and compiled versions
FROM scratch AS generate

FROM ghcr.io/siderolabs/ca-certificates:v1.2.0 AS image-ca-certificates

FROM ghcr.io/siderolabs/fhs:v1.2.0 AS image-fhs

# runs markdownlint
FROM docker.io/node:19.0.1-alpine3.16 AS lint-markdown
WORKDIR /src
RUN npm i -g markdownlint-cli@0.32.2
RUN npm i sentences-per-line@0.2.1
COPY .markdownlint.json .
COPY ./README.md ./README.md
RUN markdownlint --ignore "CHANGELOG.md" --ignore "**/node_modules/**" --ignore '**/hack/chglog/**' --rules node_modules/sentences-per-line/index.js .

# base toolchain image
FROM ${TOOLCHAIN} AS toolchain
RUN apk --update --no-cache add bash curl build-base protoc protobuf-dev

# build tools
FROM --platform=${BUILDPLATFORM} toolchain AS tools
ENV GO111MODULE on
ARG CGO_ENABLED
ENV CGO_ENABLED ${CGO_ENABLED}
ENV GOPATH /go
ARG GOLANGCILINT_VERSION
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@${GOLANGCILINT_VERSION} \
	&& mv /go/bin/golangci-lint /bin/golangci-lint
ARG GOFUMPT_VERSION
RUN go install mvdan.cc/gofumpt@${GOFUMPT_VERSION} \
	&& mv /go/bin/gofumpt /bin/gofumpt
RUN go install golang.org/x/vuln/cmd/govulncheck@latest \
	&& mv /go/bin/govulncheck /bin/govulncheck
ARG GOIMPORTS_VERSION
RUN go install golang.org/x/tools/cmd/goimports@${GOIMPORTS_VERSION} \
	&& mv /go/bin/goimports /bin/goimports
ARG DEEPCOPY_VERSION
RUN go install github.com/siderolabs/deep-copy@${DEEPCOPY_VERSION} \
	&& mv /go/bin/deep-copy /bin/deep-copy

# tools and sources
FROM tools AS base
WORKDIR /src
COPY ./go.mod .
COPY ./go.sum .
RUN --mount=type=cache,target=/go/pkg go mod download
RUN --mount=type=cache,target=/go/pkg go mod verify
COPY ./cmd ./cmd
COPY ./internal ./internal
RUN --mount=type=cache,target=/go/pkg go list -mod=readonly all >/dev/null

# builds bldr-darwin-amd64
FROM base AS bldr-darwin-amd64-build
COPY --from=generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOARCH=amd64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-darwin-amd64

# builds bldr-darwin-arm64
FROM base AS bldr-darwin-arm64-build
COPY --from=generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOARCH=arm64 GOOS=darwin go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-darwin-arm64

# builds bldr-linux-amd64
FROM base AS bldr-linux-amd64-build
COPY --from=generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOARCH=amd64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-linux-amd64

# builds bldr-linux-arm64
FROM base AS bldr-linux-arm64-build
COPY --from=generate / /
WORKDIR /src/cmd/bldr
ARG GO_BUILDFLAGS
ARG GO_LDFLAGS
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/version"
ARG SHA
ARG TAG
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg GOARCH=arm64 GOOS=linux go build ${GO_BUILDFLAGS} -ldflags "${GO_LDFLAGS} -X ${VERSION_PKG}.Name=bldr -X ${VERSION_PKG}.SHA=${SHA} -X ${VERSION_PKG}.Tag=${TAG}" -o /bldr-linux-arm64

# builds the integration test binary
FROM base AS integration-build
ARG REGISTRY
ARG USERNAME
ARG TAG
ARG VERSION_PKG="github.com/siderolabs/bldr/internal/pkg/constants"
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg go test -c -tags integration -ldflags "-s -w -X ${VERSION_PKG}.Version=${TAG} -X ${VERSION_PKG}.DefaultOrganization=${USERNAME} -X ${VERSION_PKG}.DefaultRegistry=${REGISTRY}" ./internal/pkg/integration

# runs gofumpt
FROM base AS lint-gofumpt
RUN FILES="$(gofumpt -l .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'gofumpt -w .':\n${FILES}"; exit 1)

# runs goimports
FROM base AS lint-goimports
RUN FILES="$(goimports -l -local github.com/siderolabs/bldr .)" && test -z "${FILES}" || (echo -e "Source code is not formatted with 'goimports -w -local github.com/siderolabs/bldr .':\n${FILES}"; exit 1)

# runs golangci-lint
FROM base AS lint-golangci-lint
COPY .golangci.yml .
ENV GOGC 50
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/root/.cache/golangci-lint --mount=type=cache,target=/go/pkg golangci-lint run --config .golangci.yml

# runs govulncheck
FROM base AS lint-govulncheck
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg govulncheck ./...

# runs unit-tests with race detector
FROM base AS unit-tests-race
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg --mount=type=cache,target=/tmp CGO_ENABLED=1 go test -v -race -count 1 ${TESTPKGS}

# runs unit-tests
FROM base AS unit-tests-run
ARG TESTPKGS
RUN --mount=type=cache,target=/root/.cache/go-build --mount=type=cache,target=/go/pkg --mount=type=cache,target=/tmp go test -v -covermode=atomic -coverprofile=coverage.txt -coverpkg=${TESTPKGS} -count 1 ${TESTPKGS}

FROM scratch AS bldr-darwin-amd64
COPY --from=bldr-darwin-amd64-build /bldr-darwin-amd64 /bldr-darwin-amd64

FROM scratch AS bldr-darwin-arm64
COPY --from=bldr-darwin-arm64-build /bldr-darwin-arm64 /bldr-darwin-arm64

FROM scratch AS bldr-linux-amd64
COPY --from=bldr-linux-amd64-build /bldr-linux-amd64 /bldr-linux-amd64

FROM scratch AS bldr-linux-arm64
COPY --from=bldr-linux-arm64-build /bldr-linux-arm64 /bldr-linux-arm64

# copies out the integration test binary
FROM scratch AS integration.test
COPY --from=integration-build /src/integration.test /integration.test

FROM scratch AS unit-tests
COPY --from=unit-tests-run /src/coverage.txt /coverage.txt

FROM bldr-linux-${TARGETARCH} AS bldr

FROM scratch AS bldr-all
COPY --from=bldr-darwin-amd64 / /
COPY --from=bldr-darwin-arm64 / /
COPY --from=bldr-linux-amd64 / /
COPY --from=bldr-linux-arm64 / /

FROM scratch AS image-bldr
ARG TARGETARCH
COPY --from=bldr bldr-linux-${TARGETARCH} /bldr
COPY --from=image-fhs / /
COPY --from=image-ca-certificates / /
LABEL org.opencontainers.image.source https://github.com/siderolabs/bldr
ENTRYPOINT ["/bldr","frontend"]

