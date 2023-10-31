package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	CoinType = 606

	// Stos defines the denomination displayed to users in client applications.
	Stos  = "stos"
	Gwei  = "gwei"
	Wei   = "wei"
	Utros = "utros" // reward denom

	// WeiDenomUnit defines the base denomination unit for stos.
	// 1 stos = 1x10^{WeiDenomUnit} wei
	WeiDenomUnit  = 18
	GweiDenomUnit = 9

	StosToWei  = 1e18 // 1 Stos = 1e18 wei
	StosToGwei = 1e9  // 1 Stos = 1e9 Gwei
	GweiToWei  = 1e9  // 1 Gwei = 1e9 wei

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)

func NewCoin(amount sdkmath.Int) sdk.Coin {
	return sdk.NewCoin(Wei, amount)
}

func NewDecCoin(amount sdkmath.Int) sdk.DecCoin {
	return sdk.NewDecCoin(Wei, amount)
}

func NewCoinInt64(amount int64) sdk.Coin {
	return sdk.NewInt64Coin(Wei, amount)
}
