package api

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stratosnet/stratos-chain/web3/types"
)

var _ types.Web3Service = (*Web3API)(nil)

// PublicWeb3API is the web3_ prefixed set of APIs in the Web3 JSON-RPC spec.
type Web3API struct{}

// NewAPI creates an instance of the Web3 API.
func NewAPI() *Web3API {
	return &Web3API{}
}

// ClientVersion returns the client version in the Web3 user agent format.
func (*Web3API) ClientVersion() string {
	info := version.NewInfo()
	fmt.Println("ClientVersion@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@@")
	return fmt.Sprintf("%s-%s", info.Name, info.Version)
}

// Sha3 returns the keccak-256 hash of the passed-in input.
func (*Web3API) Sha3(input hexutil.Bytes) hexutil.Bytes {
	return crypto.Keccak256(input)
}
