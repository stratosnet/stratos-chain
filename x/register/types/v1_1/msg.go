package v1_1

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	_ sdk.Msg = &MsgCreateResourceNode{}
	_ sdk.Msg = &MsgUpdateResourceNode{}
	_ sdk.Msg = &MsgCreateMetaNode{}
	_ sdk.Msg = &MsgUpdateMetaNode{}
)

// ValidateBasic validity check for the CreateResourceNode
func (msg MsgCreateResourceNode) ValidateBasic() error {
	return nil
}

func (msg MsgCreateResourceNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateResourceNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pk)
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) ValidateBasic() error {
	return nil
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateResourceNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

func (msg MsgCreateMetaNode) ValidateBasic() error {
	return nil
}

func (msg MsgCreateMetaNode) GetSigners() []sdk.AccAddress {
	// Owner pays the tx fees
	addr, err := sdk.AccAddressFromBech32(msg.GetOwnerAddress())
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}

}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (msg MsgCreateMetaNode) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var pk cryptotypes.PubKey
	return unpacker.UnpackAny(msg.Pubkey, &pk)
}

// GetSigners implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) GetSigners() []sdk.AccAddress {
	addr, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{addr.Bytes()}
}

// ValidateBasic implements the sdk.Msg interface.
func (msg MsgUpdateMetaNode) ValidateBasic() error {
	return nil
}
