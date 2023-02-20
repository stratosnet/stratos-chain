package server

import (
	"github.com/tendermint/tendermint/node"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"
	ethlog "github.com/ethereum/go-ethereum/log"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/stratosnet/stratos-chain/rpc"
	"github.com/stratosnet/stratos-chain/server/config"
	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
)

// StartJSONRPC starts the JSON-RPC server
func StartJSONRPC(ctx *server.Context, tmNode *node.Node, evmKeeper *evmkeeper.Keeper, ms storetypes.MultiStore, clientCtx client.Context, config config.Config) error {
	logger := ctx.Logger.With("module", "geth")
	ethlog.Root().SetHandler(ethlog.FuncHandler(func(r *ethlog.Record) error {
		switch r.Lvl {
		case ethlog.LvlTrace, ethlog.LvlDebug:
			logger.Debug(r.Msg, r.Ctx...)
		case ethlog.LvlInfo, ethlog.LvlWarn:
			logger.Info(r.Msg, r.Ctx...)
		case ethlog.LvlError, ethlog.LvlCrit:
			logger.Error(r.Msg, r.Ctx...)
		}
		return nil
	}))

	apis := rpc.GetRPCAPIs(ctx, tmNode, evmKeeper, ms, clientCtx, config.JSONRPC.API)
	web3Srv := rpc.NewWeb3Server(config, logger)
	err := web3Srv.StartHTTP(apis)
	if err != nil {
		return err
	}
	err = web3Srv.StartWS(apis)
	if err != nil {
		return err
	}
	return nil
}
