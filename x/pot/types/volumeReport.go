package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

type ReportInfo struct {
	Epoch     sdk.Int
	Reference string
}

func NewReportInfo(epoch sdk.Int, reference string) ReportInfo {
	return ReportInfo{
		Epoch:     epoch,
		Reference: reference,
	}
}
