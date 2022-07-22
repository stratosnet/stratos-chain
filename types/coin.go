package types

const (
	CoinType = 606

	// DisplayDenom defines the denomination displayed to users in client applications.
	DisplayDenom = "USTOS"

	USTOS string = "ustos"

	//1 eth = 1x10^18 wei
	EthDenomUnit = 18

	// BaseDenomUnit defines the base denomination unit for stos.
	// 1 stos = 1x10^{BaseDenomUnit} ustos
	BaseDenomUnit = 9

	// 1 ustos = 1x10^{DenomUnitDiff} wei
	DenomUnitDiff = EthDenomUnit - BaseDenomUnit

	// DefaultGasPrice is default gas price for evm transactions
	DefaultGasPrice = 20
)
