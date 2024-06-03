package types

import (
	"bytes"

	"cosmossdk.io/errors"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/proto"

	"github.com/ethereum/go-ethereum/common"

	prototypes "github.com/cosmos/gogoproto/types"
)

var (
	_ sdk.Msg                            = &MsgCosmosData{}
	_ codectypes.UnpackInterfacesMessage = &MsgCosmosData{}

	// Web3CosmosAddress is used as a target for resolve tx data
	Web3CosmosAddress = common.BytesToAddress([]byte("web3cosmos"))
)

func IsCosmosHandler(to *common.Address) bool {
	if to != nil && bytes.Equal(to.Bytes(), Web3CosmosAddress.Bytes()) {
		return true
	}
	return false
}

// TxDataToAny convert evm tx data to any
func TxDataToAny(data []byte) (*codectypes.Any, error) {
	var msgAny prototypes.Any
	if err := proto.Unmarshal(data, &msgAny); err != nil {
		return nil, err
	}

	var pb prototypes.DynamicAny

	if err := prototypes.UnmarshalAny(&msgAny, &pb); err != nil {
		return nil, err
	}

	any, err := codectypes.NewAnyWithValue(pb.Message)
	if err != nil {
		return nil, err
	}
	return any, nil
}

// TxDataToMsg convert evm tx data to cosmos msg
func TxDataToSdkMsg(cdc codectypes.AnyUnpacker, data []byte) (sdk.Msg, error) {
	any, err := TxDataToAny(data)
	if err != nil {
		return nil, err
	}

	cMsg, ok := any.GetCachedValue().(sdk.Msg)
	if !ok {
		return nil, err
	}

	// this should resolved nested anyies in msg
	if err := codectypes.UnpackInterfaces(cMsg, cdc); err != nil {
		return nil, err
	}

	return cMsg, nil
}

// MsgCosmosData is dynamic data from MsgEthereumTx which will be handled by the router
type MsgCosmosData struct {
	cdc  codectypes.AnyUnpacker
	from sdk.Address
	msg  sdk.Msg
}

func NewMsgCosmosData(cdc codectypes.AnyUnpacker, from sdk.Address) *MsgCosmosData {
	return &MsgCosmosData{cdc: cdc, from: from}
}

func (mcd *MsgCosmosData) IsNil() bool {
	return mcd.msg == nil
}

// SetMsg explisitly (could bypass Parse)
func (mcd *MsgCosmosData) SetMsg(msg sdk.Msg) {
	mcd.msg = msg
}

// Parse tx data to sdk.Msg
func (mcd *MsgCosmosData) Parse(data []byte) error {
	cMsg, err := TxDataToSdkMsg(mcd.cdc, data)
	if err != nil {
		return err
	}
	mcd.SetMsg(cMsg)

	return nil
}

func (mcd *MsgCosmosData) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if mcd.IsNil() {
		return nil
	}
	return codectypes.UnpackInterfaces(mcd.msg, unpacker)
}

func (mcd *MsgCosmosData) String() string {
	if mcd.IsNil() {
		return "No msg"
	}
	return mcd.msg.String()
}

func (mcd *MsgCosmosData) Reset() {}

// ProtoMessage invoke msg ProtoMessage
func (mcd *MsgCosmosData) ProtoMessage() {
	if mcd.IsNil() {
		return
	}
	mcd.msg.ProtoMessage()
}

// GetMsgs returns a single sdk.Msg.
func (mcd *MsgCosmosData) GetMsgs() []sdk.Msg {
	if mcd.IsNil() {
		return []sdk.Msg{}
	}
	return []sdk.Msg{mcd.msg}
}

// GetSigners return signers of sdk.Msg
func (mcd *MsgCosmosData) GetSigners() []sdk.AccAddress {
	if mcd.IsNil() {
		return []sdk.AccAddress{}
	}
	return mcd.msg.GetSigners()
}

// ValidateCosmosMsg validates msg (mostly for ante checks and estimates)
func (mcd *MsgCosmosData) ValidateBasic() (err error) {
	// could be occured in nested msg during check if something missed or unparsed (like signers check)
	defer func() {
		if r := recover(); r != nil {
			err = errors.Wrapf(sdkerrors.ErrUnknownRequest, "%v", r)
		}
	}()

	if mcd.IsNil() {
		return errors.Wrapf(sdkerrors.ErrUnknownRequest, "msg not found")
	}

	msg := mcd.GetMsgs()[0]

	// should be first as it is faster to compare and skip tx
	_, ok := msg.(*MsgEthereumTx)
	if ok {
		return errors.Wrapf(sdkerrors.ErrInvalidType, "circular msgs not allowed")
	}

	cSigners := msg.GetSigners()
	if len(cSigners) == 0 || mcd.from == nil || !bytes.Equal(cSigners[0], mcd.from.Bytes()) {
		return errors.Wrapf(sdkerrors.ErrorInvalidSigner, "cosmos signer is not a signer of eth msg")
	}

	if err := msg.ValidateBasic(); err != nil {
		return errors.Wrap(err, "tx basic validation failed")
	}

	return nil
}
