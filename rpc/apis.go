// Package rpc contains RPC handler methods and utilities to start
// Stratos Web3-compatibly JSON-RPC server.
package rpc

import (
	"github.com/cometbft/cometbft/node"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/server"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/debug"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/eth"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/eth/filters"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/miner"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/net"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/personal"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/txpool"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/ethereum/web3"
	"github.com/stratosnet/stratos-chain/rpc/types"

	evmkeeper "github.com/stratosnet/stratos-chain/x/evm/keeper"
)

// RPC namespaces and API version
const (
	// Cosmos namespaces

	CosmosNamespace = "cosmos"

	// Ethereum namespaces

	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"
	TxPoolNamespace   = "txpool"
	DebugNamespace    = "debug"
	MinerNamespace    = "miner"

	apiVersion = "1.0"
)

// GetRPCAPIs returns the list of all APIs
func GetRPCAPIs(ctx *server.Context, tmNode *node.Node, evmKeeper *evmkeeper.Keeper, ms storetypes.MultiStore, clientCtx client.Context, selectedAPIs []string) ([]rpc.API, error) {
	nonceLock := new(types.AddrLocker)
	evmBackend, err := backend.NewBackend(ctx, tmNode, evmKeeper, ms, ctx.Logger, clientCtx)

	if err != nil {
		return []rpc.API{}, err
	}

	return []rpc.API{
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   eth.NewPublicAPI(ctx.Logger, clientCtx, evmBackend, nonceLock),
			Public:    true,
		},
		{
			Namespace: EthNamespace,
			Version:   apiVersion,
			Service:   filters.NewPublicAPI(ctx.Logger, clientCtx, tmNode.EventBus(), evmBackend),
			Public:    true,
		},
		{
			Namespace: Web3Namespace,
			Version:   apiVersion,
			Service:   web3.NewPublicAPI(),
			Public:    true,
		},
		{
			Namespace: NetNamespace,
			Version:   apiVersion,
			Service:   net.NewPublicAPI(evmBackend),
			Public:    true,
		},
		{
			Namespace: PersonalNamespace,
			Version:   apiVersion,
			Service:   personal.NewAPI(ctx.Logger, clientCtx, evmBackend),
			Public:    false,
		},
		{
			Namespace: TxPoolNamespace,
			Version:   apiVersion,
			Service:   txpool.NewPublicAPI(ctx.Logger, clientCtx, evmBackend),
			Public:    true,
		},
		{
			Namespace: DebugNamespace,
			Version:   apiVersion,
			Service:   debug.NewAPI(ctx, evmBackend, clientCtx),
			Public:    true,
		},
		{
			Namespace: MinerNamespace,
			Version:   apiVersion,
			Service:   miner.NewPrivateAPI(ctx, clientCtx, evmBackend),
			Public:    false,
		},
	}, nil
}
