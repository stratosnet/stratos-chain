package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

// RegisterCodec registers concrete types on codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgCreateResourceNode{}, "register/CreateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveResourceNode{}, "register/RemoveResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNode{}, "register/UpdateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNodeStake{}, "register/UpdateResourceNodeStakeTx", nil)

	cdc.RegisterConcrete(MsgCreateIndexingNode{}, "register/CreateIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveIndexingNode{}, "register/RemoveIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateIndexingNode{}, "register/UpdateIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateIndexingNodeStake{}, "register/UpdateIndexingNodeStakeTx", nil)

	cdc.RegisterConcrete(MsgIndexingNodeRegistrationVote{}, "register/MsgIndexingNodeRegistrationVote", nil)
}

// ModuleCdc defines the module codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.New()
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
