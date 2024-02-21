package client

import (
	govclient "github.com/cosmos/cosmos-sdk/x/gov/client"

	"github.com/stratosnet/stratos-chain/x/evm/client/cli"
)

var (
	EVMChangeProxyImplementationHandler = govclient.NewProposalHandler(cli.NewEVMProxyImplmentationUpgrade)
)
