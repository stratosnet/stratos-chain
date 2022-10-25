package types

import "github.com/cosmos/cosmos-sdk/types"

// querier keys
const (
	QueryParams         = "params"
	QueryUploadedFile   = "uploaded_file"
	QuerySimulatePrepay = "simulate_prepay"
	QueryCurrNozPrice   = "curr_noz_price"
	QueryNozSupply      = "noz_supply"
)

type QueryUploadedFileParams struct {
	Sender types.AccAddress `json:"sender" yaml:"sender"`
}

// NewQueryUploadedFileParams creates a new instance of QueryValidatorSlashesParams
func NewQueryUploadedFileParams(sender types.AccAddress) QueryUploadedFileParams {
	return QueryUploadedFileParams{
		Sender: sender,
	}
}
