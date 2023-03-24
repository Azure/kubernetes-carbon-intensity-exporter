REGISTRY ?= ghcr.io
SERVER_IMG_NAME ?= server
SERVERIMG_TAG ?= 0.1.0

OUTPUT_TYPE ?= type=docker
BUILDPLATFORM ?= linux/amd64
BUILDX_BUILDER_NAME ?= img-builder
QEMU_VERSION ?= 5.2.0-2

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
		--tag $(REGISTRY)/$(SERVER_IMG_NAME):$(SERVERIMG_TAG) .


build:
	go build -o _output/bin/carbon-data-provider ./cmd/carbon-data-provider/
