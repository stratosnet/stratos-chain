package types

// pot module event types
const (
	EventTypeVolumeReport      = "volume_report"
	EventTypeWithdraw          = "withdraw"
	EventTypeFoundationDeposit = "foundation_deposit"

	AttributeKeyEpoch              = "report_epoch"
	AttributeKeyReportReference    = "report_reference"
	AttributeKeyAmount             = "amount"
	AttributeKeyNodeAddress        = "node_address"
	AttributeKeyOwnerAddress       = "owner_address"
	AttributeKeyTotalConsumedOzone = "total_consumed_ozone"

	AttributeValueCategory = ModuleName
)
