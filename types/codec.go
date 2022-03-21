package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	// StAccountName is the amino encoding name for StAccount
	StAccountName = "stratos/StAccount"
)

func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(SdsAddress{}, "SdsAddress", nil)
	cdc.RegisterConcrete(&StAccount{}, StAccountName, nil)
}
