package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	"github.com/cosmos/cosmos-sdk/x/authz"
)

// RegisterCodec registers concrete types on codec
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

// ModuleCdc defines the module codec
var ModuleCdc *codec.LegacyAmino

func init() {
	ModuleCdc = codec.NewLegacyAmino()
	RegisterLegacyAminoCodec(ModuleCdc)
	cryptocodec.RegisterCrypto(ModuleCdc)
	ModuleCdc.Seal()
}
