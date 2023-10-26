package types

import (
	"github.com/ipfs/go-cid"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// verify interface at compile time
var (
	_ sdk.Msg = &MsgFileUpload{}
	_ sdk.Msg = &MsgPrepay{}
	_ sdk.Msg = &MsgUpdateParams{}
)

const (
	TypeMsgFileUpload   = "FileUploadTx"
	TypeMsgPrepay       = "SdsPrepayTx"
	TypeMsgUpdateParams = "update_params"
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
		panic(err)
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
		return errors.Wrap(ErrInvalidFileHash, "failed to validate file hash")
	}

	from, err := sdk.AccAddressFromBech32(msg.GetFrom())
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "failed to validate from address")
	}

	reporter, err := stratos.SdsAddressFromBech32(msg.GetReporter())
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "failed to parse reporter address")
	}

	uploader, err := sdk.AccAddressFromBech32(msg.GetUploader())
	if err != nil {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "failed to parse uploader address")
	}

	if from.Empty() {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "missing address of tx from")
	}
	if reporter.Empty() {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "missing address of tx reporter")
	}
	if uploader.Empty() {
		return errors.Wrap(sdkerrors.ErrInvalidAddress, "missing address of file uploader")
	}
	if len(msg.FileHash) == 0 {
		return errors.Wrap(sdkerrors.ErrInvalidRequest, "missing file hash")
	}
	return nil
}

// --------------------------------------------------------------------------------------------------------------------

// NewMsgPrepay NewMsg<Action> creates a new Msg<Action> instance
func NewMsgPrepay(sender string, beneficiary string, amount sdk.Coins) *MsgPrepay {
	return &MsgPrepay{
		Sender:      sender,
		Beneficiary: beneficiary,
		Amount:      amount,
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
	_, err := sdk.AccAddressFromBech32(msg.GetSender())
	if err != nil {
		return ErrInvalidSenderAddr
	}

	_, err = sdk.AccAddressFromBech32(msg.GetBeneficiary())
	if err != nil {
		return ErrInvalidBeneficiaryAddr
	}

	if msg.Amount.Empty() {
		return errors.Wrap(sdkerrors.ErrInvalidCoins, "missing amount to send")
	}

	return nil
}

// --------------------------------------------------------------------------------------------------------------------

func NewMsgUpdateParams(params Params, authority string) *MsgUpdateParams {
	return &MsgUpdateParams{
		Params:    params,
		Authority: authority,
	}
}

// Route implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

// Type implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

// ValidateBasic implements sdk.Msg
func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return err
	}

	return msg.Params.Validate()
}

// GetSignBytes implements sdk.Msg
func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	authority := sdk.MustAccAddressFromBech32(msg.Authority)
	return []sdk.AccAddress{authority}
}

// GetSigners implements legacytx.LegacyMsg
func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}
