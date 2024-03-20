package types

import (
	sdkmath "cosmossdk.io/math"
)

// NewFileInfo constructor
func NewFileInfo(height sdkmath.Int, reporters []byte, uploader string) FileInfo {
	return FileInfo{
		Height:    height,
		Reporters: reporters,
		Uploader:  uploader,
	}
}
