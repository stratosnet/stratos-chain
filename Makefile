BUILDDIR ?= $(CURDIR)/build

APP_VER := v0.9.0
COMMIT := $(GIT_COMMIT_HASH)

ifeq ($(COMMIT),)
    VERSION := $(APP_VER)
else
	VERSION := $(APP_VER)-$(COMMIT)
endif

ldflags= -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION)

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
# 	go $@ -mod=readonly $(BUILD_ARGS) $(VERSION) ./cmd/...
	go $@ $(BUILD_ARGS) $(BUILD_FLAGS) ./cmd/...
#	CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" go $@ -mod=readonly $(BUILD_ARGS) $(VERSION) -tags "cleveldb" ./cmd/...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

build-mac: go.sum
	LEDGER_ENABLED=false GOOS=darwin GOARCH=amd64 $(MAKE) build

build-windows: go.sum
	LEDGER_ENABLED=false GOOS=windows GOARCH=amd64 $(MAKE) build

clean:
	rm -rf $(BUILDDIR)/

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

.PHONY: build-linux build-mac build clean
