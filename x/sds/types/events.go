package types

// sds module event types
const (
	EventTypeFileUpload = "FileUpload"
	EventTypePrepay     = "Prepay"

	AttributeKeyReporter = "reporter"
	AttributeKeyFileHash = "file_hash"
	AttributeKeyUploader = "uploader"

	AttributeKeyAmount       = "amount"
	AttributeKeyPurchasedNoz = "purchased_noz"

	AttributeValueCategory = ModuleName
)
