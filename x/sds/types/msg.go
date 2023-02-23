package types

import (
	"github.com/ipfs/go-cid"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// verify interface at compile time
var (
	_ sdk.Msg = &MsgFileUpload{}
	_ sdk.Msg = &MsgPrepay{}
)

const (
	TypeMsgFileUpload = "FileUploadTx"
	TypeMsgPrepay     = "SdsPrepayTx"
)

// NewMsgUpload creates a new Msg<Action> instance
func NewMsgUpload(fileHash string, from, reporter, uploader string) *MsgFileUpload {
	return &MsgFileUpload{
		FileHash: fileHash,
		From:     from,
		Reporter: reporter,
		Uploader: uploader,
	}
}

// nolint
func (msg MsgFileUpload) Route() string {
	return RouterKey
}

func (msg MsgFileUpload) Type() string {
	return TypeMsgFileUpload
}

func (msg MsgFileUpload) GetSigners() []sdk.AccAddress {
	accAddr, err := sdk.AccAddressFromBech32(msg.GetFrom())
	if err != nil {
		return []sdk.AccAddress{}
	}
	return []sdk.AccAddress{accAddr.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgFileUpload) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgFileUpload) ValidateBasic() error {
	_, err := cid.Decode(msg.FileHash)
	if err != nil {
		return sdkerrors.Wrap(ErrInvalidFileHash, "failed to validate file hash")
	}

	reporter, err := stratos.SdsAddressFromBech32(msg.GetReporter())
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "failed to parse reporter address")
	}

	uploader, err := sdk.AccAddressFromBech32(msg.GetUploader())
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "failed to parse uploader address")
	}

	if reporter.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing address of tx reporter")
	}
	if uploader.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing address of file uploader")
	}
	if len(msg.FileHash) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "missing file hash")
	}
	return nil
}

// NewMsgPrepay NewMsg<Action> creates a new Msg<Action> instance
func NewMsgPrepay(sender string, amount sdk.DecCoins) *MsgPrepay {
	return &MsgPrepay{
		Sender: sender,
		Amount: amount,
	}
}

func (msg MsgPrepay) Route() string {
	return RouterKey
}

func (msg MsgPrepay) Type() string {
	return TypeMsgPrepay
}

func (msg MsgPrepay) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.GetSender())
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{sender.Bytes()}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgPrepay) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgPrepay) ValidateBasic() error {
	sender, err := sdk.AccAddressFromBech32(msg.GetSender())
	if err != nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "invalid sender address")
	}

	if sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.Amount.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "missing amount to send")
	}
	return nil
}
