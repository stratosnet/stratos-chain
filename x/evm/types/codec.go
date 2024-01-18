package types

import (
	"cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
	tx "github.com/cosmos/cosmos-sdk/types/tx"
	"github.com/gogo/protobuf/proto"
)

var ModuleCdc = codec.NewProtoCodec(codectypes.NewInterfaceRegistry())

type (
	ExtensionOptionsEthereumTxI interface{}
)

// RegisterInterfaces registers the client interfaces to protobuf Any.
func RegisterInterfaces(registry codectypes.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgEthereumTx{},
		&MsgUpdateParams{},
	)
	registry.RegisterInterface(
		"stratos.evm.v1.ExtensionOptionsEthereumTx",
		(*ExtensionOptionsEthereumTxI)(nil),
		&ExtensionOptionsEthereumTx{},
	)
	registry.RegisterInterface(
		"stratos.evm.v1.TxData",
		(*TxData)(nil),
		&DynamicFeeTx{},
		&AccessListTx{},
		&LegacyTx{},
	)
	// for tx decoder handling
	registry.RegisterImplementations(
		(*tx.TxExtensionOptionI)(nil),
		&MsgEthereumTx{},
		&ExtensionOptionsEthereumTx{},
	)

	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

// PackTxData constructs a new Any packed with the given tx data value. It returns
// an error if the client state can't be cast to a protobuf message or if the concrete
// implementation is not registered to the protobuf codec.
func PackTxData(txData TxData) (*codectypes.Any, error) {
	msg, ok := txData.(proto.Message)
	if !ok {
		return nil, errors.Wrapf(sdkerrors.ErrPackAny, "cannot proto marshal %T", txData)
	}

	anyTxData, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return nil, errors.Wrap(sdkerrors.ErrPackAny, err.Error())
	}

	return anyTxData, nil
}

// UnpackTxData unpacks an Any into a TxData. It returns an error if the
// client state can't be unpacked into a TxData.
func UnpackTxData(any *codectypes.Any) (TxData, error) {
	if any == nil {
		return nil, errors.Wrap(sdkerrors.ErrUnpackAny, "protobuf Any message cannot be nil")
	}

	txData, ok := any.GetCachedValue().(TxData)
	if !ok {
		return nil, errors.Wrapf(sdkerrors.ErrUnpackAny, "cannot unpack Any into TxData %T", any)
	}

	return txData, nil
}
