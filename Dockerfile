# syntax = docker/dockerfile-upstream:1.2.0

FROM golang:1.18-alpine AS base
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
ENV GOCACHE /.cache/go-build
ENV GOMODCACHE /.cache/mod
RUN apk --update --no-cache add bash curl
RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | bash -s -- -b /bin v1.45.2
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN --mount=type=cache,target=/.cache go mod download
RUN --mount=type=cache,target=/.cache go mod verify
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./main.go ./main.go
RUN --mount=type=cache,target=/.cache go list -mod=readonly all >/dev/null

FROM base AS build
ARG VERSION
ARG USERNAME
ARG REGISTRY
RUN --mount=type=cache,target=/.cache GOOS=linux CGO_ENABLED=0 \
    go build \
    -ldflags "-extldflags \"-static\" -s -w -X github.com/talos-systems/bldr/internal/pkg/constants.Version=${VERSION} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultOrganization=${USERNAME} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultRegistry=${REGISTRY}" \
    -o /bldr .
RUN --mount=type=cache,target=/.cache GOOS=linux \
    go test -c \
    -ldflags "-extldflags \"-static\" -s -w -X github.com/talos-systems/bldr/internal/pkg/constants.Version=${VERSION} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultOrganization=${USERNAME} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultRegistry=${REGISTRY}" \
    ./internal/pkg/integration

FROM base AS lint
COPY hack/golang/golangci-lint.yaml .
RUN --mount=type=cache,target=/.cache golangci-lint run --config golangci-lint.yaml

FROM scratch AS bldr
LABEL org.opencontainers.image.source https://github.com/siderolabs/bldr
COPY --from=build /bldr /bldr

FROM scratch AS integration.test
COPY --from=build /src/integration.test /integration.test

FROM bldr AS frontend
LABEL org.opencontainers.image.source https://github.com/siderolabs/bldr
ENTRYPOINT ["/bldr", "frontend"]
