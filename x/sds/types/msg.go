package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ipfs/go-cid"
	stratos "github.com/stratosnet/stratos-chain/types"
)

const (
	ConstFileUpload = "FileUploadTx"
	ConstSdsPrepay  = "SdsPrepayTx"
)

// verify interface at compile time
var _ sdk.Msg = &MsgFileUpload{}

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
func (msg MsgFileUpload) Route() string { return RouterKey }
func (msg MsgFileUpload) Type() string  { return ConstFileUpload }
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
		return sdkerrors.Wrap(err, "failed to validate file hash")
	}

	reporter, err := stratos.SdsAddressFromBech32(msg.GetReporter())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to parse reporter address")
	}

	uploader, err := sdk.AccAddressFromBech32(msg.GetUploader())
	if err != nil {
		return sdkerrors.Wrap(err, "failed to parse uploader address")
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

// verify interface at compile time
var _ sdk.Msg = &MsgPrepay{}

// NewMsgPrepay NewMsg<Action> creates a new Msg<Action> instance
func NewMsgPrepay(sender string, coins sdk.Coins) *MsgPrepay {

	return &MsgPrepay{
		Sender: sender,
		Coins:  coins,
	}
}

// nolint
func (msg MsgPrepay) Route() string { return RouterKey }
func (msg MsgPrepay) Type() string  { return ConstSdsPrepay }
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
		return err
	}

	if sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.Coins.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "missing coins to send")
	}
	return nil
}
