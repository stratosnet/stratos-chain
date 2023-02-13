package net

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stratosnet/stratos-chain/rpc/backend"
)

// PublicAPI is the eth_ prefixed set of APIs in the Web3 JSON-RPC spec.
type PublicAPI struct {
	backend backend.BackendI
}

// NewPublicAPI creates an instance of the public Net Web3 API.
func NewPublicAPI(backend backend.BackendI) *PublicAPI {
	return &PublicAPI{
		backend: backend,
	}
}

// Version returns the current ethereum protocol version.
func (s *PublicAPI) Version() (string, error) {
	ctx := s.backend.GetSdkContext()
	params := s.backend.GetEVMKeeper().GetParams(ctx)
	return params.ChainConfig.ChainID.String(), nil
}

// Listening returns if client is actively listening for network connections.
func (s *PublicAPI) Listening() bool {
	return s.backend.GetNode().IsListening()
}

// PeerCount returns the number of peers currently connected to the client.
func (s *PublicAPI) PeerCount() hexutil.Big {
	return hexutil.Big(*big.NewInt(int64(len(s.backend.GetSwitch().Peers().List()))))
}
