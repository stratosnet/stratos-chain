package types

// pot module event types
const (
	EventTypeVolumeReport = "VolumeReport"

	AttributeKeyReporter            = "reporter"
	AttributeKeyEpoch               = "epoch"
	AttributeKeyReportReferenceHash = "report_reference_hash"
	AttributeKeyNodesVolume         = "nodes_volume"

	// TODO: Create your event types
	// EventType<Action>    		= "action"

	// TODO: Create keys fo your events, the values will be derived from the msg
	// AttributeKeyAddress  		= "address"

	// TODO: Some events may not have values for that reason you want to emit that something happened.
	// AttributeValueDoubleSign = "double_sign"

	AttributeValueCategory = ModuleName
)
