BUILDDIR ?= $(CURDIR)/build

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

VERSION := -ldflags="-X github.com/cosmos/cosmos-sdk/version.Version=v0.3.0"

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
# 	go $@ -mod=readonly $(BUILD_ARGS) $(VERSION) ./cmd/...
	go $@ $(BUILD_ARGS) $(VERSION) ./cmd/...
#	CGO_ENABLED=1 CGO_LDFLAGS="-lsnappy" go $@ -mod=readonly $(BUILD_ARGS) $(VERSION) -tags "cleveldb" ./cmd/...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

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
