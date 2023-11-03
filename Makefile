REGISTRY ?= ghcr.io
SERVER_IMG_NAME ?= server
SERVER_IMG_TAG ?= 0.1.0
EXPORTER_IMG_NAME ?= exporter
EXPORTER_IMG_TAG ?= 0.1.0

OUTPUT_TYPE ?= type=docker
BUILDPLATFORM ?= linux/amd64,linux/arm64
BUILDX_BUILDER_NAME ?= img-builder
QEMU_VERSION ?= 5.2.0-2

TOOLS_DIR := hack/tools
TOOLS_BIN_DIR := $(abspath $(TOOLS_DIR)/bin)

GOLANGCI_LINT_VER := v1.49.0
GOLANGCI_LINT_BIN := golangci-lint
GOLANGCI_LINT := $(abspath $(TOOLS_BIN_DIR)/$(GOLANGCI_LINT_BIN)-$(GOLANGCI_LINT_VER))

# Scripts
GO_INSTALL := ./hack/go-install.sh

## --------------------------------------
## Tooling Binaries
## --------------------------------------

$(GOLANGCI_LINT):
	GOBIN=$(TOOLS_BIN_DIR) $(GO_INSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint $(GOLANGCI_LINT_BIN) $(GOLANGCI_LINT_VER)


.PHONY: docker-buildx-builder
docker-buildx-builder:
	@if ! docker buildx ls | grep $(BUILDX_BUILDER_NAME); then \
  		docker run --rm --privileged multiarch/qemu-user-static:$(QEMU_VERSION) --reset -p yes; \
		docker buildx create --name $(BUILDX_BUILDER_NAME) --use; \
		docker buildx inspect $(BUILDX_BUILDER_NAME) --bootstrap; \
	fi

.PHONY: docker-build-server-image
docker-build-server-image: docker-buildx-builder
	 docker buildx build \
		--file docker/server/Dockerfile \
		--output=$(OUTPUT_TYPE) \
		--platform="$(BUILDPLATFORM)" \
		--pull \
		--tag $(REGISTRY)/$(SERVER_IMG_NAME):$(SERVER_IMG_TAG) .

.PHONY: docker-build-exporter-image
docker-build-exporter-image: docker-buildx-builder
	 docker buildx build \
		--file docker/exporter/Dockerfile \
		--output=$(OUTPUT_TYPE) \
		--platform="$(BUILDPLATFORM)" \
		--pull \
		--tag $(REGISTRY)/$(EXPORTER_IMG_NAME):$(EXPORTER_IMG_TAG) .


build:
	go build -o _output/bin/exporter ./cmd/exporter/

## --------------------------------------
## Linting
## --------------------------------------

.PHONY: lint
lint: $(GOLANGCI_LINT)
	$(GOLANGCI_LINT) run -v
