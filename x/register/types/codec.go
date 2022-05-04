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
	cdc.RegisterConcrete(MsgCreateResourceNode{}, "register/CreateResourceNode", nil)
	cdc.RegisterConcrete(MsgRemoveResourceNode{}, "register/RemoveResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNode{}, "register/UpdateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNodeStake{}, "register/UpdateResourceNodeStakeTx", nil)

	cdc.RegisterConcrete(MsgCreateIndexingNode{}, "register/CreateIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveIndexingNode{}, "register/RemoveIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateIndexingNode{}, "register/UpdateIndexingNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateIndexingNodeStake{}, "register/UpdateIndexingNodeStakeTx", nil)
	cdc.RegisterConcrete(MsgIndexingNodeRegistrationVote{}, "register/MsgIndexingNodeRegistrationVote", nil)
}

// RegisterInterfaces registers the x/register interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateResourceNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateResourceNode{},
		&MsgUpdateResourceNodeStake{},
		&MsgCreateIndexingNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateIndexingNode{},
		&MsgUpdateIndexingNodeStake{},
		&MsgIndexingNodeRegistrationVote{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		//&StakeAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/register module codec. Note, the codec should
	// ONLY be used in certain instances of tests and for JSON encoding as Amino is
	// still used for that purpose.
	//
	// The actual codec used for serialization should be provided to x/register and
	// defined at the application level.
	ModuleCdc = codec.NewAminoCodec(amino)
)

// ModuleCdc defines the module codec
func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
