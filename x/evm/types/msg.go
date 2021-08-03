package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// message type and route constants
const (
	// TypeMsgEthermint defines the type string of Stratos message
	TypeMsgStratosTx = "stratos"
)

type MsgStratosTx struct {
	AccountNonce uint64          `json:"nonce"`
	Price        sdk.Int         `json:"gasPrice"`
	GasLimit     uint64          `json:"gas"`
	Recipient    *sdk.AccAddress `json:"to" rlp:"nil"` // nil means contract creation
	Amount       sdk.Int         `json:"value"`
	Payload      []byte          `json:"input"`

	// From address (formerly derived from signature)
	From sdk.AccAddress `json:"from"`
}

// NewMsgEthermint returns a reference to a new Ethermint transaction
func NewMsgStratosTx(
	nonce uint64, to *sdk.AccAddress, amount sdk.Int,
	gasLimit uint64, gasPrice sdk.Int, payload []byte, from sdk.AccAddress,
) MsgStratosTx {
	return MsgStratosTx{
		AccountNonce: nonce,
		Price:        gasPrice,
		GasLimit:     gasLimit,
		Recipient:    to,
		Amount:       amount,
		Payload:      payload,
		From:         from,
	}
}

func (msg MsgStratosTx) String() string {
	return fmt.Sprintf("nonce=%d gasPrice=%d gasLimit=%d recipient=%s amount=%d data=0x%x from=%s",
		msg.AccountNonce, msg.Price, msg.GasLimit, msg.Recipient, msg.Amount, msg.Payload, msg.From)
}

// Route should return the name of the module
func (msg MsgStratosTx) Route() string { return RouterKey }

// Type returns the action of the message
func (msg MsgStratosTx) Type() string { return TypeMsgStratosTx }

// GetSignBytes encodes the message for signing
func (msg MsgStratosTx) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// ValidateBasic runs stateless checks on the message
func (msg MsgStratosTx) ValidateBasic() error {
	if msg.Price.IsZero() {
		return sdkerrors.Wrapf(ErrInvalidPrice, "gas price cannot be 0")
	}

	if msg.Price.Sign() == -1 {
		return sdkerrors.Wrapf(ErrInvalidPrice, "gas price cannot be negative %s", msg.Price)
	}

	// Amount can be 0
	if msg.Amount.Sign() == -1 {
		return sdkerrors.Wrapf(ErrInvalidPrice, "amount cannot be negative %s", msg.Amount)
	}

	return nil
}

// GetSigners defines whose signature is required
func (msg MsgStratosTx) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.From}
}

// To returns the recipient address of the transaction. It returns nil if the
// transaction is a contract creation.
func (msg MsgStratosTx) To() *sdk.AccAddress {
	if msg.Recipient == nil {
		return nil
	}

	addr := sdk.AccAddress(msg.Recipient.Bytes())
	return &addr
}
