package types

import sdk "github.com/cosmos/cosmos-sdk/types"

const (
	// USTOS defines the default coin denomination used in Stratos in:
	//
	// - Staking parameters: denomination used as stake in the dPoS chain
	// - Mint parameters: denomination minted due to fee distribution rewards
	// - Governance parameters: denomination used for spam prevention in proposal deposits
	// - Crisis parameters: constant fee denomination used for spam prevention to check broken invariant
	// - EVM parameters: denomination used for running EVM state transitions in Stratos.
	USTOS string = "ustos"

	// BaseDenomUnit defines the base denomination unit for Photons.
	// 1 photon = 1x10^{BaseDenomUnit} aphoton
	BaseDenomUnit = 18
)

// NewUstosSCoinInt64 is a utility function that returns an "aphoton" coin with the given int64 amount.
// The function will panic if the provided amount is negative.
func NewUstosSCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(USTOS, amount)
}
