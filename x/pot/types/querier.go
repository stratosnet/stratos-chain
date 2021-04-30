package types

import "github.com/cosmos/cosmos-sdk/types"

// querier keys
const (
	QueryParams           = "params"
	QueryVolumeReportHash = "volume_report"
)

// QueryVolumeReportHashParams for query 'custom/distr/validator_outstanding_rewards'
type QueryVolumeReportHashParams struct {
	Reporter types.AccAddress `json:"reporter" yaml:"reporter"`
}

// NewQueryVolumeReportParams creates a new instance of QueryVolumeReportParams
func NewQueryVolumeReportParams(reporter types.AccAddress) QueryVolumeReportHashParams {
	return QueryVolumeReportHashParams{
		Reporter: reporter,
	}
}
