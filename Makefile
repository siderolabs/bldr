REGISTRY ?= ghcr.io
USERNAME ?= siderolabs
TAG ?= $(shell git describe --tag --always --dirty)
REGISTRY_AND_USERNAME := $(REGISTRY)/$(USERNAME)
RUN_TESTS ?= TestIntegration

BUILD_PLATFORM ?= linux/amd64
PLATFORM ?= linux/amd64,linux/arm64
PROGRESS ?= auto
PUSH ?= false

BUILD := docker buildx build
COMMON_ARGS := --progress=$(PROGRESS)
COMMON_ARGS += --platform=$(PLATFORM)
COMMON_ARGS += --build-arg=VERSION=$(TAG)
COMMON_ARGS += --build-arg=USERNAME=$(USERNAME)
COMMON_ARGS += --build-arg=REGISTRY=$(REGISTRY)
COMMON_ARGS += --build-arg=http_proxy=$(http_proxy)
COMMON_ARGS += --build-arg=https_proxy=$(https_proxy)

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

.PHONY: integration.test
integration.test:
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

.PHONY: integration
integration: integration.test bldr
	cd internal/pkg/integration && PATH="$$PWD/../../../out/$(subst /,_,$(BUILD_PLATFORM)):$$PATH"  integration.test -test.v -test.run $(RUN_TESTS)

.PHONY: conformance
conformance: ## Performs policy checks against the commit and source code.
	docker run --rm -it -v $(PWD):/src -w /src ghcr.io/siderolabs/conform:v0.1.0-alpha.22 enforce
