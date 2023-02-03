package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// RegisterLegacyAminoCodec RegisterCodec registers concrete types on codec
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	// this line is used by starport scaffolding # 1
	cdc.RegisterConcrete(MsgVolumeReport{}, "pot/VolumeReportTx", nil)
	cdc.RegisterConcrete(MsgWithdraw{}, "pot/WithdrawTx", nil)
	cdc.RegisterConcrete(MsgFoundationDeposit{}, "pot/FoundationDepositTx", nil)
	cdc.RegisterConcrete(MsgSlashingResourceNode{}, "pot/SlashingResourceNodeTx", nil)
}

// RegisterInterfaces registers the x/register interfaces types with the interface registry
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations((*sdk.Msg)(nil),
		&MsgVolumeReport{},
		&MsgWithdraw{},
		&MsgFoundationDeposit{},
		&MsgSlashingResourceNode{},
	)
	registry.RegisterImplementations(
		(*authz.Authorization)(nil),
		//&StakeAuthorization{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// ModuleCdc references the global x/pot module codec. Note, the codec should
// ONLY be used in certain instances of tests and for JSON encoding as Amino is
// still used for that purpose.
//
// The actual codec used for serialization should be provided to x/pot and
// defined at the application level.
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterLegacyAminoCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
