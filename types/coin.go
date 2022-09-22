package types

const (
	/*
		1 ETH = 10^9 Gwei = 10^18 wei = 1 STOS = 10^9 nstos
		1 ETH = 1 STOS
		1 Gwei = 1 nstos
		todo: use nstos instead of ustos
	*/

	CoinType = 606

	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "USTOS"

	USTOS string = "ustos"

	//1 eth = 1x10^18 wei
	WeiDenomUnit = 18

	// BaseDenomUnit defines the base denomination unit for stos.
	// 1 stos = 1x10^{BaseDenomUnit} ustos
	BaseDenomUnit = 9

	// 1 ustos = 1x10^{WeiUstosUnitDiff} wei
	WeiUstosUnitDiff = WeiDenomUnit - BaseDenomUnit

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)
