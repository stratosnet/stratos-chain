package types

import "github.com/cosmos/cosmos-sdk/types"

// querier keys
const (
	QueryParams         = "params"
	QueryUploadedFile   = "uploaded_file"
	QueryPrepay         = "prepay"
	QuerySimulatePrepay = "simulate_prepay"
	QueryCurrUozPrice   = "curr_uoz_price"
	QueryUozSupply      = "uoz_supply"
)

// QueryUploadedFileParams for query 'custom/distr/validator_outstanding_rewards'
type QueryUploadedFileParams struct {
	Sender types.AccAddress `json:"sender" yaml:"sender"`
}

// NewQueryUploadedFileParams creates a new instance of QueryValidatorSlashesParams
func NewQueryUploadedFileParams(sender types.AccAddress) QueryUploadedFileParams {
	return QueryUploadedFileParams{
		Sender: sender,
	}
}
