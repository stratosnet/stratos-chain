package types

import "github.com/cosmos/cosmos-sdk/codec"

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgEthereumTx{}, "stratos/MsgEthereumTx", nil)
	cdc.RegisterConcrete(MsgStratosTx{}, "stratos/MsgStratosTx", nil)
	cdc.RegisterConcrete(TxData{}, "stratos/TxData", nil)
	cdc.RegisterConcrete(ChainConfig{}, "stratos/ChainConfig", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
