# syntax=docker/dockerfile:1

# Build the manager binary
ARG GO_VERSION=1.19
FROM --platform=$BUILDPLATFORM golang:${GO_VERSION} as builder

WORKDIR /workspace
# Copy the Go Modules manifests
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,source=go.sum,target=go.sum \
    --mount=type=bind,source=go.mod,target=go.mod \
    go mod download -x

# Copy the go source
COPY cmd/exporter/ cmd/exporter/
COPY pkg/ pkg/

# Build
ARG TARGETOS TARGETARCH
RUN --mount=type=cache,target=/go/pkg/mod/ \
    --mount=type=bind,target=. \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH GO111MODULE=on go build -o /bin/exporter cmd/exporter/main.go

# Use distroless as minimal base image to package the exporter binary
# Refer to https://github.com/GoogleContainerTools/distroless for more details
FROM --platform=${TARGETPLATFORM:-linux/amd64} gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /bin/exporter .

ENTRYPOINT ["/exporter"]