package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ParamSubspace defines the expected Subspace interfacace
type ParamSubspace interface {
	SetUploadFile(ctx sdk.Context, key []byte, value string)
	GetUploadFile(ctx sdk.Context, key []byte) MsgFileUpload
}
