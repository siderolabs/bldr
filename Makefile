REGISTRY ?= docker.io
USERNAME ?= autonomy
TAG = $(shell git describe --tags --dirty --always)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)

PLATFORM ?= linux/amd64
PROGRESS ?= auto
PUSH ?= false

BUILD := docker buildx build
COMMON_ARGS := --progress=$(PROGRESS)
COMMON_ARGS += --platform=$(PLATFORM)
COMMON_ARGS += --build-arg=VERSION=$(TAG)
COMMON_ARGS += --build-arg=USERNAME=$(USERNAME)
COMMON_ARGS += --build-arg=REGISTRY=$(REGISTRY)

PKGS := frontend bldr

all: $(PKGS) lint

.PHONY: bldr
bldr:
	mkdir -p out
	$(BUILD) $(COMMON_ARGS) \
	--target=$@ \
	--output=type=local,dest=./out \
	-f ./Dockerfile \
	.

.PHONY: lint
lint:
	$(BUILD) $(COMMON_ARGS) \
	--target=$@ \
	-f ./Dockerfile \
	.

.PHONY: frontend
frontend:

	$(BUILD) $(COMMON_ARGS) \
	--push=$(PUSH) \
	--target=$@ \
	--tag $(REGISTRY_AND_USERNAME)/bldr:$(TAG)-$@ \
	-f ./Dockerfile \
	.
