package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
	authzcodec "github.com/cosmos/cosmos-sdk/x/authz/codec"
	govcodec "github.com/cosmos/cosmos-sdk/x/gov/codec"
	groupcodec "github.com/cosmos/cosmos-sdk/x/group/codec"
	"github.com/stratosnet/stratos-chain/x/register/types/v1_1"
)

// RegisterLegacyAminoCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgCreateResourceNode{}, "register/CreateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveResourceNode{}, "register/RemoveResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNode{}, "register/UpdateResourceNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateResourceNodeDeposit{}, "register/UpdateResourceNodeDepositTx", nil)

	cdc.RegisterConcrete(MsgCreateMetaNode{}, "register/CreateMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgRemoveMetaNode{}, "register/RemoveMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateMetaNode{}, "register/UpdateMetaNodeTx", nil)
	cdc.RegisterConcrete(MsgUpdateMetaNodeDeposit{}, "register/UpdateMetaNodeDepositTx", nil)
	cdc.RegisterConcrete(MsgMetaNodeRegistrationVote{}, "register/MsgMetaNodeRegistrationVote", nil)

	cdc.RegisterConcrete(MsgUpdateParams{}, "register/UpdateParamsTx", nil)
}

// RegisterInterfaces registers the x/register interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgCreateResourceNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateResourceNode{},
		&MsgUpdateResourceNodeDeposit{},
		&MsgCreateMetaNode{},
		&MsgRemoveResourceNode{},
		&MsgUpdateMetaNode{},
		&MsgUpdateMetaNodeDeposit{},
		&MsgMetaNodeRegistrationVote{},
		&MsgUpdateParams{},
		&v1_1.MsgCreateResourceNode{},
		&v1_1.MsgUpdateResourceNode{},
		&v1_1.MsgCreateMetaNode{},
		&v1_1.MsgUpdateMetaNode{},
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

	// Register all Amino interfaces and concrete types on the authz and gov Amino codec so that this can later be
	// used to properly serialize MsgGrant, MsgExec and MsgSubmitProposal instances
	RegisterLegacyAminoCodec(authzcodec.Amino)
	RegisterLegacyAminoCodec(govcodec.Amino)
	RegisterLegacyAminoCodec(groupcodec.Amino)
}
