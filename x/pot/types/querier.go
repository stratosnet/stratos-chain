package types

import "github.com/cosmos/cosmos-sdk/types"

// querier keys
const (
	QueryParams           = "params"
	QueryVolumeReportHash = "volume_report"
)

// QueryVolumeReportParams for query 'custom/distr/validator_outstanding_rewards'
type QueryVolumeReportParams struct {
	Reporter types.AccAddress `json:"reporter" yaml:"reporter"`
}
