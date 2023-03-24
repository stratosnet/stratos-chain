package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewFileInfo constructor
func NewFileInfo(height sdk.Int, reporters []byte, uploader string) FileInfo {
	return FileInfo{
		Height:    height,
		Reporters: reporters,
		Uploader:  uploader,
	}
}
