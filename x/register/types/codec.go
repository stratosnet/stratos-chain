package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// RegisterLegacyAminoCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateResourceNode{}, "register/CreateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveResourceNode{}, "register/RemoveResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNode{}, "register/UpdateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNodeStake{}, "register/UpdateResourceNodeStakeTx", nil)

	cdc.RegisterConcrete(MsgCreateMetaNode{}, "register/CreateMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveMetaNode{}, "register/RemoveMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateMetaNode{}, "register/UpdateMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateMetaNodeStake{}, "register/UpdateMetaNodeStakeTx", nil)
	cdc.RegisterConcrete(MsgMetaNodeRegistrationVote{}, "register/MsgMetaNodeRegistrationVote", nil)
}

// RegisterInterfaces registers the x/register interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateResourceNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateResourceNode{},
		&MsgUpdateResourceNodeStake{},
		&MsgCreateMetaNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateMetaNode{},
		&MsgUpdateMetaNodeStake{},
		&MsgMetaNodeRegistrationVote{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// ModuleCdc references the global x/register module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding as Amino is
// still used for that purpose.
//
// The actual codec used for serialization should be provided to x/register and
// defined at the application level.

var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterLegacyAminoCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
