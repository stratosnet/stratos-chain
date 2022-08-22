package types

// pot module event types
const (
	EventTypeVolumeReport      = "volume_report"
	EventTypeWithdraw          = "withdraw"
	EventTypeFoundationDeposit = "foundation_deposit"
	EventTypeSlashing          = "slashing"

	AttributeKeyEpoch            = "epoch"
	AttributeKeyReportReference  = "report_reference"
	AttributeKeyAmount           = "amount"
	AttributeKeyWalletAddress    = "wallet_address"
	AttributeKeyNodeP2PAddress   = "p2p_address"
	AttributeKeySlashingNodeType = "slashing_type"
	AttributeKeyNodeSuspended    = "suspend"
	//AttributeKeyTotalConsumedOzone = "total_consumed_ozone"

	AttributeValueCategory = ModuleName
)
