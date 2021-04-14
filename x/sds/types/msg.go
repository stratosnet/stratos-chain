package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

type MsgFileUpload struct {
	FileHash []byte         `json:"file_hash" yaml:"file_hash"` // hash of file
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`       // sender of tx
}

// verify interface at compile time
var _ sdk.Msg = &MsgFileUpload{}

// NewMsg<Action> creates a new Msg<Action> instance
func NewMsgUpload(fileHash []byte, sender sdk.AccAddress) MsgFileUpload {
	return MsgFileUpload{
		FileHash: fileHash,
		Sender:   sender,
	}
}

const Const = "FileUploadTx"

// nolint
func (msg MsgFileUpload) Route() string { return RouterKey }
func (msg MsgFileUpload) Type() string  { return Const }
func (msg MsgFileUpload) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

// GetSignBytes gets the bytes for the message signer to sign on
func (msg MsgFileUpload) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

// ValidateBasic validity check for the AnteHandler
func (msg MsgFileUpload) ValidateBasic() error {
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	return nil
}
