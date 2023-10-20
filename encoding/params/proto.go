package params

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/std"

	cryptocodec "github.com/stratosnet/stratos-chain/crypto/codec"
	stratos "github.com/stratosnet/stratos-chain/types"
)

// RegisterInterfaces registers Interfaces from types, crypto, and SDK std.
func RegisterInterfaces(interfaceRegistry codectypes.InterfaceRegistry) {
	std.RegisterInterfaces(interfaceRegistry)
	cryptocodec.RegisterInterfaces(interfaceRegistry)
	stratos.RegisterInterfaces(interfaceRegistry)
}
