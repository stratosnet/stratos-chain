package types

import (
	"encoding/json"
	"strings"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	stratos "github.com/stratosnet/stratos-chain/types"
)

// NewGenesisState creates a new GenesisState object
func NewGenesisState(params Params,
	resourceNodes ResourceNodes,
	metaNodes MetaNodes,
	remainingNozLimit sdk.Int,
	slashingInfo []Slashing,
	depositNozRate sdk.Dec,
) *GenesisState {
	return &GenesisState{
		Params:            params,
		ResourceNodes:     resourceNodes,
		MetaNodes:         metaNodes,
		RemainingNozLimit: remainingNozLimit,
		Slashing:          slashingInfo,
		DepositNozRate:    depositNozRate,
	}
}

// DefaultGenesisState - default GenesisState used by Cosmos Hub
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Params:            DefaultParams(),
		ResourceNodes:     ResourceNodes{},
		MetaNodes:         MetaNodes{},
		RemainingNozLimit: DefaultRemainingNozLimit,
		Slashing:          make([]Slashing, 0),
		DepositNozRate:    DefaultDepositNozRate,
	}
}

// GetGenesisStateFromAppState returns x/auth GenesisState given raw application
// genesis state.
func GetGenesisStateFromAppState(cdc codec.JSONCodec, appState map[string]json.RawMessage) GenesisState {
	var genesisState GenesisState
	if appState[ModuleName] != nil {
		cdc.MustUnmarshalJSON(appState[ModuleName], &genesisState)
	}

	return genesisState
}

// ValidateGenesis validates the register genesis parameters
func ValidateGenesis(data GenesisState) error {
	if err := data.GetParams().Validate(); err != nil {
		return err
	}
	if err := data.GetResourceNodes().Validate(); err != nil {
		return err
	}
	if err := data.GetMetaNodes().Validate(); err != nil {
		return err
	}

	if (data.RemainingNozLimit).LT(sdk.ZeroInt()) {
		return ErrRemainingNozLimit
	}

	if (data.DepositNozRate).LTE(sdk.ZeroDec()) {
		return ErrDepositNozRate
	}
	return nil
}

func (v GenesisMetaNode) ToMetaNode() (MetaNode, error) {
	ownerAddress, err := sdk.AccAddressFromBech32(v.OwnerAddress)
	if err != nil {
		return MetaNode{}, sdkerrors.Wrap(ErrInvalidOwnerAddr, err.Error())
	}

	beneficiaryAddress := ownerAddress
	if len(strings.TrimSpace(v.BeneficiaryAddress)) > 0 {
		beneficiaryAddress, err = sdk.AccAddressFromBech32(v.BeneficiaryAddress)
		if err != nil {
			return MetaNode{}, sdkerrors.Wrap(ErrInvalidBeneficiaryAddr, err.Error())
		}
	}

	netAddr, err := stratos.SdsAddressFromBech32(v.GetNetworkAddress())
	if err != nil {
		return MetaNode{}, sdkerrors.Wrap(ErrInvalidNetworkAddr, err.Error())
	}

	tokens, err := sdk.ParseCoinsNormalized(v.Tokens)
	if err != nil {
		return MetaNode{}, sdkerrors.Wrap(ErrBadDenom, err.Error())
	}
	tokenAmt := tokens.AmountOf(DefaultBondDenom)

	return MetaNode{
		NetworkAddress:     netAddr.String(),
		Pubkey:             v.GetPubkey(),
		Suspend:            v.GetSuspend(),
		Status:             v.GetStatus(),
		Tokens:             tokenAmt,
		OwnerAddress:       ownerAddress.String(),
		BeneficiaryAddress: beneficiaryAddress.String(),
		Description:        v.GetDescription(),
	}, nil
}

func NewSlashing(walletAddress sdk.AccAddress, value sdk.Int) Slashing {
	return Slashing{
		WalletAddress: walletAddress.String(),
		Value:         value.Int64(),
	}
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (g GenesisState) UnpackInterfaces(c codectypes.AnyUnpacker) error {
	for i := range g.MetaNodes {
		if err := g.MetaNodes[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	for i := range g.ResourceNodes {
		if err := g.ResourceNodes[i].UnpackInterfaces(c); err != nil {
			return err
		}
	}
	return nil
}
