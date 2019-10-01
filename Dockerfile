# syntax = docker/dockerfile-upstream:1.1.2-experimental

FROM golang:1.13.1-alpine AS base
ENV GO111MODULE on
ENV GOPROXY https://proxy.golang.org
ENV CGO_ENABLED 0
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
RUN go mod download
RUN go mod verify
COPY . .
RUN go list -mod=readonly all >/dev/null

FROM base AS build
ARG VERSION
ARG USERNAME
ARG REGISTRY
RUN --mount=type=cache,target=/root/.cache/go-build GOOS=linux CGO_ENABLED=0 \
    go build \
    -ldflags "-extldflags \"-static\" -s -w -X github.com/talos-systems/bldr/internal/pkg/constants.Version=${VERSION} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultOrganization=${USERNAME} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultRegistry=${REGISTRY}" \
    -o /bldr .

FROM scratch AS bldr
COPY --from=build /bldr /bldr

FROM bldr AS frontend
ENTRYPOINT ["/bldr", "frontend"]
