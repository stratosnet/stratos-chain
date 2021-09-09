package common

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/client/context"
	registerKeeper "github.com/stratosnet/stratos-chain/x/register/keeper"
)

// QueryResourceNodeList queries the resource node list
func QueryResourceNodeList(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, registerKeeper.QueryResourceNodeList)
	return cliCtx.QueryWithData(route, nil)
}

// QueryIndexingNodeList queries the indexing node list
func QueryIndexingNodeList(cliCtx context.CLIContext, queryRoute string) ([]byte, int64, error) {
	route := fmt.Sprintf("custom/%s/%s", queryRoute, registerKeeper.QueryIndexingNodeList)
	return cliCtx.QueryWithData(route, nil)
}
