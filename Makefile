#!/usr/bin/make -f

BUILDDIR ?= $(CURDIR)/build
LEDGER_ENABLED ?= false

APP_VER := v0.12.0
COMMIT := $(GIT_COMMIT_HASH)
TEST_DOCKER_REPO=stratos-chain-e2e

ifeq ($(COMMIT),)
	VERSION := $(APP_VER)
else
	VERSION := $(APP_VER)-$(COMMIT)
endif

ldflags= -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION)
ldflags += -X github.com/cosmos/cosmos-sdk/version.Name=stchain
ifeq ($(LEDGER_ENABLED),true)
  build_tags += ledger
endif

ifeq (cleveldb,$(findstring cleveldb,$(BUILD_OPTIONS)))
  build_tags += gcc cleveldb
  ldflags += -X github.com/cosmos/cosmos-sdk/types.DBBackend=cleveldb
endif

whitespace :=
whitespace := $(whitespace) $(whitespace)
comma := ,
build_tags_comma_sep := $(subst $(whitespace),$(comma),$(build_tags))
ldflags += -X "github.com/cosmos/cosmos-sdk/version.BuildTags=$(build_tags_comma_sep)"

BUILD_FLAGS += -ldflags '$(ldflags)'
BUILD_FLAGS += -tags "$(build_tags)"


BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
# 	go $@ -mod=readonly $(BUILD_ARGS) $(VERSION) ./cmd/...
	go$(GO_VERSION) $@ $(BUILD_ARGS) $(BUILD_FLAGS) ./cmd/...
#	CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" go $@ -mod=readonly $(BUILD_ARGS) $(BUILD_FLAGS) ./cmd/...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-linux: go.sum
	GOOS=linux GOARCH=amd64 $(MAKE) build

build-mac: go.sum
	GOOS=darwin GOARCH=amd64 $(MAKE) build

build-windows: go.sum
	GOOS=windows GOARCH=amd64 $(MAKE) build

build-cleveldb: go.sum
	 CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" BUILD_OPTIONS=cleveldb $(MAKE) build

clean:
	rm -rf $(BUILDDIR)/

coverage:
	go$(GO_VERSION) test ./... -coverprofile cover.out -coverpkg=./...
	go$(GO_VERSION) tool cover -html cover.out -o cover.html
	go$(GO_VERSION) tool cover -func cover.out | grep total:
	rm cover.out

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-stchaind-node:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	$(MAKE) -C networks/local
	@if ! [ -f build/node0/stchaind/config/genesis.json ]; then cp -r networks/local/cluster-config-sample/* build/ ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down

build-docker-e2e:
	@docker build -f tests/e2e/Dockerfile -t ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) --build-arg uid=$(shell id -u) --build-arg gid=$(shell id -g) .
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:latest

build-docker:
	@docker build -f Dockerfile -t ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) --build-arg uid=$(shell id -u) --build-arg gid=$(shell id -g) .
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:$(shell git rev-parse --abbrev-ref HEAD | sed 's#/#_#g')
	@docker tag ${TEST_DOCKER_REPO}:$(shell git rev-parse --short HEAD) ${TEST_DOCKER_REPO}:latest

.PHONY: build-linux build-mac build-cleveldb build clean
