REGISTRY ?= docker.io
USERNAME ?= autonomy
TAG = $(shell gitmeta image tag)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)

PLATFORM ?= linux/amd64,linux/arm64
PROGRESS ?= auto
PUSH ?= false

BUILD := docker buildx build
COMMON_ARGS := --progress=$(PROGRESS)
COMMON_ARGS += --platform=$(PLATFORM)
COMMON_ARGS += --build-arg=VERSION=$(TAG)

PKGS := alpine scratch bldr

all: $(PKGS)

.PHONY: bldr
bldr:
	mkdir -p out
	$(BUILD) $(COMMON_ARGS) \
	--target=$@ \
	--output=type=local,dest=./out \
	-f ./Dockerfile \
	.

.PHONY: scratch
scratch:

	$(BUILD) $(COMMON_ARGS) \
	--push=$(PUSH) \
	--target=$@ \
	--tag $(REGISTRY_AND_USERNAME)/bldr:$(TAG)-$@ \
	-f ./Dockerfile \
	.

.PHONY: alpine
alpine:
	$(BUILD) $(COMMON_ARGS) \
	--push=$(PUSH) \
	--target=$@ \
	--tag $(REGISTRY_AND_USERNAME)/bldr:$(TAG)-$@ \
	-f ./Dockerfile \
	.
