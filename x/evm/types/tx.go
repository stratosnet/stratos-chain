package types

import (
	bytes "bytes"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/vm"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	govv1beta1 "github.com/cosmos/cosmos-sdk/x/gov/types/v1beta1"
)

// Failed returns if the contract execution failed in vm errors
func (m *MsgEthereumTxResponse) Failed() bool {
	return len(m.VmError) > 0
}

// Return is a helper function to help caller distinguish between revert reason
// and function return. Return returns the data after execution if no error occurs.
func (m *MsgEthereumTxResponse) Return() []byte {
	if m.Failed() {
		return nil
	}
	return common.CopyBytes(m.Ret)
}

// Revert returns the concrete revert reason if the execution is aborted by `REVERT`
// opcode. Note the reason can be nil if no data supplied with revert opcode.
func (m *MsgEthereumTxResponse) Revert() []byte {
	if m.VmError != vm.ErrExecutionReverted.Error() {
		return nil
	}
	return common.CopyBytes(m.Ret)
}

// AsLegacyV0 converts to legacy structure before v012
func (m *MsgUpdateImplmentationProposal) AsLegacyV0() govv1beta1.Content {
	return NewUpdateImplmentationProposal(
		m.ProxyAddress,
		m.ImplementationAddress,
		m.Data,
		m.Amount,
	)
}

// GetSigners returns the expected signers for MsgUpdateImplmentationProposal.
func (m *MsgUpdateImplmentationProposal) GetSigners() []sdk.AccAddress {
	addr, _ := sdk.AccAddressFromBech32(m.Authority)
	return []sdk.AccAddress{addr}
}

// ValidateBasic runs basic stateless validity checks
func (m *MsgUpdateImplmentationProposal) ValidateBasic() error {
	err := govv1beta1.ValidateAbstract(m.AsLegacyV0())
	if err != nil {
		return err
	}

	if len(m.Authority) > 0 {
		if _, err := sdk.AccAddressFromBech32(m.Authority); err != nil {
			return sdkerrors.ErrInvalidAddress.Wrapf("invalid authority address: %s", err)
		}
	}

	if !common.IsHexAddress(m.ProxyAddress) {
		return fmt.Errorf("address '%s' is not valid", m.ProxyAddress)
	}

	if !common.IsHexAddress(m.ImplementationAddress) {
		return fmt.Errorf("address '%s' is not valid", m.ImplementationAddress)
	}

	if bytes.Equal(common.HexToAddress(m.ImplementationAddress).Bytes(), common.Address{}.Bytes()) {
		return fmt.Errorf("implementation address could not be zero address")
	}

	if m.Amount == nil {
		return fmt.Errorf("amount should be zero or greater")
	}

	return nil
}
