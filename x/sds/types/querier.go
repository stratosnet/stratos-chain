package types

// querier keys
const (
	QueryParams         = "params"
	QueryFileUpload     = "file_upload"
	QuerySimulatePrepay = "simulate_prepay"
	QueryCurrNozPrice   = "curr_noz_price"
	QueryNozSupply      = "noz_supply"
)

type QueryFileUploadParams struct {
	FileHash string `json:"file_hash" yaml:"file_hash"`
}

// NewQueryUploadedFileParams creates a new instance of QueryValidatorSlashesParams
func NewQueryFileUploadParams(fileHash string) QueryFileUploadParams {
	return QueryFileUploadParams{
		FileHash: fileHash,
	}
}
