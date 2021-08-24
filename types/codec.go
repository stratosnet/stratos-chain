package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

const (
	// StAccountName is the amino encoding name for StAccount
	StAccountName = "stratos/StAccount"
)

// RegisterCodec registers the account interfaces and concrete types on the
// provided Amino codec.
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(&StAccount{}, StAccountName, nil)
}
