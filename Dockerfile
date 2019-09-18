FROM alpine:3.9 AS base
RUN apk add build-base ca-certificates go
WORKDIR /src
COPY ./go.mod ./
COPY ./go.sum ./
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
COPY . .
ARG VERSION
ARG USERNAME
ARG REGISTRY
RUN GOOS=linux CGO_ENABLED=0 go build  -a -ldflags "-extldflags \"-static\" -s -w -X github.com/talos-systems/bldr/internal/pkg/constants.Version=${VERSION} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultOrganization=${USERNAME} -X github.com/talos-systems/bldr/internal/pkg/constants.DefaultRegistry=${REGISTRY}" -o /bldr .
FROM scratch AS bldr
COPY --from=build /bldr /bldr

FROM scratch AS scratch
ARG TARGETPLATFORM
ENV TARGETPLATFORM ${TARGETPLATFORM}
ARG BUILDPLATFORM
ENV BUILDPLATFORM ${BUILDPLATFORM}
COPY --from=bldr /bldr /bldr
WORKDIR /pkg
ENTRYPOINT ["/bldr"]
CMD ["build"]

FROM alpine:3.9 AS alpine
ARG TARGETPLATFORM
ENV TARGETPLATFORM ${TARGETPLATFORM}
ARG BUILDPLATFORM
ENV BUILDPLATFORM ${BUILDPLATFORM}
RUN apk --no-cache add bash ca-certificates
RUN [ "ln", "-svf", "/bin/bash", "/bin/sh" ]
COPY --from=bldr /bldr /bldr
WORKDIR /pkg
ENTRYPOINT ["/bldr"]
CMD ["build"]
