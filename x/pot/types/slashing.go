package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type Slashing struct {
	P2PAddress    sdk.AccAddress
	SlashingCoins sdk.Coins
}

func NewSlashing(p2pAddress sdk.AccAddress, slashingCoins sdk.Coins) Slashing {
	return Slashing{
		P2PAddress:    p2pAddress,
		SlashingCoins: slashingCoins,
	}
}
