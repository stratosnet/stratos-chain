package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/tendermint/tendermint/crypto"
)

// ensure Msg interface compliance at compile time
var (
	_ sdk.Msg = &MsgCreateResourceNode{}
	_ sdk.Msg = &MsgCreateIndexingNode{}
)

// NewMsgCreateResourceNode NewMsg<Action> creates a new Msg<Action> instance
func NewMsgCreateResourceNode(
	resourceNodeAddress sdk.ValAddress, pubKey crypto.PubKey, value sdk.Coin, delegatorAddress sdk.AccAddress,
) MsgCreateResourceNode {
	return MsgCreateResourceNode{
		ResourceNodeAddress: resourceNodeAddress,
		PubKey:              pubKey,
		Value:               value,
		DelegatorAddress:    delegatorAddress,
	}
}

type MsgCreateResourceNode struct {
	ResourceNodeAddress sdk.ValAddress `json:"resource_node_address" yaml:"resource_node_address"`
	PubKey              crypto.PubKey  `json:"pubkey" yaml:"pubkey"`
	Value               sdk.Coin       `json:"value" yaml:"value"`
	DelegatorAddress    sdk.AccAddress `json:"delegator_address" yaml:"delegator_address"`
}

func (msg MsgCreateResourceNode) Route() string {
	return RouterKey
}

func (msg MsgCreateResourceNode) Type() string {
	return "create_resource_node"
}

// ValidateBasic validity check for the CreateResourceNode
func (msg MsgCreateResourceNode) ValidateBasic() error {
	if msg.ResourceNodeAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing resource node address")
	}
	if msg.DelegatorAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing delegator address")
	}
	if !msg.Value.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "value is not positive")
	}
	return nil
}

func (msg MsgCreateResourceNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateResourceNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{sdk.AccAddress(msg.ResourceNodeAddress)}
}

type MsgCreateIndexingNode struct {
	IndexingNodeAddress sdk.ValAddress `json:"indexing_node_address" yaml:"indexing_node_address"`
	Sender              sdk.AccAddress `json:"sender" yaml:"sender"`
	PubKey              crypto.PubKey  `json:"pubkey" yaml:"pubkey"`
	Value               sdk.Coin       `json:"value" yaml:"value"`
}

func (msg MsgCreateIndexingNode) Route() string {
	return RouterKey
}

func (msg MsgCreateIndexingNode) Type() string {
	return "create_indexing_node"
}

func (msg MsgCreateIndexingNode) ValidateBasic() error {
	if msg.IndexingNodeAddress.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing indexing node address")
	}
	if msg.Sender.Empty() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if !msg.Value.IsPositive() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "value is not positive")
	}
	return nil
}

func (msg MsgCreateIndexingNode) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgCreateIndexingNode) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
