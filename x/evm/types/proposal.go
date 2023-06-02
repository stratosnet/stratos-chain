package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/ethereum/go-ethereum/common"
)

const (
	ProposalTypeUpdateImplementation = "UpdateImplementation"
)

var (
	_ govtypes.Content = &UpdateImplmentationProposal{}
)

func init() {
	govtypes.RegisterProposalType(ProposalTypeUpdateImplementation)
}

// NewUpdateImplmentationProposal creates a new proxy update proposal.
func NewUpdateImplmentationProposal(proxy, impl common.Address, data []byte, amt *sdk.Int) govtypes.Content {
	return &UpdateImplmentationProposal{
		ProxyAddress:          proxy.Hex(),
		ImplementationAddress: impl.Hex(),
		Data:                  data,
		Amount:                amt,
	}
}

// GetTitle returns the title of a new proxy update proposal.
func (epd *UpdateImplmentationProposal) GetTitle() string {
	return "New proxy upgrade function call"
}

// GetDescription returns the description of a new proxy update proposal.
func (epd *UpdateImplmentationProposal) GetDescription() string {
	return fmt.Sprintf(
		"This is upgrade for proxy '%s' address with a new implementation '%s'",
		epd.ProxyAddress, epd.ImplementationAddress,
	)
}

// ProposalRoute returns the routing key of a new proxy update proposal.
func (epd *UpdateImplmentationProposal) ProposalRoute() string { return RouterKey }

// ProposalType returns the type of a new proxy update proposal.
func (epd *UpdateImplmentationProposal) ProposalType() string {
	return ProposalTypeUpdateImplementation
}

// ValidateBasic runs basic stateless validity checks
func (epd *UpdateImplmentationProposal) ValidateBasic() error {
	err := govtypes.ValidateAbstract(epd)
	if err != nil {
		return err
	}

	if !common.IsHexAddress(epd.ProxyAddress) {
		return fmt.Errorf("address '%s' is not valid", epd.ProxyAddress)
	}

	if !common.IsHexAddress(epd.ImplementationAddress) {
		return fmt.Errorf("address '%s' is not valid", epd.ImplementationAddress)
	}

	if bytes.Equal(common.HexToAddress(epd.ImplementationAddress).Bytes(), common.Address{}.Bytes()) {
		return fmt.Errorf("implementation address could not be zero address")
	}

	if epd.Amount == nil {
		return fmt.Errorf("amount should be zero or greater")
	}

	return nil
}
