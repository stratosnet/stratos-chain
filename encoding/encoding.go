package encoding

import (
	simappparams "cosmossdk.io/simapp/params"
	"github.com/cosmos/cosmos-sdk/types/module"

	"github.com/stratosnet/stratos-chain/encoding/params"
)

// MakeEncodingConfig creates an EncodingConfig for testing. This function
// should be used only in tests or when creating a new app instance (NewApp*()).
// App user shouldn't create new codecs - use the app.AppCodec instead.
func MakeEncodingConfig(mb module.BasicManager) simappparams.EncodingConfig {
	encodingConfig := params.MakeEncodingConfig()

	params.RegisterLegacyAminoCodec(encodingConfig.Amino)
	mb.RegisterLegacyAminoCodec(encodingConfig.Amino)
	params.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	mb.RegisterInterfaces(encodingConfig.InterfaceRegistry)
	return encodingConfig
}
