package rpc

import (
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/cosmos/cosmos-sdk/client/context"

	"github.com/stratosnet/stratos-chain/crypto/ethsecp256k1"
	"github.com/stratosnet/stratos-chain/rpc/backend"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/eth"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/eth/filters"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/net"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/personal"
	"github.com/stratosnet/stratos-chain/rpc/namespaces/web3"
	rpctypes "github.com/stratosnet/stratos-chain/rpc/types"
)

// RPC namespaces and API version
const (
	Web3Namespace     = "web3"
	EthNamespace      = "eth"
	PersonalNamespace = "personal"
	NetNamespace      = "net"
	flagRPCAPI        = "rpc-api"

	apiVersion = "1.0"
)

// GetAPIs returns the list of all APIs from the Ethereum namespaces
func GetAPIs(clientCtx context.CLIContext, selectedApis []string, keys ...ethsecp256k1.PrivKey) []rpc.API {
	nonceLock := new(rpctypes.AddrLocker)
	backend := backend.New(clientCtx)
	ethAPI := eth.NewAPI(clientCtx, backend, nonceLock, keys...)

	var apis []rpc.API

	for _, api := range selectedApis {
		switch api {
		case Web3Namespace:
			apis = append(apis,
				rpc.API{
					Namespace: Web3Namespace,
					Version:   apiVersion,
					Service:   web3.NewAPI(),
					Public:    true,
				},
			)
		case EthNamespace:
			apis = append(apis,
				rpc.API{
					Namespace: EthNamespace,
					Version:   apiVersion,
					Service:   ethAPI,
					Public:    true,
				},
			)
			apis = append(apis,
				rpc.API{
					Namespace: EthNamespace,
					Version:   apiVersion,
					Service:   filters.NewAPI(clientCtx, backend),
					Public:    true,
				},
			)
		case PersonalNamespace:
			apis = append(apis,
				rpc.API{
					Namespace: PersonalNamespace,
					Version:   apiVersion,
					Service:   personal.NewAPI(ethAPI),
					Public:    false,
				},
			)
		case NetNamespace:
			apis = append(apis,
				rpc.API{
					Namespace: NetNamespace,
					Version:   apiVersion,
					Service:   net.NewAPI(clientCtx),
					Public:    true,
				},
			)
		}
	}

	return apis
}
