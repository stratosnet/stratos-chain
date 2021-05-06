BUILDDIR ?= $(CURDIR)/build

BUILD_FLAGS := -tags "$(build_tags)" -ldflags '$(ldflags)'

BUILD_TARGETS := build install

build: BUILD_ARGS=-o $(BUILDDIR)/

$(BUILD_TARGETS): go.sum $(BUILDDIR)/
	go $@ -mod=readonly $(BUILD_ARGS) ./...

$(BUILDDIR)/:
	mkdir -p $(BUILDDIR)/

build-linux: go.sum
	LEDGER_ENABLED=false GOOS=linux GOARCH=amd64 $(MAKE) build

###############################################################################
###                                Localnet                                 ###
###############################################################################

build-docker-stratos-chaind-node:
	$(MAKE) -C networks/local

# Run a 4-node testnet locally
localnet-start: build-linux localnet-stop
	$(MAKE) -C networks/local
	@if ! [ -f build/node0/stratos-chaind/config/genesis.json ]; then cp -r networks/local/cluster-config-sample/* build/ ; fi
	docker-compose up -d

# Stop testnet
localnet-stop:
	docker-compose down
